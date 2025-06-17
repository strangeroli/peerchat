package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Xelvra/peerchat/internal/message"
	"github.com/Xelvra/peerchat/internal/user"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/core/routing"
	libp2pquic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/sirupsen/logrus"
)

const (
	// Protocol ID for Xelvra messaging
	XelvraProtocolID = protocol.ID("/xelvra/1.0.0")

	// Performance targets from README
	MaxIdleMemoryMB   = 20
	MaxIdleCPUPercent = 1
	MaxLatencyMs      = 50
)

// NetworkTransport represents active network transport information
type NetworkTransport struct {
	Type       string  `json:"type"` // "tcp", "quic", "relay"
	LocalAddr  string  `json:"local_addr"`
	RemoteAddr string  `json:"remote_addr,omitempty"`
	IsActive   bool    `json:"is_active"`
	Latency    int     `json:"latency_ms,omitempty"`
	PacketLoss float64 `json:"packet_loss,omitempty"`
}

// NATInfo represents NAT traversal information
type NATInfo struct {
	Type        string   `json:"type"` // "none", "full_cone", "restricted", "port_restricted", "symmetric"
	PublicIP    string   `json:"public_ip"`
	PublicPort  int      `json:"public_port"`
	LocalIP     string   `json:"local_ip"`
	LocalPort   int      `json:"local_port"`
	STUNServers []string `json:"stun_servers"`
	UsingRelay  bool     `json:"using_relay"`
	RelayAddr   string   `json:"relay_addr,omitempty"`
}

// DiscoveryStatus represents peer discovery status
type DiscoveryStatus struct {
	MDNSActive     bool      `json:"mdns_active"`
	DHTActive      bool      `json:"dht_active"`
	UDPBroadcast   bool      `json:"udp_broadcast"`
	BootstrapPeers []string  `json:"bootstrap_peers"`
	KnownPeers     int       `json:"known_peers"`
	LastDiscovery  time.Time `json:"last_discovery"`
}

// NodeStatus represents the current status of a running node
type NodeStatus struct {
	PeerID            string    `json:"peer_id"`
	ListenAddrs       []string  `json:"listen_addrs"`
	ConnectedPeers    int       `json:"connected_peers"`
	UptimeSeconds     float64   `json:"uptime_seconds"`
	MessagesProcessed int64     `json:"messages_processed"`
	StartTime         time.Time `json:"start_time"`
	LastUpdate        time.Time `json:"last_update"`
	ProcessID         int       `json:"process_id"`
	IsRunning         bool      `json:"is_running"`

	// Extended network information
	Transports     []NetworkTransport `json:"transports"`
	NATInfo        *NATInfo           `json:"nat_info,omitempty"`
	Discovery      *DiscoveryStatus   `json:"discovery,omitempty"`
	NetworkQuality string             `json:"network_quality"` // "excellent", "good", "poor", "offline"
}

// PeerChatNode represents the main P2P node for Xelvra messenger
type PeerChatNode struct {
	host   host.Host
	ctx    context.Context
	cancel context.CancelFunc
	logger *logrus.Logger

	// Performance monitoring
	startTime    time.Time
	messageCount int64
	mu           sync.RWMutex

	// Configuration
	config *NodeConfig

	// Message handling
	messageManager *message.MessageManager
	identity       *user.MessengerID

	// Network components
	stunClient       *LegacySTUNClient
	discoveryManager *DiscoveryManager
	energyManager    *EnergyManager
	natInfo          *NATInfo
}

// NodeConfig holds configuration for the P2P node
type NodeConfig struct {
	ListenAddrs    []string
	BootstrapPeers []peer.AddrInfo
	EnableQUIC     bool
	EnableTCP      bool
	LogLevel       logrus.Level
	Logger         *logrus.Logger // External logger to use
}

// DefaultNodeConfig returns a default configuration optimized for performance
func DefaultNodeConfig() *NodeConfig {
	return &NodeConfig{
		ListenAddrs: []string{
			"/ip4/0.0.0.0/tcp/0",
			"/ip4/0.0.0.0/udp/0/quic-v1",
		},
		EnableQUIC: true,
		EnableTCP:  true,
		LogLevel:   logrus.InfoLevel,
	}
}

// NewPeerChatNode creates a new P2P node with optimized settings
func NewPeerChatNode(ctx context.Context, config *NodeConfig) (*PeerChatNode, error) {
	if config == nil {
		config = DefaultNodeConfig()
	}

	// Use provided logger or create default one
	var logger *logrus.Logger
	if config.Logger != nil {
		logger = config.Logger
	} else {
		logger = logrus.New()
		logger.SetLevel(config.LogLevel)
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	}

	// Generate MessengerID (which includes Ed25519 keys)
	identity, err := user.GenerateMessengerID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate messenger identity: %w", err)
	}

	// Convert to libp2p private key
	privKey, err := crypto.UnmarshalEd25519PrivateKey(identity.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert private key: %w", err)
	}

	// Create context with cancellation
	nodeCtx, cancel := context.WithCancel(ctx)

	// Configure libp2p options for optimal performance
	opts := []libp2p.Option{
		libp2p.Identity(privKey),
		libp2p.ListenAddrStrings(config.ListenAddrs...),
		libp2p.Ping(false),   // Disable built-in ping to save resources
		libp2p.EnableRelay(), // Enable relay for NAT traversal (basic relay support)
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			// Create DHT for routing
			dht, err := dual.New(nodeCtx, h)
			if err != nil {
				return nil, err
			}
			return dht, nil
		}),
	}

	// Add TCP transport
	if config.EnableTCP {
		opts = append(opts, libp2p.Transport(tcp.NewTCPTransport))
		logger.Info("TCP transport enabled")
	}

	// Add QUIC transport with buffer size configuration
	if config.EnableQUIC {
		// Try to increase UDP buffer sizes for QUIC
		if err := increaseUDPBufferSizes(logger); err != nil {
			logger.WithError(err).Warn("Failed to increase UDP buffer sizes, QUIC performance may be reduced")
		}

		// Redirect QUIC logs to our logger
		if err := redirectQUICLogs(logger); err != nil {
			logger.WithError(err).Warn("Failed to redirect QUIC logs")
		}

		opts = append(opts, libp2p.Transport(libp2pquic.NewTransport))
		logger.Info("QUIC transport enabled")
	}

	// Create the libp2p host
	h, err := libp2p.New(opts...)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create libp2p host: %w", err)
	}

	node := &PeerChatNode{
		host:      h,
		ctx:       nodeCtx,
		cancel:    cancel,
		logger:    logger,
		startTime: time.Now(),
		config:    config,
		identity:  identity,
	}

	// Create network components
	node.stunClient = NewLegacySTUNClient(logger)
	node.discoveryManager = NewDiscoveryManager(h, logger)
	node.energyManager = NewEnergyManager(nodeCtx, logger)

	// Create message manager
	node.messageManager = message.NewMessageManager(h, identity, logger)

	// Set up stream handler for Xelvra protocol
	h.SetStreamHandler(XelvraProtocolID, node.handleStream)

	logger.WithFields(logrus.Fields{
		"peer_id": h.ID().String(),
		"addrs":   h.Addrs(),
	}).Info("PeerChatNode initialized successfully")

	return node, nil
}

// Start begins the P2P node operations
func (n *PeerChatNode) Start() error {
	n.logger.Info("Starting PeerChatNode...")

	// Log performance targets
	n.logger.WithFields(logrus.Fields{
		"max_idle_memory_mb":   MaxIdleMemoryMB,
		"max_idle_cpu_percent": MaxIdleCPUPercent,
		"max_latency_ms":       MaxLatencyMs,
	}).Info("Performance targets set")

	// Start message manager
	n.logger.Debug("About to start MessageManager...")
	if err := n.messageManager.Start(); err != nil {
		return fmt.Errorf("failed to start message manager: %w", err)
	}
	n.logger.Debug("MessageManager started, registering handlers...")

	// Register console message handler for text messages
	n.logger.Debug("Creating console message handler...")
	consoleHandler := message.NewConsoleMessageHandler(n.logger)
	n.messageManager.RegisterHandler(message.MessageTypeText, consoleHandler)
	n.messageManager.RegisterHandler(message.MessageTypeSystem, consoleHandler)
	n.logger.Debug("Message handlers registered, writing status file...")

	// Start NAT discovery
	n.logger.Debug("Starting NAT discovery...")
	go n.discoverNAT()

	// Start energy management
	n.logger.Debug("Starting energy management...")
	if err := n.energyManager.Start(); err != nil {
		n.logger.WithError(err).Warn("Failed to start energy management")
	}

	// Start peer discovery
	n.logger.Debug("Starting peer discovery...")
	if err := n.discoveryManager.Start(); err != nil {
		n.logger.WithError(err).Warn("Failed to start peer discovery")
	}

	// Write initial status file
	if err := n.writeStatusFile(); err != nil {
		n.logger.WithError(err).Warn("Failed to write status file")
	}
	n.logger.Debug("Status file written successfully")

	n.logger.Info("PeerChatNode started successfully")
	return nil
}

// Stop gracefully shuts down the P2P node
func (n *PeerChatNode) Stop() error {
	n.logger.Info("Stopping PeerChatNode...")

	// Stop energy manager
	if n.energyManager != nil {
		if err := n.energyManager.Stop(); err != nil {
			n.logger.WithError(err).Error("Failed to stop energy manager")
		}
	}

	// Stop discovery manager
	if n.discoveryManager != nil {
		if err := n.discoveryManager.Stop(); err != nil {
			n.logger.WithError(err).Error("Failed to stop discovery manager")
		}
	}

	// Stop message manager
	if n.messageManager != nil {
		if err := n.messageManager.Stop(); err != nil {
			n.logger.WithError(err).Error("Failed to stop message manager")
		}
	}

	// Remove status file
	if err := n.removeStatusFile(); err != nil {
		n.logger.WithError(err).Warn("Failed to remove status file")
	}

	// Cancel context to stop all operations
	n.cancel()

	// Close the libp2p host
	if err := n.host.Close(); err != nil {
		n.logger.WithError(err).Error("Error closing libp2p host")
		return err
	}

	uptime := time.Since(n.startTime)
	n.mu.RLock()
	msgCount := n.messageCount
	n.mu.RUnlock()

	n.logger.WithFields(logrus.Fields{
		"uptime_seconds":     uptime.Seconds(),
		"messages_processed": msgCount,
	}).Info("PeerChatNode stopped successfully")

	return nil
}

// GetHost returns the underlying libp2p host
func (n *PeerChatNode) GetHost() host.Host {
	return n.host
}

// GetPeerID returns the node's peer ID
func (n *PeerChatNode) GetPeerID() peer.ID {
	return n.host.ID()
}

// SendMessage sends a message to a peer
func (n *PeerChatNode) SendMessage(to string, content []byte, msgType message.MessageType) error {
	if n.messageManager == nil {
		return fmt.Errorf("message manager not initialized")
	}
	return n.messageManager.SendMessage(to, content, msgType)
}

// SendFile sends a file to a peer
func (n *PeerChatNode) SendFile(peerID peer.ID, filePath string) error {
	if n.messageManager == nil {
		return fmt.Errorf("message manager not initialized")
	}
	return n.messageManager.SendFile(peerID, filePath)
}

// GetIdentity returns the node's identity
func (n *PeerChatNode) GetIdentity() *user.MessengerID {
	return n.identity
}

// GetEnergyProfile returns the current energy profile
func (n *PeerChatNode) GetEnergyProfile() *EnergyProfile {
	if n.energyManager != nil {
		return n.energyManager.GetEnergyProfile()
	}
	return nil
}

// SetBatteryLevel updates the battery level for energy optimization
func (n *PeerChatNode) SetBatteryLevel(level float64) {
	if n.energyManager != nil {
		n.energyManager.SetBatteryLevel(level)
	}
}

// GetStats returns basic node statistics
func (n *PeerChatNode) GetStats() map[string]interface{} {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return map[string]interface{}{
		"peer_id":            n.host.ID().String(),
		"uptime_seconds":     time.Since(n.startTime).Seconds(),
		"messages_processed": n.messageCount,
		"connected_peers":    len(n.host.Network().Peers()),
		"listen_addrs":       n.host.Addrs(),
	}
}

// handleStream handles incoming streams on the Xelvra protocol
func (n *PeerChatNode) handleStream(stream network.Stream) {
	defer func() {
		if err := stream.Close(); err != nil {
			n.logger.WithError(err).Error("Failed to close stream")
		}
	}()

	n.mu.Lock()
	n.messageCount++
	n.mu.Unlock()

	remotePeer := stream.Conn().RemotePeer()
	n.logger.WithFields(logrus.Fields{
		"remote_peer": remotePeer.String(),
		"protocol":    stream.Protocol(),
	}).Debug("Handling incoming stream")

	// TODO: Implement message handling
	// For now, just acknowledge the connection
	_, err := stream.Write([]byte("ACK"))
	if err != nil {
		n.logger.WithError(err).Error("Failed to write ACK to stream")
	}
}

// getStatusFilePath returns the path to the status file
func getStatusFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".xelvra", "node_status.json"), nil
}

// writeStatusFile writes the current node status to a file
func (n *PeerChatNode) writeStatusFile() error {
	statusPath, err := getStatusFilePath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(statusPath), 0700); err != nil {
		return err
	}

	// Create status object
	addrs := make([]string, len(n.host.Addrs()))
	for i, addr := range n.host.Addrs() {
		addrs[i] = addr.String()
	}

	// Collect transport information
	transports := []NetworkTransport{}
	for _, addr := range n.host.Addrs() {
		transport := NetworkTransport{
			Type:      "tcp", // Simplified - would need to detect actual transport
			LocalAddr: addr.String(),
			IsActive:  true,
		}
		transports = append(transports, transport)
	}

	// Get discovery status
	var discoveryStatus *DiscoveryStatus
	if n.discoveryManager != nil {
		discoveryStatus = n.discoveryManager.GetStatus()
	}

	n.mu.RLock()
	natInfo := n.natInfo
	n.mu.RUnlock()

	status := NodeStatus{
		PeerID:            n.host.ID().String(),
		ListenAddrs:       addrs,
		ConnectedPeers:    len(n.host.Network().Peers()),
		UptimeSeconds:     time.Since(n.startTime).Seconds(),
		MessagesProcessed: n.messageCount,
		StartTime:         n.startTime,
		LastUpdate:        time.Now(),
		ProcessID:         os.Getpid(),
		IsRunning:         true,
		Transports:        transports,
		NATInfo:           natInfo,
		Discovery:         discoveryStatus,
		NetworkQuality:    n.GetNetworkQuality(),
	}

	// Write to file
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(statusPath, data, 0600)
}

// removeStatusFile removes the status file when node stops
func (n *PeerChatNode) removeStatusFile() error {
	statusPath, err := getStatusFilePath()
	if err != nil {
		return err
	}

	// Mark as not running before removing
	if err := n.markStatusNotRunning(); err != nil {
		n.logger.WithError(err).Warn("Failed to mark status as not running")
	}

	return os.Remove(statusPath)
}

// markStatusNotRunning updates the status file to indicate node is not running
func (n *PeerChatNode) markStatusNotRunning() error {
	statusPath, err := getStatusFilePath()
	if err != nil {
		return err
	}

	// Read current status
	data, err := os.ReadFile(statusPath)
	if err != nil {
		return err
	}

	var status NodeStatus
	if err := json.Unmarshal(data, &status); err != nil {
		return err
	}

	// Update status
	status.IsRunning = false
	status.LastUpdate = time.Now()

	// Write back
	data, err = json.MarshalIndent(status, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(statusPath, data, 0600)
}

// ReadNodeStatus reads the current node status from file
func ReadNodeStatus() (*NodeStatus, error) {
	statusPath, err := getStatusFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(statusPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No status file means no running node
		}
		return nil, err
	}

	var status NodeStatus
	if err := json.Unmarshal(data, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

// discoverNAT performs NAT discovery in background
func (n *PeerChatNode) discoverNAT() {
	ctx, cancel := context.WithTimeout(n.ctx, 30*time.Second)
	defer cancel()

	// Get local port from first listen address
	localPort := 0
	if len(n.host.Addrs()) > 0 {
		addr := n.host.Addrs()[0]
		if tcpAddr, err := net.ResolveTCPAddr("tcp", addr.String()); err == nil {
			localPort = tcpAddr.Port
		}
	}

	natInfo, err := n.stunClient.DiscoverNAT(ctx, localPort)
	if err != nil {
		n.logger.WithError(err).Warn("NAT discovery failed")
		return
	}

	n.mu.Lock()
	n.natInfo = natInfo
	n.mu.Unlock()

	n.logger.WithFields(logrus.Fields{
		"nat_type":    natInfo.Type,
		"public_ip":   natInfo.PublicIP,
		"public_port": natInfo.PublicPort,
	}).Info("NAT discovery completed")

	// Update status file with new NAT info
	if err := n.writeStatusFile(); err != nil {
		n.logger.WithError(err).Warn("Failed to update status file with NAT info")
	}
}

// GetNetworkQuality determines network quality based on various factors
func (n *PeerChatNode) GetNetworkQuality() string {
	n.mu.RLock()
	defer n.mu.RUnlock()

	connectedPeers := len(n.host.Network().Peers())

	if n.natInfo == nil {
		return "unknown"
	}

	// Determine quality based on NAT type and connectivity
	switch n.natInfo.Type {
	case "none":
		if connectedPeers > 0 {
			return "excellent"
		}
		return "good"
	case "full_cone", "port_restricted":
		if connectedPeers > 0 {
			return "good"
		}
		return "fair"
	case "symmetric":
		if n.natInfo.UsingRelay {
			return "poor"
		}
		if connectedPeers > 0 {
			return "fair"
		}
		return "poor"
	default:
		return "unknown"
	}
}

// increaseUDPBufferSizes attempts to increase UDP buffer sizes for QUIC
func increaseUDPBufferSizes(logger *logrus.Logger) error {
	// This is a best-effort attempt to increase UDP buffer sizes
	// The actual implementation would need to use syscalls or external tools
	// For now, we'll just log the recommendation

	logger.Info("QUIC UDP buffer size optimization:")
	logger.Info("  Recommended: sudo sysctl -w net.core.rmem_max=7340032")
	logger.Info("  Recommended: sudo sysctl -w net.core.wmem_max=7340032")
	logger.Info("  Current settings may limit QUIC performance")

	// In a production implementation, you might:
	// 1. Check current buffer sizes
	// 2. Attempt to increase them programmatically (requires root)
	// 3. Fall back to TCP if QUIC performance is poor

	return nil
}

// redirectQUICLogs redirects QUIC library logs to our logger
func redirectQUICLogs(logger *logrus.Logger) error {
	// QUIC-go logs using Go's standard log package
	// We'll redirect only the log package, not stderr, to avoid breaking console output

	// Redirect Go's standard log package to our logger
	// This catches QUIC-go logs that use log.Printf
	log.SetOutput(&logrusWriter{logger: logger})
	log.SetFlags(0) // Remove timestamp since logrus handles it

	logger.Info("QUIC logging redirected to log file")
	logger.Debug("Go log package redirected to capture QUIC messages")

	return nil
}

// logrusWriter implements io.Writer to redirect Go's log package to logrus
type logrusWriter struct {
	logger *logrus.Logger
}

func (w *logrusWriter) Write(p []byte) (n int, err error) {
	logLine := strings.TrimSpace(string(p))
	if logLine != "" {
		// Check if it's a QUIC-related log
		if strings.Contains(logLine, "quic") ||
			strings.Contains(logLine, "UDP") ||
			strings.Contains(logLine, "buffer size") ||
			strings.Contains(logLine, "receive buffer") ||
			strings.Contains(logLine, "failed to sufficiently increase") {
			// Log QUIC messages as debug level to reduce console noise
			w.logger.Debug("QUIC: " + logLine)
		} else {
			// Log other messages as info
			w.logger.Info("log: " + logLine)
		}
	}
	return len(p), nil
}
