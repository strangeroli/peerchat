# Xelvra P2P Messenger - User Guide

## Table of Contents
- [Getting Started](#getting-started)
- [Installation](#installation)
- [First Time Setup](#first-time-setup)
- [Basic Usage](#basic-usage)
- [Interactive Chat](#interactive-chat)
- [Peer Discovery](#peer-discovery)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)

## Getting Started

Xelvra is a secure, decentralized peer-to-peer messaging platform that allows you to communicate directly with other users without relying on central servers. Your messages are end-to-end encrypted and your identity is self-sovereign.

### Key Features
- **Decentralized**: No central servers required
- **Secure**: End-to-end encryption for all messages
- **Private**: Self-sovereign identity management
- **Cross-platform**: Works on Linux, macOS, and Windows
- **Real-time**: Instant message delivery over P2P networks

## Installation

### Prerequisites
- Go 1.21 or later
- Git
- Network connectivity (for P2P communication)

### Build from Source
```bash
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go
```

### Verify Installation
```bash
./bin/peerchat-cli version
./bin/peerchat-cli doctor
```

## First Time Setup

### 1. Initialize Your Identity
```bash
./bin/peerchat-cli init
```
This creates your cryptographic identity and configuration files in `~/.xelvra/`.

### 2. Test Network Connectivity
```bash
./bin/peerchat-cli doctor
```
This diagnoses your network setup and ensures P2P communication will work.

### 3. Get Your Identity
```bash
./bin/peerchat-cli id
```
This displays your Peer ID and DID that others can use to connect to you.

## Basic Usage

### Starting Interactive Chat
```bash
./bin/peerchat-cli start
```
This starts the interactive chat mode where you can:
- Send messages to connected peers
- Discover other users on your network
- Connect to specific peers
- View your connection status

### Passive Listening (Debug Mode)
```bash
./bin/peerchat-cli listen
```
This starts passive listening mode that displays all logs and network activity. Useful for debugging and monitoring.

### Check Node Status
```bash
./bin/peerchat-cli status
```
Shows detailed information about your node, connections, and network status.

## Interactive Chat

When you run `peerchat-cli start`, you enter interactive chat mode with the following commands:

### Chat Commands
- `/help` - Show available commands
- `/peers` - List currently connected peers
- `/discover` - Discover peers on your network
- `/connect <peer_id>` - Connect to a specific peer
- `/status` - Show your node status
- `/quit` - Exit the chat

### Sending Messages
Simply type your message and press Enter. It will be sent to all connected peers.

```
> Hello, world!
📤 Sending: Hello, world!
✅ Message sent to 2 peer(s): 'Hello, world!'
```

### Discovering Peers
Use the `/discover` command to find other Xelvra users on your network:

```
> /discover
🔍 Starting peer discovery...
⏳ Scanning for 10 seconds...
..........
📊 Discovery completed
👥 Total discovered peers: 2
📋 Discovered peers:
  1. 12D3KooWExample1...
  2. 12D3KooWExample2...
💡 Use '/connect <peer_id>' to connect to a peer
```

### Connecting to Peers
Once you've discovered peers, connect to them:

```
> /connect 12D3KooWExample1...
🔗 Attempting to connect to peer: 12D3KooWExample1...
✅ Successfully connected to peer: 12D3KooWExample1...
```

## Peer Discovery

Xelvra uses multiple methods to discover peers:

### 1. mDNS (Local Network)
Automatically discovers peers on your local network using multicast DNS.

### 2. UDP Broadcast
Sends broadcast messages on your local network to find nearby peers.

### 3. Manual Connection
Connect directly to peers using their Peer ID if you know it.

### Network Requirements
- **Firewall**: Allow UDP port 42424 for peer discovery
- **Multicast**: Your network must support multicast/broadcast
- **NAT**: Most NAT configurations work automatically

## Troubleshooting

### Common Issues

#### No Peers Found
**Problem**: `/discover` finds no peers even though others are running Xelvra.

**Solutions**:
1. Check firewall settings - allow UDP port 42424
2. Ensure you're on the same network
3. Try different network (mobile hotspot)
4. Run `peerchat-cli doctor` for diagnostics

#### Connection Failed
**Problem**: Cannot connect to discovered peers.

**Solutions**:
1. Verify the peer is still online
2. Check network connectivity
3. Ensure both nodes are using compatible versions
4. Try discovering again - peer addresses may have changed

#### Simulation Mode
**Problem**: Node starts in simulation mode instead of real P2P.

**Solutions**:
1. Check network connectivity
2. Verify firewall settings
3. Try different network interface
4. Run `peerchat-cli doctor` for detailed diagnostics

### Getting Help
- Run `peerchat-cli doctor` for automated diagnostics
- Check logs in `~/.xelvra/peerchat.log`
- Use `peerchat-cli listen` to see real-time network activity
- Consult the [Developer Guide](DEVELOPER_GUIDE.md) for technical details

## FAQ

### Q: Is Xelvra secure?
A: Yes, all messages are end-to-end encrypted using modern cryptographic protocols. Your identity is self-sovereign and not controlled by any central authority.

### Q: Do I need to open ports in my firewall?
A: For local network discovery, allow UDP port 42424. For direct connections, Xelvra uses dynamic ports and works with most NAT configurations.

### Q: Can I use Xelvra over the internet?
A: Currently, Xelvra is optimized for local network communication. Internet-wide communication is planned for future releases.

### Q: Where are my keys stored?
A: Your cryptographic keys and configuration are stored in `~/.xelvra/`. Keep this directory secure and backed up.

### Q: How do I backup my identity?
A: Backup the entire `~/.xelvra/` directory. This contains your keys, configuration, and identity information.

### Q: Can I run multiple instances?
A: Yes, but each instance needs its own configuration directory. Use the `--config` flag to specify different directories.

---

For technical details and development information, see the [Developer Guide](DEVELOPER_GUIDE.md).
