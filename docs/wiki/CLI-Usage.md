# CLI Usage Guide

Comprehensive guide for using the Xelvra command-line interface effectively.

## 📋 Table of Contents

- [Quick Start](#quick-start)
- [Command Overview](#command-overview)
- [Interactive Mode](#interactive-mode)
- [Common Workflows](#common-workflows)
- [Advanced Usage](#advanced-usage)
- [Troubleshooting](#troubleshooting)

## 🚀 Quick Start

### Basic Setup and First Chat

```bash
# 1. Initialize your identity
peerchat-cli init

# 2. Test your setup
peerchat-cli doctor

# 3. Start interactive chat
peerchat-cli start
```

Once in interactive mode:
```
> /discover
> /connect 12D3KooWExample...
> Hello! This is my first message on Xelvra! 👋
```

## 📖 Command Overview

### Essential Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `init` | Initialize identity | `peerchat-cli init` |
| `start` | Interactive chat | `peerchat-cli start` |
| `discover` | Find peers | `peerchat-cli discover` |
| `connect` | Connect to peer | `peerchat-cli connect 12D3KooW...` |
| `send` | Send message | `peerchat-cli send 12D3KooW... "Hello"` |
| `status` | Node status | `peerchat-cli status` |
| `doctor` | Diagnostics | `peerchat-cli doctor` |

### File Operations

| Command | Purpose | Example |
|---------|---------|---------|
| `send-file` | Send file | `peerchat-cli send-file 12D3KooW... file.pdf` |

### Information Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `id` | Show identity | `peerchat-cli id` |
| `version` | Show version | `peerchat-cli version` |
| `manual` | Show help | `peerchat-cli manual` |

## 💬 Interactive Mode

### Starting Interactive Mode

```bash
peerchat-cli start
```

You'll see:
```
🚀 Xelvra P2P Messenger - Interactive Chat
Your Peer ID: 12D3KooWYourPeerID...
Your DID: did:xelvra:yourhash...

Network Status: ✅ Online
Listening on: /ip4/192.168.1.50/tcp/4001

Type /help for commands, or just start typing to send messages!
> 
```

### Interactive Commands

#### Essential Commands
```bash
/help           # Show all available commands
/peers          # List connected peers
/discover       # Find peers on network
/connect <id>   # Connect to specific peer
/status         # Show node status
/quit           # Exit chat mode
```

#### Discovery and Connection
```bash
# Discover peers
> /discover
🔍 Starting peer discovery...
⏳ Scanning for 10 seconds...
📊 Discovery completed
👥 Total discovered peers: 2
📋 Discovered peers:
  1. 12D3KooWAlice... (Alice's Node)
  2. 12D3KooWBob... (Bob's Node)

# Connect to a peer
> /connect 12D3KooWAlice...
🔗 Attempting to connect to peer: Alice's Node
✅ Successfully connected to peer: Alice's Node
```

#### Messaging
```bash
# Send message (just type and press Enter)
> Hello Alice! How are you today? 😊
📤 Sending: Hello Alice! How are you today? 😊
✅ Message sent to 1 peer(s)

# Receive message
📨 Message from Alice: Hi there! I'm doing great, thanks for asking! 🎉
```

#### File Transfer
```bash
# Send file
> /send-file 12D3KooWAlice... ~/Documents/report.pdf
📤 Sending file: report.pdf (2.5 MB)
⏳ Transfer progress: [████████████████████] 100%
✅ File sent successfully to Alice's Node

# File received notification
📨 Incoming file from Alice: presentation.pptx (5.2 MB)
📁 Saved to: ~/.xelvra/downloads/presentation.pptx
✅ File received successfully
```

#### Peer Management
```bash
# List current connections
> /peers
👥 Connected peers (2):
  1. Alice's Node (12D3KooWAlice...) - Connected 15m ago
  2. Bob's Node (12D3KooWBob...) - Connected 8m ago

# Disconnect from peer
> /disconnect 12D3KooWBob...
🔌 Disconnected from peer: Bob's Node
```

### Navigation and Shortcuts

- **Arrow Keys** - Navigate command history
- **Tab** - Auto-complete commands and peer IDs
- **Ctrl+C** - Exit chat mode
- **Ctrl+L** - Clear screen
- **Ctrl+A** - Move to beginning of line
- **Ctrl+E** - Move to end of line

## 🔄 Common Workflows

### Workflow 1: First Time Setup

```bash
# Step 1: Initialize
peerchat-cli init
# Creates ~/.xelvra/ directory with your identity

# Step 2: Verify setup
peerchat-cli doctor
# Checks network connectivity and configuration

# Step 3: View your identity
peerchat-cli id
# Shows your Peer ID and DID for sharing

# Step 4: Start chatting
peerchat-cli start
# Enters interactive mode
```

### Workflow 2: Quick Message to Known Peer

```bash
# If you know the peer ID, connect and send directly
peerchat-cli connect 12D3KooWExample...
peerchat-cli send 12D3KooWExample... "Quick message!"

# Or do it all in interactive mode
peerchat-cli start
> /connect 12D3KooWExample...
> Quick message!
```

### Workflow 3: File Sharing Session

```bash
# Start interactive mode
peerchat-cli start

# Discover and connect
> /discover
> /connect 12D3KooWFriend...

# Send file
> /send-file 12D3KooWFriend... ~/Documents/important.pdf

# Continue chatting
> File sent! Let me know when you receive it.
```

### Workflow 4: Network Troubleshooting

```bash
# Check system health
peerchat-cli doctor

# If issues found, try detailed diagnostics
peerchat-cli doctor --detailed

# Test specific components
peerchat-cli doctor --test network

# Monitor network activity
peerchat-cli listen --verbose

# Check current status
peerchat-cli status --format json
```

### Workflow 5: Multiple Identity Management

```bash
# Create work identity
peerchat-cli init --config-dir ~/.xelvra-work

# Create personal identity  
peerchat-cli init --config-dir ~/.xelvra-personal

# Use work identity
peerchat-cli start --config-dir ~/.xelvra-work

# Use personal identity
peerchat-cli start --config-dir ~/.xelvra-personal
```

## 🔧 Advanced Usage

### Custom Configuration

```bash
# Use custom config directory
peerchat-cli start --config-dir ~/.xelvra-custom

# Specify network interface
peerchat-cli start --interface eth0

# Use specific port
peerchat-cli start --port 8080

# Disable discovery
peerchat-cli start --no-discovery
```

### Output Formatting

```bash
# JSON output for scripting
peerchat-cli status --format json
peerchat-cli discover --format json

# YAML output
peerchat-cli status --format yaml

# Watch mode for monitoring
peerchat-cli status --watch --interval 5
```

### Batch Operations

```bash
# Discover and save peer list
peerchat-cli discover --format json > peers.json

# Connect to multiple peers from file
cat peers.json | jq -r '.peers[].id' | while read peer; do
    peerchat-cli connect "$peer"
done

# Send message to all connected peers
peerchat-cli status --format json | jq -r '.connections.peers[].id' | while read peer; do
    peerchat-cli send "$peer" "Broadcast message"
done
```

### Debugging and Monitoring

```bash
# Enable debug logging
peerchat-cli start --log-level debug

# Monitor all network activity
peerchat-cli listen --verbose

# Save debug output to file
peerchat-cli listen --verbose --output debug.log

# Filter log messages
peerchat-cli listen --filter error
```

### Automation Scripts

#### Auto-discovery Script
```bash
#!/bin/bash
# auto-discover.sh - Automatically discover and connect to peers

echo "Starting auto-discovery..."
PEERS=$(peerchat-cli discover --timeout 30 --format json | jq -r '.peers[].id')

for peer in $PEERS; do
    echo "Connecting to $peer..."
    peerchat-cli connect "$peer" --timeout 10
done

echo "Auto-discovery complete!"
```

#### Status Monitor Script
```bash
#!/bin/bash
# monitor.sh - Monitor node status

while true; do
    STATUS=$(peerchat-cli status --format json)
    PEERS=$(echo "$STATUS" | jq -r '.connections.total')
    echo "$(date): Connected to $PEERS peers"
    sleep 60
done
```

## 🛠️ Troubleshooting

### Common Issues and Solutions

#### "Command not found"
```bash
# Check if binary is in PATH
which peerchat-cli

# If not, use full path or add to PATH
export PATH=$PATH:/path/to/peerchat-cli
```

#### "No peers found"
```bash
# Check network connectivity
peerchat-cli doctor

# Try different discovery methods
peerchat-cli discover --method mdns
peerchat-cli discover --method udp

# Check firewall settings
sudo ufw status
sudo ufw allow 42424/udp
```

#### "Connection failed"
```bash
# Verify peer is still online
peerchat-cli discover | grep <peer_id>

# Try with longer timeout
peerchat-cli connect <peer_id> --timeout 60

# Check detailed status
peerchat-cli status --detailed
```

#### "Permission denied"
```bash
# Check file permissions
ls -la ~/.xelvra/

# Fix permissions if needed
chmod 600 ~/.xelvra/identity.key
chmod 755 ~/.xelvra/
```

### Debug Information Collection

```bash
# Collect system information
peerchat-cli version > debug_info.txt
peerchat-cli doctor --detailed >> debug_info.txt
uname -a >> debug_info.txt

# Collect recent logs
tail -100 ~/.xelvra/peerchat.log >> debug_info.txt

# Test network connectivity
peerchat-cli discover --timeout 30 --format json >> debug_info.txt
```

### Performance Optimization

```bash
# Reduce resource usage
peerchat-cli start --max-peers 10

# Optimize for mobile/low-power
export XELVRA_MAX_PEERS=5
export XELVRA_CONNECTION_TIMEOUT=15s
peerchat-cli start
```

## 📞 Getting Help

### Built-in Help

```bash
# General help
peerchat-cli --help

# Command-specific help
peerchat-cli connect --help
peerchat-cli send-file --help

# Interactive help
peerchat-cli start
> /help
```

### Manual Pages

```bash
# View built-in manual
peerchat-cli manual

# View specific command manual
peerchat-cli manual connect
peerchat-cli manual send-file
```

### Community Resources

- **[GitHub Wiki](https://github.com/Xelvra/peerchat/wiki)** - Complete documentation
- **[GitHub Discussions](https://github.com/Xelvra/peerchat/discussions)** - Community Q&A
- **[Troubleshooting Guide](Troubleshooting)** - Detailed troubleshooting
- **[FAQ](FAQ)** - Frequently asked questions

---

**Master the CLI and become a Xelvra power user!** 🚀 *#XelvraFree*
