package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Xelvra/peerchat/internal/p2p"
)

// MonitorLogFileRealTime monitors log file and sends new entries to channel
func MonitorLogFileRealTime(logChan chan<- string) {
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
				logChan <- FormatLogEntry(line)
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

// FormatLogEntry formats JSON log entry for console display
func FormatLogEntry(jsonLine string) string {
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

// RunInlinePeerDiscovery runs peer discovery within the chat interface
func RunInlinePeerDiscovery(wrapper *p2p.P2PWrapper) {
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
