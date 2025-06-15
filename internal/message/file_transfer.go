package message

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sirupsen/logrus"
)

const (
	// File transfer protocol constants
	FileChunkSize     = 32 * 1024 // 32KB chunks for optimal performance
	FileHeaderSize    = 1024      // Maximum size for file metadata header
	FileTransferMagic = 0x58454C56 // "XELV" magic number for file transfers
)

// FileTransferStatus represents the status of a file transfer
type FileTransferStatus int

const (
	FileTransferPending FileTransferStatus = iota
	FileTransferActive
	FileTransferCompleted
	FileTransferFailed
	FileTransferCancelled
)

// String returns string representation of FileTransferStatus
func (fts FileTransferStatus) String() string {
	switch fts {
	case FileTransferPending:
		return "pending"
	case FileTransferActive:
		return "active"
	case FileTransferCompleted:
		return "completed"
	case FileTransferFailed:
		return "failed"
	case FileTransferCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

// FileMetadata contains information about a file being transferred
type FileMetadata struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	Hash        string    `json:"hash"`        // SHA256 hash for integrity verification
	MimeType    string    `json:"mime_type"`
	Timestamp   time.Time `json:"timestamp"`
	ChunkCount  int       `json:"chunk_count"`
	ChunkSize   int       `json:"chunk_size"`
}

// FileTransferRequest represents a file transfer request
type FileTransferRequest struct {
	Magic    uint32       `json:"magic"`
	Type     string       `json:"type"` // "request", "accept", "reject", "chunk", "complete"
	Metadata FileMetadata `json:"metadata,omitempty"`
	ChunkID  int          `json:"chunk_id,omitempty"`
	Data     []byte       `json:"data,omitempty"`
	Error    string       `json:"error,omitempty"`
}

// FileTransfer represents an active file transfer session
type FileTransfer struct {
	ID           string
	PeerID       peer.ID
	Metadata     FileMetadata
	Status       FileTransferStatus
	Progress     float64 // 0.0 to 1.0
	BytesTotal   int64
	BytesSent    int64
	BytesReceived int64
	StartTime    time.Time
	EndTime      time.Time
	Error        error
	
	// File handling
	file         *os.File
	isOutgoing   bool
	chunks       map[int]bool // Track received chunks
	logger       *logrus.Logger
}

// NewFileTransfer creates a new file transfer session
func NewFileTransfer(id string, peerID peer.ID, metadata FileMetadata, isOutgoing bool, logger *logrus.Logger) *FileTransfer {
	return &FileTransfer{
		ID:         id,
		PeerID:     peerID,
		Metadata:   metadata,
		Status:     FileTransferPending,
		Progress:   0.0,
		BytesTotal: metadata.Size,
		StartTime:  time.Now(),
		isOutgoing: isOutgoing,
		chunks:     make(map[int]bool),
		logger:     logger,
	}
}

// CalculateFileHash calculates SHA256 hash of a file
func CalculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// CreateFileMetadata creates metadata for a file
func CreateFileMetadata(filePath string) (*FileMetadata, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	hash, err := CalculateFileHash(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate file hash: %w", err)
	}

	chunkCount := int((fileInfo.Size() + FileChunkSize - 1) / FileChunkSize)

	metadata := &FileMetadata{
		ID:         fmt.Sprintf("file_%d", time.Now().UnixNano()),
		Name:       filepath.Base(filePath),
		Size:       fileInfo.Size(),
		Hash:       hash,
		MimeType:   detectMimeType(filePath),
		Timestamp:  fileInfo.ModTime(),
		ChunkCount: chunkCount,
		ChunkSize:  FileChunkSize,
	}

	return metadata, nil
}

// detectMimeType detects MIME type based on file extension
func detectMimeType(filePath string) string {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".txt":
		return "text/plain"
	case ".pdf":
		return "application/pdf"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".mp4":
		return "video/mp4"
	case ".mp3":
		return "audio/mpeg"
	case ".zip":
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}

// SendFileRequest sends a file transfer request
func (ft *FileTransfer) SendFileRequest(stream network.Stream) error {
	request := FileTransferRequest{
		Magic:    FileTransferMagic,
		Type:     "request",
		Metadata: ft.Metadata,
	}

	return ft.sendRequest(stream, request)
}

// SendFileChunk sends a file chunk
func (ft *FileTransfer) SendFileChunk(stream network.Stream, chunkID int, data []byte) error {
	request := FileTransferRequest{
		Magic:   FileTransferMagic,
		Type:    "chunk",
		ChunkID: chunkID,
		Data:    data,
	}

	return ft.sendRequest(stream, request)
}

// SendFileComplete sends file transfer completion notification
func (ft *FileTransfer) SendFileComplete(stream network.Stream) error {
	request := FileTransferRequest{
		Magic: FileTransferMagic,
		Type:  "complete",
	}

	return ft.sendRequest(stream, request)
}

// sendRequest sends a file transfer request over the stream
func (ft *FileTransfer) sendRequest(stream network.Stream, request FileTransferRequest) error {
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Write length prefix (4 bytes)
	length := uint32(len(data))
	if err := binary.Write(stream, binary.BigEndian, length); err != nil {
		return fmt.Errorf("failed to write length: %w", err)
	}

	// Write data
	if _, err := stream.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	return nil
}

// UpdateProgress updates the transfer progress
func (ft *FileTransfer) UpdateProgress() {
	if ft.BytesTotal > 0 {
		if ft.isOutgoing {
			ft.Progress = float64(ft.BytesSent) / float64(ft.BytesTotal)
		} else {
			ft.Progress = float64(ft.BytesReceived) / float64(ft.BytesTotal)
		}
	}
}

// Close closes the file transfer and cleans up resources
func (ft *FileTransfer) Close() error {
	if ft.file != nil {
		return ft.file.Close()
	}
	return nil
}

// IsOutgoing returns true if this is an outgoing file transfer
func (ft *FileTransfer) IsOutgoing() bool {
	return ft.isOutgoing
}

// FileTransferManager manages file transfers
type FileTransferManager struct {
	transfers map[string]*FileTransfer
	logger    *logrus.Logger
}

// NewFileTransferManager creates a new file transfer manager
func NewFileTransferManager(logger *logrus.Logger) *FileTransferManager {
	return &FileTransferManager{
		transfers: make(map[string]*FileTransfer),
		logger:    logger,
	}
}

// StartFileTransfer initiates a file transfer
func (ftm *FileTransferManager) StartFileTransfer(ctx context.Context, stream network.Stream, filePath string, peerID peer.ID) error {
	// Create file metadata
	metadata, err := CreateFileMetadata(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file metadata: %w", err)
	}

	// Create file transfer session
	transfer := NewFileTransfer(metadata.ID, peerID, *metadata, true, ftm.logger)
	ftm.transfers[metadata.ID] = transfer

	ftm.logger.WithFields(logrus.Fields{
		"transfer_id": metadata.ID,
		"file_name":   metadata.Name,
		"file_size":   metadata.Size,
		"peer_id":     peerID.String(),
	}).Info("Starting file transfer")

	// Send file transfer request
	if err := transfer.SendFileRequest(stream); err != nil {
		transfer.Status = FileTransferFailed
		transfer.Error = err
		return fmt.Errorf("failed to send file request: %w", err)
	}

	// Wait for response (simplified - in production would be async)
	response, err := ftm.readResponse(stream)
	if err != nil {
		transfer.Status = FileTransferFailed
		transfer.Error = err
		return fmt.Errorf("failed to read response: %w", err)
	}

	if response.Type == "reject" {
		transfer.Status = FileTransferFailed
		transfer.Error = fmt.Errorf("file transfer rejected: %s", response.Error)
		return transfer.Error
	}

	if response.Type != "accept" {
		transfer.Status = FileTransferFailed
		transfer.Error = fmt.Errorf("unexpected response type: %s", response.Type)
		return transfer.Error
	}

	// Start sending file chunks
	return ftm.sendFileChunks(ctx, stream, transfer, filePath)
}

// sendFileChunks sends file data in chunks
func (ftm *FileTransferManager) sendFileChunks(ctx context.Context, stream network.Stream, transfer *FileTransfer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	transfer.file = file
	transfer.Status = FileTransferActive

	buffer := make([]byte, FileChunkSize)
	chunkID := 0

	for {
		select {
		case <-ctx.Done():
			transfer.Status = FileTransferCancelled
			return ctx.Err()
		default:
		}

		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			transfer.Status = FileTransferFailed
			transfer.Error = err
			return fmt.Errorf("failed to read file chunk: %w", err)
		}

		// Send chunk
		chunkData := buffer[:n]
		if err := transfer.SendFileChunk(stream, chunkID, chunkData); err != nil {
			transfer.Status = FileTransferFailed
			transfer.Error = err
			return fmt.Errorf("failed to send chunk %d: %w", chunkID, err)
		}

		transfer.BytesSent += int64(n)
		transfer.UpdateProgress()

		ftm.logger.WithFields(logrus.Fields{
			"transfer_id": transfer.ID,
			"chunk_id":    chunkID,
			"chunk_size":  n,
			"progress":    fmt.Sprintf("%.1f%%", transfer.Progress*100),
		}).Debug("Sent file chunk")

		chunkID++
	}

	// Send completion notification
	if err := transfer.SendFileComplete(stream); err != nil {
		transfer.Status = FileTransferFailed
		transfer.Error = err
		return fmt.Errorf("failed to send completion: %w", err)
	}

	transfer.Status = FileTransferCompleted
	transfer.EndTime = time.Now()

	ftm.logger.WithFields(logrus.Fields{
		"transfer_id": transfer.ID,
		"file_name":   transfer.Metadata.Name,
		"bytes_sent":  transfer.BytesSent,
		"duration":    transfer.EndTime.Sub(transfer.StartTime),
	}).Info("File transfer completed successfully")

	return nil
}

// readResponse reads a file transfer response from the stream
func (ftm *FileTransferManager) readResponse(stream network.Stream) (*FileTransferRequest, error) {
	// Read length prefix
	var length uint32
	if err := binary.Read(stream, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("failed to read length: %w", err)
	}

	if length > FileHeaderSize {
		return nil, fmt.Errorf("response too large: %d bytes", length)
	}

	// Read data
	data := make([]byte, length)
	if _, err := io.ReadFull(stream, data); err != nil {
		return nil, fmt.Errorf("failed to read response data: %w", err)
	}

	// Parse response
	var response FileTransferRequest
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Magic != FileTransferMagic {
		return nil, fmt.Errorf("invalid magic number: %x", response.Magic)
	}

	return &response, nil
}

// GetTransfer returns a file transfer by ID
func (ftm *FileTransferManager) GetTransfer(id string) (*FileTransfer, bool) {
	transfer, exists := ftm.transfers[id]
	return transfer, exists
}

// ListTransfers returns all active transfers
func (ftm *FileTransferManager) ListTransfers() []*FileTransfer {
	transfers := make([]*FileTransfer, 0, len(ftm.transfers))
	for _, transfer := range ftm.transfers {
		transfers = append(transfers, transfer)
	}
	return transfers
}

// CleanupTransfer removes a completed or failed transfer
func (ftm *FileTransferManager) CleanupTransfer(id string) {
	if transfer, exists := ftm.transfers[id]; exists {
		transfer.Close()
		delete(ftm.transfers, id)
	}
}
