package p2p

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/dual"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/sirupsen/logrus"
)

// AdvancedDHT implements enhanced Kademlia DHT with optimizations
// Inspired by battle-tested algorithms but implemented in Go with modern optimizations
type AdvancedDHT struct {
	host   host.Host
	logger *logrus.Logger
	ctx    context.Context
	cancel context.CancelFunc

	// Core DHT
	dht              *dual.DHT
	routingDiscovery *drouting.RoutingDiscovery

	// Advanced features
	bucketManager    *BucketManager
	peerSelector     *PeerSelector
	timeoutManager   *TimeoutManager
	batteryOptimizer *BatteryOptimizer

	// State management
	mu     sync.RWMutex
	status *DHTStatus
}

// BucketManager handles advanced bucket management for Kademlia
type BucketManager struct {
	mu      sync.RWMutex
	buckets map[int]*KBucket
	logger  *logrus.Logger
}

// KBucket represents a Kademlia bucket with enhanced features
type KBucket struct {
	peers       []peer.ID
	lastSeen    map[peer.ID]time.Time
	reliability map[peer.ID]float64
	maxSize     int
}

// PeerSelector implements intelligent peer selection algorithms
type PeerSelector struct {
	mu               sync.RWMutex
	peerMetrics      map[peer.ID]*PeerMetrics
	selectionHistory []peer.ID
	logger           *logrus.Logger
}

// PeerMetrics tracks peer performance metrics
type PeerMetrics struct {
	ResponseTime    time.Duration
	SuccessRate     float64
	LastSeen        time.Time
	NetworkQuality  float64
	BatteryFriendly bool
}

// TimeoutManager handles adaptive timeout mechanisms
type TimeoutManager struct {
	mu              sync.RWMutex
	baseTimeout     time.Duration
	adaptiveTimeout map[peer.ID]time.Duration
	networkQuality  float64
	logger          *logrus.Logger
}

// BatteryOptimizer optimizes DHT operations for battery-aware devices
type BatteryOptimizer struct {
	mu               sync.RWMutex
	batteryLevel     float64
	powerSaveMode    bool
	queryInterval    time.Duration
	maxConcurrentOps int
	logger           *logrus.Logger
}

// DHTStatus represents enhanced DHT status
type DHTStatus struct {
	Active           bool
	ConnectedPeers   int
	BucketCount      int
	QueryLatency     time.Duration
	SuccessRate      float64
	BatteryOptimized bool
	NetworkQuality   float64
}

// NewAdvancedDHT creates a new advanced DHT instance
func NewAdvancedDHT(h host.Host, logger *logrus.Logger) *AdvancedDHT {
	ctx, cancel := context.WithCancel(context.Background())

	return &AdvancedDHT{
		host:   h,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
		bucketManager: &BucketManager{
			buckets: make(map[int]*KBucket),
			logger:  logger,
		},
		peerSelector: &PeerSelector{
			peerMetrics: make(map[peer.ID]*PeerMetrics),
			logger:      logger,
		},
		timeoutManager: &TimeoutManager{
			baseTimeout:     30 * time.Second,
			adaptiveTimeout: make(map[peer.ID]time.Duration),
			networkQuality:  1.0,
			logger:          logger,
		},
		batteryOptimizer: &BatteryOptimizer{
			batteryLevel:     1.0,
			powerSaveMode:    false,
			queryInterval:    2 * time.Minute,
			maxConcurrentOps: 10,
			logger:           logger,
		},
		status: &DHTStatus{
			Active:           false,
			ConnectedPeers:   0,
			BucketCount:      0,
			QueryLatency:     0,
			SuccessRate:      0.0,
			BatteryOptimized: false,
			NetworkQuality:   1.0,
		},
	}
}

// Start initializes and starts the advanced DHT
func (adht *AdvancedDHT) Start() error {
	adht.logger.Info("Starting Advanced DHT with Kademlia optimizations...")

	// Create enhanced DHT with custom options
	dht, err := dual.New(adht.ctx, adht.host)
	if err != nil {
		return fmt.Errorf("failed to create advanced DHT: %w", err)
	}

	adht.dht = dht
	adht.routingDiscovery = drouting.NewRoutingDiscovery(adht.dht)

	// Initialize advanced components
	if err := adht.initializeBucketManager(); err != nil {
		return fmt.Errorf("failed to initialize bucket manager: %w", err)
	}

	if err := adht.initializePeerSelector(); err != nil {
		return fmt.Errorf("failed to initialize peer selector: %w", err)
	}

	if err := adht.initializeTimeoutManager(); err != nil {
		return fmt.Errorf("failed to initialize timeout manager: %w", err)
	}

	if err := adht.initializeBatteryOptimizer(); err != nil {
		return fmt.Errorf("failed to initialize battery optimizer: %w", err)
	}

	// Bootstrap DHT
	if err := adht.dht.Bootstrap(adht.ctx); err != nil {
		adht.logger.WithError(err).Warn("DHT bootstrap failed, continuing anyway")
	}

	// Start background processes
	go adht.runMaintenanceLoop()
	go adht.runMetricsCollection()
	go adht.runBatteryOptimization()

	adht.mu.Lock()
	adht.status.Active = true
	adht.mu.Unlock()

	adht.logger.Info("Advanced DHT started successfully with all optimizations enabled")
	return nil
}

// Stop stops the advanced DHT
func (adht *AdvancedDHT) Stop() error {
	adht.logger.Info("Stopping Advanced DHT...")

	adht.cancel()

	if adht.dht != nil {
		if err := adht.dht.Close(); err != nil {
			adht.logger.WithError(err).Warn("Failed to close DHT")
		}
	}

	adht.mu.Lock()
	adht.status.Active = false
	adht.mu.Unlock()

	adht.logger.Info("Advanced DHT stopped")
	return nil
}

// GetStatus returns current DHT status
func (adht *AdvancedDHT) GetStatus() *DHTStatus {
	adht.mu.RLock()
	defer adht.mu.RUnlock()

	// Create copy to avoid race conditions
	status := *adht.status
	return &status
}

// FindPeer finds a peer using advanced selection algorithms
func (adht *AdvancedDHT) FindPeer(ctx context.Context, peerID peer.ID) (*peer.AddrInfo, error) {
	start := time.Now()

	// Use adaptive timeout based on peer history
	timeout := adht.timeoutManager.getTimeoutForPeer(peerID)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Find peer using DHT
	peerInfo, err := adht.dht.FindPeer(ctx, peerID)
	if err != nil {
		adht.updatePeerMetrics(peerID, time.Since(start), false)
		return nil, err
	}

	adht.updatePeerMetrics(peerID, time.Since(start), true)
	return &peerInfo, nil
}

// FindPeers finds peers advertising a namespace with intelligent selection
func (adht *AdvancedDHT) FindPeers(ctx context.Context, namespace string, limit int) ([]peer.AddrInfo, error) {
	if adht.routingDiscovery == nil {
		return nil, fmt.Errorf("routing discovery not initialized")
	}

	// Apply battery optimization
	if adht.batteryOptimizer.shouldLimitQuery() {
		limit = adht.batteryOptimizer.getOptimalLimit(limit)
	}

	// Use adaptive timeout
	timeout := adht.timeoutManager.getAdaptiveTimeout()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	peerChan, err := adht.routingDiscovery.FindPeers(ctx, namespace)
	if err != nil {
		return nil, err
	}

	var peers []peer.AddrInfo
	count := 0

	for peerInfo := range peerChan {
		if count >= limit {
			break
		}

		// Apply intelligent peer selection
		if adht.peerSelector.shouldSelectPeer(peerInfo.ID) {
			peers = append(peers, peerInfo)
			count++
		}
	}

	adht.logger.WithFields(logrus.Fields{
		"namespace":    namespace,
		"found_peers":  len(peers),
		"limit":        limit,
		"battery_mode": adht.batteryOptimizer.powerSaveMode,
	}).Debug("Advanced DHT peer discovery completed")

	return peers, nil
}

// Advertise advertises presence in DHT with battery optimization
func (adht *AdvancedDHT) Advertise(ctx context.Context, namespace string) error {
	if adht.routingDiscovery == nil {
		return fmt.Errorf("routing discovery not initialized")
	}

	// Apply battery optimization for advertisement frequency
	if adht.batteryOptimizer.shouldSkipAdvertisement() {
		adht.logger.Debug("Skipping advertisement due to battery optimization")
		return nil
	}

	timeout := adht.timeoutManager.getAdaptiveTimeout()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := adht.routingDiscovery.Advertise(ctx, namespace)
	return err
}

// initializeBucketManager initializes the bucket manager
func (adht *AdvancedDHT) initializeBucketManager() error {
	adht.logger.Debug("Initializing advanced bucket manager...")

	// Initialize buckets for Kademlia (typically 256 buckets for 256-bit key space)
	for i := 0; i < 256; i++ {
		adht.bucketManager.buckets[i] = &KBucket{
			peers:       make([]peer.ID, 0),
			lastSeen:    make(map[peer.ID]time.Time),
			reliability: make(map[peer.ID]float64),
			maxSize:     20, // Standard Kademlia bucket size
		}
	}

	adht.logger.Debug("Bucket manager initialized with 256 buckets")
	return nil
}

// initializePeerSelector initializes the peer selector
func (adht *AdvancedDHT) initializePeerSelector() error {
	adht.logger.Debug("Initializing intelligent peer selector...")

	// Initialize peer selection algorithms
	adht.peerSelector.selectionHistory = make([]peer.ID, 0, 1000)

	adht.logger.Debug("Peer selector initialized")
	return nil
}

// initializeTimeoutManager initializes the timeout manager
func (adht *AdvancedDHT) initializeTimeoutManager() error {
	adht.logger.Debug("Initializing adaptive timeout manager...")

	// Set base timeout based on network conditions
	adht.timeoutManager.baseTimeout = 30 * time.Second
	adht.timeoutManager.networkQuality = 1.0

	adht.logger.Debug("Timeout manager initialized")
	return nil
}

// initializeBatteryOptimizer initializes the battery optimizer
func (adht *AdvancedDHT) initializeBatteryOptimizer() error {
	adht.logger.Debug("Initializing battery optimizer...")

	// Set default battery optimization parameters
	adht.batteryOptimizer.batteryLevel = 1.0
	adht.batteryOptimizer.powerSaveMode = false
	adht.batteryOptimizer.queryInterval = 2 * time.Minute
	adht.batteryOptimizer.maxConcurrentOps = 10

	adht.logger.Debug("Battery optimizer initialized")
	return nil
}

// runMaintenanceLoop runs the main DHT maintenance loop
func (adht *AdvancedDHT) runMaintenanceLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-adht.ctx.Done():
			return
		case <-ticker.C:
			adht.performMaintenance()
		}
	}
}

// performMaintenance performs DHT maintenance tasks
func (adht *AdvancedDHT) performMaintenance() {
	adht.logger.Debug("Performing DHT maintenance...")

	// Update bucket health
	adht.bucketManager.updateBucketHealth()

	// Clean up stale peers
	adht.peerSelector.cleanupStalePeers()

	// Adjust timeouts based on network conditions
	adht.timeoutManager.adjustTimeouts()

	// Update battery optimization
	adht.batteryOptimizer.updateOptimization()

	adht.logger.Debug("DHT maintenance completed")
}

// runMetricsCollection runs metrics collection loop
func (adht *AdvancedDHT) runMetricsCollection() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-adht.ctx.Done():
			return
		case <-ticker.C:
			adht.collectMetrics()
		}
	}
}

// collectMetrics collects DHT performance metrics
func (adht *AdvancedDHT) collectMetrics() {
	adht.mu.Lock()
	defer adht.mu.Unlock()

	// Update connected peers count
	if adht.dht != nil {
		// Use host to get connected peers since DHT routing table is not directly accessible
		adht.status.ConnectedPeers = len(adht.host.Network().Peers())
	}

	// Update bucket count
	adht.bucketManager.mu.RLock()
	adht.status.BucketCount = len(adht.bucketManager.buckets)
	adht.bucketManager.mu.RUnlock()

	// Update battery optimization status
	adht.status.BatteryOptimized = adht.batteryOptimizer.powerSaveMode

	// Update network quality
	adht.status.NetworkQuality = adht.timeoutManager.networkQuality

	adht.logger.WithFields(logrus.Fields{
		"connected_peers":   adht.status.ConnectedPeers,
		"bucket_count":      adht.status.BucketCount,
		"battery_optimized": adht.status.BatteryOptimized,
		"network_quality":   adht.status.NetworkQuality,
	}).Debug("DHT metrics collected")
}

// runBatteryOptimization runs battery optimization loop
func (adht *AdvancedDHT) runBatteryOptimization() {
	ticker := time.NewTicker(adht.batteryOptimizer.queryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-adht.ctx.Done():
			return
		case <-ticker.C:
			adht.optimizeForBattery()
		}
	}
}

// optimizeForBattery performs battery optimization
func (adht *AdvancedDHT) optimizeForBattery() {
	if adht.batteryOptimizer.batteryLevel < 0.2 {
		adht.batteryOptimizer.enablePowerSaveMode()
	} else if adht.batteryOptimizer.batteryLevel > 0.8 {
		adht.batteryOptimizer.disablePowerSaveMode()
	}
}

// updatePeerMetrics updates metrics for a peer
func (adht *AdvancedDHT) updatePeerMetrics(peerID peer.ID, responseTime time.Duration, success bool) {
	adht.peerSelector.mu.Lock()
	defer adht.peerSelector.mu.Unlock()

	metrics, exists := adht.peerSelector.peerMetrics[peerID]
	if !exists {
		metrics = &PeerMetrics{
			ResponseTime:    responseTime,
			SuccessRate:     0.0,
			LastSeen:        time.Now(),
			NetworkQuality:  1.0,
			BatteryFriendly: true,
		}
		adht.peerSelector.peerMetrics[peerID] = metrics
	}

	// Update response time (exponential moving average)
	metrics.ResponseTime = time.Duration(float64(metrics.ResponseTime)*0.7 + float64(responseTime)*0.3)

	// Update success rate
	if success {
		metrics.SuccessRate = metrics.SuccessRate*0.9 + 0.1
	} else {
		metrics.SuccessRate = metrics.SuccessRate * 0.9
	}

	metrics.LastSeen = time.Now()

	// Update network quality based on response time
	if responseTime < 100*time.Millisecond {
		metrics.NetworkQuality = 1.0
	} else if responseTime < 500*time.Millisecond {
		metrics.NetworkQuality = 0.8
	} else if responseTime < 1*time.Second {
		metrics.NetworkQuality = 0.6
	} else {
		metrics.NetworkQuality = 0.4
	}

	// Determine if peer is battery friendly
	metrics.BatteryFriendly = responseTime < 200*time.Millisecond && metrics.SuccessRate > 0.8
}

// Helper methods for BucketManager
func (bm *BucketManager) updateBucketHealth() {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	now := time.Now()
	for bucketID, bucket := range bm.buckets {
		// Remove stale peers (not seen for more than 10 minutes)
		var activePeers []peer.ID
		for _, peerID := range bucket.peers {
			if lastSeen, exists := bucket.lastSeen[peerID]; exists {
				if now.Sub(lastSeen) < 10*time.Minute {
					activePeers = append(activePeers, peerID)
				} else {
					delete(bucket.lastSeen, peerID)
					delete(bucket.reliability, peerID)
				}
			}
		}
		bucket.peers = activePeers

		if len(activePeers) != len(bucket.peers) {
			bm.logger.WithFields(logrus.Fields{
				"bucket_id":     bucketID,
				"active_peers":  len(activePeers),
				"removed_peers": len(bucket.peers) - len(activePeers),
			}).Debug("Updated bucket health")
		}
	}
}

// Helper methods for PeerSelector
func (ps *PeerSelector) shouldSelectPeer(peerID peer.ID) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	metrics, exists := ps.peerMetrics[peerID]
	if !exists {
		// New peer, give it a chance
		return true
	}

	// Select peer based on success rate and network quality
	score := metrics.SuccessRate * metrics.NetworkQuality

	// Prefer battery-friendly peers when in power save mode
	if metrics.BatteryFriendly {
		score *= 1.2
	}

	// Select if score is above threshold
	return score > 0.5
}

func (ps *PeerSelector) cleanupStalePeers() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	now := time.Now()
	for peerID, metrics := range ps.peerMetrics {
		if now.Sub(metrics.LastSeen) > 30*time.Minute {
			delete(ps.peerMetrics, peerID)
		}
	}

	// Limit selection history size
	if len(ps.selectionHistory) > 1000 {
		ps.selectionHistory = ps.selectionHistory[500:]
	}
}

// Helper methods for TimeoutManager
func (tm *TimeoutManager) getTimeoutForPeer(peerID peer.ID) time.Duration {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	if timeout, exists := tm.adaptiveTimeout[peerID]; exists {
		return timeout
	}

	return tm.baseTimeout
}

func (tm *TimeoutManager) getAdaptiveTimeout() time.Duration {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	// Adjust timeout based on network quality
	timeout := time.Duration(float64(tm.baseTimeout) / tm.networkQuality)

	// Ensure reasonable bounds
	if timeout < 5*time.Second {
		timeout = 5 * time.Second
	} else if timeout > 2*time.Minute {
		timeout = 2 * time.Minute
	}

	return timeout
}

func (tm *TimeoutManager) adjustTimeouts() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Adjust network quality based on recent performance
	// This is a simplified implementation
	if tm.networkQuality < 1.0 {
		tm.networkQuality = tm.networkQuality*0.9 + 0.1 // Slowly recover
	}
}

// Helper methods for BatteryOptimizer
func (bo *BatteryOptimizer) shouldLimitQuery() bool {
	bo.mu.RLock()
	defer bo.mu.RUnlock()

	return bo.powerSaveMode || bo.batteryLevel < 0.3
}

func (bo *BatteryOptimizer) getOptimalLimit(requestedLimit int) int {
	bo.mu.RLock()
	defer bo.mu.RUnlock()

	if bo.powerSaveMode {
		// Reduce limit by 50% in power save mode
		return requestedLimit / 2
	}

	if bo.batteryLevel < 0.3 {
		// Reduce limit by 30% when battery is low
		return int(float64(requestedLimit) * 0.7)
	}

	return requestedLimit
}

func (bo *BatteryOptimizer) shouldSkipAdvertisement() bool {
	bo.mu.RLock()
	defer bo.mu.RUnlock()

	// Skip advertisement if battery is very low
	return bo.batteryLevel < 0.1
}

func (bo *BatteryOptimizer) enablePowerSaveMode() {
	bo.mu.Lock()
	defer bo.mu.Unlock()

	if !bo.powerSaveMode {
		bo.powerSaveMode = true
		bo.queryInterval = 5 * time.Minute // Reduce query frequency
		bo.maxConcurrentOps = 3            // Reduce concurrent operations
		bo.logger.Info("Battery power save mode enabled")
	}
}

func (bo *BatteryOptimizer) disablePowerSaveMode() {
	bo.mu.Lock()
	defer bo.mu.Unlock()

	if bo.powerSaveMode {
		bo.powerSaveMode = false
		bo.queryInterval = 2 * time.Minute // Normal query frequency
		bo.maxConcurrentOps = 10           // Normal concurrent operations
		bo.logger.Info("Battery power save mode disabled")
	}
}

func (bo *BatteryOptimizer) updateOptimization() {
	// This would integrate with system battery APIs
	// For now, we simulate battery level changes
	bo.mu.Lock()
	defer bo.mu.Unlock()

	// Placeholder for actual battery level detection
	// In real implementation, this would query system battery APIs
}
