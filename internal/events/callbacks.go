package events

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// CallbackType represents different types of callbacks
type CallbackType string

const (
	CallbackTypeUI       CallbackType = "ui"
	CallbackTypeAPI      CallbackType = "api"
	CallbackTypeInternal CallbackType = "internal"
	CallbackTypeMetrics  CallbackType = "metrics"
)

// CallbackPriority represents callback execution priority
type CallbackPriority int

const (
	PriorityLow    CallbackPriority = 1
	PriorityNormal CallbackPriority = 5
	PriorityHigh   CallbackPriority = 10
	PriorityCritical CallbackPriority = 15
)

// CallbackConfig represents callback configuration
type CallbackConfig struct {
	Type        CallbackType
	Priority    CallbackPriority
	Timeout     time.Duration
	Retries     int
	Async       bool
	Debounce    time.Duration // Minimum time between callback executions
	Filter      func(Event) bool
}

// CallbackResult represents the result of a callback execution
type CallbackResult struct {
	Success   bool
	Error     error
	Duration  time.Duration
	Timestamp time.Time
}

// CallbackInfo represents information about a registered callback
type CallbackInfo struct {
	ID          string
	EventType   EventType
	Config      CallbackConfig
	Handler     EventHandler
	LastCall    time.Time
	CallCount   int64
	ErrorCount  int64
	TotalTime   time.Duration
}

// CallbackManager manages event callbacks with advanced features
type CallbackManager struct {
	mu           sync.RWMutex
	callbacks    map[string]*CallbackInfo
	eventBus     *EventBus
	logger       *logrus.Logger
	ctx          context.Context
	cancel       context.CancelFunc
	debounceMap  map[string]time.Time
	debounceMu   sync.RWMutex
}

// NewCallbackManager creates a new callback manager
func NewCallbackManager(eventBus *EventBus, logger *logrus.Logger) *CallbackManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	cm := &CallbackManager{
		callbacks:   make(map[string]*CallbackInfo),
		eventBus:    eventBus,
		logger:      logger,
		ctx:         ctx,
		cancel:      cancel,
		debounceMap: make(map[string]time.Time),
	}
	
	return cm
}

// RegisterCallback registers a new callback with configuration
func (cm *CallbackManager) RegisterCallback(eventType EventType, handler EventHandler, config CallbackConfig) string {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	callbackID := fmt.Sprintf("%s_%s_%d", config.Type, eventType, time.Now().UnixNano())
	
	// Wrap handler with callback management logic
	wrappedHandler := cm.wrapHandler(callbackID, handler, config)
	
	// Subscribe to event bus
	subscriptionID := cm.eventBus.SubscribeWithOptions(
		eventType,
		wrappedHandler,
		config.Filter,
		int(config.Priority),
	)
	
	// Store callback info
	cm.callbacks[callbackID] = &CallbackInfo{
		ID:        callbackID,
		EventType: eventType,
		Config:    config,
		Handler:   handler,
	}
	
	cm.logger.WithFields(logrus.Fields{
		"callback_id":     callbackID,
		"subscription_id": subscriptionID,
		"event_type":      eventType,
		"callback_type":   config.Type,
		"priority":        config.Priority,
	}).Debug("Callback registered")
	
	return callbackID
}

// wrapHandler wraps the original handler with callback management features
func (cm *CallbackManager) wrapHandler(callbackID string, handler EventHandler, config CallbackConfig) EventHandler {
	return func(event Event) error {
		start := time.Now()
		
		// Check debounce
		if config.Debounce > 0 {
			cm.debounceMu.Lock()
			lastCall, exists := cm.debounceMap[callbackID]
			if exists && time.Since(lastCall) < config.Debounce {
				cm.debounceMu.Unlock()
				cm.logger.WithFields(logrus.Fields{
					"callback_id": callbackID,
					"event_type":  event.Type,
				}).Debug("Callback debounced")
				return nil
			}
			cm.debounceMap[callbackID] = start
			cm.debounceMu.Unlock()
		}
		
		// Update statistics
		cm.mu.Lock()
		if info, exists := cm.callbacks[callbackID]; exists {
			info.LastCall = start
			info.CallCount++
		}
		cm.mu.Unlock()
		
		var err error
		
		if config.Async {
			// Execute asynchronously
			go func() {
				err = cm.executeWithRetries(callbackID, handler, event, config)
				cm.updateStats(callbackID, start, err)
			}()
		} else {
			// Execute synchronously
			err = cm.executeWithRetries(callbackID, handler, event, config)
			cm.updateStats(callbackID, start, err)
		}
		
		return err
	}
}

// executeWithRetries executes a handler with retry logic
func (cm *CallbackManager) executeWithRetries(callbackID string, handler EventHandler, event Event, config CallbackConfig) error {
	var lastErr error
	
	for attempt := 0; attempt <= config.Retries; attempt++ {
		// Create context with timeout if specified
		ctx := cm.ctx
		if config.Timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(cm.ctx, config.Timeout)
			defer cancel()
		}
		
		// Execute handler
		done := make(chan error, 1)
		go func() {
			done <- handler(event)
		}()
		
		select {
		case err := <-done:
			if err == nil {
				return nil // Success
			}
			lastErr = err
			
			cm.logger.WithFields(logrus.Fields{
				"callback_id": callbackID,
				"event_type":  event.Type,
				"attempt":     attempt + 1,
				"error":       err,
			}).Warn("Callback execution failed")
			
		case <-ctx.Done():
			lastErr = fmt.Errorf("callback timeout: %v", ctx.Err())
			
			cm.logger.WithFields(logrus.Fields{
				"callback_id": callbackID,
				"event_type":  event.Type,
				"attempt":     attempt + 1,
				"timeout":     config.Timeout,
			}).Warn("Callback execution timed out")
		}
		
		// Wait before retry (exponential backoff)
		if attempt < config.Retries {
			backoff := time.Duration(attempt+1) * 100 * time.Millisecond
			time.Sleep(backoff)
		}
	}
	
	return fmt.Errorf("callback failed after %d attempts: %v", config.Retries+1, lastErr)
}

// updateStats updates callback statistics
func (cm *CallbackManager) updateStats(callbackID string, start time.Time, err error) {
	duration := time.Since(start)
	
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	if info, exists := cm.callbacks[callbackID]; exists {
		info.TotalTime += duration
		if err != nil {
			info.ErrorCount++
		}
	}
}

// UnregisterCallback removes a callback
func (cm *CallbackManager) UnregisterCallback(callbackID string) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	if info, exists := cm.callbacks[callbackID]; exists {
		// Remove from event bus (we need to track subscription IDs for this)
		// For now, we'll just remove from our tracking
		delete(cm.callbacks, callbackID)
		
		// Clean up debounce map
		cm.debounceMu.Lock()
		delete(cm.debounceMap, callbackID)
		cm.debounceMu.Unlock()
		
		cm.logger.WithFields(logrus.Fields{
			"callback_id": callbackID,
			"event_type":  info.EventType,
		}).Debug("Callback unregistered")
		
		return true
	}
	
	return false
}

// GetCallbackInfo returns information about a callback
func (cm *CallbackManager) GetCallbackInfo(callbackID string) (*CallbackInfo, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	if info, exists := cm.callbacks[callbackID]; exists {
		// Return a copy to avoid race conditions
		infoCopy := *info
		return &infoCopy, true
	}
	
	return nil, false
}

// ListCallbacks returns all registered callbacks
func (cm *CallbackManager) ListCallbacks() map[string]*CallbackInfo {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	result := make(map[string]*CallbackInfo)
	for id, info := range cm.callbacks {
		infoCopy := *info
		result[id] = &infoCopy
	}
	
	return result
}

// GetStats returns callback manager statistics
func (cm *CallbackManager) GetStats() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	stats := map[string]interface{}{
		"total_callbacks": len(cm.callbacks),
	}
	
	// Group by type
	typeStats := make(map[CallbackType]int)
	priorityStats := make(map[CallbackPriority]int)
	eventTypeStats := make(map[EventType]int)
	
	var totalCalls int64
	var totalErrors int64
	var totalTime time.Duration
	
	for _, info := range cm.callbacks {
		typeStats[info.Config.Type]++
		priorityStats[info.Config.Priority]++
		eventTypeStats[info.EventType]++
		
		totalCalls += info.CallCount
		totalErrors += info.ErrorCount
		totalTime += info.TotalTime
	}
	
	stats["by_type"] = typeStats
	stats["by_priority"] = priorityStats
	stats["by_event_type"] = eventTypeStats
	stats["total_calls"] = totalCalls
	stats["total_errors"] = totalErrors
	stats["total_time"] = totalTime.String()
	
	if totalCalls > 0 {
		stats["error_rate"] = float64(totalErrors) / float64(totalCalls) * 100
		stats["avg_duration"] = (totalTime / time.Duration(totalCalls)).String()
	}
	
	return stats
}

// Stop stops the callback manager
func (cm *CallbackManager) Stop() {
	cm.logger.Info("Stopping callback manager...")
	cm.cancel()
	
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	// Clear all callbacks
	cm.callbacks = make(map[string]*CallbackInfo)
	
	cm.debounceMu.Lock()
	cm.debounceMap = make(map[string]time.Time)
	cm.debounceMu.Unlock()
	
	cm.logger.Info("Callback manager stopped")
}
