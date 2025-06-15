# User Manual

Complete guide for using Xelvra P2P Messenger effectively and securely.

## 📋 Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
- [Command Reference](#command-reference)
- [Interactive Chat](#interactive-chat)
- [Peer Management](#peer-management)
- [File Transfer](#file-transfer)
- [Configuration](#configuration)
- [Security Features](#security-features)
- [Advanced Usage](#advanced-usage)

## 🌟 Overview

Xelvra P2P Messenger is a secure, decentralized communication platform that enables direct peer-to-peer messaging without central servers. This manual covers all features and functionality available in the current CLI version.

### Key Concepts

**Peer-to-Peer (P2P)**: Messages travel directly between devices without intermediary servers.

**End-to-End Encryption**: Only you and your intended recipient can read your messages.

**Decentralized Identity**: Your identity is self-sovereign and stored locally on your device.

**Peer ID**: A unique cryptographic identifier used for connections (e.g., `12D3KooW...`).

**DID**: Decentralized Identifier in the format `did:xelvra:<hash>`.

## 🚀 Getting Started

### First-Time Setup

1. **Initialize your identity:**
```bash
peerchat-cli init
```

2. **Test your setup:**
```bash
peerchat-cli doctor
```

3. **View your identity:**
```bash
peerchat-cli id
```

4. **Start interactive chat:**
```bash
peerchat-cli start
```

## 📖 Command Reference

### Core Commands

#### `init`
Initialize your cryptographic identity and configuration.
```bash
peerchat-cli init [--config-dir PATH]
```

**Options:**
- `--config-dir`: Specify custom configuration directory (default: `~/.xelvra`)

#### `start`
Start interactive chat mode.
```bash
peerchat-cli start [--port PORT] [--interface INTERFACE]
```

**Options:**
- `--port`: Specify listening port (0 for auto-select)
- `--interface`: Specify network interface to use

#### `listen`
Start passive listening mode (shows all logs and network activity).
```bash
peerchat-cli listen [--verbose]
```

**Options:**
- `--verbose`: Show detailed debug information

#### `status`
Display current node status and connections.
```bash
peerchat-cli status [--detailed]
```

**Options:**
- `--detailed`: Show comprehensive status information

#### `discover`
Discover peers on your network.
```bash
peerchat-cli discover [--timeout SECONDS] [--method METHOD]
```

**Options:**
- `--timeout`: Discovery timeout in seconds (default: 10)
- `--method`: Discovery method (mdns, udp, all)

#### `connect`
Connect to a specific peer.
```bash
peerchat-cli connect <PEER_ID> [--timeout SECONDS]
```

**Arguments:**
- `PEER_ID`: Target peer's ID or multiaddress

**Options:**
- `--timeout`: Connection timeout in seconds

#### `send`
Send a message to connected peers.
```bash
peerchat-cli send <PEER_ID> <MESSAGE>
```

**Arguments:**
- `PEER_ID`: Target peer's ID
- `MESSAGE`: Message text to send

#### `send-file`
Send a file to a peer.
```bash
peerchat-cli send-file <PEER_ID> <FILE_PATH>
```

**Arguments:**
- `PEER_ID`: Target peer's ID
- `FILE_PATH`: Path to file to send

#### `doctor`
Run network diagnostics and system checks.
```bash
peerchat-cli doctor [--detailed] [--fix]
```

**Options:**
- `--detailed`: Show detailed diagnostic information
- `--fix`: Attempt to fix common issues automatically

#### `id`
Display your identity information.
```bash
peerchat-cli id [--qr] [--export]
```

**Options:**
- `--qr`: Display identity as QR code
- `--export`: Export identity for sharing

#### `version`
Show version information.
```bash
peerchat-cli version [--detailed]
```

#### `manual`
Display built-in manual pages.
```bash
peerchat-cli manual [COMMAND]
```

## 💬 Interactive Chat

When you run `peerchat-cli start`, you enter interactive chat mode with these features:

### Chat Commands

All interactive commands start with `/`:

- `/help` - Show available commands
- `/peers` - List currently connected peers
- `/discover` - Discover peers on your network
- `/connect <peer_id>` - Connect to a specific peer
- `/disconnect <peer_id>` - Disconnect from a peer
- `/status` - Show your node status
- `/clear` - Clear the chat screen
- `/quit` or `/exit` - Exit the chat

### Sending Messages

Simply type your message and press Enter:
```
> Hello, world! 👋
📤 Sending: Hello, world! 👋
✅ Message sent to 2 peer(s)
```

### Receiving Messages

Incoming messages are displayed with sender information:
```
📨 Message from Alice (12D3KooW...): Hi there! How are you?
```

### Navigation

- **Arrow Keys**: Navigate command history
- **Tab**: Auto-complete commands and peer IDs
- **Ctrl+C**: Exit chat mode
- **Ctrl+L**: Clear screen

## 👥 Peer Management

### Discovering Peers

Use the discovery system to find other Xelvra users:

```bash
# In CLI mode
peerchat-cli discover

# In interactive mode
> /discover
🔍 Starting peer discovery...
⏳ Scanning for 10 seconds...
📊 Discovery completed
👥 Total discovered peers: 3
📋 Discovered peers:
  1. 12D3KooWExample1... (Alice)
  2. 12D3KooWExample2... (Bob)
  3. 12D3KooWExample3... (Charlie)
```

### Connecting to Peers

```bash
# Connect to a specific peer
> /connect 12D3KooWExample1...
🔗 Attempting to connect to peer: Alice
✅ Successfully connected to peer: Alice
```

### Managing Connections

```bash
# List current connections
> /peers
👥 Connected peers (2):
  1. Alice (12D3KooWExample1...)
  2. Bob (12D3KooWExample2...)

# Disconnect from a peer
> /disconnect 12D3KooWExample1...
🔌 Disconnected from peer: Alice
```

## 📁 File Transfer

### Sending Files

```bash
# Send a file to a connected peer
peerchat-cli send-file 12D3KooWExample1... /path/to/document.pdf

# In interactive mode
> /send-file 12D3KooWExample1... ~/Documents/report.pdf
📤 Sending file: report.pdf (2.5 MB)
⏳ Transfer progress: [████████████████████] 100%
✅ File sent successfully to Alice
```

### Receiving Files

Files are automatically received and saved to your downloads directory:

```
📨 Incoming file from Bob: presentation.pptx (5.2 MB)
📁 Saved to: ~/.xelvra/downloads/presentation.pptx
✅ File received successfully
```

### File Transfer Features

- **Chunked Transfer**: Large files are split into chunks for reliability
- **Progress Tracking**: Real-time progress indicators
- **Integrity Verification**: Automatic checksum verification
- **Resume Support**: Interrupted transfers can be resumed
- **Encryption**: All file transfers are end-to-end encrypted

## ⚙️ Configuration

### Configuration File

Edit `~/.xelvra/config.yaml` to customize Xelvra:

```yaml
# Network settings
network:
  listen_port: 0          # 0 = auto-select port
  discovery_port: 42424   # UDP discovery port
  enable_mdns: true       # Local network discovery
  enable_udp_broadcast: true
  max_peers: 50          # Maximum concurrent connections
  connection_timeout: 30s

# User settings
user:
  display_name: "Your Name"
  auto_accept_files: false
  download_directory: "~/.xelvra/downloads"

# Logging settings
logging:
  level: "info"          # debug, info, warn, error
  file: "peerchat.log"
  max_size: 10           # MB
  max_backups: 3
  max_age: 30            # days

# Security settings
security:
  key_rotation_days: 60  # Automatic key rotation
  max_message_size: 1048576  # 1MB
  enable_forward_secrecy: true
```

### Environment Variables

You can also configure Xelvra using environment variables:

```bash
export XELVRA_CONFIG_DIR="~/.xelvra"
export XELVRA_LOG_LEVEL="debug"
export XELVRA_LISTEN_PORT="0"
export XELVRA_DISCOVERY_PORT="42424"
```

## 🔒 Security Features

### Identity Security

- **Private Key Protection**: Your private key never leaves your device
- **Key Backup**: Always backup your `~/.xelvra/` directory
- **Key Rotation**: Automatic key rotation every 60 days (configurable)

### Message Security

- **End-to-End Encryption**: All messages encrypted with Signal Protocol
- **Forward Secrecy**: Past messages remain secure even if keys are compromised
- **Message Integrity**: Digital signatures verify message authenticity

### Network Security

- **Peer Authentication**: Cryptographic verification of peer identities
- **Transport Encryption**: All network traffic is encrypted
- **Metadata Protection**: Communication patterns are obfuscated

### Best Practices

1. **Regular Backups**: Backup your identity regularly
2. **Secure Networks**: Use trusted networks when possible
3. **Update Regularly**: Keep Xelvra updated to the latest version
4. **Verify Peers**: Only connect to trusted peers
5. **Monitor Logs**: Review logs for unusual activity

## 🔧 Advanced Usage

### Multiple Instances

Run multiple Xelvra instances with different identities:

```bash
# Create separate config directories
mkdir ~/.xelvra-work ~/.xelvra-personal

# Initialize separate identities
peerchat-cli init --config-dir ~/.xelvra-work
peerchat-cli init --config-dir ~/.xelvra-personal

# Run with specific config
peerchat-cli start --config-dir ~/.xelvra-work
```

### Custom Network Interfaces

Specify which network interface to use:

```bash
# List available interfaces
ip addr show  # Linux
ifconfig      # macOS
ipconfig      # Windows

# Use specific interface
peerchat-cli start --interface eth0
peerchat-cli start --interface wlan0
```

### Debugging and Diagnostics

```bash
# Enable debug logging
peerchat-cli start --log-level debug

# Monitor network activity
peerchat-cli listen --verbose

# Run comprehensive diagnostics
peerchat-cli doctor --detailed

# Test specific features
peerchat-cli doctor --test-discovery
peerchat-cli doctor --test-nat-traversal
```

### Automation and Scripting

```bash
# Automated peer discovery
peerchat-cli discover --timeout 30 --format json > peers.json

# Batch file sending
for file in *.txt; do
    peerchat-cli send-file $PEER_ID "$file"
done

# Status monitoring
while true; do
    peerchat-cli status --format json | jq '.connected_peers'
    sleep 60
done
```

## 🆘 Getting Help

### Built-in Help

```bash
# General help
peerchat-cli --help

# Command-specific help
peerchat-cli connect --help

# Interactive help
> /help
```

### Troubleshooting

1. **Run diagnostics**: `peerchat-cli doctor`
2. **Check logs**: `tail -f ~/.xelvra/peerchat.log`
3. **Test connectivity**: `peerchat-cli listen`
4. **Reset configuration**: `peerchat-cli init --reset`

### Community Support

- **[GitHub Wiki](https://github.com/Xelvra/peerchat/wiki)** - Comprehensive documentation
- **[GitHub Discussions](https://github.com/Xelvra/peerchat/discussions)** - Community Q&A
- **[GitHub Issues](https://github.com/Xelvra/peerchat/issues)** - Bug reports and feature requests
- **[FAQ](https://github.com/Xelvra/peerchat/wiki/FAQ)** - Frequently asked questions

---

**Welcome to the decentralized future of communication!** 🌐 *#XelvraFree*
