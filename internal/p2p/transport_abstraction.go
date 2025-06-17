package p2p

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/sirupsen/logrus"
)

// NetworkInterface defines the abstract network interface for testing and flexibility
// Inspired by battle-tested network abstraction patterns
type NetworkInterface interface {
	// Connection management
	Connect(ctx context.Context, peerID peer.ID, addrs []string) (Connection, error)
	Listen(ctx context.Context, addr string) (Listener, error)
	Close() error

	// Stream operations
	NewStream(ctx context.Context, conn Connection, protocol string) (Stream, error)
	AcceptStream(listener Listener) (Stream, error)

	// Network information
	LocalAddresses() []string
	RemoteAddresses(conn Connection) []string
	ConnectionStatus(conn Connection) ConnectionStatus
}

// Connection represents an abstract network connection
type Connection interface {
	// Basic connection operations
	ID() string
	RemotePeer() peer.ID
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close() error

	// Connection properties
	IsConnected() bool
	Latency() time.Duration
	Bandwidth() (upload, download int64)

	// Stream management
	OpenStream(protocol string) (Stream, error)
	AcceptStream() (Stream, error)
}

// Stream represents an abstract network stream
type Stream interface {
	io.ReadWriteCloser

	// Stream properties
	ID() string
	Protocol() string
	Connection() Connection

	// Flow control
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

// Listener represents an abstract network listener
type Listener interface {
	Accept() (Connection, error)
	Close() error
	Addr() net.Addr
}

// ConnectionStatus represents connection status
type ConnectionStatus int

const (
	StatusDisconnected ConnectionStatus = iota
	StatusConnecting
	StatusConnected
	StatusError
)

// TransportManager manages different transport implementations
type TransportManager struct {
	mu     sync.RWMutex
	logger *logrus.Logger

	// Transport implementations
	transports map[string]NetworkInterface
	primary    string
	fallbacks  []string

	// Connection management
	connections map[string]Connection
	connPool    *ConnectionPool

	// Metrics and monitoring
	metrics *TransportMetrics
}

// ConnectionPool manages connection reuse and pooling
type ConnectionPool struct {
	mu          sync.RWMutex
	connections map[string]*PooledConnection
	maxSize     int
	maxAge      time.Duration
	logger      *logrus.Logger
}

// PooledConnection represents a pooled connection
type PooledConnection struct {
	conn     Connection
	created  time.Time
	lastUsed time.Time
	useCount int
	inUse    bool
}

// TransportMetrics tracks transport performance
type TransportMetrics struct {
	mu                sync.RWMutex
	connectionsTotal  map[string]int64
	connectionsActive map[string]int64
	bytesTransferred  map[string]int64
	latencyHistogram  map[string][]time.Duration
	errorCounts       map[string]int64
}

// LibP2PTransport implements NetworkInterface using libp2p
type LibP2PTransport struct {
	host   host.Host
	logger *logrus.Logger
}

// MockTransport implements NetworkInterface for testing
type MockTransport struct {
	mu          sync.RWMutex
	connections map[string]Connection
	listeners   map[string]Listener
	logger      *logrus.Logger

	// Test configuration
	simulateLatency bool
	latencyRange    time.Duration
	dropRate        float64
	errorRate       float64
}

// NewTransportManager creates a new transport manager
func NewTransportManager(logger *logrus.Logger) *TransportManager {
	return &TransportManager{
		logger:      logger,
		transports:  make(map[string]NetworkInterface),
		connections: make(map[string]Connection),
		connPool: &ConnectionPool{
			connections: make(map[string]*PooledConnection),
			maxSize:     100,
			maxAge:      30 * time.Minute,
			logger:      logger,
		},
		metrics: &TransportMetrics{
			connectionsTotal:  make(map[string]int64),
			connectionsActive: make(map[string]int64),
			bytesTransferred:  make(map[string]int64),
			latencyHistogram:  make(map[string][]time.Duration),
			errorCounts:       make(map[string]int64),
		},
	}
}

// RegisterTransport registers a transport implementation
func (tm *TransportManager) RegisterTransport(name string, transport NetworkInterface) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.transports[name] = transport
	tm.logger.WithField("transport", name).Info("Transport registered")
}

// SetPrimaryTransport sets the primary transport
func (tm *TransportManager) SetPrimaryTransport(name string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, exists := tm.transports[name]; !exists {
		return fmt.Errorf("transport %s not registered", name)
	}

	tm.primary = name
	tm.logger.WithField("transport", name).Info("Primary transport set")
	return nil
}

// SetFallbackTransports sets fallback transports in order of preference
func (tm *TransportManager) SetFallbackTransports(names []string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, name := range names {
		if _, exists := tm.transports[name]; !exists {
			return fmt.Errorf("transport %s not registered", name)
		}
	}

	tm.fallbacks = names
	tm.logger.WithField("fallbacks", names).Info("Fallback transports set")
	return nil
}

// Connect attempts to connect using primary transport with fallbacks
func (tm *TransportManager) Connect(ctx context.Context, peerID peer.ID, addrs []string) (Connection, error) {
	// Try to get connection from pool first
	if conn := tm.connPool.getConnection(peerID.String()); conn != nil {
		tm.logger.WithField("peer_id", peerID.String()).Debug("Using pooled connection")
		return conn, nil
	}

	// Try primary transport first
	if tm.primary != "" {
		if conn, err := tm.tryConnect(ctx, tm.primary, peerID, addrs); err == nil {
			tm.connPool.addConnection(peerID.String(), conn)
			return conn, nil
		}
	}

	// Try fallback transports
	for _, transportName := range tm.fallbacks {
		if conn, err := tm.tryConnect(ctx, transportName, peerID, addrs); err == nil {
			tm.connPool.addConnection(peerID.String(), conn)
			return conn, nil
		}
	}

	return nil, fmt.Errorf("failed to connect using any transport")
}

// tryConnect attempts connection using specific transport
func (tm *TransportManager) tryConnect(ctx context.Context, transportName string, peerID peer.ID, addrs []string) (Connection, error) {
	tm.mu.RLock()
	transport, exists := tm.transports[transportName]
	tm.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("transport %s not found", transportName)
	}

	start := time.Now()
	conn, err := transport.Connect(ctx, peerID, addrs)
	duration := time.Since(start)

	// Update metrics
	tm.updateMetrics(transportName, duration, err == nil)

	if err != nil {
		tm.logger.WithFields(logrus.Fields{
			"transport": transportName,
			"peer_id":   peerID.String(),
			"duration":  duration,
			"error":     err,
		}).Debug("Connection attempt failed")
		return nil, err
	}

	tm.logger.WithFields(logrus.Fields{
		"transport": transportName,
		"peer_id":   peerID.String(),
		"duration":  duration,
	}).Info("Connection established")

	return conn, nil
}

// NewLibP2PTransport creates a new libp2p transport
func NewLibP2PTransport(h host.Host, logger *logrus.Logger) *LibP2PTransport {
	return &LibP2PTransport{
		host:   h,
		logger: logger,
	}
}

// Connect implements NetworkInterface for libp2p
func (lt *LibP2PTransport) Connect(ctx context.Context, peerID peer.ID, addrs []string) (Connection, error) {
	// Convert addresses to multiaddrs and create peer.AddrInfo
	// This is a simplified implementation

	// Try to connect using libp2p host
	err := lt.host.Connect(ctx, peer.AddrInfo{ID: peerID})
	if err != nil {
		return nil, err
	}

	// Return wrapped connection
	return &LibP2PConnection{
		peerID: peerID,
		host:   lt.host,
		logger: lt.logger,
	}, nil
}

// Listen implements NetworkInterface for libp2p
func (lt *LibP2PTransport) Listen(ctx context.Context, addr string) (Listener, error) {
	// libp2p handles listening automatically
	return &LibP2PListener{
		host:   lt.host,
		logger: lt.logger,
	}, nil
}

// Close implements NetworkInterface for libp2p
func (lt *LibP2PTransport) Close() error {
	return lt.host.Close()
}

// NewStream implements NetworkInterface for libp2p
func (lt *LibP2PTransport) NewStream(ctx context.Context, conn Connection, protocolStr string) (Stream, error) {
	libp2pConn, ok := conn.(*LibP2PConnection)
	if !ok {
		return nil, fmt.Errorf("invalid connection type")
	}

	stream, err := lt.host.NewStream(ctx, libp2pConn.peerID, protocol.ID(protocolStr))
	if err != nil {
		return nil, err
	}

	return &LibP2PStream{
		stream: stream,
		logger: lt.logger,
	}, nil
}

// AcceptStream implements NetworkInterface for libp2p
func (lt *LibP2PTransport) AcceptStream(listener Listener) (Stream, error) {
	// libp2p uses stream handlers instead of explicit accept
	return nil, fmt.Errorf("libp2p uses stream handlers")
}

// LocalAddresses implements NetworkInterface for libp2p
func (lt *LibP2PTransport) LocalAddresses() []string {
	addrs := lt.host.Addrs()
	result := make([]string, len(addrs))
	for i, addr := range addrs {
		result[i] = addr.String()
	}
	return result
}

// RemoteAddresses implements NetworkInterface for libp2p
func (lt *LibP2PTransport) RemoteAddresses(conn Connection) []string {
	libp2pConn, ok := conn.(*LibP2PConnection)
	if !ok {
		return nil
	}

	peerInfo := lt.host.Peerstore().PeerInfo(libp2pConn.peerID)
	result := make([]string, len(peerInfo.Addrs))
	for i, addr := range peerInfo.Addrs {
		result[i] = addr.String()
	}
	return result
}

// ConnectionStatus implements NetworkInterface for libp2p
func (lt *LibP2PTransport) ConnectionStatus(conn Connection) ConnectionStatus {
	libp2pConn, ok := conn.(*LibP2PConnection)
	if !ok {
		return StatusError
	}

	connectedness := lt.host.Network().Connectedness(libp2pConn.peerID)
	switch connectedness {
	case network.Connected:
		return StatusConnected
	case network.CanConnect:
		return StatusConnecting
	default:
		return StatusDisconnected
	}
}

// LibP2PConnection implements Connection for libp2p
type LibP2PConnection struct {
	peerID peer.ID
	host   host.Host
	logger *logrus.Logger
}

func (lc *LibP2PConnection) ID() string {
	return lc.peerID.String()
}

func (lc *LibP2PConnection) RemotePeer() peer.ID {
	return lc.peerID
}

func (lc *LibP2PConnection) LocalAddr() net.Addr {
	// Return first local address
	addrs := lc.host.Addrs()
	if len(addrs) > 0 {
		// This is simplified - real implementation would parse multiaddr
		return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}
	}
	return nil
}

func (lc *LibP2PConnection) RemoteAddr() net.Addr {
	// This is simplified - real implementation would parse peer addresses
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}
}

func (lc *LibP2PConnection) Close() error {
	return lc.host.Network().ClosePeer(lc.peerID)
}

func (lc *LibP2PConnection) IsConnected() bool {
	return lc.host.Network().Connectedness(lc.peerID) == network.Connected
}

func (lc *LibP2PConnection) Latency() time.Duration {
	// Get latency from libp2p if available
	return 50 * time.Millisecond // Placeholder
}

func (lc *LibP2PConnection) Bandwidth() (upload, download int64) {
	// Get bandwidth info from libp2p if available
	return 1000000, 1000000 // Placeholder: 1MB/s
}

func (lc *LibP2PConnection) OpenStream(protocolStr string) (Stream, error) {
	stream, err := lc.host.NewStream(context.Background(), lc.peerID, protocol.ID(protocolStr))
	if err != nil {
		return nil, err
	}

	return &LibP2PStream{
		stream: stream,
		logger: lc.logger,
	}, nil
}

func (lc *LibP2PConnection) AcceptStream() (Stream, error) {
	// libp2p uses stream handlers
	return nil, fmt.Errorf("libp2p uses stream handlers")
}

// LibP2PStream implements Stream for libp2p
type LibP2PStream struct {
	stream network.Stream
	logger *logrus.Logger
}

func (ls *LibP2PStream) Read(p []byte) (n int, err error) {
	return ls.stream.Read(p)
}

func (ls *LibP2PStream) Write(p []byte) (n int, err error) {
	return ls.stream.Write(p)
}

func (ls *LibP2PStream) Close() error {
	return ls.stream.Close()
}

func (ls *LibP2PStream) ID() string {
	return ls.stream.ID()
}

func (ls *LibP2PStream) Protocol() string {
	return string(ls.stream.Protocol())
}

func (ls *LibP2PStream) Connection() Connection {
	// Return wrapped connection
	return &LibP2PConnection{
		peerID: ls.stream.Conn().RemotePeer(),
		host:   nil, // Would need reference to host
		logger: ls.logger,
	}
}

func (ls *LibP2PStream) SetDeadline(t time.Time) error {
	return ls.stream.SetDeadline(t)
}

func (ls *LibP2PStream) SetReadDeadline(t time.Time) error {
	return ls.stream.SetReadDeadline(t)
}

func (ls *LibP2PStream) SetWriteDeadline(t time.Time) error {
	return ls.stream.SetWriteDeadline(t)
}

// LibP2PListener implements Listener for libp2p
type LibP2PListener struct {
	host   host.Host
	logger *logrus.Logger
}

func (ll *LibP2PListener) Accept() (Connection, error) {
	// libp2p handles connections automatically
	return nil, fmt.Errorf("libp2p handles connections automatically")
}

func (ll *LibP2PListener) Close() error {
	return ll.host.Close()
}

func (ll *LibP2PListener) Addr() net.Addr {
	// Return first address
	addrs := ll.host.Addrs()
	if len(addrs) > 0 {
		return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}
	}
	return nil
}

// Helper methods for TransportManager
func (tm *TransportManager) updateMetrics(transportName string, duration time.Duration, success bool) {
	tm.metrics.mu.Lock()
	defer tm.metrics.mu.Unlock()

	tm.metrics.connectionsTotal[transportName]++

	if success {
		tm.metrics.connectionsActive[transportName]++

		// Update latency histogram
		if tm.metrics.latencyHistogram[transportName] == nil {
			tm.metrics.latencyHistogram[transportName] = make([]time.Duration, 0, 100)
		}

		latencies := tm.metrics.latencyHistogram[transportName]
		latencies = append(latencies, duration)

		// Keep only last 100 measurements
		if len(latencies) > 100 {
			latencies = latencies[1:]
		}

		tm.metrics.latencyHistogram[transportName] = latencies
	} else {
		tm.metrics.errorCounts[transportName]++
	}
}

// Helper methods for ConnectionPool
func (cp *ConnectionPool) getConnection(key string) Connection {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	pooled, exists := cp.connections[key]
	if !exists {
		return nil
	}

	// Check if connection is still valid and not too old
	if time.Since(pooled.created) > cp.maxAge || pooled.inUse || !pooled.conn.IsConnected() {
		delete(cp.connections, key)
		return nil
	}

	pooled.inUse = true
	pooled.lastUsed = time.Now()
	pooled.useCount++

	return pooled.conn
}

func (cp *ConnectionPool) addConnection(key string, conn Connection) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	// Check pool size limit
	if len(cp.connections) >= cp.maxSize {
		// Remove oldest connection
		var oldestKey string
		var oldestTime time.Time

		for k, v := range cp.connections {
			if oldestKey == "" || v.lastUsed.Before(oldestTime) {
				oldestKey = k
				oldestTime = v.lastUsed
			}
		}

		if oldestKey != "" {
			if oldConn := cp.connections[oldestKey]; oldConn != nil {
				oldConn.conn.Close()
			}
			delete(cp.connections, oldestKey)
		}
	}

	cp.connections[key] = &PooledConnection{
		conn:     conn,
		created:  time.Now(),
		lastUsed: time.Now(),
		useCount: 0,
		inUse:    false,
	}
}

func (cp *ConnectionPool) releaseConnection(key string) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if pooled, exists := cp.connections[key]; exists {
		pooled.inUse = false
	}
}

func (cp *ConnectionPool) cleanup() {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	now := time.Now()
	for key, pooled := range cp.connections {
		if now.Sub(pooled.created) > cp.maxAge || !pooled.conn.IsConnected() {
			pooled.conn.Close()
			delete(cp.connections, key)
		}
	}
}
