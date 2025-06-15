package p2p

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Xelvra/peerchat/internal/message"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sirupsen/logrus"
)

// P2PWrapper provides a safe interface to P2P functionality
// It can fallback to simulation if real P2P fails
type P2PWrapper struct {
	useSimulation bool
	realNode      *PeerChatNode
	ctx           context.Context
	logger        *logrus.Logger
}

// NodeInfo contains basic node information
type NodeInfo struct {
	PeerID      string
	DID         string
	ListenAddrs []string
	IsRunning   bool
}

// NewP2PWrapper creates a new P2P wrapper
func NewP2PWrapper(ctx context.Context, useSimulation bool) *P2PWrapper {
	logger := setupLogger()
	return &P2PWrapper{
		useSimulation: useSimulation,
		ctx:           ctx,
		logger:        logger,
	}
}

// setupLogger configures logging to file with rotation
func setupLogger() *logrus.Logger {
	logger := logrus.New()

	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home not available
		logger.SetOutput(os.Stderr)
		return logger
	}

	// Create .xelvra directory if it doesn't exist
	xelvraDir := filepath.Join(home, ".xelvra")
	if err := os.MkdirAll(xelvraDir, 0700); err != nil {
		logger.SetOutput(os.Stderr)
		return logger
	}

	// Setup log rotation
	logFile := filepath.Join(xelvraDir, "peerchat.log")
	if err := rotateLogIfNeeded(logFile); err != nil {
		// If rotation fails, continue with stderr
		logger.SetOutput(os.Stderr)
		return logger
	}

	// Open log file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		logger.SetOutput(os.Stderr)
		return logger
	}

	// Set JSON formatter for structured logging
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000000000Z07:00",
	})

	logger.SetOutput(file)
	logger.SetLevel(logrus.InfoLevel)

	return logger
}

// Start starts the P2P node (real or simulated)
func (w *P2PWrapper) Start() error {
	if w.useSimulation {
		return w.startSimulation()
	}
	return w.startReal()
}

// Stop stops the P2P node
func (w *P2PWrapper) Stop() error {
	if w.useSimulation {
		return nil // Nothing to stop in simulation
	}
	if w.realNode != nil {
		return w.realNode.Stop()
	}
	return nil
}

// GetNodeInfo returns basic node information
func (w *P2PWrapper) GetNodeInfo() *NodeInfo {
	if w.useSimulation {
		return &NodeInfo{
			PeerID:      "12D3KooWSimulatedPeerID...",
			DID:         "did:xelvra:simulated...",
			ListenAddrs: []string{"/ip4/127.0.0.1/tcp/4001"},
			IsRunning:   true,
		}
	}

	if w.realNode == nil {
		return &NodeInfo{
			IsRunning: false,
		}
	}

	// Get real node info
	addrs := make([]string, 0)
	for _, addr := range w.realNode.GetHost().Addrs() {
		addrs = append(addrs, addr.String())
	}

	return &NodeInfo{
		PeerID:      w.realNode.GetPeerID().String(),
		DID:         w.realNode.GetIdentity().GetDID(),
		ListenAddrs: addrs,
		IsRunning:   true,
	}
}

// SendMessage sends a message to a peer
func (w *P2PWrapper) SendMessage(peerID, messageText string) error {
	if w.useSimulation {
		// Simulate message sending
		w.logger.WithFields(logrus.Fields{
			"peer_id": peerID,
			"message": messageText,
		}).Info("Simulated message sent")
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	if w.realNode == nil {
		return fmt.Errorf("node not started")
	}

	// Send real message using the message manager
	w.logger.WithFields(logrus.Fields{
		"peer_id": peerID,
		"message": messageText,
	}).Info("Sending real message")

	// Send message using proper message type
	return w.realNode.SendMessage(peerID, []byte(messageText), message.MessageTypeText)
}

// startSimulation starts simulation mode
func (w *P2PWrapper) startSimulation() error {
	// Simulate startup delay
	time.Sleep(200 * time.Millisecond)
	return nil
}

// startReal starts real P2P node
func (w *P2PWrapper) startReal() error {
	w.logger.Info("Attempting to start real P2P node...")

	// Create context with timeout to prevent hanging
	ctx, cancel := context.WithTimeout(w.ctx, 2*time.Second)
	defer cancel()

	// Try to create real P2P node with timeout
	config := DefaultNodeConfig()
	config.LogLevel = w.logger.Level // Use our log level
	config.Logger = w.logger         // Use our file logger

	// Use a channel to handle timeout
	type result struct {
		node *PeerChatNode
		err  error
	}

	resultChan := make(chan result, 1)

	go func() {
		node, err := NewPeerChatNode(ctx, config)
		resultChan <- result{node: node, err: err}
	}()

	select {
	case res := <-resultChan:
		if res.err != nil {
			w.logger.WithError(res.err).Warn("Failed to create real P2P node, falling back to simulation")
			// Fallback to simulation if real P2P fails
			w.useSimulation = true
			return w.startSimulation()
		}

		// Try to start the node with timeout
		startChan := make(chan error, 1)
		go func() {
			startChan <- res.node.Start()
		}()

		select {
		case err := <-startChan:
			if err != nil {
				w.logger.WithError(err).Warn("Failed to start real P2P node, falling back to simulation")
				// Fallback to simulation if start fails
				w.useSimulation = true
				return w.startSimulation()
			}
			w.logger.Info("Real P2P node started successfully")
			w.realNode = res.node
			return nil
		case <-ctx.Done():
			w.logger.Warn("P2P node start timed out, falling back to simulation")
			w.useSimulation = true
			return w.startSimulation()
		}

	case <-ctx.Done():
		w.logger.Warn("P2P node creation timed out, falling back to simulation")
		w.useSimulation = true
		return w.startSimulation()
	}
}

// IsUsingSimulation returns true if using simulation mode
func (w *P2PWrapper) IsUsingSimulation() bool {
	return w.useSimulation
}

// GetDiscoveredPeers returns list of discovered peers
func (w *P2PWrapper) GetDiscoveredPeers() []string {
	if w.useSimulation {
		return []string{} // No peers in simulation
	}

	if w.realNode == nil || w.realNode.discoveryManager == nil {
		return []string{}
	}

	peers := w.realNode.discoveryManager.GetDiscoveredPeers()
	result := make([]string, len(peers))
	for i, peer := range peers {
		result[i] = peer.String()
	}
	return result
}

// GetConnectedPeers returns list of currently connected peers
func (w *P2PWrapper) GetConnectedPeers() []string {
	if w.useSimulation {
		return []string{} // No peers in simulation
	}

	if w.realNode == nil {
		return []string{}
	}

	// Get connected peers from libp2p host
	connectedPeers := w.realNode.host.Network().Peers()
	result := make([]string, len(connectedPeers))
	for i, peer := range connectedPeers {
		result[i] = peer.String()
	}
	return result
}

// ConnectToPeer attempts to connect to a specific peer
func (w *P2PWrapper) ConnectToPeer(peerIDStr string) bool {
	if w.useSimulation {
		return false // Cannot connect in simulation
	}

	if w.realNode == nil {
		return false
	}

	// Parse peer ID
	peerID, err := peer.Decode(peerIDStr)
	if err != nil {
		w.logger.WithError(err).Error("Invalid peer ID format")
		return false
	}

	// Get peer addresses from discovery manager
	peerAddrs := w.realNode.discoveryManager.GetPeerAddresses(peerID)
	if len(peerAddrs) == 0 {
		w.logger.WithField("peer_id", peerIDStr).Warn("No addresses found for peer")
		return false
	}

	// Create peer info with addresses
	peerInfo := peer.AddrInfo{
		ID:    peerID,
		Addrs: peerAddrs,
	}

	// Try to connect
	ctx, cancel := context.WithTimeout(w.ctx, 10*time.Second)
	defer cancel()

	if err := w.realNode.host.Connect(ctx, peerInfo); err != nil {
		w.logger.WithError(err).WithField("peer_id", peerIDStr).Error("Failed to connect to peer")
		return false
	}

	w.logger.WithField("peer_id", peerIDStr).Info("Successfully connected to peer")
	return true
}

// SendMessageToMultiplePeers sends a message to specified peers
func (w *P2PWrapper) SendMessageToMultiplePeers(message string, peerIDs []string) bool {
	if w.useSimulation {
		return false // Cannot send in simulation
	}

	if w.realNode == nil {
		return false
	}

	success := true
	for _, peerIDStr := range peerIDs {
		// Use the existing SendMessage method
		if err := w.SendMessage(peerIDStr, message); err != nil {
			w.logger.WithError(err).WithField("peer_id", peerIDStr).Error("Failed to send message")
			success = false
		} else {
			w.logger.WithField("peer_id", peerIDStr).WithField("message", message).Info("Message sent successfully")
		}
	}

	return success
}

// rotateLogIfNeeded checks if log rotation is needed and performs it
func rotateLogIfNeeded(logFile string) error {
	// Check if log file exists
	info, err := os.Stat(logFile)
	if os.IsNotExist(err) {
		return nil // No rotation needed for non-existent file
	}
	if err != nil {
		return err
	}

	const maxSizeBytes = 5 * 1024 * 1024 // 5MB
	const maxLines = 10000               // 10,000 lines

	// Check file size
	needsRotation := info.Size() > maxSizeBytes

	// Check line count if size check didn't trigger rotation
	if !needsRotation {
		lineCount, err := countLines(logFile)
		if err == nil && lineCount > maxLines {
			needsRotation = true
		}
	}

	if !needsRotation {
		return nil
	}

	// Perform rotation
	return performLogRotation(logFile)
}

// countLines counts the number of lines in a file
func countLines(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			// Log error but don't fail the function - line counting can continue
			_ = err // Explicitly ignore error
		}
	}()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	return lineCount, scanner.Err()
}

// performLogRotation rotates log files (keeps 3 old versions)
func performLogRotation(logFile string) error {
	const maxBackups = 3

	// Remove oldest backup if it exists
	oldestBackup := logFile + "." + strconv.Itoa(maxBackups)
	if _, err := os.Stat(oldestBackup); err == nil {
		if err := os.Remove(oldestBackup); err != nil {
			// Log error but continue with rotation - backup cleanup is not critical
			_ = err // Explicitly ignore error
		}
	}

	// Shift existing backups
	for i := maxBackups - 1; i >= 1; i-- {
		oldName := logFile + "." + strconv.Itoa(i)
		newName := logFile + "." + strconv.Itoa(i+1)

		if _, err := os.Stat(oldName); err == nil {
			if err := os.Rename(oldName, newName); err != nil {
				// Log error but continue with rotation
				fmt.Printf("Warning: Failed to rotate log backup %s to %s: %v\n", oldName, newName, err)
			}
		}
	}

	// Move current log to .1
	backup := logFile + ".1"
	if err := os.Rename(logFile, backup); err != nil {
		return fmt.Errorf("failed to rotate log: %v", err)
	}

	return nil
}
