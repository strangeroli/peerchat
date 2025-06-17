package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Xelvra/peerchat/internal/p2p"
	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
)

// RunInteractiveChat starts the P2P node with interactive chat
func RunInteractiveChat(cmd *cobra.Command, args []string) {
	fmt.Println("🚀 Starting Xelvra P2P Messenger CLI")
	fmt.Printf("Version: %s\n", "0.2.0-alpha")
	fmt.Println("💬 Interactive Chat Mode")
	fmt.Println("📝 Logs are written to ~/.xelvra/peerchat.log")
	fmt.Println()

	// Create P2P wrapper (try real P2P first, fallback to simulation)
	ctx := context.Background()
	wrapper := p2p.NewP2PWrapper(ctx, false)

	fmt.Println("🔧 Initializing P2P node...")
	if err := wrapper.Start(); err != nil {
		fmt.Printf("❌ Failed to start real P2P node: %v\n", err)
		fmt.Println("🔄 Falling back to simulation mode...")

		// Try simulation mode
		wrapper = p2p.NewP2PWrapper(ctx, true)
		if err := wrapper.Start(); err != nil {
			fmt.Printf("❌ Failed to start simulation mode: %v\n", err)
			return
		}
	}
	defer func() {
		if err := wrapper.Stop(); err != nil {
			fmt.Printf("Warning: Failed to stop wrapper: %v\n", err)
		}
	}()

	// Get node information
	nodeInfo := wrapper.GetNodeInfo()

	fmt.Println("✅ P2P node started successfully!")
	fmt.Printf("🆔 Your Peer ID: %s\n", nodeInfo.PeerID)
	fmt.Printf("🌐 Your DID: %s\n", nodeInfo.DID)
	fmt.Printf("📡 Listening on: %v\n", nodeInfo.ListenAddrs)
	fmt.Println()

	if wrapper.IsUsingSimulation() {
		fmt.Println("⚠️  Note: Using simulation mode (real P2P failed to start)")
		fmt.Println("💡 Messages will be simulated. For real P2P, check network settings.")
	} else {
		fmt.Println("✅ Using real P2P networking")
		fmt.Println("💡 Share your Peer ID with others to receive messages")
	}

	fmt.Println()
	fmt.Println("💬 Interactive chat started! Type /help for commands.")
	fmt.Println("🎯 Features: Tab completion, command history, arrow keys")
	fmt.Println()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Try to create readline instance for advanced input
	rl, completer, err := CreateReadlineInstance()
	if err != nil {
		fmt.Printf("⚠️  Failed to initialize advanced input: %v\n", err)
		fmt.Println("💡 Falling back to basic input mode")
		// Fallback to basic input mode would go here
		return
	}
	defer func() {
		if err := rl.Close(); err != nil {
			fmt.Printf("Warning: Failed to close readline: %v\n", err)
		}
	}()

	// Create input channel
	inputChan := make(chan string)

	// Start advanced input handler with readline
	go func() {
		defer close(inputChan)

		for {
			// Update peer completions periodically
			completer.UpdatePeers(wrapper)

			line, err := rl.Readline()
			if err != nil {
				switch err {
				case readline.ErrInterrupt:
					inputChan <- "/quit"
					return
				case io.EOF:
					inputChan <- "/quit"
					return
				default:
					return
				}
			}

			input := strings.TrimSpace(line)
			if input != "" {
				inputChan <- input
			}
		}
	}()

	// Main event loop
	for {
		select {
		case <-sigChan:
			fmt.Println("\n👋 Shutdown signal received, stopping node...")
			fmt.Println("✅ Node stopped successfully")
			fmt.Println("👋 Goodbye!")
			return

		case input, ok := <-inputChan:
			if !ok {
				fmt.Println("\n👋 Input closed, shutting down...")
				return
			}

			if input == "" {
				continue
			}

			// Handle commands
			if strings.HasPrefix(input, "/") {
				if input == "/quit" || input == "/exit" {
					fmt.Println("👋 Goodbye!")
					return
				}
				HandleChatCommand(input, wrapper, nodeInfo)
			} else {
				// Send message to all connected peers
				HandleChatMessage(input, wrapper)
			}

		default:
			// Check for incoming messages (placeholder)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// RunDaemonMode runs the P2P node as a background daemon
func RunDaemonMode(cmd *cobra.Command, args []string) {
	fmt.Println("🔧 Starting Xelvra P2P Messenger in daemon mode...")
	fmt.Printf("Version: %s\n", "0.2.0-alpha")
	fmt.Println("📝 All logs will be written to ~/.xelvra/peerchat.log")
	fmt.Println()

	// Create P2P wrapper
	ctx := context.Background()
	wrapper := p2p.NewP2PWrapper(ctx, false)

	fmt.Println("🔧 Initializing P2P node...")
	if err := wrapper.Start(); err != nil {
		fmt.Printf("❌ Failed to start P2P node: %v\n", err)
		return
	}
	defer func() {
		if err := wrapper.Stop(); err != nil {
			fmt.Printf("Warning: Failed to stop wrapper: %v\n", err)
		}
	}()

	// Get node information
	nodeInfo := wrapper.GetNodeInfo()

	fmt.Println("✅ P2P node started in daemon mode!")
	fmt.Printf("🆔 Your Peer ID: %s\n", nodeInfo.PeerID)
	fmt.Printf("🌐 Your DID: %s\n", nodeInfo.DID)
	fmt.Printf("📡 Listening on: %v\n", nodeInfo.ListenAddrs)
	fmt.Println()
	fmt.Println("🔄 Running in background... Press Ctrl+C to stop")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	fmt.Println("\n👋 Shutdown signal received, stopping daemon...")
	fmt.Println("✅ Daemon stopped successfully")
}
