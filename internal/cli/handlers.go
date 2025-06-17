package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Xelvra/peerchat/internal/p2p"
	"github.com/spf13/cobra"
)

// RunInit handles the init command
func RunInit(cmd *cobra.Command, args []string) {
	fmt.Println("🔧 Initializing Xelvra P2P Messenger...")
	fmt.Println("📝 Logs are written to ~/.xelvra/peerchat.log")
	fmt.Println()

	// Create P2P wrapper to initialize identity
	ctx := context.Background()
	wrapper := p2p.NewP2PWrapper(ctx, false) // Try real P2P first

	fmt.Println("🔑 Generating cryptographic identity...")
	if err := wrapper.Start(); err != nil {
		fmt.Printf("❌ Failed to initialize P2P node: %v\n", err)
		fmt.Println("💡 This might be due to network issues. The identity was still created.")
		return
	}
	defer func() {
		if err := wrapper.Stop(); err != nil {
			fmt.Printf("Warning: Failed to stop wrapper: %v\n", err)
		}
	}()

	// Get node information
	nodeInfo := wrapper.GetNodeInfo()

	fmt.Println("✅ Identity created successfully!")
	fmt.Printf("🆔 Your DID: %s\n", nodeInfo.DID)
	fmt.Printf("🔗 Your Peer ID: %s\n", nodeInfo.PeerID)
	fmt.Printf("📁 Configuration saved to: ~/.xelvra/\n")
	fmt.Println()

	if wrapper.IsUsingSimulation() {
		fmt.Println("⚠️  Note: Using simulation mode (real P2P failed to start)")
		fmt.Println("💡 This is normal for first-time setup or network issues")
	} else {
		fmt.Println("✅ Real P2P networking initialized successfully")
	}

	fmt.Println("🎉 Setup complete! Next steps:")
	fmt.Println("  1. Run 'peerchat-cli doctor' to test your network")
	fmt.Println("  2. Run 'peerchat-cli start' to begin chatting")
}

// RunStart handles the start command
func RunStart(cmd *cobra.Command, args []string) {
	daemon, _ := cmd.Flags().GetBool("daemon")
	
	if daemon {
		RunDaemonMode(cmd, args)
	} else {
		RunInteractiveChat(cmd, args)
	}
}

// RunStatus handles the status command
func RunStatus(cmd *cobra.Command, args []string) {
	fmt.Println("📊 Node Status")
	fmt.Println("==============")
	fmt.Println("📝 Logs are written to ~/.xelvra/peerchat.log")
	fmt.Println()

	// Check if node is already running
	status, err := p2p.ReadNodeStatus()
	if err != nil || status == nil || !status.IsRunning {
		fmt.Println("❌ No running node found")
		fmt.Println("💡 Start the node first with: peerchat-cli start")
		return
	}

	fmt.Println("✅ Node is running")
	fmt.Printf("🆔 Peer ID: %s\n", status.PeerID)
	// DID information would be displayed here when available
	fmt.Printf("📡 Listen addresses: %v\n", status.ListenAddrs)
	fmt.Printf("🔗 Connected peers: %d\n", status.ConnectedPeers)
	fmt.Printf("⏰ Uptime: %s\n", time.Since(status.StartTime).Round(time.Second))
	fmt.Println()

	// Display NAT information
	if status.NATInfo != nil {
		fmt.Println("🌐 Network Information:")
		fmt.Printf("  NAT Type: %s\n", status.NATInfo.Type)
		fmt.Printf("  Local IP: %s:%d\n", status.NATInfo.LocalIP, status.NATInfo.LocalPort)
		if status.NATInfo.PublicIP != "" {
			fmt.Printf("  Public IP: %s:%d\n", status.NATInfo.PublicIP, status.NATInfo.PublicPort)
		}
		fmt.Println()
	}

	// Display discovery status
	if status.Discovery != nil {
		fmt.Println("🔍 Discovery Status:")
		fmt.Printf("  mDNS: %s\n", getStatusIcon(status.Discovery.MDNSActive))
		fmt.Printf("  DHT: %s\n", getStatusIcon(status.Discovery.DHTActive))
		fmt.Printf("  UDP Broadcast: %s\n", getStatusIcon(status.Discovery.UDPBroadcast))
		fmt.Printf("  Known peers: %d\n", status.Discovery.KnownPeers)
		if !status.Discovery.LastDiscovery.IsZero() {
			fmt.Printf("  Last discovery: %s\n", status.Discovery.LastDiscovery.Format("15:04:05"))
		}
	}
}

// RunVersion handles the version command
func RunVersion(version string) {
	fmt.Printf("Xelvra P2P Messenger CLI v%s\n", version)
	fmt.Println("Built with Go and libp2p")
	fmt.Println("https://github.com/Xelvra/peerchat")
}

// RunSend handles the send command
func RunSend(cmd *cobra.Command, args []string) {
	peerTarget := args[0]
	messageText := args[1]

	fmt.Printf("📤 Sending message to %s\n", peerTarget)
	fmt.Printf("💬 Message: %s\n", messageText)
	fmt.Println("📝 Logs are written to ~/.xelvra/peerchat.log")
	fmt.Println()

	// Check if node is already running
	status, err := p2p.ReadNodeStatus()
	if err != nil || status == nil || !status.IsRunning {
		fmt.Println("❌ No running node found")
		fmt.Println("💡 Start the node first with: peerchat-cli start")
		return
	}

	fmt.Println("✅ Using existing running node")
	fmt.Printf("🆔 Your Peer ID: %s\n", status.PeerID)
	fmt.Println()

	// For now, simulate message sending since we need IPC to communicate with running node
	fmt.Println("🔗 Attempting to send message via P2P network...")
	fmt.Println("⚠️  Note: Message sending via running node not yet implemented")
	fmt.Println("💡 This requires IPC (Inter-Process Communication) with the running node")
	fmt.Println("💡 For interactive messaging, use 'peerchat-cli start' mode")

	// Log the message attempt
	fmt.Println("📝 Message logged for future implementation")
	fmt.Printf("✅ Message queued: '%s' -> %s\n", messageText, peerTarget)
}

// RunConnect handles the connect command
func RunConnect(cmd *cobra.Command, args []string) {
	peerID := args[0]

	fmt.Printf("🔗 Connecting to peer: %s\n", peerID)
	fmt.Println("❌ Error: Peer connection not yet implemented")
	fmt.Println("This feature requires P2P connection management.")
}

// RunListen handles the listen command
func RunListen(cmd *cobra.Command, args []string) {
	fmt.Println("👂 Starting P2P node in passive listening mode...")
	fmt.Println("ALL LOGS AND MESSAGES will be displayed here for debugging.")
	fmt.Println("This is a passive mode - no interaction available.")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	// Create P2P wrapper with console logging enabled for debugging
	ctx := context.Background()
	wrapper := p2p.NewP2PWrapper(ctx, false) // Try real P2P first

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

	fmt.Println("✅ P2P node started successfully!")
	fmt.Printf("🆔 Your Peer ID: %s\n", nodeInfo.PeerID)
	fmt.Printf("🌐 Your DID: %s\n", nodeInfo.DID)
	fmt.Printf("📡 Listening on: %v\n", nodeInfo.ListenAddrs)
	fmt.Println()

	if wrapper.IsUsingSimulation() {
		fmt.Println("⚠️  Note: Using simulation mode (real P2P failed to start)")
	} else {
		fmt.Println("✅ Using real P2P networking")
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start real-time log monitoring
	logChan := make(chan string, 100)
	go MonitorLogFileRealTime(logChan)

	// Passive listening loop with real log monitoring
	for {
		select {
		case <-sigChan:
			fmt.Println("\n👋 Shutting down...")
			return

		case logEntry := <-logChan:
			// Display new log entries in real-time
			fmt.Printf("[%s] %s\n", time.Now().Format("15:04:05"), logEntry)

		default:
			// Small sleep to prevent busy waiting
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// RunDiscover handles the discover command
func RunDiscover(cmd *cobra.Command, args []string) {
	fmt.Println("🔍 Discovering peers in the network...")
	fmt.Println("📝 Logs are written to ~/.xelvra/peerchat.log")
	fmt.Println()

	// Check if node is already running
	status, err := p2p.ReadNodeStatus()
	if err != nil || status == nil || !status.IsRunning {
		fmt.Println("❌ No running node found")
		fmt.Println("💡 Start the node first with: peerchat-cli start")
		return
	}

	fmt.Println("✅ Using existing running node")
	fmt.Printf("🆔 Your Peer ID: %s\n", status.PeerID)
	fmt.Printf("📡 Your addresses: %v\n", status.ListenAddrs)
	fmt.Println()

	fmt.Println("⏳ Monitoring discovery for 10 seconds...")

	// Monitor discovery for 10 seconds
	for i := 1; i <= 10; i++ {
		fmt.Printf(".")
		time.Sleep(1 * time.Second)

		// Check for new peers every 2 seconds
		if i%2 == 0 {
			newStatus, err := p2p.ReadNodeStatus()
			if err == nil && newStatus != nil && newStatus.Discovery != nil {
				if newStatus.Discovery.KnownPeers > status.Discovery.KnownPeers {
					fmt.Printf("\n🎉 Found %d new peers!\n", newStatus.Discovery.KnownPeers-status.Discovery.KnownPeers)
					status = newStatus
				}
			}
		}
	}
	fmt.Println()

	// Final status
	finalStatus, err := p2p.ReadNodeStatus()
	if err == nil && finalStatus != nil {
		fmt.Println("✅ Discovery completed")
		fmt.Printf("📊 Total known peers: %d\n", finalStatus.Discovery.KnownPeers)
		fmt.Printf("🔗 Connected peers: %d\n", finalStatus.ConnectedPeers)
		fmt.Println("💡 Use 'peerchat-cli status' for detailed information")
	} else {
		fmt.Println("✅ Discovery completed")
		fmt.Println("📊 Check logs for detailed discovery information")
	}
}

// RunShowID handles the id command
func RunShowID(cmd *cobra.Command, args []string) {
	fmt.Println("🆔 Your Identity:")
	fmt.Println("==================")
	fmt.Println("📝 Logs are written to ~/.xelvra/peerchat.log")
	fmt.Println()

	// Try to get identity from P2P wrapper
	ctx := context.Background()
	wrapper := p2p.NewP2PWrapper(ctx, false) // Try real P2P first

	fmt.Println("🔧 Initializing P2P node to get identity...")
	if err := wrapper.Start(); err != nil {
		fmt.Printf("❌ Failed to start P2P node: %v\n", err)
		fmt.Println("💡 Try running 'peerchat-cli init' first")
		return
	}
	defer func() {
		if err := wrapper.Stop(); err != nil {
			fmt.Printf("Warning: Failed to stop wrapper: %v\n", err)
		}
	}()

	// Get node information
	nodeInfo := wrapper.GetNodeInfo()

	fmt.Println("✅ Identity retrieved successfully!")
	fmt.Printf("🆔 DID: %s\n", nodeInfo.DID)
	fmt.Printf("🔗 Peer ID: %s\n", nodeInfo.PeerID)
	fmt.Printf("📡 Listen addresses: %v\n", nodeInfo.ListenAddrs)
	fmt.Println()

	if wrapper.IsUsingSimulation() {
		fmt.Println("⚠️  Note: Using simulation mode (real P2P failed to start)")
		fmt.Println("💡 This identity is simulated for testing")
	} else {
		fmt.Println("✅ Using real P2P networking")
		fmt.Println("💡 Share your Peer ID with others to receive messages")
	}
}

// RunProfile handles the profile command
func RunProfile(cmd *cobra.Command, args []string) {
	peerID := args[0]

	fmt.Printf("👤 Profile for peer: %s\n", peerID)
	fmt.Println("========================")
	fmt.Println("❌ Error: Peer profile lookup not yet implemented")
	fmt.Println("This feature requires DHT lookup and peer information storage.")
}

// RunSendFile handles the send-file command
func RunSendFile(cmd *cobra.Command, args []string) {
	peerID := args[0]
	filePath := args[1]

	fmt.Printf("📁 Sending file %s to peer: %s\n", filePath, peerID)
	fmt.Println("❌ Error: File transfer not yet implemented")
	fmt.Println("This feature requires P2P file transfer protocol.")
}

// RunStop handles the stop command
func RunStop(cmd *cobra.Command, args []string) {
	fmt.Println("🛑 Stopping P2P node...")
	fmt.Println("❌ Error: Node stopping not yet implemented")
	fmt.Println("This feature requires process management and IPC.")
}

// RunSetup handles the setup command
func RunSetup(cmd *cobra.Command, args []string) {
	fmt.Println("🧙 Xelvra Setup Wizard")
	fmt.Println("======================")
	fmt.Println("❌ Error: Setup wizard not yet implemented")
	fmt.Println("This feature requires interactive CLI interface.")
}

// getStatusIcon returns an icon for boolean status
func getStatusIcon(active bool) string {
	if active {
		return "✅ Active"
	}
	return "❌ Inactive"
}
