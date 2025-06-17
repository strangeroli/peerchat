package events

import (
	"time"

	"github.com/sirupsen/logrus"
)

// EventEmitter provides a convenient interface for emitting events
type EventEmitter struct {
	bus    *EventBus
	source string
	logger *logrus.Logger
}

// NewEventEmitter creates a new event emitter
func NewEventEmitter(bus *EventBus, source string, logger *logrus.Logger) *EventEmitter {
	return &EventEmitter{
		bus:    bus,
		source: source,
		logger: logger,
	}
}

// EmitPeerConnected emits a peer connected event
func (ee *EventEmitter) EmitPeerConnected(peerID string, address string) error {
	return ee.bus.Publish(Event{
		Type:   EventPeerConnected,
		Source: ee.source,
		Data: map[string]interface{}{
			"peer_id": peerID,
			"address": address,
		},
	})
}

// EmitPeerDisconnected emits a peer disconnected event
func (ee *EventEmitter) EmitPeerDisconnected(peerID string, reason string) error {
	return ee.bus.Publish(Event{
		Type:   EventPeerDisconnected,
		Source: ee.source,
		Data: map[string]interface{}{
			"peer_id": peerID,
			"reason":  reason,
		},
	})
}

// EmitPeerDiscovered emits a peer discovered event
func (ee *EventEmitter) EmitPeerDiscovered(peerID string, addresses []string, method string) error {
	return ee.bus.Publish(Event{
		Type:   EventPeerDiscovered,
		Source: ee.source,
		Data: map[string]interface{}{
			"peer_id":   peerID,
			"addresses": addresses,
			"method":    method,
		},
	})
}

// EmitMessageReceived emits a message received event
func (ee *EventEmitter) EmitMessageReceived(fromPeerID string, message string, messageType string) error {
	return ee.bus.Publish(Event{
		Type:   EventMessageReceived,
		Source: ee.source,
		Data: map[string]interface{}{
			"from_peer_id":  fromPeerID,
			"message":       message,
			"message_type":  messageType,
			"received_at":   time.Now(),
		},
	})
}

// EmitMessageSent emits a message sent event
func (ee *EventEmitter) EmitMessageSent(toPeerID string, message string, messageID string) error {
	return ee.bus.Publish(Event{
		Type:   EventMessageSent,
		Source: ee.source,
		Data: map[string]interface{}{
			"to_peer_id": toPeerID,
			"message":    message,
			"message_id": messageID,
			"sent_at":    time.Now(),
		},
	})
}

// EmitMessageFailed emits a message failed event
func (ee *EventEmitter) EmitMessageFailed(toPeerID string, message string, error string) error {
	return ee.bus.Publish(Event{
		Type:   EventMessageFailed,
		Source: ee.source,
		Data: map[string]interface{}{
			"to_peer_id": toPeerID,
			"message":    message,
			"error":      error,
			"failed_at":  time.Now(),
		},
	})
}

// EmitFileTransferStarted emits a file transfer started event
func (ee *EventEmitter) EmitFileTransferStarted(transferID string, filename string, size int64, peerID string) error {
	return ee.bus.Publish(Event{
		Type:   EventFileTransferStarted,
		Source: ee.source,
		Data: map[string]interface{}{
			"transfer_id": transferID,
			"filename":    filename,
			"size":        size,
			"peer_id":     peerID,
			"started_at":  time.Now(),
		},
	})
}

// EmitFileTransferProgress emits a file transfer progress event
func (ee *EventEmitter) EmitFileTransferProgress(transferID string, bytesTransferred int64, totalBytes int64) error {
	progress := float64(bytesTransferred) / float64(totalBytes) * 100
	
	return ee.bus.Publish(Event{
		Type:   EventFileTransferProgress,
		Source: ee.source,
		Data: map[string]interface{}{
			"transfer_id":        transferID,
			"bytes_transferred":  bytesTransferred,
			"total_bytes":        totalBytes,
			"progress_percent":   progress,
			"updated_at":         time.Now(),
		},
	})
}

// EmitFileTransferCompleted emits a file transfer completed event
func (ee *EventEmitter) EmitFileTransferCompleted(transferID string, filename string, totalBytes int64) error {
	return ee.bus.Publish(Event{
		Type:   EventFileTransferCompleted,
		Source: ee.source,
		Data: map[string]interface{}{
			"transfer_id":   transferID,
			"filename":      filename,
			"total_bytes":   totalBytes,
			"completed_at":  time.Now(),
		},
	})
}

// EmitFileTransferFailed emits a file transfer failed event
func (ee *EventEmitter) EmitFileTransferFailed(transferID string, filename string, error string) error {
	return ee.bus.Publish(Event{
		Type:   EventFileTransferFailed,
		Source: ee.source,
		Data: map[string]interface{}{
			"transfer_id": transferID,
			"filename":    filename,
			"error":       error,
			"failed_at":   time.Now(),
		},
	})
}

// EmitNodeStarted emits a node started event
func (ee *EventEmitter) EmitNodeStarted(nodeID string, addresses []string) error {
	return ee.bus.Publish(Event{
		Type:   EventNodeStarted,
		Source: ee.source,
		Data: map[string]interface{}{
			"node_id":    nodeID,
			"addresses":  addresses,
			"started_at": time.Now(),
		},
	})
}

// EmitNodeStopped emits a node stopped event
func (ee *EventEmitter) EmitNodeStopped(nodeID string, reason string) error {
	return ee.bus.Publish(Event{
		Type:   EventNodeStopped,
		Source: ee.source,
		Data: map[string]interface{}{
			"node_id":    nodeID,
			"reason":     reason,
			"stopped_at": time.Now(),
		},
	})
}

// EmitNodeError emits a node error event
func (ee *EventEmitter) EmitNodeError(nodeID string, error string, severity string) error {
	return ee.bus.Publish(Event{
		Type:   EventNodeError,
		Source: ee.source,
		Data: map[string]interface{}{
			"node_id":   nodeID,
			"error":     error,
			"severity":  severity,
			"error_at":  time.Now(),
		},
	})
}

// EmitNetworkConnected emits a network connected event
func (ee *EventEmitter) EmitNetworkConnected(networkType string, address string) error {
	return ee.bus.Publish(Event{
		Type:   EventNetworkConnected,
		Source: ee.source,
		Data: map[string]interface{}{
			"network_type":  networkType,
			"address":       address,
			"connected_at":  time.Now(),
		},
	})
}

// EmitNetworkDisconnected emits a network disconnected event
func (ee *EventEmitter) EmitNetworkDisconnected(networkType string, address string, reason string) error {
	return ee.bus.Publish(Event{
		Type:   EventNetworkDisconnected,
		Source: ee.source,
		Data: map[string]interface{}{
			"network_type":     networkType,
			"address":          address,
			"reason":           reason,
			"disconnected_at":  time.Now(),
		},
	})
}

// EmitNetworkError emits a network error event
func (ee *EventEmitter) EmitNetworkError(networkType string, error string, severity string) error {
	return ee.bus.Publish(Event{
		Type:   EventNetworkError,
		Source: ee.source,
		Data: map[string]interface{}{
			"network_type": networkType,
			"error":        error,
			"severity":     severity,
			"error_at":     time.Now(),
		},
	})
}

// EmitCustom emits a custom event
func (ee *EventEmitter) EmitCustom(eventType EventType, data map[string]interface{}) error {
	return ee.bus.Publish(Event{
		Type:   eventType,
		Source: ee.source,
		Data:   data,
	})
}

// EmitCustomSync emits a custom event synchronously
func (ee *EventEmitter) EmitCustomSync(eventType EventType, data map[string]interface{}) error {
	return ee.bus.PublishSync(Event{
		Type:   eventType,
		Source: ee.source,
		Data:   data,
	})
}

// SetSource updates the source identifier for this emitter
func (ee *EventEmitter) SetSource(source string) {
	ee.source = source
}

// GetSource returns the current source identifier
func (ee *EventEmitter) GetSource() string {
	return ee.source
}
