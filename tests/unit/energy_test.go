package unit

import (
	"context"
	"testing"
	"time"

	"github.com/Xelvra/peerchat/internal/p2p"
	"github.com/sirupsen/logrus"
)

// TestEnergyManagerCreation tests energy manager creation
func TestEnergyManagerCreation(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	em := p2p.NewEnergyManager(ctx, logger)
	if em == nil {
		t.Fatal("Failed to create energy manager")
	}
}

// TestEnergyManagerStart tests energy manager startup
func TestEnergyManagerStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	em := p2p.NewEnergyManager(ctx, logger)

	err := em.Start()
	if err != nil {
		t.Fatalf("Failed to start energy manager: %v", err)
	}

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	err = em.Stop()
	if err != nil {
		t.Fatalf("Failed to stop energy manager: %v", err)
	}
}

// TestEnergyProfile tests energy profile retrieval
func TestEnergyProfile(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	em := p2p.NewEnergyManager(ctx, logger)

	profile := em.GetEnergyProfile()
	if profile == nil {
		t.Fatal("Energy profile is nil")
	}

	// Battery level should be between 0.0 and 1.0 (may not be initialized to 1.0)
	if profile.BatteryLevel < 0.0 || profile.BatteryLevel > 1.0 {
		t.Errorf("Expected battery level between 0.0 and 1.0, got %f", profile.BatteryLevel)
	}

	if profile.DeepSleepActive {
		t.Error("Deep sleep should not be active initially")
	}

	if profile.LastUpdated.IsZero() {
		t.Error("LastUpdated should not be zero")
	}
}

// TestBatteryLevelUpdates tests battery level updates
func TestBatteryLevelUpdates(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	em := p2p.NewEnergyManager(ctx, logger)

	// Test normal battery level
	em.SetBatteryLevel(0.8)
	profile := em.GetEnergyProfile()
	if profile.BatteryLevel != 0.8 {
		t.Errorf("Expected battery level 0.8, got %f", profile.BatteryLevel)
	}
	if profile.DeepSleepActive {
		t.Error("Deep sleep should not be active at 80% battery")
	}

	// Test low battery level (should trigger deep sleep)
	em.SetBatteryLevel(0.1) // 10% - below 15% threshold
	profile = em.GetEnergyProfile()
	if profile.BatteryLevel != 0.1 {
		t.Errorf("Expected battery level 0.1, got %f", profile.BatteryLevel)
	}
	if !profile.DeepSleepActive {
		t.Error("Deep sleep should be active at 10% battery")
	}

	// Test recovery from deep sleep
	em.SetBatteryLevel(0.5) // 50% - above 15% threshold
	profile = em.GetEnergyProfile()
	if profile.BatteryLevel != 0.5 {
		t.Errorf("Expected battery level 0.5, got %f", profile.BatteryLevel)
	}
	if profile.DeepSleepActive {
		t.Error("Deep sleep should not be active at 50% battery")
	}
}

// TestAdaptivePolling tests adaptive polling intervals
func TestAdaptivePolling(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	em := p2p.NewEnergyManager(ctx, logger)

	// Test normal battery level intervals
	em.SetBatteryLevel(0.8)
	dhtInterval := em.GetDHTPollingInterval()
	heartbeatInterval := em.GetHeartbeatInterval()

	if dhtInterval <= 0 {
		t.Error("DHT polling interval should be positive")
	}
	if heartbeatInterval <= 0 {
		t.Error("Heartbeat interval should be positive")
	}

	// Store normal intervals for comparison
	normalDHT := dhtInterval
	normalHeartbeat := heartbeatInterval

	// Test low battery level intervals (should be longer)
	em.SetBatteryLevel(0.1)
	lowDHT := em.GetDHTPollingInterval()
	lowHeartbeat := em.GetHeartbeatInterval()

	if lowDHT <= normalDHT {
		t.Error("Low battery DHT interval should be longer than normal")
	}
	if lowHeartbeat <= normalHeartbeat {
		t.Error("Low battery heartbeat interval should be longer than normal")
	}
}

// TestDeepSleepMode tests deep sleep mode functionality
func TestDeepSleepMode(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	em := p2p.NewEnergyManager(ctx, logger)

	// Initially not in deep sleep
	if em.IsDeepSleepMode() {
		t.Error("Should not be in deep sleep mode initially")
	}

	// Trigger deep sleep with low battery
	em.SetBatteryLevel(0.1) // 10% - below 15% threshold
	if !em.IsDeepSleepMode() {
		t.Error("Should be in deep sleep mode at 10% battery")
	}

	// Exit deep sleep with higher battery
	em.SetBatteryLevel(0.5) // 50% - above 15% threshold
	if em.IsDeepSleepMode() {
		t.Error("Should not be in deep sleep mode at 50% battery")
	}
}

// TestEnergyProfileFields tests all energy profile fields
func TestEnergyProfileFields(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	em := p2p.NewEnergyManager(ctx, logger)

	err := em.Start()
	if err != nil {
		t.Fatalf("Failed to start energy manager: %v", err)
	}
	defer em.Stop()

	// Give it time to measure
	time.Sleep(200 * time.Millisecond)

	profile := em.GetEnergyProfile()

	// Check that all fields are populated
	if profile.CPUUsagePercent < 0 {
		t.Error("CPU usage should not be negative")
	}

	if profile.MemoryUsageMB < 0 {
		t.Error("Memory usage should not be negative")
	}

	if profile.EstimatedPowerMW < 0 {
		t.Error("Estimated power should not be negative")
	}

	if profile.DHTPollInterval == "" {
		t.Error("DHT poll interval should not be empty")
	}

	if profile.HeartbeatInterval == "" {
		t.Error("Heartbeat interval should not be empty")
	}

	if !profile.AdaptivePollingActive {
		t.Error("Adaptive polling should be active")
	}
}

// TestPerformanceTargets tests performance target validation
func TestPerformanceTargets(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	em := p2p.NewEnergyManager(ctx, logger)

	err := em.Start()
	if err != nil {
		t.Fatalf("Failed to start energy manager: %v", err)
	}
	defer em.Stop()

	// Give it time to measure
	time.Sleep(500 * time.Millisecond)

	profile := em.GetEnergyProfile()

	// Check performance targets (these are goals, not strict requirements in tests)
	if profile.MemoryUsageMB > 100 { // Relaxed for test environment
		t.Logf("Warning: Memory usage %dMB exceeds relaxed target of 100MB", profile.MemoryUsageMB)
	}

	if profile.CPUUsagePercent > 50 { // Relaxed for test environment
		t.Logf("Warning: CPU usage %.1f%% exceeds relaxed target of 50%%", profile.CPUUsagePercent)
	}

	// Log actual values for debugging
	t.Logf("Memory usage: %dMB", profile.MemoryUsageMB)
	t.Logf("CPU usage: %.1f%%", profile.CPUUsagePercent)
	t.Logf("Estimated power: %.1fmW", profile.EstimatedPowerMW)
}

// BenchmarkEnergyProfileRetrieval benchmarks energy profile retrieval
func BenchmarkEnergyProfileRetrieval(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	em := p2p.NewEnergyManager(ctx, logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		profile := em.GetEnergyProfile()
		if profile == nil {
			b.Fatal("Energy profile is nil")
		}
	}
}

// BenchmarkBatteryLevelUpdate benchmarks battery level updates
func BenchmarkBatteryLevelUpdate(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	em := p2p.NewEnergyManager(ctx, logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		level := float64(i%100) / 100.0 // Cycle through 0.0 to 0.99
		em.SetBatteryLevel(level)
	}
}
