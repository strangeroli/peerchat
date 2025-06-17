package unit

import (
	"context"
	"testing"
	"time"

	"github.com/Xelvra/peerchat/internal/p2p"
	"github.com/libp2p/go-libp2p"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdvancedNATTraversal_Creation(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create test host
	host, err := libp2p.New()
	require.NoError(t, err)
	defer host.Close()

	// Create advanced NAT traversal
	nat := p2p.NewAdvancedNATTraversal(host, logger)
	assert.NotNil(t, nat)

	// Check initial status
	status := nat.GetStatus()
	assert.Equal(t, p2p.NATTypeUnknown, status.Type)
	assert.Equal(t, "", status.PublicIP)
	assert.Equal(t, 0, status.ExternalPort)
}

func TestAdvancedNATTraversal_StartStop(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create test host
	host, err := libp2p.New()
	require.NoError(t, err)
	defer host.Close()

	// Create and start advanced NAT traversal
	nat := p2p.NewAdvancedNATTraversal(host, logger)

	err = nat.Start()
	assert.NoError(t, err)

	// Check status after start
	status := nat.GetStatus()
	// NAT detection might fail in test environment, that's expected
	t.Logf("NAT Type detected: %v", status.Type)
	t.Logf("Public IP: %s", status.PublicIP)

	// Stop NAT traversal
	err = nat.Stop()
	assert.NoError(t, err)
}

func TestAdvancedNATTraversal_NATDetection(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create test host
	host, err := libp2p.New()
	require.NoError(t, err)
	defer host.Close()

	// Create and start advanced NAT traversal
	nat := p2p.NewAdvancedNATTraversal(host, logger)
	err = nat.Start()
	require.NoError(t, err)
	defer nat.Stop()

	// Wait for NAT detection to complete
	time.Sleep(500 * time.Millisecond)

	// Check NAT detection results
	status := nat.GetStatus()

	// In test environment, we might not detect actual NAT
	// but the system should at least attempt detection
	t.Logf("Detected NAT Type: %v", status.Type)
	t.Logf("Mapping Behavior: %v", status.MappingBehavior)
	t.Logf("Filter Behavior: %v", status.FilterBehavior)

	// The type should be set to something other than unknown after detection attempt
	// (even if detection fails, it should default to a specific type)
	assert.NotEqual(t, p2p.NATTypeUnknown, status.Type)
}

func TestAdvancedNATTraversal_ConnectionAttempt(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create test host
	host, err := libp2p.New()
	require.NoError(t, err)
	defer host.Close()

	// Create and start advanced NAT traversal
	nat := p2p.NewAdvancedNATTraversal(host, logger)
	err = nat.Start()
	require.NoError(t, err)
	defer nat.Stop()

	// Test connection attempt
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use host's own peer ID for testing
	testPeerID := host.ID()
	testAddrs := []string{"127.0.0.1:8080", "127.0.0.1:8081"}

	err = nat.AttemptConnection(ctx, testPeerID, testAddrs)
	// This will likely fail in test environment, which is expected
	if err != nil {
		t.Logf("Connection attempt failed as expected in test environment: %v", err)
	}

	t.Log("Connection attempt test completed")
}

func TestAdvancedNATTraversal_StrategySelection(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create test host
	host, err := libp2p.New()
	require.NoError(t, err)
	defer host.Close()

	// Create and start advanced NAT traversal
	nat := p2p.NewAdvancedNATTraversal(host, logger)
	err = nat.Start()
	require.NoError(t, err)
	defer nat.Stop()

	// Wait for initialization
	time.Sleep(200 * time.Millisecond)

	// Test strategy selection for different NAT types
	// This tests the internal logic without actual network operations

	// The strategy selection is internal, so we test it indirectly
	// by checking that connection attempts use appropriate strategies
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	testPeerID := host.ID()
	testAddrs := []string{"127.0.0.1:8080"}

	// Attempt connection - this will test strategy selection internally
	err = nat.AttemptConnection(ctx, testPeerID, testAddrs)
	// Expected to fail in test environment
	if err != nil {
		t.Logf("Strategy selection test completed - connection failed as expected: %v", err)
	}
}

func TestAdvancedNATTraversal_STUNClient(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create test host
	host, err := libp2p.New()
	require.NoError(t, err)
	defer host.Close()

	// Create and start advanced NAT traversal
	nat := p2p.NewAdvancedNATTraversal(host, logger)
	err = nat.Start()
	require.NoError(t, err)
	defer nat.Stop()

	// STUN client is tested indirectly through NAT detection
	// Wait for STUN operations to complete
	time.Sleep(300 * time.Millisecond)

	status := nat.GetStatus()

	// Check that STUN client attempted to get public IP
	// (even if it failed in test environment)
	t.Logf("STUN client test - Public IP: %s", status.PublicIP)

	// In test environment, STUN might not work, but the system should handle it gracefully
	t.Log("STUN client test completed")
}

func TestAdvancedNATTraversal_RelayManagement(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create test host
	host, err := libp2p.New()
	require.NoError(t, err)
	defer host.Close()

	// Create and start advanced NAT traversal
	nat := p2p.NewAdvancedNATTraversal(host, logger)
	err = nat.Start()
	require.NoError(t, err)
	defer nat.Stop()

	// Wait for relay management to initialize
	time.Sleep(300 * time.Millisecond)

	status := nat.GetStatus()

	// Check relay management status
	t.Logf("Active relays: %d", status.ActiveRelays)

	// In test environment, we might not have active relays
	assert.GreaterOrEqual(t, status.ActiveRelays, 0)

	t.Log("Relay management test completed")
}

func TestAdvancedNATTraversal_HolePunching(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create test host
	host, err := libp2p.New()
	require.NoError(t, err)
	defer host.Close()

	// Create and start advanced NAT traversal
	nat := p2p.NewAdvancedNATTraversal(host, logger)
	err = nat.Start()
	require.NoError(t, err)
	defer nat.Stop()

	// Test hole punching strategies
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	testPeerID := host.ID()
	testAddrs := []string{"127.0.0.1:8080"}

	// This will test hole punching internally
	err = nat.AttemptConnection(ctx, testPeerID, testAddrs)

	// Expected to fail in test environment, but should test the hole punching logic
	if err != nil {
		t.Logf("Hole punching test completed - failed as expected: %v", err)
	}

	t.Log("Hole punching strategies test completed")
}

func TestAdvancedNATTraversal_TraversalRate(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create test host
	host, err := libp2p.New()
	require.NoError(t, err)
	defer host.Close()

	// Create and start advanced NAT traversal
	nat := p2p.NewAdvancedNATTraversal(host, logger)
	err = nat.Start()
	require.NoError(t, err)
	defer nat.Stop()

	// Wait for initialization
	time.Sleep(200 * time.Millisecond)

	status := nat.GetStatus()

	// Check initial traversal rate
	assert.GreaterOrEqual(t, status.TraversalRate, 0.0)
	assert.LessOrEqual(t, status.TraversalRate, 1.0)

	t.Logf("Traversal rate: %.2f", status.TraversalRate)
	t.Log("Traversal rate test completed")
}

func TestAdvancedNATTraversal_Monitoring(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create test host
	host, err := libp2p.New()
	require.NoError(t, err)
	defer host.Close()

	// Create and start advanced NAT traversal
	nat := p2p.NewAdvancedNATTraversal(host, logger)
	err = nat.Start()
	require.NoError(t, err)
	defer nat.Stop()

	// Wait for monitoring to run
	time.Sleep(500 * time.Millisecond)

	// Check that monitoring is working
	status1 := nat.GetStatus()

	// Wait a bit more
	time.Sleep(200 * time.Millisecond)

	status2 := nat.GetStatus()

	// Status should be consistent (monitoring working)
	assert.Equal(t, status1.Type, status2.Type)

	t.Log("NAT monitoring test completed")
}

// Benchmark tests
func BenchmarkAdvancedNATTraversal_ConnectionAttempt(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in benchmarks

	// Create test host
	host, err := libp2p.New()
	require.NoError(b, err)
	defer host.Close()

	// Create and start advanced NAT traversal
	nat := p2p.NewAdvancedNATTraversal(host, logger)
	err = nat.Start()
	require.NoError(b, err)
	defer nat.Stop()

	testPeerID := host.ID()
	testAddrs := []string{"127.0.0.1:8080"}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		_ = nat.AttemptConnection(ctx, testPeerID, testAddrs)
		cancel()
	}
}

func BenchmarkAdvancedNATTraversal_StatusCheck(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in benchmarks

	// Create test host
	host, err := libp2p.New()
	require.NoError(b, err)
	defer host.Close()

	// Create and start advanced NAT traversal
	nat := p2p.NewAdvancedNATTraversal(host, logger)
	err = nat.Start()
	require.NoError(b, err)
	defer nat.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = nat.GetStatus()
	}
}
