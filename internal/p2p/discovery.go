package p2p

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/dual"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	multiaddr "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

// DiscoveryManager handles peer discovery using hierarchical approach
// Implements the hierarchical discovery from tmp/Návrhy.md:
// 1. Local Discovery (BLE, Wi-Fi Direct, mDNS) - fastest and most efficient
// 2. Global Discovery (DHT) - for distributed peer finding
type DiscoveryManager struct {
	host   host.Host
	logger *logrus.Logger
	ctx    context.Context
	cancel context.CancelFunc

	// Discovery methods
	mdnsService      mdns.Service
	dht              *dual.DHT
	routingDiscovery *drouting.RoutingDiscovery

	// Bootstrap peers for DHT
	bootstrapPeers []peer.AddrInfo

	// Local discovery cache (LRU)
	localPeerCache map[peer.ID]*peer.AddrInfo
	cacheMaxSize   int
	cacheOrder     []peer.ID // For LRU eviction

	// Status tracking
	mu              sync.RWMutex
	discoveredPeers map[peer.ID]*peer.AddrInfo
	status          *DiscoveryStatus

	// Hierarchical discovery priorities
	localDiscoveryActive  bool
	globalDiscoveryActive bool
}

// NewDiscoveryManager creates a new discovery manager
func NewDiscoveryManager(h host.Host, logger *logrus.Logger) *DiscoveryManager {
	ctx, cancel := context.WithCancel(context.Background())

	// Define bootstrap peers for DHT
	bootstrapPeers := getBootstrapPeers()

	return &DiscoveryManager{
		host:            h,
		logger:          logger,
		ctx:             ctx,
		cancel:          cancel,
		bootstrapPeers:  bootstrapPeers,
		discoveredPeers: make(map[peer.ID]*peer.AddrInfo),
		localPeerCache:  make(map[peer.ID]*peer.AddrInfo),
		cacheMaxSize:    100, // LRU cache for 100 local peers
		cacheOrder:      make([]peer.ID, 0),
		status: &DiscoveryStatus{
			MDNSActive:     false,
			DHTActive:      false,
			UDPBroadcast:   false,
			BootstrapPeers: make([]string, len(bootstrapPeers)),
			KnownPeers:     0,
			LastDiscovery:  time.Now(),
		},
		localDiscoveryActive:  false,
		globalDiscoveryActive: false,
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

	// Start DHT discovery
	if err := dm.startDHT(); err != nil {
		dm.logger.WithError(err).Warn("Failed to start DHT discovery")
	} else {
		dm.mu.Lock()
		dm.status.DHTActive = true
		dm.mu.Unlock()
		dm.logger.Info("DHT discovery started")
	}

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

	if dm.dht != nil {
		if err := dm.dht.Close(); err != nil {
			dm.logger.WithError(err).Warn("Failed to close DHT")
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
	defer func() {
		if err := conn.Close(); err != nil {
			dm.logger.WithError(err).Error("Failed to close UDP connection")
		}
	}()

	dm.logger.Info("Listening for UDP broadcasts on :42424")

	buffer := make([]byte, 1024)
	for {
		select {
		case <-dm.ctx.Done():
			return
		default:
			if err := conn.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
				dm.logger.WithError(err).Debug("Failed to set UDP read deadline")
				continue
			}
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
	defer func() {
		if err := conn.Close(); err != nil {
			dm.logger.WithError(err).Error("Failed to close UDP broadcast connection")
		}
	}()

	_, err = conn.Write([]byte(message))
	if err != nil {
		dm.logger.WithError(err).Warn("Failed to send UDP broadcast")
		return
	}

	dm.logger.WithFields(logrus.Fields{
		"message": message,
		"target":  "255.255.255.255:42424",
		"peer_id": dm.host.ID().String(),
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

// getBootstrapPeers returns a list of bootstrap peers for DHT
func getBootstrapPeers() []peer.AddrInfo {
	// Use IPFS bootstrap peers as they run compatible DHT and relay services
	bootstrapAddrs := []string{
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zp9FCS47PpbUANZBTokb6BPWjkp8Bk",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
		// Additional relay servers for better NAT traversal
		"/ip4/147.75.77.187/tcp/4001/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/ip4/147.75.195.153/tcp/4001/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
	}

	var bootstrapPeers []peer.AddrInfo
	for _, addrStr := range bootstrapAddrs {
		addr, err := multiaddr.NewMultiaddr(addrStr)
		if err != nil {
			continue
		}

		peerInfo, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			continue
		}

		bootstrapPeers = append(bootstrapPeers, *peerInfo)
	}

	return bootstrapPeers
}

// startDHT initializes and starts the DHT for peer discovery
func (dm *DiscoveryManager) startDHT() error {
	dm.logger.Info("Starting DHT for global peer discovery...")

	// Create DHT with bootstrap peers
	dht, err := dual.New(dm.ctx, dm.host)
	if err != nil {
		return fmt.Errorf("failed to create DHT: %w", err)
	}

	dm.dht = dht

	// Bootstrap the DHT
	if err := dm.dht.Bootstrap(dm.ctx); err != nil {
		dm.logger.WithError(err).Warn("DHT bootstrap failed, continuing anyway")
	}

	// Connect to bootstrap peers
	go dm.connectToBootstrapPeers()

	// Create routing discovery
	dm.routingDiscovery = drouting.NewRoutingDiscovery(dm.dht)

	// Start advertising our presence
	go dm.advertisePresence()

	// Start discovering peers via DHT
	go dm.discoverViaDHT()

	dm.logger.WithFields(logrus.Fields{
		"bootstrap_peers": len(dm.bootstrapPeers),
	}).Info("DHT started successfully for global peer discovery")

	return nil
}

// connectToBootstrapPeers connects to bootstrap peers for DHT
func (dm *DiscoveryManager) connectToBootstrapPeers() {
	dm.logger.Info("Connecting to bootstrap peers...")

	connected := 0
	for _, peerInfo := range dm.bootstrapPeers {
		ctx, cancel := context.WithTimeout(dm.ctx, 10*time.Second)

		if err := dm.host.Connect(ctx, peerInfo); err != nil {
			dm.logger.WithError(err).WithField("peer_id", peerInfo.ID.String()).Debug("Failed to connect to bootstrap peer")
		} else {
			dm.logger.WithField("peer_id", peerInfo.ID.String()).Info("Connected to bootstrap peer")
			connected++
		}

		cancel()
	}

	dm.logger.WithFields(logrus.Fields{
		"connected": connected,
		"total":     len(dm.bootstrapPeers),
	}).Info("Bootstrap peer connection completed")

	// Update status
	dm.mu.Lock()
	for i, peer := range dm.bootstrapPeers {
		if i < len(dm.status.BootstrapPeers) {
			dm.status.BootstrapPeers[i] = peer.ID.String()
		}
	}
	dm.mu.Unlock()
}

// advertisePresence advertises our presence in the DHT
func (dm *DiscoveryManager) advertisePresence() {
	if dm.routingDiscovery == nil {
		return
	}

	ticker := time.NewTicker(5 * time.Minute) // Advertise every 5 minutes
	defer ticker.Stop()

	// Advertise immediately
	dm.doAdvertise()

	for {
		select {
		case <-dm.ctx.Done():
			return
		case <-ticker.C:
			dm.doAdvertise()
		}
	}
}

// doAdvertise performs the actual advertisement
func (dm *DiscoveryManager) doAdvertise() {
	ctx, cancel := context.WithTimeout(dm.ctx, 30*time.Second)
	defer cancel()

	// Advertise with Xelvra namespace
	_, err := dm.routingDiscovery.Advertise(ctx, "xelvra-p2p")
	if err != nil {
		dm.logger.WithError(err).Debug("Failed to advertise presence in DHT")
	} else {
		dm.logger.Debug("Successfully advertised presence in DHT")
	}
}

// discoverViaDHT discovers peers using DHT
func (dm *DiscoveryManager) discoverViaDHT() {
	if dm.routingDiscovery == nil {
		return
	}

	ticker := time.NewTicker(2 * time.Minute) // Discover every 2 minutes
	defer ticker.Stop()

	// Discover immediately after a short delay to allow DHT to bootstrap
	time.Sleep(30 * time.Second)
	dm.doDHTDiscovery()

	for {
		select {
		case <-dm.ctx.Done():
			return
		case <-ticker.C:
			dm.doDHTDiscovery()
		}
	}
}

// doDHTDiscovery performs the actual DHT peer discovery
func (dm *DiscoveryManager) doDHTDiscovery() {
	ctx, cancel := context.WithTimeout(dm.ctx, 60*time.Second)
	defer cancel()

	dm.logger.Debug("Discovering peers via DHT...")

	// Find peers advertising the Xelvra namespace
	peerChan, err := dm.routingDiscovery.FindPeers(ctx, "xelvra-p2p")
	if err != nil {
		dm.logger.WithError(err).Debug("Failed to start DHT peer discovery")
		return
	}

	discovered := 0
	for peerInfo := range peerChan {
		// Skip ourselves
		if peerInfo.ID == dm.host.ID() {
			continue
		}

		dm.logger.WithFields(logrus.Fields{
			"peer_id": peerInfo.ID.String(),
			"addrs":   peerInfo.Addrs,
		}).Info("Discovered peer via DHT")

		// Store discovered peer
		dm.mu.Lock()
		dm.discoveredPeers[peerInfo.ID] = &peerInfo
		dm.status.LastDiscovery = time.Now()
		dm.mu.Unlock()

		// Try to connect to the discovered peer
		go func(pi peer.AddrInfo) {
			connectCtx, connectCancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer connectCancel()

			if err := dm.host.Connect(connectCtx, pi); err != nil {
				dm.logger.WithError(err).WithField("peer_id", pi.ID.String()).Debug("Failed to connect to DHT-discovered peer")
			} else {
				dm.logger.WithField("peer_id", pi.ID.String()).Info("Successfully connected to DHT-discovered peer")
			}
		}(peerInfo)

		discovered++
	}

	if discovered > 0 {
		dm.logger.WithField("discovered", discovered).Info("DHT peer discovery completed")
	} else {
		dm.logger.Debug("No new peers discovered via DHT")
	}
}

// addToLocalCache adds a peer to the local LRU cache
func (dm *DiscoveryManager) addToLocalCache(peerID peer.ID, peerInfo *peer.AddrInfo) {
	// Check if peer already exists in cache
	if _, exists := dm.localPeerCache[peerID]; exists {
		// Move to front of LRU order
		dm.moveToFront(peerID)
		return
	}

	// Add new peer
	dm.localPeerCache[peerID] = peerInfo
	dm.cacheOrder = append([]peer.ID{peerID}, dm.cacheOrder...)

	// Evict oldest if cache is full
	if len(dm.localPeerCache) > dm.cacheMaxSize {
		oldest := dm.cacheOrder[len(dm.cacheOrder)-1]
		delete(dm.localPeerCache, oldest)
		dm.cacheOrder = dm.cacheOrder[:len(dm.cacheOrder)-1]
	}
}

// moveToFront moves a peer to the front of the LRU order
func (dm *DiscoveryManager) moveToFront(peerID peer.ID) {
	// Find and remove peer from current position
	for i, id := range dm.cacheOrder {
		if id == peerID {
			dm.cacheOrder = append(dm.cacheOrder[:i], dm.cacheOrder[i+1:]...)
			break
		}
	}
	// Add to front
	dm.cacheOrder = append([]peer.ID{peerID}, dm.cacheOrder...)
}

// getFromLocalCache retrieves a peer from local cache
func (dm *DiscoveryManager) getFromLocalCache(peerID peer.ID) (*peer.AddrInfo, bool) {
	if peerInfo, exists := dm.localPeerCache[peerID]; exists {
		dm.moveToFront(peerID)
		return peerInfo, true
	}
	return nil, false
}

// isLocalPeer checks if a peer is in the local network (same subnet)
func (dm *DiscoveryManager) isLocalPeer(peerInfo *peer.AddrInfo) bool {
	for _, addr := range peerInfo.Addrs {
		addrStr := addr.String()
		// Check for local network addresses
		if strings.Contains(addrStr, "192.168.") ||
			strings.Contains(addrStr, "10.") ||
			strings.Contains(addrStr, "172.16.") ||
			strings.Contains(addrStr, "172.17.") ||
			strings.Contains(addrStr, "172.18.") ||
			strings.Contains(addrStr, "172.19.") ||
			strings.Contains(addrStr, "172.2") ||
			strings.Contains(addrStr, "172.30.") ||
			strings.Contains(addrStr, "172.31.") ||
			strings.Contains(addrStr, "127.0.0.1") {
			return true
		}
	}
	return false
}
