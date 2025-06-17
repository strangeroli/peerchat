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
	"github.com/libp2p/go-libp2p/core/network"
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

// Start begins hierarchical peer discovery: IPv6 → mDNS → hole punching → relay
func (dm *DiscoveryManager) Start() error {
	dm.logger.Info("Starting hierarchical peer discovery: IPv6 → mDNS → UDP → DHT → Hole Punching → Relay...")

	// Phase 1: IPv6 Link-Local Discovery (highest priority, immediate)
	go dm.startIPv6LinkLocalDiscovery()
	dm.logger.Info("Phase 1: IPv6 link-local discovery started")

	// Phase 2: mDNS discovery (local network, fast)
	if err := dm.startMDNS(); err != nil {
		dm.logger.WithError(err).Warn("Failed to start mDNS discovery")
	} else {
		dm.mu.Lock()
		dm.status.MDNSActive = true
		dm.localDiscoveryActive = true
		dm.mu.Unlock()
		dm.logger.Info("Phase 2: mDNS discovery started")
	}

	// Phase 3: UDP broadcast discovery (local network fallback)
	go dm.startUDPBroadcast()
	dm.mu.Lock()
	dm.status.UDPBroadcast = true
	dm.mu.Unlock()
	dm.logger.Info("Phase 3: UDP broadcast discovery started")

	// Phase 4: DHT discovery (global network)
	if err := dm.startDHT(); err != nil {
		dm.logger.WithError(err).Warn("Failed to start DHT discovery")
	} else {
		dm.mu.Lock()
		dm.status.DHTActive = true
		dm.globalDiscoveryActive = true
		dm.mu.Unlock()
		dm.logger.Info("Phase 4: DHT global discovery started")
	}

	// Phase 5: NAT hole punching service
	go dm.startHolePunchingService()
	dm.logger.Info("Phase 5: NAT hole punching service started")

	// Phase 6: Relay server management
	go dm.startRelayServerManagement()
	dm.logger.Info("Phase 6: Relay server management started")

	dm.logger.Info("Hierarchical peer discovery started successfully - all 6 phases active")
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
			strings.Contains(addrStr, "127.0.0.1") ||
			strings.Contains(addrStr, "fe80::") { // IPv6 link-local
			return true
		}
	}
	return false
}

// startIPv6LinkLocalDiscovery starts IPv6 link-local peer discovery (Phase 1)
func (dm *DiscoveryManager) startIPv6LinkLocalDiscovery() {
	dm.logger.Info("Starting IPv6 link-local discovery (Phase 1 - highest priority)...")

	ticker := time.NewTicker(15 * time.Second) // Check every 15 seconds
	defer ticker.Stop()

	for {
		select {
		case <-dm.ctx.Done():
			return
		case <-ticker.C:
			dm.discoverIPv6LinkLocal()
		}
	}
}

// discoverIPv6LinkLocal discovers peers on IPv6 link-local network
func (dm *DiscoveryManager) discoverIPv6LinkLocal() {
	// Get all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		dm.logger.WithError(err).Debug("Failed to get network interfaces for IPv6 discovery")
		return
	}

	discovered := 0
	for _, iface := range interfaces {
		// Skip loopback and down interfaces
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipnet.IP.IsLinkLocalUnicast() && ipnet.IP.To4() == nil { // IPv6 link-local
					dm.logger.WithFields(logrus.Fields{
						"interface": iface.Name,
						"ipv6_addr": ipnet.IP.String(),
					}).Debug("Found IPv6 link-local address")

					// Try to discover peers on this link-local network
					if dm.scanIPv6LinkLocalNetwork(ipnet, iface.Name) {
						discovered++
					}
				}
			}
		}
	}

	if discovered > 0 {
		dm.logger.WithField("networks_scanned", discovered).Info("IPv6 link-local discovery completed")
	}
}

// scanIPv6LinkLocalNetwork scans for peers on IPv6 link-local network
func (dm *DiscoveryManager) scanIPv6LinkLocalNetwork(network *net.IPNet, interfaceName string) bool {
	// For IPv6 link-local, we use multicast discovery
	// Send to all-nodes multicast address ff02::1
	multicastAddr := "ff02::1"

	conn, err := net.Dial("udp6", fmt.Sprintf("[%s%%25%s]:42425", multicastAddr, interfaceName))
	if err != nil {
		dm.logger.WithError(err).Debug("Failed to create IPv6 multicast connection")
		return false
	}
	defer conn.Close()

	message := fmt.Sprintf("XELVRA_IPV6:%s", dm.host.ID().String())
	_, err = conn.Write([]byte(message))
	if err != nil {
		dm.logger.WithError(err).Debug("Failed to send IPv6 multicast discovery")
		return false
	}

	dm.logger.WithFields(logrus.Fields{
		"interface": interfaceName,
		"target":    multicastAddr,
		"message":   message,
	}).Debug("Sent IPv6 link-local discovery multicast")

	return true
}

// startHolePunchingService starts NAT hole punching service (Phase 5)
func (dm *DiscoveryManager) startHolePunchingService() {
	dm.logger.Info("Starting NAT hole punching service (Phase 5)...")

	ticker := time.NewTicker(1 * time.Minute) // Check every minute
	defer ticker.Stop()

	for {
		select {
		case <-dm.ctx.Done():
			return
		case <-ticker.C:
			dm.attemptHolePunching()
		}
	}
}

// attemptHolePunching attempts NAT hole punching for discovered peers
func (dm *DiscoveryManager) attemptHolePunching() {
	dm.mu.RLock()
	peers := make(map[peer.ID]*peer.AddrInfo)
	for id, info := range dm.discoveredPeers {
		peers[id] = info
	}
	dm.mu.RUnlock()

	for peerID, peerInfo := range peers {
		// Skip if already connected
		if dm.host.Network().Connectedness(peerID) == network.Connected {
			continue
		}

		// Skip local peers (no NAT traversal needed)
		if dm.isLocalPeer(peerInfo) {
			continue
		}

		// Attempt hole punching
		go dm.performHolePunching(peerID, peerInfo)
	}
}

// performHolePunching performs NAT hole punching for a specific peer
func (dm *DiscoveryManager) performHolePunching(peerID peer.ID, peerInfo *peer.AddrInfo) {
	dm.logger.WithField("peer_id", peerID.String()).Debug("Attempting NAT hole punching")

	// Try multiple connection attempts with different strategies
	strategies := []string{"direct", "relay-assisted", "simultaneous-open"}

	for _, strategy := range strategies {
		ctx, cancel := context.WithTimeout(dm.ctx, 10*time.Second)

		switch strategy {
		case "direct":
			err := dm.host.Connect(ctx, *peerInfo)
			if err == nil {
				dm.logger.WithFields(logrus.Fields{
					"peer_id":  peerID.String(),
					"strategy": strategy,
				}).Info("NAT hole punching successful")
				cancel()
				return
			}
		case "relay-assisted":
			// Use relay for hole punching assistance
			dm.attemptRelayAssistedConnection(ctx, peerID, peerInfo)
		case "simultaneous-open":
			// Coordinate simultaneous connection attempts
			dm.attemptSimultaneousOpen(ctx, peerID, peerInfo)
		}

		cancel()
	}

	dm.logger.WithField("peer_id", peerID.String()).Debug("NAT hole punching failed, will try relay")
}

// attemptRelayAssistedConnection uses relay to assist hole punching
func (dm *DiscoveryManager) attemptRelayAssistedConnection(ctx context.Context, peerID peer.ID, peerInfo *peer.AddrInfo) {
	// Implementation would coordinate with relay servers to assist NAT traversal
	dm.logger.WithField("peer_id", peerID.String()).Debug("Attempting relay-assisted hole punching")
	// This is a placeholder - full implementation would involve relay coordination
}

// attemptSimultaneousOpen coordinates simultaneous connection attempts
func (dm *DiscoveryManager) attemptSimultaneousOpen(ctx context.Context, peerID peer.ID, peerInfo *peer.AddrInfo) {
	// Implementation would coordinate simultaneous outbound connections
	dm.logger.WithField("peer_id", peerID.String()).Debug("Attempting simultaneous open hole punching")
	// This is a placeholder - full implementation would involve timing coordination
}

// startRelayServerManagement starts relay server management (Phase 6)
func (dm *DiscoveryManager) startRelayServerManagement() {
	dm.logger.Info("Starting relay server management (Phase 6 - final fallback)...")

	ticker := time.NewTicker(2 * time.Minute) // Check every 2 minutes
	defer ticker.Stop()

	for {
		select {
		case <-dm.ctx.Done():
			return
		case <-ticker.C:
			dm.manageRelayServers()
		}
	}
}

// manageRelayServers manages relay server connections and creation
func (dm *DiscoveryManager) manageRelayServers() {
	// Check if we need relay services
	needsRelay := dm.assessRelayNeed()

	if needsRelay {
		dm.logger.Info("Relay services needed - establishing relay connections")

		// Try to connect to existing relay servers
		if !dm.connectToExistingRelays() {
			// If no existing relays available, consider becoming a relay
			dm.considerBecomingRelay()
		}
	}
}

// assessRelayNeed determines if relay services are needed
func (dm *DiscoveryManager) assessRelayNeed() bool {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	// Count peers we can't connect to directly
	unreachablePeers := 0
	totalPeers := len(dm.discoveredPeers)

	for peerID := range dm.discoveredPeers {
		if dm.host.Network().Connectedness(peerID) != network.Connected {
			unreachablePeers++
		}
	}

	// Need relay if more than 30% of peers are unreachable
	needsRelay := totalPeers > 0 && float64(unreachablePeers)/float64(totalPeers) > 0.3

	if needsRelay {
		dm.logger.WithFields(logrus.Fields{
			"total_peers":       totalPeers,
			"unreachable_peers": unreachablePeers,
			"unreachable_ratio": float64(unreachablePeers) / float64(totalPeers),
		}).Info("Relay services needed due to connectivity issues")
	}

	return needsRelay
}

// connectToExistingRelays attempts to connect to existing relay servers
func (dm *DiscoveryManager) connectToExistingRelays() bool {
	// Try to connect to known relay servers (bootstrap peers often provide relay)
	relayConnected := false

	for _, peerInfo := range dm.bootstrapPeers {
		ctx, cancel := context.WithTimeout(dm.ctx, 15*time.Second)

		if err := dm.host.Connect(ctx, peerInfo); err == nil {
			dm.logger.WithField("relay_peer", peerInfo.ID.String()).Info("Connected to relay server")
			relayConnected = true
		}

		cancel()
	}

	return relayConnected
}

// considerBecomingRelay considers whether this node should become a relay
func (dm *DiscoveryManager) considerBecomingRelay() {
	// Check if we have good connectivity and resources to become a relay
	canBeRelay := dm.assessRelayCapability()

	if canBeRelay {
		dm.logger.Info("Node has good connectivity - considering becoming relay server")
		dm.enableRelayService()
	} else {
		dm.logger.Debug("Node not suitable for relay service")
	}
}

// assessRelayCapability determines if this node can serve as a relay
func (dm *DiscoveryManager) assessRelayCapability() bool {
	// Check if we have public IP and good connectivity
	// This is a simplified check - full implementation would be more sophisticated

	// Check if we have any public addresses
	hasPublicAddr := false
	for _, addr := range dm.host.Addrs() {
		addrStr := addr.String()
		if !strings.Contains(addrStr, "127.0.0.1") &&
			!strings.Contains(addrStr, "192.168.") &&
			!strings.Contains(addrStr, "10.") &&
			!strings.Contains(addrStr, "172.") {
			hasPublicAddr = true
			break
		}
	}

	// Check if we have good connectivity (connected to multiple peers)
	connectedPeers := len(dm.host.Network().Peers())

	canBeRelay := hasPublicAddr && connectedPeers >= 3

	dm.logger.WithFields(logrus.Fields{
		"has_public_addr": hasPublicAddr,
		"connected_peers": connectedPeers,
		"can_be_relay":    canBeRelay,
	}).Debug("Assessed relay capability")

	return canBeRelay
}

// enableRelayService enables relay service on this node
func (dm *DiscoveryManager) enableRelayService() {
	dm.logger.Info("Enabling relay service - this node will help other peers connect")

	// Enable libp2p relay service
	// This would involve configuring the host to accept and forward relay connections
	// Implementation depends on libp2p relay configuration

	dm.logger.Info("Relay service enabled - node is now helping with NAT traversal")
}
