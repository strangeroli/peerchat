package message

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Xelvra/peerchat/internal/crypto"
	"github.com/Xelvra/peerchat/internal/user"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/sirupsen/logrus"
)

const (
	// Protocol IDs for different message types
	MessageProtocolID = protocol.ID("/xelvra/message/1.0.0")
	FileProtocolID    = protocol.ID("/xelvra/file/1.0.0")
	GroupProtocolID   = protocol.ID("/xelvra/group/1.0.0")

	// Message limits
	MaxMessageSize = 64 * 1024         // 64KB max message size
	MaxFileSize    = 100 * 1024 * 1024 // 100MB max file size

	// Timeouts
	MessageTimeout = 30 * time.Second
	FileTimeout    = 5 * time.Minute
)

// MessageType represents different types of messages
type MessageType int

const (
	MessageTypeText MessageType = iota
	MessageTypeFile
	MessageTypeImage
	MessageTypeAudio
	MessageTypeVideo
	MessageTypeSystem
)

// String returns string representation of MessageType
func (mt MessageType) String() string {
	switch mt {
	case MessageTypeText:
		return "text"
	case MessageTypeFile:
		return "file"
	case MessageTypeImage:
		return "image"
	case MessageTypeAudio:
		return "audio"
	case MessageTypeVideo:
		return "video"
	case MessageTypeSystem:
		return "system"
	default:
		return "unknown"
	}
}

// Message represents a message in the Xelvra network
type Message struct {
	ID          string                 `json:"id"`
	Type        MessageType            `json:"type"`
	From        string                 `json:"from"` // Sender's DID
	To          string                 `json:"to"`   // Recipient's DID
	GroupID     string                 `json:"group_id,omitempty"`
	Content     []byte                 `json:"content"` // Encrypted content
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Signature   []byte                 `json:"signature"`
	IsEncrypted bool                   `json:"is_encrypted"`
}

// MessageManager handles message processing and routing
type MessageManager struct {
	host     host.Host
	identity *user.MessengerID
	logger   *logrus.Logger

	// Message storage and routing
	incomingMessages chan *Message
	outgoingMessages chan *Message
	messageHandlers  map[MessageType]MessageHandler

	// File transfer management
	fileTransferManager *FileTransferManager

	// Peer connections and encryption states
	peerSessions map[peer.ID]*crypto.DoubleRatchetState
	sessionMutex sync.RWMutex

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// MessageHandler defines the interface for handling different message types
type MessageHandler interface {
	HandleMessage(ctx context.Context, msg *Message) error
}

// NewMessageManager creates a new message manager
func NewMessageManager(h host.Host, identity *user.MessengerID, logger *logrus.Logger) *MessageManager {
	ctx, cancel := context.WithCancel(context.Background())

	mm := &MessageManager{
		host:                h,
		identity:            identity,
		logger:              logger,
		incomingMessages:    make(chan *Message, 100),
		outgoingMessages:    make(chan *Message, 100),
		messageHandlers:     make(map[MessageType]MessageHandler),
		fileTransferManager: NewFileTransferManager(logger),
		peerSessions:        make(map[peer.ID]*crypto.DoubleRatchetState),
		ctx:                 ctx,
		cancel:              cancel,
	}

	// Set up stream handlers
	h.SetStreamHandler(MessageProtocolID, mm.handleMessageStream)
	h.SetStreamHandler(FileProtocolID, mm.handleFileStream)
	h.SetStreamHandler(GroupProtocolID, mm.handleGroupStream)

	return mm
}

// Start begins message processing
func (mm *MessageManager) Start() error {
	mm.logger.Info("Starting MessageManager...")

	// Start message processing goroutines
	mm.logger.Debug("Adding goroutines to wait group...")
	mm.wg.Add(2)
	mm.logger.Debug("Starting processIncomingMessages goroutine...")
	go mm.processIncomingMessages()
	mm.logger.Debug("Starting processOutgoingMessages goroutine...")
	go mm.processOutgoingMessages()

	mm.logger.Info("MessageManager started successfully")
	return nil
}

// Stop gracefully stops the message manager
func (mm *MessageManager) Stop() error {
	mm.logger.Info("Stopping MessageManager...")

	mm.cancel()
	mm.wg.Wait()

	// Close channels
	close(mm.incomingMessages)
	close(mm.outgoingMessages)

	mm.logger.Info("MessageManager stopped successfully")
	return nil
}

// SendMessage sends a message to a peer
func (mm *MessageManager) SendMessage(to string, content []byte, msgType MessageType) error {
	// Create message
	msg := &Message{
		ID:          uuid.New().String(),
		Type:        msgType,
		From:        mm.identity.GetDID(),
		To:          to,
		Content:     content,
		Timestamp:   time.Now(),
		IsEncrypted: false,
	}

	// Sign the message
	if err := mm.signMessage(msg); err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	// Queue for sending
	select {
	case mm.outgoingMessages <- msg:
		return nil
	case <-mm.ctx.Done():
		return fmt.Errorf("message manager stopped")
	default:
		return fmt.Errorf("outgoing message queue full")
	}
}

// RegisterHandler registers a handler for a specific message type
func (mm *MessageManager) RegisterHandler(msgType MessageType, handler MessageHandler) {
	mm.messageHandlers[msgType] = handler
}

// processIncomingMessages processes incoming messages
func (mm *MessageManager) processIncomingMessages() {
	defer mm.wg.Done()

	for {
		select {
		case msg := <-mm.incomingMessages:
			if err := mm.handleIncomingMessage(msg); err != nil {
				mm.logger.WithError(err).Error("Failed to handle incoming message")
			}
		case <-mm.ctx.Done():
			return
		}
	}
}

// processOutgoingMessages processes outgoing messages
func (mm *MessageManager) processOutgoingMessages() {
	defer mm.wg.Done()

	for {
		select {
		case msg := <-mm.outgoingMessages:
			if err := mm.handleOutgoingMessage(msg); err != nil {
				mm.logger.WithError(err).Error("Failed to handle outgoing message")
			}
		case <-mm.ctx.Done():
			return
		}
	}
}

// handleIncomingMessage processes an incoming message
func (mm *MessageManager) handleIncomingMessage(msg *Message) error {
	mm.logger.WithFields(logrus.Fields{
		"message_id": msg.ID,
		"from":       msg.From,
		"type":       msg.Type.String(),
	}).Debug("Processing incoming message")

	// Verify message signature
	if !mm.verifyMessage(msg) {
		return fmt.Errorf("message signature verification failed")
	}

	// Decrypt message if encrypted
	if msg.IsEncrypted {
		if err := mm.decryptMessage(msg); err != nil {
			return fmt.Errorf("failed to decrypt message: %w", err)
		}
	}

	// Route to appropriate handler
	if handler, exists := mm.messageHandlers[msg.Type]; exists {
		return handler.HandleMessage(mm.ctx, msg)
	}

	// Default handling for unregistered message types
	mm.logger.WithField("type", msg.Type.String()).Warn("No handler registered for message type")
	return nil
}

// handleOutgoingMessage processes an outgoing message
func (mm *MessageManager) handleOutgoingMessage(msg *Message) error {
	mm.logger.WithFields(logrus.Fields{
		"message_id": msg.ID,
		"to":         msg.To,
		"type":       msg.Type.String(),
	}).Debug("Processing outgoing message")

	// For now, try to parse the recipient as a peer ID directly
	// TODO: Implement proper DID to Peer ID resolution
	recipientPeerID, err := peer.Decode(msg.To)
	if err != nil {
		mm.logger.WithError(err).Error("Failed to decode recipient peer ID")
		return fmt.Errorf("invalid recipient peer ID: %w", err)
	}

	// Open a stream to the recipient
	stream, err := mm.host.NewStream(context.Background(), recipientPeerID, MessageProtocolID)
	if err != nil {
		mm.logger.WithError(err).Error("Failed to open stream to recipient")
		return fmt.Errorf("failed to connect to recipient: %w", err)
	}
	defer stream.Close()

	// Serialize and send the message
	msgData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	// Write message length first (4 bytes)
	msgLen := uint32(len(msgData))
	lenBytes := make([]byte, 4)
	lenBytes[0] = byte(msgLen >> 24)
	lenBytes[1] = byte(msgLen >> 16)
	lenBytes[2] = byte(msgLen >> 8)
	lenBytes[3] = byte(msgLen)

	if _, err := stream.Write(lenBytes); err != nil {
		return fmt.Errorf("failed to write message length: %w", err)
	}

	// Write message data
	if _, err := stream.Write(msgData); err != nil {
		return fmt.Errorf("failed to write message data: %w", err)
	}

	mm.logger.WithFields(logrus.Fields{
		"message_id": msg.ID,
		"to":         msg.To,
		"size":       len(msgData),
	}).Info("Message sent successfully")

	return nil
}

// handleMessageStream handles incoming message streams
func (mm *MessageManager) handleMessageStream(stream network.Stream) {
	defer stream.Close()

	remotePeer := stream.Conn().RemotePeer()
	mm.logger.WithField("peer", remotePeer.String()).Debug("Handling message stream")

	// Read message length (4 bytes)
	lenBytes := make([]byte, 4)
	if _, err := stream.Read(lenBytes); err != nil {
		mm.logger.WithError(err).Error("Failed to read message length")
		return
	}

	msgLen := uint32(lenBytes[0])<<24 | uint32(lenBytes[1])<<16 | uint32(lenBytes[2])<<8 | uint32(lenBytes[3])
	if msgLen > MaxMessageSize {
		mm.logger.WithField("size", msgLen).Error("Message too large")
		return
	}

	// Read message data
	msgData := make([]byte, msgLen)
	if _, err := stream.Read(msgData); err != nil {
		mm.logger.WithError(err).Error("Failed to read message data")
		return
	}

	// Parse message
	var msg Message
	if err := json.Unmarshal(msgData, &msg); err != nil {
		mm.logger.WithError(err).Error("Failed to parse message")
		return
	}

	mm.logger.WithFields(logrus.Fields{
		"message_id": msg.ID,
		"from":       msg.From,
		"type":       msg.Type.String(),
		"size":       len(msgData),
	}).Info("Message received")

	// Queue message for processing
	select {
	case mm.incomingMessages <- &msg:
		// Message queued successfully
	case <-mm.ctx.Done():
		return
	default:
		mm.logger.Warn("Incoming message queue full, dropping message")
	}
}

// handleFileStream handles incoming file streams
func (mm *MessageManager) handleFileStream(stream network.Stream) {
	defer stream.Close()

	remotePeer := stream.Conn().RemotePeer()
	mm.logger.WithField("peer", remotePeer.String()).Debug("Handling file stream")

	// Handle file transfer protocol
	if err := mm.processFileTransferStream(stream, remotePeer); err != nil {
		mm.logger.WithError(err).Error("Failed to process file transfer stream")
	}
}

// handleGroupStream handles incoming group message streams
func (mm *MessageManager) handleGroupStream(stream network.Stream) {
	defer stream.Close()

	remotePeer := stream.Conn().RemotePeer()
	mm.logger.WithField("peer", remotePeer.String()).Debug("Handling group stream")

	// TODO: Implement group message handling
}

// signMessage signs a message with the identity key
func (mm *MessageManager) signMessage(msg *Message) error {
	// Serialize message for signing (excluding signature field)
	msgData, err := json.Marshal(struct {
		ID        string                 `json:"id"`
		Type      MessageType            `json:"type"`
		From      string                 `json:"from"`
		To        string                 `json:"to"`
		GroupID   string                 `json:"group_id,omitempty"`
		Content   []byte                 `json:"content"`
		Metadata  map[string]interface{} `json:"metadata,omitempty"`
		Timestamp time.Time              `json:"timestamp"`
	}{
		ID:        msg.ID,
		Type:      msg.Type,
		From:      msg.From,
		To:        msg.To,
		GroupID:   msg.GroupID,
		Content:   msg.Content,
		Metadata:  msg.Metadata,
		Timestamp: msg.Timestamp,
	})
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	// Sign the message
	signature, err := mm.identity.Sign(msgData)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	msg.Signature = signature
	return nil
}

// verifyMessage verifies a message signature
func (mm *MessageManager) verifyMessage(msg *Message) bool {
	// TODO: Get sender's public key from DID
	// TODO: Verify signature
	return true // Placeholder
}

// encryptMessage encrypts a message using Signal Protocol
func (mm *MessageManager) encryptMessage(msg *Message) error {
	// TODO: Implement Signal Protocol encryption
	return nil
}

// SendFile initiates a file transfer to a peer
func (mm *MessageManager) SendFile(peerID peer.ID, filePath string) error {
	mm.logger.WithFields(logrus.Fields{
		"peer_id":   peerID.String(),
		"file_path": filePath,
	}).Info("Initiating file transfer")

	// Open a stream to the peer for file transfer
	stream, err := mm.host.NewStream(context.Background(), peerID, FileProtocolID)
	if err != nil {
		return fmt.Errorf("failed to open file stream to peer: %w", err)
	}
	defer stream.Close()

	// Start file transfer
	return mm.fileTransferManager.StartFileTransfer(mm.ctx, stream, filePath, peerID)
}

// processFileTransferStream processes incoming file transfer streams
func (mm *MessageManager) processFileTransferStream(stream network.Stream, remotePeer peer.ID) error {
	mm.logger.WithField("peer", remotePeer.String()).Debug("Processing file transfer stream")

	// Read the initial request
	request, err := mm.readFileTransferRequest(stream)
	if err != nil {
		return fmt.Errorf("failed to read file transfer request: %w", err)
	}

	switch request.Type {
	case "request":
		return mm.handleFileTransferRequest(stream, remotePeer, request)
	case "chunk":
		return mm.handleFileChunk(stream, remotePeer, request)
	case "complete":
		return mm.handleFileComplete(stream, remotePeer, request)
	default:
		return fmt.Errorf("unknown file transfer request type: %s", request.Type)
	}
}

// readFileTransferRequest reads a file transfer request from stream
func (mm *MessageManager) readFileTransferRequest(stream network.Stream) (*FileTransferRequest, error) {
	// Read length prefix
	var length uint32
	if err := binary.Read(stream, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("failed to read length: %w", err)
	}

	if length > FileHeaderSize {
		return nil, fmt.Errorf("request too large: %d bytes", length)
	}

	// Read data
	data := make([]byte, length)
	if _, err := io.ReadFull(stream, data); err != nil {
		return nil, fmt.Errorf("failed to read request data: %w", err)
	}

	// Parse request
	var request FileTransferRequest
	if err := json.Unmarshal(data, &request); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	if request.Magic != FileTransferMagic {
		return nil, fmt.Errorf("invalid magic number: %x", request.Magic)
	}

	return &request, nil
}

// handleFileTransferRequest handles incoming file transfer requests
func (mm *MessageManager) handleFileTransferRequest(stream network.Stream, remotePeer peer.ID, request *FileTransferRequest) error {
	mm.logger.WithFields(logrus.Fields{
		"peer":      remotePeer.String(),
		"file_name": request.Metadata.Name,
		"file_size": request.Metadata.Size,
	}).Info("Received file transfer request")

	// For now, automatically accept all file transfers
	// In production, this would prompt the user or check policies
	response := FileTransferRequest{
		Magic: FileTransferMagic,
		Type:  "accept",
	}

	// Send acceptance response
	if err := mm.sendFileTransferResponse(stream, response); err != nil {
		return fmt.Errorf("failed to send acceptance: %w", err)
	}

	// Create download directory if it doesn't exist
	downloadDir := filepath.Join(os.Getenv("HOME"), ".xelvra", "downloads")
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return fmt.Errorf("failed to create download directory: %w", err)
	}

	// Create file transfer session for receiving
	transfer := NewFileTransfer(request.Metadata.ID, remotePeer, request.Metadata, false, mm.logger)
	mm.fileTransferManager.transfers[request.Metadata.ID] = transfer

	// Create destination file
	destPath := filepath.Join(downloadDir, request.Metadata.Name)
	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}

	transfer.file = file
	transfer.Status = FileTransferActive

	mm.logger.WithFields(logrus.Fields{
		"transfer_id": transfer.ID,
		"dest_path":   destPath,
	}).Info("File transfer accepted, ready to receive")

	return nil
}

// handleFileChunk handles incoming file chunks
func (mm *MessageManager) handleFileChunk(stream network.Stream, remotePeer peer.ID, request *FileTransferRequest) error {
	// Find the active transfer (simplified - would need better lookup)
	var transfer *FileTransfer
	for _, t := range mm.fileTransferManager.transfers {
		if t.PeerID == remotePeer && t.Status == FileTransferActive && !t.isOutgoing {
			transfer = t
			break
		}
	}

	if transfer == nil {
		return fmt.Errorf("no active file transfer found for peer %s", remotePeer.String())
	}

	// Write chunk to file
	if _, err := transfer.file.Write(request.Data); err != nil {
		transfer.Status = FileTransferFailed
		transfer.Error = err
		return fmt.Errorf("failed to write chunk: %w", err)
	}

	transfer.BytesReceived += int64(len(request.Data))
	transfer.UpdateProgress()
	transfer.chunks[request.ChunkID] = true

	mm.logger.WithFields(logrus.Fields{
		"transfer_id": transfer.ID,
		"chunk_id":    request.ChunkID,
		"chunk_size":  len(request.Data),
		"progress":    fmt.Sprintf("%.1f%%", transfer.Progress*100),
	}).Debug("Received file chunk")

	return nil
}

// handleFileComplete handles file transfer completion
func (mm *MessageManager) handleFileComplete(stream network.Stream, remotePeer peer.ID, request *FileTransferRequest) error {
	// Find the active transfer
	var transfer *FileTransfer
	for _, t := range mm.fileTransferManager.transfers {
		if t.PeerID == remotePeer && t.Status == FileTransferActive && !t.isOutgoing {
			transfer = t
			break
		}
	}

	if transfer == nil {
		return fmt.Errorf("no active file transfer found for peer %s", remotePeer.String())
	}

	// Close the file
	if err := transfer.file.Close(); err != nil {
		mm.logger.WithError(err).Warn("Failed to close received file")
	}

	transfer.Status = FileTransferCompleted
	transfer.EndTime = time.Now()

	mm.logger.WithFields(logrus.Fields{
		"transfer_id":    transfer.ID,
		"file_name":      transfer.Metadata.Name,
		"bytes_received": transfer.BytesReceived,
		"duration":       transfer.EndTime.Sub(transfer.StartTime),
	}).Info("File transfer completed successfully")

	// TODO: Verify file hash
	// TODO: Send completion acknowledgment

	return nil
}

// sendFileTransferResponse sends a file transfer response
func (mm *MessageManager) sendFileTransferResponse(stream network.Stream, response FileTransferRequest) error {
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	// Write length prefix
	length := uint32(len(data))
	if err := binary.Write(stream, binary.BigEndian, length); err != nil {
		return fmt.Errorf("failed to write length: %w", err)
	}

	// Write data
	if _, err := stream.Write(data); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

// decryptMessage decrypts a message using Signal Protocol
func (mm *MessageManager) decryptMessage(msg *Message) error {
	// TODO: Implement Signal Protocol decryption
	return nil
}
