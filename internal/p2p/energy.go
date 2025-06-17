package p2p

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// EnergyManager handles energy optimization strategies
// Implements energy optimization from tmp/Energetická Optimalizace.md
type EnergyManager struct {
	logger *logrus.Logger
	ctx    context.Context
	cancel context.CancelFunc

	// Energy monitoring
	mu              sync.RWMutex
	cpuUsage        float64
	memoryUsage     int64
	networkActivity int64
	lastMeasurement time.Time
	energyProfile   *EnergyProfile

	// Adaptive polling
	dhtPollInterval   time.Duration
	heartbeatInterval time.Duration
	batteryLevel      float64 // 0.0 - 1.0

	// Deep sleep mode
	deepSleepMode      bool
	deepSleepThreshold float64 // Battery level threshold for deep sleep

	// Performance targets from README
	targetIdleMemoryMB   int
	targetIdleCPUPercent float64
	targetLatencyMs      int
}

// EnergyProfile represents current energy consumption profile
type EnergyProfile struct {
	CPUUsagePercent   float64   `json:"cpu_usage_percent"`
	MemoryUsageMB     int64     `json:"memory_usage_mb"`
	NetworkActivityMB float64   `json:"network_activity_mb"`
	EstimatedPowerMW  float64   `json:"estimated_power_mw"`
	BatteryLevel      float64   `json:"battery_level"`
	DeepSleepActive   bool      `json:"deep_sleep_active"`
	LastUpdated       time.Time `json:"last_updated"`

	// Optimization status
	AdaptivePollingActive bool   `json:"adaptive_polling_active"`
	DHTPollInterval       string `json:"dht_poll_interval"`
	HeartbeatInterval     string `json:"heartbeat_interval"`
}

// NewEnergyManager creates a new energy manager
func NewEnergyManager(ctx context.Context, logger *logrus.Logger) *EnergyManager {
	energyCtx, cancel := context.WithCancel(ctx)

	return &EnergyManager{
		logger:               logger,
		ctx:                  energyCtx,
		cancel:               cancel,
		dhtPollInterval:      2 * time.Minute,  // Default DHT polling
		heartbeatInterval:    30 * time.Second, // Default heartbeat
		batteryLevel:         1.0,              // Assume full battery initially
		deepSleepThreshold:   0.15,             // 15% battery threshold
		targetIdleMemoryMB:   20,               // From README
		targetIdleCPUPercent: 1.0,              // From README
		targetLatencyMs:      50,               // From README
		energyProfile: &EnergyProfile{
			LastUpdated: time.Now(),
		},
	}
}

// Start begins energy monitoring and optimization
func (em *EnergyManager) Start() error {
	em.logger.Info("Starting energy optimization manager...")

	// Start monitoring goroutine
	go em.monitorEnergyUsage()

	// Start adaptive polling optimization
	go em.adaptivePollingOptimization()

	em.logger.WithFields(logrus.Fields{
		"target_idle_memory_mb":   em.targetIdleMemoryMB,
		"target_idle_cpu_percent": em.targetIdleCPUPercent,
		"target_latency_ms":       em.targetLatencyMs,
		"deep_sleep_threshold":    em.deepSleepThreshold,
	}).Info("Energy optimization started with performance targets")

	return nil
}

// Stop stops energy monitoring
func (em *EnergyManager) Stop() error {
	em.logger.Info("Stopping energy optimization manager...")
	em.cancel()
	return nil
}

// GetEnergyProfile returns current energy profile
func (em *EnergyManager) GetEnergyProfile() *EnergyProfile {
	em.mu.RLock()
	defer em.mu.RUnlock()

	// Create a copy to avoid race conditions
	profile := *em.energyProfile
	return &profile
}

// SetBatteryLevel updates the battery level (0.0 - 1.0)
func (em *EnergyManager) SetBatteryLevel(level float64) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.batteryLevel = level
	em.energyProfile.BatteryLevel = level

	// Check if we should enter deep sleep mode
	if level <= em.deepSleepThreshold && !em.deepSleepMode {
		em.enterDeepSleepMode()
	} else if level > em.deepSleepThreshold && em.deepSleepMode {
		em.exitDeepSleepMode()
	}
}

// monitorEnergyUsage continuously monitors energy usage
func (em *EnergyManager) monitorEnergyUsage() {
	ticker := time.NewTicker(10 * time.Second) // Monitor every 10 seconds
	defer ticker.Stop()

	for {
		select {
		case <-em.ctx.Done():
			return
		case <-ticker.C:
			em.measureEnergyUsage()
		}
	}
}

// measureEnergyUsage measures current energy usage
func (em *EnergyManager) measureEnergyUsage() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	em.mu.Lock()
	defer em.mu.Unlock()

	// Update memory usage (convert bytes to MB)
	em.memoryUsage = int64(memStats.Alloc / 1024 / 1024)

	// Estimate CPU usage (simplified - would need more sophisticated monitoring in production)
	// For now, we'll use a basic heuristic based on goroutines and memory allocation rate
	numGoroutines := runtime.NumGoroutine()
	em.cpuUsage = float64(numGoroutines) * 0.1 // Rough estimate

	// Update energy profile
	em.energyProfile.CPUUsagePercent = em.cpuUsage
	em.energyProfile.MemoryUsageMB = em.memoryUsage
	em.energyProfile.LastUpdated = time.Now()

	// Estimate power consumption (simplified model)
	// Base consumption + CPU factor + Memory factor
	basePowerMW := 5.0                             // Base power consumption
	cpuPowerMW := em.cpuUsage * 2.0                // CPU contribution
	memoryPowerMW := float64(em.memoryUsage) * 0.1 // Memory contribution

	em.energyProfile.EstimatedPowerMW = basePowerMW + cpuPowerMW + memoryPowerMW

	// Log if usage exceeds targets
	if em.memoryUsage > int64(em.targetIdleMemoryMB) {
		em.logger.WithFields(logrus.Fields{
			"current_memory_mb": em.memoryUsage,
			"target_memory_mb":  em.targetIdleMemoryMB,
		}).Warn("Memory usage exceeds target")
	}

	if em.cpuUsage > em.targetIdleCPUPercent {
		em.logger.WithFields(logrus.Fields{
			"current_cpu_percent": em.cpuUsage,
			"target_cpu_percent":  em.targetIdleCPUPercent,
		}).Warn("CPU usage exceeds target")
	}
}

// adaptivePollingOptimization adjusts polling intervals based on battery and activity
func (em *EnergyManager) adaptivePollingOptimization() {
	ticker := time.NewTicker(1 * time.Minute) // Adjust every minute
	defer ticker.Stop()

	for {
		select {
		case <-em.ctx.Done():
			return
		case <-ticker.C:
			em.optimizePollingIntervals()
		}
	}
}

// optimizePollingIntervals adjusts polling intervals based on current conditions
func (em *EnergyManager) optimizePollingIntervals() {
	em.mu.Lock()
	defer em.mu.Unlock()

	// Adjust intervals based on battery level
	batteryFactor := em.batteryLevel

	// Base intervals
	baseDHTInterval := 2 * time.Minute
	baseHeartbeatInterval := 30 * time.Second

	// Adjust based on battery level
	if batteryFactor < 0.2 { // Low battery
		em.dhtPollInterval = baseDHTInterval * 5         // 10 minutes
		em.heartbeatInterval = baseHeartbeatInterval * 4 // 2 minutes
	} else if batteryFactor < 0.5 { // Medium battery
		em.dhtPollInterval = baseDHTInterval * 2         // 4 minutes
		em.heartbeatInterval = baseHeartbeatInterval * 2 // 1 minute
	} else { // Good battery
		em.dhtPollInterval = baseDHTInterval
		em.heartbeatInterval = baseHeartbeatInterval
	}

	// Update energy profile
	em.energyProfile.AdaptivePollingActive = true
	em.energyProfile.DHTPollInterval = em.dhtPollInterval.String()
	em.energyProfile.HeartbeatInterval = em.heartbeatInterval.String()

	em.logger.WithFields(logrus.Fields{
		"battery_level":      em.batteryLevel,
		"dht_poll_interval":  em.dhtPollInterval,
		"heartbeat_interval": em.heartbeatInterval,
	}).Debug("Adaptive polling intervals updated")
}

// enterDeepSleepMode activates deep sleep mode for energy conservation
func (em *EnergyManager) enterDeepSleepMode() {
	em.deepSleepMode = true
	em.energyProfile.DeepSleepActive = true

	// Drastically reduce polling intervals
	em.dhtPollInterval = 10 * time.Minute  // Very infrequent DHT queries
	em.heartbeatInterval = 5 * time.Minute // Very infrequent heartbeats

	em.logger.WithFields(logrus.Fields{
		"battery_level": em.batteryLevel,
		"threshold":     em.deepSleepThreshold,
	}).Warn("Entering deep sleep mode for energy conservation")
}

// exitDeepSleepMode deactivates deep sleep mode
func (em *EnergyManager) exitDeepSleepMode() {
	em.deepSleepMode = false
	em.energyProfile.DeepSleepActive = false

	// Restore normal polling intervals
	em.optimizePollingIntervals()

	em.logger.WithField("battery_level", em.batteryLevel).Info("Exiting deep sleep mode")
}

// GetDHTPollingInterval returns current DHT polling interval
func (em *EnergyManager) GetDHTPollingInterval() time.Duration {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return em.dhtPollInterval
}

// GetHeartbeatInterval returns current heartbeat interval
func (em *EnergyManager) GetHeartbeatInterval() time.Duration {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return em.heartbeatInterval
}

// IsDeepSleepMode returns true if in deep sleep mode
func (em *EnergyManager) IsDeepSleepMode() bool {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return em.deepSleepMode
}
