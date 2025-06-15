package message

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

// ConsoleMessageHandler handles messages by printing them to console
type ConsoleMessageHandler struct {
	logger *logrus.Logger
}

// NewConsoleMessageHandler creates a new console message handler
func NewConsoleMessageHandler(logger *logrus.Logger) *ConsoleMessageHandler {
	return &ConsoleMessageHandler{
		logger: logger,
	}
}

// HandleMessage handles a message by printing it to console
func (h *ConsoleMessageHandler) HandleMessage(ctx context.Context, msg *Message) error {
	switch msg.Type {
	case MessageTypeText:
		fmt.Printf("\n📨 Message from %s:\n", msg.From)
		fmt.Printf("   %s\n", string(msg.Content))
		fmt.Printf("   [%s]\n\n", msg.Timestamp.Format("15:04:05"))
		
	case MessageTypeSystem:
		fmt.Printf("\n🔧 System message from %s:\n", msg.From)
		fmt.Printf("   %s\n", string(msg.Content))
		fmt.Printf("   [%s]\n\n", msg.Timestamp.Format("15:04:05"))
		
	default:
		fmt.Printf("\n📦 %s message from %s:\n", msg.Type.String(), msg.From)
		fmt.Printf("   Size: %d bytes\n", len(msg.Content))
		fmt.Printf("   [%s]\n\n", msg.Timestamp.Format("15:04:05"))
	}
	
	h.logger.WithFields(logrus.Fields{
		"message_id": msg.ID,
		"from":       msg.From,
		"type":       msg.Type.String(),
	}).Debug("Message handled by console handler")
	
	return nil
}
