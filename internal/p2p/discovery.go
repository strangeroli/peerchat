package p2p

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

// DiscoveryManager handles peer discovery using multiple methods
type DiscoveryManager struct {
	host     host.Host
	logger   *logrus.Logger
	ctx      context.Context
	cancel   context.CancelFunc
	
	// Discovery methods
	mdnsService mdns.Service
	
	// Status tracking
	mu             sync.RWMutex
	discoveredPeers map[peer.ID]*peer.AddrInfo
	status         *DiscoveryStatus
}

// NewDiscoveryManager creates a new discovery manager
func NewDiscoveryManager(h host.Host, logger *logrus.Logger) *DiscoveryManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &DiscoveryManager{
		host:            h,
		logger:          logger,
		ctx:             ctx,
		cancel:          cancel,
		discoveredPeers: make(map[peer.ID]*peer.AddrInfo),
		status: &DiscoveryStatus{
			MDNSActive:     false,
			DHTActive:      false,
			UDPBroadcast:   false,
			BootstrapPeers: []string{},
			KnownPeers:     0,
			LastDiscovery:  time.Now(),
		},
	}
}

// Start begins peer discovery
func (dm *DiscoveryManager) Start() error {
	dm.logger.Info("Starting peer discovery...")

	// Start mDNS discovery
	if err := dm.startMDNS(); err != nil {
		dm.logger.WithError(err).Warn("Failed to start mDNS discovery")
	} else {
		dm.mu.Lock()
		dm.status.MDNSActive = true
		dm.mu.Unlock()
		dm.logger.Info("mDNS discovery started")
	}

	// Start UDP broadcast discovery
	go dm.startUDPBroadcast()
	dm.mu.Lock()
	dm.status.UDPBroadcast = true
	dm.mu.Unlock()

	// TODO: Start DHT discovery
	// TODO: Connect to bootstrap peers

	dm.logger.Info("Peer discovery started successfully")
	return nil
}

// Stop stops peer discovery
func (dm *DiscoveryManager) Stop() error {
	dm.logger.Info("Stopping peer discovery...")
	
	dm.cancel()
	
	if dm.mdnsService != nil {
		if err := dm.mdnsService.Close(); err != nil {
			dm.logger.WithError(err).Warn("Failed to close mDNS service")
		}
	}
	
	dm.mu.Lock()
	dm.status.MDNSActive = false
	dm.status.DHTActive = false
	dm.status.UDPBroadcast = false
	dm.mu.Unlock()
	
	dm.logger.Info("Peer discovery stopped")
	return nil
}

// GetStatus returns current discovery status
func (dm *DiscoveryManager) GetStatus() *DiscoveryStatus {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	
	// Update known peers count
	dm.status.KnownPeers = len(dm.discoveredPeers)
	
	// Copy status to avoid race conditions
	status := *dm.status
	return &status
}

// GetDiscoveredPeers returns list of discovered peers
func (dm *DiscoveryManager) GetDiscoveredPeers() []peer.ID {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	peers := make([]peer.ID, 0, len(dm.discoveredPeers))
	for peerID := range dm.discoveredPeers {
		peers = append(peers, peerID)
	}
	return peers
}

// GetPeerAddresses returns addresses for a specific peer
func (dm *DiscoveryManager) GetPeerAddresses(peerID peer.ID) []multiaddr.Multiaddr {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	if peerInfo, exists := dm.discoveredPeers[peerID]; exists {
		return peerInfo.Addrs
	}
	return nil
}

// startMDNS starts mDNS peer discovery
func (dm *DiscoveryManager) startMDNS() error {
	dm.logger.Info("Starting mDNS service with service name: xelvra-p2p")

	service := mdns.NewMdnsService(dm.host, "xelvra-p2p", &discoveryNotifee{dm: dm})
	if err := service.Start(); err != nil {
		return fmt.Errorf("failed to start mDNS service: %w", err)
	}

	dm.mdnsService = service
	dm.logger.WithFields(logrus.Fields{
		"service_name": "xelvra-p2p",
		"peer_id":      dm.host.ID().String(),
		"addresses":    dm.host.Addrs(),
	}).Info("mDNS service started successfully for local peer discovery")
	return nil
}

// startUDPBroadcast starts UDP broadcast discovery for local network
func (dm *DiscoveryManager) startUDPBroadcast() {
	dm.logger.Info("Starting UDP broadcast discovery...")
	
	// Listen for broadcasts
	go dm.listenUDPBroadcast()
	
	// Send periodic broadcasts
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-dm.ctx.Done():
			return
		case <-ticker.C:
			dm.sendUDPBroadcast()
		}
	}
}

// listenUDPBroadcast listens for UDP broadcast messages
func (dm *DiscoveryManager) listenUDPBroadcast() {
	addr, err := net.ResolveUDPAddr("udp", ":42424")
	if err != nil {
		dm.logger.WithError(err).Error("Failed to resolve UDP broadcast address")
		return
	}
	
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		dm.logger.WithError(err).Error("Failed to listen on UDP broadcast")
		return
	}
	defer conn.Close()
	
	dm.logger.Info("Listening for UDP broadcasts on :42424")
	
	buffer := make([]byte, 1024)
	for {
		select {
		case <-dm.ctx.Done():
			return
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, remoteAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				dm.logger.WithError(err).Debug("UDP broadcast read error")
				continue
			}
			
			dm.handleUDPBroadcast(buffer[:n], remoteAddr)
		}
	}
}

// sendUDPBroadcast sends UDP broadcast message
func (dm *DiscoveryManager) sendUDPBroadcast() {
	message := fmt.Sprintf("XELVRA_PEER:%s", dm.host.ID().String())

	conn, err := net.Dial("udp", "255.255.255.255:42424")
	if err != nil {
		dm.logger.WithError(err).Warn("Failed to create UDP broadcast connection - check firewall/network")
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	if err != nil {
		dm.logger.WithError(err).Warn("Failed to send UDP broadcast")
		return
	}

	dm.logger.WithFields(logrus.Fields{
		"message":   message,
		"target":    "255.255.255.255:42424",
		"peer_id":   dm.host.ID().String(),
	}).Info("Sent UDP broadcast for peer discovery")
}

// handleUDPBroadcast handles received UDP broadcast messages
func (dm *DiscoveryManager) handleUDPBroadcast(data []byte, remoteAddr *net.UDPAddr) {
	message := string(data)
	if len(message) < 12 || message[:12] != "XELVRA_PEER:" {
		return
	}
	
	peerIDStr := message[12:]
	peerID, err := peer.Decode(peerIDStr)
	if err != nil {
		dm.logger.WithError(err).Debug("Invalid peer ID in UDP broadcast")
		return
	}
	
	// Don't discover ourselves
	if peerID == dm.host.ID() {
		return
	}
	
	dm.logger.WithFields(logrus.Fields{
		"peer_id":     peerID.String(),
		"remote_addr": remoteAddr.String(),
	}).Info("Discovered peer via UDP broadcast")
	
	// Create basic peer info (UDP broadcast doesn't provide addresses)
	peerInfo := &peer.AddrInfo{
		ID:    peerID,
		Addrs: []multiaddr.Multiaddr{}, // No addresses from UDP broadcast
	}

	dm.mu.Lock()
	dm.discoveredPeers[peerID] = peerInfo
	dm.status.LastDiscovery = time.Now()
	dm.mu.Unlock()
}

// discoveryNotifee handles mDNS discovery notifications
type discoveryNotifee struct {
	dm *DiscoveryManager
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.dm.logger.WithFields(logrus.Fields{
		"peer_id": pi.ID.String(),
		"addrs":   pi.Addrs,
	}).Info("Discovered peer via mDNS")
	
	n.dm.mu.Lock()
	n.dm.discoveredPeers[pi.ID] = &pi
	n.dm.status.LastDiscovery = time.Now()
	n.dm.mu.Unlock()
	
	// Try to connect to the discovered peer
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := n.dm.host.Connect(ctx, pi); err != nil {
		n.dm.logger.WithError(err).WithField("peer_id", pi.ID.String()).Debug("Failed to connect to discovered peer")
	} else {
		n.dm.logger.WithField("peer_id", pi.ID.String()).Info("Successfully connected to discovered peer")
	}
}
