package events

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// EventType represents the type of event
type EventType string

const (
	// P2P Events
	EventPeerConnected    EventType = "peer.connected"
	EventPeerDisconnected EventType = "peer.disconnected"
	EventPeerDiscovered   EventType = "peer.discovered"
	
	// Message Events
	EventMessageReceived EventType = "message.received"
	EventMessageSent     EventType = "message.sent"
	EventMessageFailed   EventType = "message.failed"
	
	// File Transfer Events
	EventFileTransferStarted   EventType = "file.transfer.started"
	EventFileTransferProgress  EventType = "file.transfer.progress"
	EventFileTransferCompleted EventType = "file.transfer.completed"
	EventFileTransferFailed    EventType = "file.transfer.failed"
	
	// Node Events
	EventNodeStarted EventType = "node.started"
	EventNodeStopped EventType = "node.stopped"
	EventNodeError   EventType = "node.error"
	
	// Network Events
	EventNetworkConnected    EventType = "network.connected"
	EventNetworkDisconnected EventType = "network.disconnected"
	EventNetworkError        EventType = "network.error"
)

// Event represents a system event
type Event struct {
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	Data      map[string]interface{} `json:"data"`
	ID        string                 `json:"id"`
}

// EventHandler is a function that handles events
type EventHandler func(event Event) error

// EventSubscription represents a subscription to events
type EventSubscription struct {
	ID       string
	Type     EventType
	Handler  EventHandler
	Filter   func(Event) bool // Optional filter function
	Priority int              // Higher priority handlers are called first
}

// EventBus is the central event bus for the application
type EventBus struct {
	mu            sync.RWMutex
	subscriptions map[EventType][]*EventSubscription
	logger        *logrus.Logger
	ctx           context.Context
	cancel        context.CancelFunc
	eventQueue    chan Event
	workers       int
	bufferSize    int
}

// NewEventBus creates a new event bus
func NewEventBus(logger *logrus.Logger, workers int, bufferSize int) *EventBus {
	ctx, cancel := context.WithCancel(context.Background())
	
	bus := &EventBus{
		subscriptions: make(map[EventType][]*EventSubscription),
		logger:        logger,
		ctx:           ctx,
		cancel:        cancel,
		eventQueue:    make(chan Event, bufferSize),
		workers:       workers,
		bufferSize:    bufferSize,
	}
	
	// Start worker goroutines
	for i := 0; i < workers; i++ {
		go bus.worker(i)
	}
	
	return bus
}

// Subscribe adds a new event subscription
func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) string {
	return eb.SubscribeWithOptions(eventType, handler, nil, 0)
}

// SubscribeWithOptions adds a new event subscription with options
func (eb *EventBus) SubscribeWithOptions(eventType EventType, handler EventHandler, filter func(Event) bool, priority int) string {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	
	subscription := &EventSubscription{
		ID:       fmt.Sprintf("%s_%d", eventType, time.Now().UnixNano()),
		Type:     eventType,
		Handler:  handler,
		Filter:   filter,
		Priority: priority,
	}
	
	eb.subscriptions[eventType] = append(eb.subscriptions[eventType], subscription)
	
	// Sort by priority (higher first)
	subs := eb.subscriptions[eventType]
	for i := len(subs) - 1; i > 0; i-- {
		if subs[i].Priority > subs[i-1].Priority {
			subs[i], subs[i-1] = subs[i-1], subs[i]
		} else {
			break
		}
	}
	
	eb.logger.WithFields(logrus.Fields{
		"event_type":      eventType,
		"subscription_id": subscription.ID,
		"priority":        priority,
	}).Debug("Event subscription added")
	
	return subscription.ID
}

// Unsubscribe removes an event subscription
func (eb *EventBus) Unsubscribe(subscriptionID string) bool {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	
	for eventType, subs := range eb.subscriptions {
		for i, sub := range subs {
			if sub.ID == subscriptionID {
				// Remove subscription
				eb.subscriptions[eventType] = append(subs[:i], subs[i+1:]...)
				
				eb.logger.WithFields(logrus.Fields{
					"event_type":      eventType,
					"subscription_id": subscriptionID,
				}).Debug("Event subscription removed")
				
				return true
			}
		}
	}
	
	return false
}

// Publish publishes an event to the bus
func (eb *EventBus) Publish(event Event) error {
	// Set timestamp if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	
	// Generate ID if not provided
	if event.ID == "" {
		event.ID = fmt.Sprintf("%s_%d", event.Type, time.Now().UnixNano())
	}
	
	select {
	case eb.eventQueue <- event:
		eb.logger.WithFields(logrus.Fields{
			"event_type": event.Type,
			"event_id":   event.ID,
			"source":     event.Source,
		}).Debug("Event published to queue")
		return nil
	case <-eb.ctx.Done():
		return fmt.Errorf("event bus is shutting down")
	default:
		eb.logger.WithFields(logrus.Fields{
			"event_type": event.Type,
			"event_id":   event.ID,
		}).Warn("Event queue is full, dropping event")
		return fmt.Errorf("event queue is full")
	}
}

// PublishSync publishes an event synchronously
func (eb *EventBus) PublishSync(event Event) error {
	// Set timestamp if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	
	// Generate ID if not provided
	if event.ID == "" {
		event.ID = fmt.Sprintf("%s_%d", event.Type, time.Now().UnixNano())
	}
	
	return eb.processEvent(event)
}

// worker processes events from the queue
func (eb *EventBus) worker(workerID int) {
	eb.logger.WithField("worker_id", workerID).Debug("Event bus worker started")
	
	for {
		select {
		case event := <-eb.eventQueue:
			if err := eb.processEvent(event); err != nil {
				eb.logger.WithFields(logrus.Fields{
					"worker_id":  workerID,
					"event_type": event.Type,
					"event_id":   event.ID,
					"error":      err,
				}).Error("Failed to process event")
			}
		case <-eb.ctx.Done():
			eb.logger.WithField("worker_id", workerID).Debug("Event bus worker stopped")
			return
		}
	}
}

// processEvent processes a single event
func (eb *EventBus) processEvent(event Event) error {
	eb.mu.RLock()
	subscriptions := eb.subscriptions[event.Type]
	eb.mu.RUnlock()
	
	if len(subscriptions) == 0 {
		eb.logger.WithField("event_type", event.Type).Debug("No subscribers for event type")
		return nil
	}
	
	var errors []error
	
	for _, sub := range subscriptions {
		// Apply filter if provided
		if sub.Filter != nil && !sub.Filter(event) {
			continue
		}
		
		// Call handler
		if err := sub.Handler(event); err != nil {
			eb.logger.WithFields(logrus.Fields{
				"event_type":      event.Type,
				"subscription_id": sub.ID,
				"error":           err,
			}).Error("Event handler failed")
			errors = append(errors, err)
		}
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("event processing failed with %d errors", len(errors))
	}
	
	return nil
}

// Stop stops the event bus
func (eb *EventBus) Stop() {
	eb.logger.Info("Stopping event bus...")
	eb.cancel()
	
	// Drain remaining events
	close(eb.eventQueue)
	for event := range eb.eventQueue {
		eb.processEvent(event)
	}
	
	eb.logger.Info("Event bus stopped")
}

// GetStats returns event bus statistics
func (eb *EventBus) GetStats() map[string]interface{} {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	
	stats := map[string]interface{}{
		"workers":     eb.workers,
		"buffer_size": eb.bufferSize,
		"queue_size":  len(eb.eventQueue),
	}
	
	subscriptionCounts := make(map[EventType]int)
	for eventType, subs := range eb.subscriptions {
		subscriptionCounts[eventType] = len(subs)
	}
	stats["subscriptions"] = subscriptionCounts
	
	return stats
}
