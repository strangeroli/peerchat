package unit

import (
	"testing"
	"time"

	"github.com/Xelvra/peerchat/internal/events"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventBus_Creation(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	
	bus := events.NewEventBus(logger, 2, 100)
	assert.NotNil(t, bus)
	
	stats := bus.GetStats()
	assert.Equal(t, 2, stats["workers"])
	assert.Equal(t, 100, stats["buffer_size"])
	assert.Equal(t, 0, stats["queue_size"])
	
	bus.Stop()
}

func TestEventBus_SubscribeAndPublish(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	
	bus := events.NewEventBus(logger, 1, 10)
	defer bus.Stop()
	
	// Test subscription
	received := make(chan events.Event, 1)
	subscriptionID := bus.Subscribe(events.EventPeerConnected, func(event events.Event) error {
		received <- event
		return nil
	})
	
	assert.NotEmpty(t, subscriptionID)
	
	// Test publishing
	testEvent := events.Event{
		Type:   events.EventPeerConnected,
		Source: "test",
		Data: map[string]interface{}{
			"peer_id": "test-peer",
		},
	}
	
	err := bus.Publish(testEvent)
	require.NoError(t, err)
	
	// Wait for event to be processed
	select {
	case receivedEvent := <-received:
		assert.Equal(t, events.EventPeerConnected, receivedEvent.Type)
		assert.Equal(t, "test", receivedEvent.Source)
		assert.Equal(t, "test-peer", receivedEvent.Data["peer_id"])
		assert.NotEmpty(t, receivedEvent.ID)
		assert.False(t, receivedEvent.Timestamp.IsZero())
	case <-time.After(1 * time.Second):
		t.Fatal("Event not received within timeout")
	}
}

func TestEventBus_Unsubscribe(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	
	bus := events.NewEventBus(logger, 1, 10)
	defer bus.Stop()
	
	received := make(chan events.Event, 1)
	subscriptionID := bus.Subscribe(events.EventPeerConnected, func(event events.Event) error {
		received <- event
		return nil
	})
	
	// Unsubscribe
	success := bus.Unsubscribe(subscriptionID)
	assert.True(t, success)
	
	// Publish event - should not be received
	testEvent := events.Event{
		Type:   events.EventPeerConnected,
		Source: "test",
	}
	
	err := bus.Publish(testEvent)
	require.NoError(t, err)
	
	// Should not receive event
	select {
	case <-received:
		t.Fatal("Event received after unsubscribe")
	case <-time.After(100 * time.Millisecond):
		// Expected - no event received
	}
}

func TestEventBus_Priority(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	
	bus := events.NewEventBus(logger, 1, 10)
	defer bus.Stop()
	
	var order []int
	
	// Subscribe with different priorities
	bus.SubscribeWithOptions(events.EventPeerConnected, func(event events.Event) error {
		order = append(order, 1)
		return nil
	}, nil, 1)
	
	bus.SubscribeWithOptions(events.EventPeerConnected, func(event events.Event) error {
		order = append(order, 3)
		return nil
	}, nil, 3)
	
	bus.SubscribeWithOptions(events.EventPeerConnected, func(event events.Event) error {
		order = append(order, 2)
		return nil
	}, nil, 2)
	
	// Publish event
	testEvent := events.Event{
		Type:   events.EventPeerConnected,
		Source: "test",
	}
	
	err := bus.PublishSync(testEvent)
	require.NoError(t, err)
	
	// Check order (higher priority first)
	assert.Equal(t, []int{3, 2, 1}, order)
}

func TestEventEmitter(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	
	bus := events.NewEventBus(logger, 1, 10)
	defer bus.Stop()
	
	emitter := events.NewEventEmitter(bus, "test-component", logger)
	
	// Test peer connected event
	received := make(chan events.Event, 1)
	bus.Subscribe(events.EventPeerConnected, func(event events.Event) error {
		received <- event
		return nil
	})
	
	err := emitter.EmitPeerConnected("test-peer", "127.0.0.1:8080")
	require.NoError(t, err)
	
	select {
	case event := <-received:
		assert.Equal(t, events.EventPeerConnected, event.Type)
		assert.Equal(t, "test-component", event.Source)
		assert.Equal(t, "test-peer", event.Data["peer_id"])
		assert.Equal(t, "127.0.0.1:8080", event.Data["address"])
	case <-time.After(1 * time.Second):
		t.Fatal("Event not received within timeout")
	}
}

func TestEventEmitter_MessageEvents(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	
	bus := events.NewEventBus(logger, 1, 10)
	defer bus.Stop()
	
	emitter := events.NewEventEmitter(bus, "test-component", logger)
	
	// Test message received event
	received := make(chan events.Event, 1)
	bus.Subscribe(events.EventMessageReceived, func(event events.Event) error {
		received <- event
		return nil
	})
	
	err := emitter.EmitMessageReceived("sender-peer", "Hello, World!", "text")
	require.NoError(t, err)
	
	select {
	case event := <-received:
		assert.Equal(t, events.EventMessageReceived, event.Type)
		assert.Equal(t, "sender-peer", event.Data["from_peer_id"])
		assert.Equal(t, "Hello, World!", event.Data["message"])
		assert.Equal(t, "text", event.Data["message_type"])
		assert.NotNil(t, event.Data["received_at"])
	case <-time.After(1 * time.Second):
		t.Fatal("Event not received within timeout")
	}
}

func TestEventEmitter_FileTransferEvents(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	
	bus := events.NewEventBus(logger, 1, 10)
	defer bus.Stop()
	
	emitter := events.NewEventEmitter(bus, "test-component", logger)
	
	// Test file transfer progress event
	received := make(chan events.Event, 1)
	bus.Subscribe(events.EventFileTransferProgress, func(event events.Event) error {
		received <- event
		return nil
	})
	
	err := emitter.EmitFileTransferProgress("transfer-123", 500, 1000)
	require.NoError(t, err)
	
	select {
	case event := <-received:
		assert.Equal(t, events.EventFileTransferProgress, event.Type)
		assert.Equal(t, "transfer-123", event.Data["transfer_id"])
		assert.Equal(t, int64(500), event.Data["bytes_transferred"])
		assert.Equal(t, int64(1000), event.Data["total_bytes"])
		assert.Equal(t, 50.0, event.Data["progress_percent"])
	case <-time.After(1 * time.Second):
		t.Fatal("Event not received within timeout")
	}
}

func TestCallbackManager(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	
	bus := events.NewEventBus(logger, 1, 10)
	defer bus.Stop()
	
	cm := events.NewCallbackManager(bus, logger)
	defer cm.Stop()
	
	// Register callback
	received := make(chan events.Event, 1)
	callbackID := cm.RegisterCallback(
		events.EventPeerConnected,
		func(event events.Event) error {
			received <- event
			return nil
		},
		events.CallbackConfig{
			Type:     events.CallbackTypeUI,
			Priority: events.PriorityNormal,
			Async:    false,
		},
	)
	
	assert.NotEmpty(t, callbackID)
	
	// Check callback info
	info, exists := cm.GetCallbackInfo(callbackID)
	assert.True(t, exists)
	assert.Equal(t, events.EventPeerConnected, info.EventType)
	assert.Equal(t, events.CallbackTypeUI, info.Config.Type)
	
	// Publish event
	testEvent := events.Event{
		Type:   events.EventPeerConnected,
		Source: "test",
	}
	
	err := bus.Publish(testEvent)
	require.NoError(t, err)
	
	// Wait for event
	select {
	case <-received:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Event not received within timeout")
	}
	
	// Check updated stats
	info, exists = cm.GetCallbackInfo(callbackID)
	assert.True(t, exists)
	assert.Equal(t, int64(1), info.CallCount)
	assert.Equal(t, int64(0), info.ErrorCount)
}

func TestCallbackManager_Stats(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	
	bus := events.NewEventBus(logger, 1, 10)
	defer bus.Stop()
	
	cm := events.NewCallbackManager(bus, logger)
	defer cm.Stop()
	
	// Register multiple callbacks
	cm.RegisterCallback(events.EventPeerConnected, func(event events.Event) error {
		return nil
	}, events.CallbackConfig{Type: events.CallbackTypeUI, Priority: events.PriorityHigh})
	
	cm.RegisterCallback(events.EventMessageReceived, func(event events.Event) error {
		return nil
	}, events.CallbackConfig{Type: events.CallbackTypeAPI, Priority: events.PriorityNormal})
	
	stats := cm.GetStats()
	assert.Equal(t, 2, stats["total_callbacks"])
	
	byType := stats["by_type"].(map[events.CallbackType]int)
	assert.Equal(t, 1, byType[events.CallbackTypeUI])
	assert.Equal(t, 1, byType[events.CallbackTypeAPI])
}
