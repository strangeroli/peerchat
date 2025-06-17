package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Xelvra/peerchat/internal/p2p"
	"github.com/chzyer/readline"
)

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

// UpdatePeers updates the list of available peers for completion
func (c *InteractiveCompleter) UpdatePeers(wrapper *p2p.P2PWrapper) {
	if wrapper == nil {
		return
	}

	// Get connected peers using the wrapper method
	connectedPeers := wrapper.GetConnectedPeers()
	c.peers = append(c.peers[:0], connectedPeers...)
}

// CreateReadlineInstance creates a readline instance with completion and history
func CreateReadlineInstance() (*readline.Instance, *InteractiveCompleter, error) {
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
	if err := os.MkdirAll(xelvraDir, 0700); err != nil {
		return nil, nil, fmt.Errorf("failed to create xelvra directory: %w", err)
	}

	config := &readline.Config{
		Prompt:            "> ",
		HistoryFile:       filepath.Join(xelvraDir, "chat_history"),
		AutoComplete:      completer,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	}

	rl, err := readline.NewEx(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create readline instance: %w", err)
	}

	return rl, completer, nil
}

// HandleChatCommand processes chat commands like /help, /peers, etc.
func HandleChatCommand(input string, wrapper *p2p.P2PWrapper, nodeInfo *p2p.NodeInfo) {
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
		RunInlinePeerDiscovery(wrapper)

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

// HandleChatMessage sends a message to connected peers
func HandleChatMessage(message string, wrapper *p2p.P2PWrapper) {
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
