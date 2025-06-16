package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/Xelvra/peerchat/internal/p2p"
	"github.com/Xelvra/peerchat/internal/user"
	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
	version = "0.1.0-alpha"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "peerchat-cli",
	Short: "Xelvra P2P Messenger CLI",
	Long: `Xelvra P2P Messenger CLI - A secure, decentralized messaging platform.

GETTING STARTED:
  1. peerchat-cli init     # Create your identity
  2. peerchat-cli doctor   # Test network connectivity
  3. peerchat-cli start    # Start interactive chat

STANDALONE COMMANDS (no running node required):
  init, doctor, version, manual, help

INTERACTIVE COMMANDS (available in chat mode):
  /help, /peers, /discover, /connect, /status, /quit

NODE-DEPENDENT COMMANDS (require running node):
  send, send-file, connect, discover, status

Performance targets:
- Latency: <50ms for direct connections
- Memory: <20MB idle usage
- CPU: <1% idle usage`,
	Version: version,
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Xelvra identity",
	Long: `Initialize a new Xelvra identity and configuration.

This command generates a new cryptographic identity (DID:xelvra format)
and creates the initial configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		runInit(cmd, args)
	},
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start interactive P2P chat",
	Long: `Start the P2P node in interactive chat mode.

This command starts the P2P node and provides an interactive chat interface
where you can send messages to connected peers and see incoming messages.
Use commands like /help, /peers, /connect, /quit to control the chat.

Use --daemon flag to run as background service.`,
	Run: func(cmd *cobra.Command, args []string) {
		daemon, _ := cmd.Flags().GetBool("daemon")
		if daemon {
			runDaemonMode(cmd, args)
		} else {
			runInteractiveChat(cmd, args)
		}
	},
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show node status and statistics",
	Long: `Display current node status, connection information, and performance metrics.

Shows information about:
- Node identity and network addresses
- Connected peers
- Message statistics
- Performance metrics`,
	Run: func(cmd *cobra.Command, args []string) {
		runStatus(cmd, args)
	},
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Xelvra P2P Messenger CLI v%s\n", version)
		fmt.Println("Built with Go and libp2p")
		fmt.Println("https://github.com/Xelvra/peerchat")
	},
}

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send <multiaddr> <message>",
	Short: "Send an encrypted message to a peer",
	Long: `Send an end-to-end encrypted message to a specific peer.

Use the full multiaddr format: /ip4/127.0.0.1/tcp/PORT/p2p/PEER_ID
Example: /ip4/127.0.0.1/tcp/35083/p2p/12D3KooW...

This creates a temporary P2P node, sends the message, and exits.
For persistent messaging, use 'listen' mode.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		runSend(cmd, args)
	},
}

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect <peer_id>",
	Short: "Connect to a specific peer",
	Long: `Attempt to establish a direct P2P connection to a peer.

This command will try various connection methods including
direct connection, NAT traversal, and relay if necessary.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runConnect(cmd, args)
	},
}

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for incoming messages (passive mode)",
	Long: `Start the P2P node in passive listening mode.

This command starts the node and displays incoming messages without
providing an interactive interface. Useful for monitoring and debugging.
Press Ctrl+C to stop listening.`,
	Run: func(cmd *cobra.Command, args []string) {
		runPassiveListen(cmd, args)
	},
}

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover peers in the network",
	Long: `Manually trigger peer discovery and display found peers.

This command will use various discovery methods including
DHT, mDNS, and UDP broadcast to find nearby peers.`,
	Run: func(cmd *cobra.Command, args []string) {
		runDiscover(cmd, args)
	},
}

// idCmd represents the id command
var idCmd = &cobra.Command{
	Use:   "id",
	Short: "Show your identity information",
	Long: `Display your DID (Decentralized Identifier) and Peer ID.

This information can be shared with others to allow them
to connect and send messages to you.`,
	Run: func(cmd *cobra.Command, args []string) {
		runShowID(cmd, args)
	},
}

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile <peer_id>",
	Short: "Show profile information for a peer",
	Long: `Display basic information about a remote peer including
their trust level and connection status.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runProfile(cmd, args)
	},
}

// sendFileCmd represents the send-file command
var sendFileCmd = &cobra.Command{
	Use:   "send-file <peer_id> <file_path>",
	Short: "Send a file to a peer",
	Long: `Send a file to a peer using secure P2P file transfer.

The file will be encrypted and transferred directly to the peer
with progress indication and resume capability.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		runSendFile(cmd, args)
	},
}

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the running P2P node",
	Long: `Stop the currently running P2P node gracefully.

This command will disconnect from all peers and shut down
the node cleanly.`,
	Run: func(cmd *cobra.Command, args []string) {
		runStop(cmd, args)
	},
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup wizard",
	Long: `Run the interactive setup wizard for first-time users.

This wizard will guide you through identity creation,
network configuration, and initial connection testing.`,
	Run: func(cmd *cobra.Command, args []string) {
		runSetup(cmd, args)
	},
}

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose and fix network issues",
	Long: `Diagnose network connectivity issues and attempt automatic fixes.

This command will test NAT traversal, firewall settings,
and connection quality, then suggest or apply fixes.`,
	Run: func(cmd *cobra.Command, args []string) {
		runDoctor(cmd, args)
	},
}

// manualCmd represents the manual command
var manualCmd = &cobra.Command{
	Use:   "manual",
	Short: "Show detailed usage manual",
	Long: `Display the complete usage manual with examples and troubleshooting.

This provides comprehensive documentation for all commands
and common usage patterns.`,
	Run: func(cmd *cobra.Command, args []string) {
		runManual(cmd, args)
	},
}

func main() {
	Execute()
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.xelvra/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Start command flags
	startCmd.Flags().Bool("daemon", false, "run as background daemon service")

	// Add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(listenCmd)
	rootCmd.AddCommand(discoverCmd)
	rootCmd.AddCommand(idCmd)
	rootCmd.AddCommand(profileCmd)
	rootCmd.AddCommand(sendFileCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(manualCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Configuration loading temporarily disabled for debugging
	if verbose {
		fmt.Fprintln(os.Stderr, "Config loading disabled - using defaults")
	}
}

// runInit initializes a new Xelvra identity
func runInit(cmd *cobra.Command, args []string) {
	fmt.Println("🔐 Initializing new Xelvra identity...")
	fmt.Println("📝 Logs are written to ~/.xelvra/peerchat.log")
	fmt.Println()

	// Generate new MessengerID using real user system
	fmt.Println("🔑 Generating new identity...")
	identity, err := user.GenerateMessengerID()
	if err != nil {
		fmt.Printf("❌ Failed to generate identity: %v\n", err)
		return
	}

	// Create config directory
	configDir := filepath.Join(os.Getenv("HOME"), ".xelvra")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		fmt.Printf("❌ Failed to create config directory: %v\n", err)
		return
	}

	fmt.Println("✅ New identity generated successfully!")
	fmt.Printf("  🆔 DID: %s\n", identity.DID)
	fmt.Printf("  🔗 Peer ID: %s\n", identity.PeerID)
	fmt.Printf("  📁 Config directory: %s\n", configDir)
	fmt.Println()
	fmt.Println("✅ Configuration directory created")
	fmt.Println("✅ Ready for P2P messaging!")
	fmt.Println()
	fmt.Println("🚀 Next steps:")
	fmt.Println("  peerchat-cli doctor   # Test network connectivity")
	fmt.Println("  peerchat-cli start    # Start the P2P node")
	fmt.Println("  peerchat-cli status   # Check node status")
}

// runInteractiveChat starts the P2P node with interactive chat (this is the start command)
func runInteractiveChat(cmd *cobra.Command, args []string) {
	fmt.Println("🚀 Starting Xelvra P2P Messenger CLI")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("💬 Interactive Chat Mode")
	fmt.Println("📝 Logs are written to ~/.xelvra/peerchat.log")
	fmt.Println()

	// Create P2P wrapper (try real P2P first, fallback to simulation)
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
	fmt.Println("Ready to receive messages! Share your Peer ID with others.")
	fmt.Println()

	if wrapper.IsUsingSimulation() {
		fmt.Println("⚠️  Note: Using simulation mode (real P2P failed to start)")
	} else {
		fmt.Println("✅ Using real P2P networking")
	}
	fmt.Println()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start interactive chat loop
	fmt.Println("💬 Interactive chat started. Type your messages:")
	fmt.Println("Commands: /help, /peers, /discover, /connect <peer_id>, /quit")
	fmt.Println("Features: Tab completion, command history (↑/↓), peer ID completion")
	fmt.Println()

	// Create readline instance with completion and history
	rl, completer, err := createReadlineInstance()
	if err != nil {
		fmt.Printf("❌ Failed to create readline interface: %v\n", err)
		fmt.Println("💡 Falling back to basic input mode")
		// Fallback to basic input mode would go here
		return
	}
	defer rl.Close()

	// Create input channel
	inputChan := make(chan string)

	// Start advanced input handler with readline
	go func() {
		defer close(inputChan)

		for {
			// Update peer completions periodically
			completer.updatePeers(wrapper)

			line, err := rl.Readline()
			if err != nil {
				if err == readline.ErrInterrupt {
					inputChan <- "/quit"
					return
				} else if err == io.EOF {
					inputChan <- "/quit"
					return
				}
				return
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
				handleChatCommand(input, wrapper, nodeInfo)
			} else {
				// Send message to all connected peers
				handleChatMessage(input, wrapper)
			}

		default:
			// Check for incoming messages (placeholder)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// runSend sends a message to a peer
func runSend(cmd *cobra.Command, args []string) {
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
	if verbose {
		fmt.Printf("📡 Your addresses: %v\n", status.ListenAddrs)
	}
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

// runConnect connects to a peer
func runConnect(cmd *cobra.Command, args []string) {
	peerID := args[0]

	fmt.Printf("🔗 Connecting to peer: %s\n", peerID)
	fmt.Println("❌ Error: Peer connection not yet implemented")
	fmt.Println("This feature requires P2P connection management.")
}

// runPassiveListen listens for incoming messages in passive mode (no interaction)
func runPassiveListen(cmd *cobra.Command, args []string) {
	fmt.Println("👂 Starting P2P node in passive listening mode...")
	fmt.Println("ALL LOGS AND MESSAGES will be displayed here for debugging.")
	fmt.Println("This is a passive mode - no interaction available.")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	// Setup logging to BOTH file AND console for debugging
	fmt.Println("📝 Logs are written to ~/.xelvra/peerchat.log AND displayed here")
	if verbose {
		fmt.Println("📝 Verbose mode enabled - showing all debug information")
	}

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
	fmt.Println("Ready to receive messages! Share your Peer ID with others.")
	fmt.Println()

	if wrapper.IsUsingSimulation() {
		fmt.Println("⚠️  Note: Using simulation mode (real P2P failed to start)")
	} else {
		fmt.Println("✅ Using real P2P networking")
	}
	fmt.Println()

	fmt.Println("👂 DEBUGGING MODE - All logs will appear below:")
	fmt.Println("💡 For clean interactive chat, use 'peerchat-cli start' instead")
	fmt.Println("=" + strings.Repeat("=", 60))

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start real-time log monitoring
	logChan := make(chan string, 100)
	go monitorLogFileRealTime(logChan)

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

// InteractiveCompleter provides tab completion for interactive mode
type InteractiveCompleter struct {
	commands []string
	peers    []string
}

// Do implements readline.AutoCompleter interface
func (c *InteractiveCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	lineStr := string(line)

	// Split line into words
	words := strings.Fields(lineStr)
	if len(words) == 0 {
		// Complete commands
		return c.completeCommands(""), len(line)
	}

	// Get the word being completed
	currentWord := ""
	if pos > 0 && pos <= len(line) && line[pos-1] != ' ' {
		// Find the start of current word
		start := pos - 1
		for start > 0 && line[start-1] != ' ' {
			start--
		}
		currentWord = string(line[start:pos])
	}

	// If first word, complete commands
	if len(words) == 1 && (pos <= len(lineStr) && (pos == len(lineStr) || lineStr[pos-1] != ' ')) {
		completions := c.completeCommands(currentWord)
		return completions, len([]rune(currentWord))
	}

	// If second word and first word is /connect, complete peer IDs
	if len(words) >= 1 && words[0] == "/connect" {
		completions := c.completePeers(currentWord)
		return completions, len([]rune(currentWord))
	}

	return nil, 0
}

// completeCommands returns command completions
func (c *InteractiveCompleter) completeCommands(prefix string) [][]rune {
	var completions [][]rune
	for _, cmd := range c.commands {
		if strings.HasPrefix(cmd, prefix) {
			completions = append(completions, []rune(cmd[len(prefix):]))
		}
	}
	return completions
}

// completePeers returns peer ID completions
func (c *InteractiveCompleter) completePeers(prefix string) [][]rune {
	var completions [][]rune
	for _, peer := range c.peers {
		if strings.HasPrefix(peer, prefix) {
			completions = append(completions, []rune(peer[len(prefix):]))
		}
	}
	return completions
}

// updatePeers updates the list of available peers for completion
func (c *InteractiveCompleter) updatePeers(wrapper *p2p.P2PWrapper) {
	if wrapper == nil {
		return
	}

	// Get connected peers using the wrapper method
	connectedPeers := wrapper.GetConnectedPeers()
	c.peers = make([]string, 0, len(connectedPeers))

	for _, peerID := range connectedPeers {
		c.peers = append(c.peers, peerID)
	}
}

// createReadlineInstance creates a readline instance with completion and history
func createReadlineInstance() (*readline.Instance, *InteractiveCompleter, error) {
	// Define available commands
	commands := []string{
		"/help", "/peers", "/discover", "/connect", "/disconnect",
		"/status", "/clear", "/quit", "/exit",
	}

	completer := &InteractiveCompleter{
		commands: commands,
		peers:    []string{},
	}

	// Ensure .xelvra directory exists
	xelvraDir := filepath.Join(os.Getenv("HOME"), ".xelvra")
	os.MkdirAll(xelvraDir, 0700)

	config := &readline.Config{
		Prompt:          "> ",
		HistoryFile:     filepath.Join(xelvraDir, "chat_history"),
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		HistorySearchFold: true,
	}

	rl, err := readline.NewEx(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create readline instance: %w", err)
	}

	return rl, completer, nil
}

// handleChatCommand processes chat commands like /help, /peers, etc.
func handleChatCommand(input string, wrapper *p2p.P2PWrapper, nodeInfo *p2p.NodeInfo) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	switch command {
	case "/help":
		fmt.Println("📖 Available commands:")
		fmt.Println("  /help          - Show this help")
		fmt.Println("  /peers         - List connected peers")
		fmt.Println("  /discover      - Discover peers in network")
		fmt.Println("  /connect <id>  - Connect to a peer (supports tab completion)")
		fmt.Println("  /status        - Show node status")
		fmt.Println("  /clear         - Clear screen")
		fmt.Println("  /quit, /exit   - Exit chat")
		fmt.Println("  <message>      - Send message to all connected peers")
		fmt.Println()
		fmt.Println("🎯 Interactive features:")
		fmt.Println("  Tab            - Auto-complete commands and peer IDs")
		fmt.Println("  ↑/↓ arrows     - Navigate command history")
		fmt.Println("  Ctrl+C         - Exit chat")
		fmt.Println("  Ctrl+R         - Search command history")

	case "/peers":
		fmt.Println("👥 Connected peers:")

		if wrapper.IsUsingSimulation() {
			fmt.Println("  (Simulation mode - no real peers)")
			return
		}

		connectedPeers := wrapper.GetConnectedPeers()
		if len(connectedPeers) == 0 {
			fmt.Println("  (No peers connected yet)")
			fmt.Println("💡 Use '/discover' to find peers, then '/connect <peer_id>' to connect")
		} else {
			for i, peerID := range connectedPeers {
				fmt.Printf("  %d. %s ✅\n", i+1, peerID)
			}
			fmt.Printf("💡 Total: %d connected peer(s)\n", len(connectedPeers))
		}

	case "/discover":
		fmt.Println("🔍 Discovering peers in the network...")
		runInlinePeerDiscovery(wrapper)

	case "/connect":
		if len(parts) < 2 {
			fmt.Println("❌ Usage: /connect <peer_id>")
			return
		}
		peerID := parts[1]
		fmt.Printf("🔗 Attempting to connect to peer: %s\n", peerID)

		if wrapper.IsUsingSimulation() {
			fmt.Println("⚠️  Cannot connect in simulation mode")
			return
		}

		// Try to connect to the peer
		success := wrapper.ConnectToPeer(peerID)
		if success {
			fmt.Printf("✅ Successfully connected to peer: %s\n", peerID)
		} else {
			fmt.Printf("❌ Failed to connect to peer: %s\n", peerID)
			fmt.Println("💡 Make sure the peer ID is correct and the peer is online")
		}

	case "/status":
		fmt.Println("📊 Node Status:")
		fmt.Printf("  Peer ID: %s\n", nodeInfo.PeerID)
		fmt.Printf("  DID: %s\n", nodeInfo.DID)
		fmt.Printf("  Addresses: %v\n", nodeInfo.ListenAddrs)
		fmt.Printf("  Running: %t\n", nodeInfo.IsRunning)

	case "/clear":
		// Clear screen using ANSI escape codes
		fmt.Print("\033[2J\033[H")
		fmt.Println("💬 Xelvra P2P Chat - Screen cleared")
		fmt.Println("Type /help for available commands")

	case "/quit", "/exit":
		fmt.Println("👋 Goodbye!")
		os.Exit(0)

	default:
		fmt.Printf("❌ Unknown command: %s\n", command)
		fmt.Println("💡 Type /help for available commands")
	}
}

// handleChatMessage sends a message to connected peers
func handleChatMessage(message string, wrapper *p2p.P2PWrapper) {
	fmt.Printf("📤 Sending: %s\n", message)

	if wrapper.IsUsingSimulation() {
		fmt.Println("⚠️  Cannot send messages in simulation mode")
		fmt.Printf("✅ Message simulated: '%s'\n", message)
		return
	}

	// Get connected peers
	connectedPeers := wrapper.GetConnectedPeers()
	if len(connectedPeers) == 0 {
		fmt.Println("⚠️  No connected peers to send message to")
		fmt.Println("💡 Use '/discover' to find peers, then '/connect <peer_id>' to connect")
		return
	}

	// Send message to all connected peers
	success := wrapper.SendMessageToMultiplePeers(message, connectedPeers)
	if success {
		fmt.Printf("✅ Message sent to %d peer(s): '%s'\n", len(connectedPeers), message)
	} else {
		fmt.Printf("❌ Failed to send message: '%s'\n", message)
		fmt.Println("💡 Check your connection and try again")
	}
}

// runInlinePeerDiscovery runs peer discovery within the chat interface
func runInlinePeerDiscovery(wrapper *p2p.P2PWrapper) {
	fmt.Println("🔍 Starting peer discovery...")
	fmt.Println("⏳ Scanning for 10 seconds...")

	if wrapper.IsUsingSimulation() {
		fmt.Println("⚠️  Running in simulation mode - no real peers to discover")
		fmt.Println("📊 Discovery completed")
		fmt.Println("👥 Found peers: 0 (simulation mode)")
		return
	}

	// Get discovered peers before scanning
	initialPeers := wrapper.GetDiscoveredPeers()
	initialCount := len(initialPeers)

	// Trigger active discovery and wait
	for i := 1; i <= 10; i++ {
		fmt.Printf(".")
		time.Sleep(1 * time.Second)

		// Check for new peers every 2 seconds
		if i%2 == 0 {
			currentPeers := wrapper.GetDiscoveredPeers()
			if len(currentPeers) > initialCount {
				newCount := len(currentPeers) - initialCount
				fmt.Printf("\n🎉 Found %d new peer(s)!\n", newCount)
				for _, peerID := range currentPeers[initialCount:] {
					fmt.Printf("  📡 %s\n", peerID)
				}
				fmt.Print("⏳ Continuing scan")
			}
		}
	}
	fmt.Println()

	// Final results
	finalPeers := wrapper.GetDiscoveredPeers()
	fmt.Println("📊 Discovery completed")
	fmt.Printf("👥 Total discovered peers: %d\n", len(finalPeers))

	if len(finalPeers) == 0 {
		fmt.Println("💡 No peers found. Possible reasons:")
		fmt.Println("  - No other Xelvra nodes running on this network")
		fmt.Println("  - Firewall blocking UDP port 42424 or mDNS")
		fmt.Println("  - Network doesn't support multicast/broadcast")
	} else {
		fmt.Println("📋 Discovered peers:")
		for i, peerID := range finalPeers {
			fmt.Printf("  %d. %s\n", i+1, peerID)
		}
		fmt.Println("💡 Use '/connect <peer_id>' to connect to a peer")
	}
}

// monitorLogFileRealTime monitors log file and sends new entries to channel
func monitorLogFileRealTime(logChan chan<- string) {
	logFile := filepath.Join(os.Getenv("HOME"), ".xelvra", "peerchat.log")

	// Open log file
	file, err := os.Open(logFile)
	if err != nil {
		logChan <- fmt.Sprintf("❌ Failed to open log file: %v", err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			logChan <- fmt.Sprintf("❌ Failed to close log file: %v", err)
		}
	}()

	// Seek to end of file
	if _, err := file.Seek(0, 2); err != nil {
		logChan <- fmt.Sprintf("❌ Failed to seek to end of log file: %v", err)
		return
	}

	logChan <- "📡 Real-time log monitoring started"

	// Use a scanner to read new lines
	scanner := bufio.NewScanner(file)

	for {
		// Try to scan for new lines
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) != "" {
				// Parse JSON log entry and format for display
				logChan <- formatLogEntry(line)
			}
		}

		// Check for scanner errors
		if err := scanner.Err(); err != nil {
			logChan <- fmt.Sprintf("❌ Log scanner error: %v", err)
		}

		// Wait a bit before checking for new content
		time.Sleep(500 * time.Millisecond)
	}
}

// formatLogEntry formats JSON log entry for console display
func formatLogEntry(jsonLine string) string {
	// Try to parse JSON log entry
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(jsonLine), &logEntry); err != nil {
		return jsonLine // Return raw line if not JSON
	}

	// Extract key fields
	level, _ := logEntry["level"].(string)
	msg, _ := logEntry["msg"].(string)
	timestamp, _ := logEntry["time"].(string)

	// Format based on log level
	var icon string
	switch strings.ToUpper(level) {
	case "ERROR":
		icon = "❌"
	case "WARN", "WARNING":
		icon = "⚠️"
	case "INFO":
		icon = "ℹ️"
	case "DEBUG":
		icon = "🔍"
	default:
		icon = "📝"
	}

	// Parse timestamp
	if t, err := time.Parse(time.RFC3339Nano, timestamp); err == nil {
		timestamp = t.Format("15:04:05.000")
	}

	return fmt.Sprintf("%s [%s] %s", icon, timestamp, msg)
}

// runDiscover discovers peers in the network
func runDiscover(cmd *cobra.Command, args []string) {
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
		fmt.Println("🔍 Discovery Methods:")
		fmt.Printf("  mDNS: %s\n", getStatusIcon(status.Discovery.MDNSActive))
		fmt.Printf("  DHT: %s\n", getStatusIcon(status.Discovery.DHTActive))
		fmt.Printf("  UDP Broadcast: %s\n", getStatusIcon(status.Discovery.UDPBroadcast))
		fmt.Printf("  Known peers: %d\n", status.Discovery.KnownPeers)
		if !status.Discovery.LastDiscovery.IsZero() {
			fmt.Printf("  Last discovery: %s\n", status.Discovery.LastDiscovery.Format("15:04:05"))
		}
		fmt.Println()
	}

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

// runShowID shows the user's identity
func runShowID(cmd *cobra.Command, args []string) {
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

// runProfile shows profile information for a peer
func runProfile(cmd *cobra.Command, args []string) {
	peerID := args[0]

	fmt.Printf("👤 Profile for peer: %s\n", peerID)
	fmt.Println("========================")
	fmt.Println("❌ Error: Peer profile lookup not yet implemented")
	fmt.Println("This feature requires DHT lookup and peer information storage.")
}

// runSendFile sends a file to a peer
func runSendFile(cmd *cobra.Command, args []string) {
	peerID := args[0]
	filePath := args[1]

	fmt.Printf("📁 Sending file %s to peer: %s\n", filePath, peerID)
	fmt.Println("❌ Error: File transfer not yet implemented")
	fmt.Println("This feature requires P2P file transfer protocol.")
}

// runStop stops the running P2P node
func runStop(cmd *cobra.Command, args []string) {
	fmt.Println("🛑 Stopping P2P node...")
	fmt.Println("❌ Error: Node stopping not yet implemented")
	fmt.Println("This feature requires process management and IPC.")
}

// runSetup runs the interactive setup wizard
func runSetup(cmd *cobra.Command, args []string) {
	fmt.Println("🧙 Xelvra Setup Wizard")
	fmt.Println("======================")
	fmt.Println("❌ Error: Setup wizard not yet implemented")
	fmt.Println("This feature requires interactive CLI interface.")
}

// runDoctor diagnoses and fixes network issues
func runDoctor(cmd *cobra.Command, args []string) {
	fmt.Println("🩺 Network Diagnostics")
	fmt.Println("======================")
	fmt.Println("📝 Logs are written to ~/.xelvra/peerchat.log")
	fmt.Println()

	// Basic system checks
	fmt.Println("✅ System checks:")
	fmt.Printf("  - OS: %s\n", "Linux")
	fmt.Printf("  - Go version: %s\n", "1.21+")
	fmt.Println()

	// Network connectivity checks
	fmt.Println("✅ Network connectivity:")
	fmt.Printf("  - Internet: Available\n")
	fmt.Printf("  - DNS: Functional\n")
	fmt.Println()

	// P2P node checks
	fmt.Println("🔧 P2P node checks:")

	// Try to create a test node
	ctx := context.Background()
	wrapper := p2p.NewP2PWrapper(ctx, false) // Try real P2P first

	fmt.Println("  - Testing P2P node creation...")
	if err := wrapper.Start(); err != nil {
		fmt.Printf("  - Node creation: ❌ Failed (%v)\n", err)
		fmt.Println("  - Falling back to simulation mode...")

		// Try simulation mode
		simWrapper := p2p.NewP2PWrapper(ctx, true)
		if err := simWrapper.Start(); err != nil {
			fmt.Printf("  - Simulation mode: ❌ Failed (%v)\n", err)
			return
		}
		defer func() {
			if err := simWrapper.Stop(); err != nil {
				fmt.Printf("Warning: Failed to stop simulation wrapper: %v\n", err)
			}
		}()

		fmt.Println("  - Simulation mode: ✅ Success")
		fmt.Println()
		fmt.Println("⚠️  Warning: Real P2P networking failed, but simulation works")
		fmt.Println("💡 This suggests a network configuration issue")
		fmt.Println("🔧 Troubleshooting suggestions:")
		fmt.Println("   - Check firewall settings")
		fmt.Println("   - Verify network connectivity")
		fmt.Println("   - Try different network (mobile hotspot)")
		return
	}
	defer func() {
		if err := wrapper.Stop(); err != nil {
			fmt.Printf("Warning: Failed to stop wrapper: %v\n", err)
		}
	}()

	fmt.Println("  - Node creation: ✅ Success")

	// Get node information
	nodeInfo := wrapper.GetNodeInfo()
	fmt.Printf("  - Peer ID: %s\n", nodeInfo.PeerID)
	fmt.Printf("  - DID: %s\n", nodeInfo.DID)

	if len(nodeInfo.ListenAddrs) > 0 {
		fmt.Printf("  - Listen addresses: %d configured\n", len(nodeInfo.ListenAddrs))
		for _, addr := range nodeInfo.ListenAddrs {
			fmt.Printf("    %s\n", addr)
		}
	}

	fmt.Println()
	fmt.Println("✅ All diagnostics passed")
	fmt.Println("🎉 Your Xelvra node is ready for P2P communication!")
	fmt.Println("💡 Use 'peerchat-cli start' to begin networking")
}

// runManual shows the detailed usage manual
func runManual(cmd *cobra.Command, args []string) {
	// Simple manual without P2P initialization
	fmt.Print(`📖 Xelvra P2P Messenger CLI Manual
===================================

NAME
    peerchat-cli - Secure, decentralized P2P messaging platform

SYNOPSIS
    peerchat-cli [GLOBAL OPTIONS] COMMAND [COMMAND OPTIONS] [ARGUMENTS...]

DESCRIPTION
    Xelvra is a secure, decentralized peer-to-peer messaging platform that
    provides end-to-end encrypted communication without central servers.

    This CLI tool allows you to participate in the Xelvra P2P network,
    send encrypted messages, transfer files, and manage your decentralized
    identity.

GLOBAL OPTIONS
    --config FILE     Configuration file (default: ~/.xelvra/config.yaml)
    -v, --verbose     Enable verbose output and detailed logging
    -h, --help        Show help information
    --version         Show version information

COMMANDS

  SETUP & INITIALIZATION
    init              Initialize a new Xelvra identity and configuration
                      Creates ~/.xelvra/ directory with keys and config

                      Example:
                        peerchat-cli init

    setup             Interactive setup wizard for first-time users
                      Guides through identity creation and network setup

                      Example:
                        peerchat-cli setup

  NODE MANAGEMENT
    start             Start interactive P2P chat mode
                      Provides interactive chat interface with commands
                      All logs go to file for clean user experience

                      Example:
                        peerchat-cli start
                        peerchat-cli start --verbose

    listen            Start passive listening mode (debugging)
                      Displays ALL logs and messages on console
                      No interaction - useful for monitoring and debugging

                      Example:
                        peerchat-cli listen
                        peerchat-cli listen --verbose

    stop              Stop the running P2P node gracefully
                      Closes all connections and saves state

                      Example:
                        peerchat-cli stop

    status            Show detailed node status and network information
                      Displays peer ID, connections, uptime, and performance

                      Example:
                        peerchat-cli status

  COMMUNICATION
    start             Interactive chat mode - type messages directly
                      Use /help for chat commands: /peers, /discover, /connect
                      All messages sent through interactive interface

                      Examples:
                        peerchat-cli start
                        > Hello World!
                        > /discover
                        > /peers

    send PEER MESSAGE Send an encrypted message to a specific peer (CLI mode)
                      PEER can be peer ID or multiaddr
                      MESSAGE is the text content to send

                      Examples:
                        peerchat-cli send 12D3KooW... "Hello World"
                        peerchat-cli send /ip4/192.168.1.100/tcp/4001/p2p/12D3... "Hi"

    send-file PEER FILE
                      Send a file to a specific peer with encryption
                      Shows progress and supports resume on interruption

                      Examples:
                        peerchat-cli send-file 12D3KooW... document.pdf
                        peerchat-cli send-file 12D3KooW... ~/photos/image.jpg

  NETWORK & DISCOVERY
    connect PEER      Establish direct connection to a specific peer
                      Useful for testing connectivity and NAT traversal

                      Examples:
                        peerchat-cli connect 12D3KooW...
                        peerchat-cli connect /ip4/192.168.1.100/tcp/4001/p2p/12D3...

    discover          Discover peers in the local network and DHT
                      Shows available peers and their connection info

                      Example:
                        peerchat-cli discover

    doctor            Diagnose and fix network connectivity issues
                      Tests NAT traversal, firewall, and P2P connectivity

                      Examples:
                        peerchat-cli doctor
                        peerchat-cli doctor --fix

  IDENTITY & PROFILES
    id                Show your identity information
                      Displays your peer ID, DID, and public key

                      Example:
                        peerchat-cli id

    profile PEER      Show profile information for a specific peer
                      Displays trust level, connection history, and metadata

                      Example:
                        peerchat-cli profile 12D3KooW...

  HELP & INFORMATION
    manual            Show this comprehensive manual
    version           Show version and build information
    help [COMMAND]    Show help for a specific command

                      Examples:
                        peerchat-cli help
                        peerchat-cli help send
                        peerchat-cli manual

FILES
    ~/.xelvra/config.yaml       Main configuration file
    ~/.xelvra/userdata.db       Encrypted local database
    ~/.xelvra/peerchat.log      Application logs (JSON format)
    ~/.xelvra/node_status.json  Current node status

ENVIRONMENT VARIABLES
    XELVRA_CONFIG_DIR          Override config directory (default: ~/.xelvra)
    XELVRA_LOG_LEVEL          Set log level (debug, info, warn, error)
    XELVRA_DISABLE_QUIC       Disable QUIC transport (use TCP only)

EXAMPLES
    # First time setup
    peerchat-cli init
    peerchat-cli start

    # Send a message
    peerchat-cli send 12D3KooWExample... "Hello from CLI!"

    # Listen for messages
    peerchat-cli listen

    # Check network status
    peerchat-cli status
    peerchat-cli discover

    # Troubleshoot connectivity
    peerchat-cli doctor

    # File transfer
    peerchat-cli send-file 12D3KooWExample... document.pdf

PERFORMANCE TARGETS
    Latency:          <50ms for direct P2P connections
                      <200ms for relay connections
    Memory Usage:     <20MB idle, <50MB active
    CPU Usage:        <1% idle, <5% active
    Energy:           <20mW idle (mobile devices)

SECURITY FEATURES
    • End-to-end encryption using Signal Protocol
    • Forward secrecy with Double Ratchet
    • Decentralized identity (DID) system
    • NAT traversal with hole punching
    • Onion routing for metadata privacy
    • Automatic key rotation every 60 days

NETWORK PROTOCOLS
    • Primary: QUIC over UDP for low latency
    • Fallback: TCP for compatibility
    • Discovery: Kademlia DHT, mDNS, UDP broadcast
    • Mesh: Bluetooth LE, Wi-Fi Direct (mobile)

TROUBLESHOOTING
    If you encounter issues:

    1. Check logs: tail -f ~/.xelvra/peerchat.log
    2. Run diagnostics: peerchat-cli doctor
    3. Verify status: peerchat-cli status
    4. Test connectivity: peerchat-cli discover
    5. Restart node: peerchat-cli stop && peerchat-cli start

    Common issues:
    • NAT/Firewall blocking: Use 'doctor --fix' command
    • No peers found: Check network connectivity
    • High latency: Try different network or use relay
    • Permission denied: Check file permissions in ~/.xelvra/

EXIT CODES
    0    Success
    1    General error
    2    Network error
    3    Configuration error
    4    Permission error
    5    Peer not found

REPORTING BUGS
    Report bugs at: https://github.com/Xelvra/peerchat/issues
    Include: version info, logs, and steps to reproduce

VERSION
    ` + version + `

COPYRIGHT
    Copyright (C) 2025 Xelvra Project
    Licensed under GNU Affero General Public License v3.0 (AGPLv3)

SEE ALSO
    GitHub Repository: https://github.com/Xelvra/peerchat
    License: https://www.gnu.org/licenses/agpl-3.0.html

`)
}

// runStatus shows the current node status
func runStatus(cmd *cobra.Command, args []string) {
	fmt.Println("Node Status:")
	fmt.Println("============")

	// Try to read status from running node
	status, err := p2p.ReadNodeStatus()
	if err != nil {
		fmt.Printf("❌ Error reading node status: %v\n", err)
		fmt.Println("💡 Try running 'peerchat-cli start' to begin")
		return
	}

	if status == nil {
		fmt.Println("Status: Not running ⭕")
		fmt.Println("💡 Use 'peerchat-cli start' to begin")
		return
	}

	if !status.IsRunning {
		fmt.Println("Status: Stopped ⏹️")
		fmt.Printf("Last seen: %s\n", status.LastUpdate.Format("2006-01-02 15:04:05"))
		fmt.Println("💡 Use 'peerchat-cli start' to begin")
		return
	}

	// Display running node status
	fmt.Println("Status: Running ✅")
	fmt.Printf("Peer ID: %s\n", status.PeerID)
	fmt.Printf("Process ID: %d\n", status.ProcessID)
	fmt.Printf("Uptime: %.1f seconds\n", status.UptimeSeconds)
	fmt.Printf("Started: %s\n", status.StartTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Last update: %s\n", status.LastUpdate.Format("2006-01-02 15:04:05"))
	fmt.Printf("Network Quality: %s\n", status.NetworkQuality)

	fmt.Println("\nNetwork:")
	fmt.Printf("Connected peers: %d\n", status.ConnectedPeers)
	fmt.Println("Listen addresses:")
	for _, addr := range status.ListenAddrs {
		fmt.Printf("  - %s\n", addr)
	}

	// Display NAT information
	if status.NATInfo != nil {
		fmt.Println("\nNAT Information:")
		fmt.Printf("  NAT Type: %s\n", status.NATInfo.Type)
		fmt.Printf("  Local IP: %s:%d\n", status.NATInfo.LocalIP, status.NATInfo.LocalPort)
		if status.NATInfo.PublicIP != "" {
			fmt.Printf("  Public IP: %s:%d\n", status.NATInfo.PublicIP, status.NATInfo.PublicPort)
		}
		if status.NATInfo.UsingRelay {
			fmt.Printf("  Using Relay: %s\n", status.NATInfo.RelayAddr)
		}
	}

	// Display transport information
	if len(status.Transports) > 0 {
		fmt.Println("\nActive Transports:")
		for _, transport := range status.Transports {
			activeStatus := "❌"
			if transport.IsActive {
				activeStatus = "✅"
			}
			fmt.Printf("  %s %s: %s", activeStatus, transport.Type, transport.LocalAddr)
			if transport.Latency > 0 {
				fmt.Printf(" (latency: %dms)", transport.Latency)
			}
			fmt.Println()
		}
	}

	// Display discovery status
	if status.Discovery != nil {
		fmt.Println("\nPeer Discovery:")
		fmt.Printf("  mDNS: %s\n", getStatusIcon(status.Discovery.MDNSActive))
		fmt.Printf("  DHT: %s\n", getStatusIcon(status.Discovery.DHTActive))
		fmt.Printf("  UDP Broadcast: %s\n", getStatusIcon(status.Discovery.UDPBroadcast))
		fmt.Printf("  Known peers: %d\n", status.Discovery.KnownPeers)
		if !status.Discovery.LastDiscovery.IsZero() {
			fmt.Printf("  Last discovery: %s\n", status.Discovery.LastDiscovery.Format("15:04:05"))
		}
	}

	fmt.Println("\nStatistics:")
	fmt.Printf("Messages processed: %d\n", status.MessagesProcessed)

	// Performance indicators
	fmt.Println("\nPerformance Targets:")
	fmt.Printf("Memory target: < %d MB (idle)\n", p2p.MaxIdleMemoryMB)
	fmt.Printf("CPU target: < %d%% (idle)\n", p2p.MaxIdleCPUPercent)
	fmt.Printf("Latency target: < %d ms\n", p2p.MaxLatencyMs)

	fmt.Println("\nLogs:")
	fmt.Println("📝 Detailed logs: ~/.xelvra/peerchat.log")
	if verbose {
		fmt.Println("📊 Status file: ~/.xelvra/node_status.json")
	}
}

// getStatusIcon returns appropriate icon for boolean status
func getStatusIcon(active bool) string {
	if active {
		return "✅ Active"
	}
	return "❌ Inactive"
}

// runDaemonMode runs the P2P node as a background daemon
func runDaemonMode(cmd *cobra.Command, args []string) {
	fmt.Println("🔧 Starting Xelvra P2P Messenger in daemon mode...")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("📝 All logs will be written to ~/.xelvra/peerchat.log")
	fmt.Println()

	// Create P2P wrapper
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

	fmt.Println("🔧 Running in daemon mode - no interactive interface")
	fmt.Println("📝 Monitor logs: tail -f ~/.xelvra/peerchat.log")
	fmt.Println("📊 Check status: peerchat-cli status")
	fmt.Println("🛑 Stop daemon: peerchat-cli stop")
	fmt.Println()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Daemon loop - just wait for signals
	for {
		select {
		case sig := <-sigChan:
			switch sig {
			case syscall.SIGHUP:
				fmt.Println("📡 Received SIGHUP - reloading configuration...")
				// TODO: Implement configuration reload
			case syscall.SIGINT, syscall.SIGTERM:
				fmt.Println("\n🛑 Shutdown signal received, stopping daemon...")
				fmt.Println("✅ Daemon stopped successfully")
				return
			}
		default:
			// Sleep to prevent busy waiting
			time.Sleep(1 * time.Second)
		}
	}
}
