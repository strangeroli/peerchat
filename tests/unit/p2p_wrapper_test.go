package unit

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Xelvra/peerchat/internal/p2p"
)

func TestP2PWrapper_SimulationMode(t *testing.T) {
	ctx := context.Background()
	wrapper := p2p.NewP2PWrapper(ctx, true) // Force simulation mode

	// Test start in simulation mode
	err := wrapper.Start()
	if err != nil {
		t.Fatalf("Failed to start wrapper in simulation mode: %v", err)
	}
	defer wrapper.Stop()

	// Verify simulation mode is active
	if !wrapper.IsUsingSimulation() {
		t.Error("Expected wrapper to be in simulation mode")
	}

	// Test node info in simulation mode
	nodeInfo := wrapper.GetNodeInfo()
	if nodeInfo == nil {
		t.Fatal("Expected node info to be available")
	}

	if nodeInfo.PeerID != "12D3KooWSimulatedPeerID..." {
		t.Errorf("Expected simulated peer ID, got: %s", nodeInfo.PeerID)
	}

	if nodeInfo.DID != "did:xelvra:simulated..." {
		t.Errorf("Expected simulated DID, got: %s", nodeInfo.DID)
	}

	if !nodeInfo.IsRunning {
		t.Error("Expected node to be running in simulation mode")
	}

	// Test message sending in simulation mode
	err = wrapper.SendMessage("test-peer", "Hello World")
	if err != nil {
		t.Errorf("Failed to send message in simulation mode: %v", err)
	}
}

func TestP2PWrapper_RealMode(t *testing.T) {
	ctx := context.Background()
	wrapper := p2p.NewP2PWrapper(ctx, false) // Try real P2P first

	// Test start (may fallback to simulation if real P2P fails)
	err := wrapper.Start()
	if err != nil {
		t.Fatalf("Failed to start wrapper: %v", err)
	}
	defer wrapper.Stop()

	// Test node info
	nodeInfo := wrapper.GetNodeInfo()
	if nodeInfo == nil {
		t.Fatal("Expected node info to be available")
	}

	if !nodeInfo.IsRunning {
		t.Error("Expected node to be running")
	}

	// Verify we have some peer ID (real or simulated)
	if nodeInfo.PeerID == "" {
		t.Error("Expected peer ID to be set")
	}

	// Verify we have some DID (real or simulated)
	if nodeInfo.DID == "" {
		t.Error("Expected DID to be set")
	}

	// Verify we have listen addresses
	if len(nodeInfo.ListenAddrs) == 0 {
		t.Error("Expected at least one listen address")
	}
}

func TestP2PWrapper_StartStop(t *testing.T) {
	ctx := context.Background()
	wrapper := p2p.NewP2PWrapper(ctx, true) // Use simulation for reliable test

	// Test multiple start/stop cycles
	for i := 0; i < 3; i++ {
		err := wrapper.Start()
		if err != nil {
			t.Fatalf("Failed to start wrapper on iteration %d: %v", i, err)
		}

		// Verify running state
		nodeInfo := wrapper.GetNodeInfo()
		if !nodeInfo.IsRunning {
			t.Errorf("Expected node to be running on iteration %d", i)
		}

		err = wrapper.Stop()
		if err != nil {
			t.Errorf("Failed to stop wrapper on iteration %d: %v", i, err)
		}
	}
}

func TestP2PWrapper_MessageSending(t *testing.T) {
	ctx := context.Background()
	wrapper := p2p.NewP2PWrapper(ctx, true) // Use simulation for reliable test

	err := wrapper.Start()
	if err != nil {
		t.Fatalf("Failed to start wrapper: %v", err)
	}
	defer wrapper.Stop()

	// Test sending various messages
	testMessages := []struct {
		peerID  string
		message string
	}{
		{"peer1", "Hello"},
		{"peer2", "World"},
		{"peer3", "Test message with special chars: !@#$%^&*()"},
		{"peer4", ""},
	}

	for _, tm := range testMessages {
		err := wrapper.SendMessage(tm.peerID, tm.message)
		if err != nil {
			t.Errorf("Failed to send message to %s: %v", tm.peerID, err)
		}
	}
}

func TestP2PWrapper_SimulationDelay(t *testing.T) {
	ctx := context.Background()
	wrapper := p2p.NewP2PWrapper(ctx, true) // Use simulation mode

	// Measure startup time
	start := time.Now()
	err := wrapper.Start()
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Failed to start wrapper: %v", err)
	}
	defer wrapper.Stop()

	// Simulation should have some delay but not too much
	if duration < 100*time.Millisecond {
		t.Error("Expected simulation to have some startup delay")
	}
	if duration > 1*time.Second {
		t.Error("Simulation startup took too long")
	}

	// Measure message sending time
	start = time.Now()
	err = wrapper.SendMessage("test-peer", "test message")
	duration = time.Since(start)

	if err != nil {
		t.Errorf("Failed to send message: %v", err)
	}

	// Message sending should have some delay but not too much
	if duration < 50*time.Millisecond {
		t.Error("Expected message sending to have some delay")
	}
	if duration > 500*time.Millisecond {
		t.Error("Message sending took too long")
	}
}

func TestP2PWrapper_LoggingToFile(t *testing.T) {
	ctx := context.Background()
	wrapper := p2p.NewP2PWrapper(ctx, true) // Use simulation for reliable test

	err := wrapper.Start()
	if err != nil {
		t.Fatalf("Failed to start wrapper: %v", err)
	}
	defer wrapper.Stop()

	// Check if log file exists
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get home directory")
	}

	logFile := filepath.Join(home, ".xelvra", "peerchat.log")
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Expected log file to be created")
	}
}

func TestP2PWrapper_NodeInfoConsistency(t *testing.T) {
	ctx := context.Background()
	wrapper := p2p.NewP2PWrapper(ctx, true) // Use simulation for consistent test

	err := wrapper.Start()
	if err != nil {
		t.Fatalf("Failed to start wrapper: %v", err)
	}
	defer wrapper.Stop()

	// Get node info multiple times and verify consistency
	info1 := wrapper.GetNodeInfo()
	info2 := wrapper.GetNodeInfo()

	if info1.PeerID != info2.PeerID {
		t.Error("Peer ID should be consistent between calls")
	}

	if info1.DID != info2.DID {
		t.Error("DID should be consistent between calls")
	}

	if info1.IsRunning != info2.IsRunning {
		t.Error("Running status should be consistent between calls")
	}
}
