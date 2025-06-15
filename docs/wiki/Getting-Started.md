# Getting Started with Xelvra P2P Messenger

Welcome to Xelvra! This guide will help you get up and running with secure, decentralized peer-to-peer messaging in just a few minutes.

## 🎯 What is Xelvra?

Xelvra is a **secure, decentralized P2P communication platform** that allows you to:
- Send messages directly to other users without central servers
- Maintain complete privacy with end-to-end encryption
- Own your data and identity
- Communicate even in restrictive network environments

## 🚀 Quick Start (5 Minutes)

### Step 1: Installation

#### Option A: Download Pre-built Binary (Recommended)
```bash
# Download the latest release
curl -L https://github.com/Xelvra/peerchat/releases/latest/download/peerchat-cli-linux -o peerchat-cli
chmod +x peerchat-cli
```

#### Option B: Build from Source
```bash
# Prerequisites: Go 1.21+, Git
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go
```

### Step 2: Initialize Your Identity
```bash
./peerchat-cli init
```
This creates your cryptographic identity and stores it securely in `~/.xelvra/`.

### Step 3: Test Your Setup
```bash
./peerchat-cli doctor
```
This runs network diagnostics to ensure everything is working correctly.

### Step 4: Start Chatting!
```bash
./peerchat-cli start
```
This launches the interactive chat interface where you can discover and connect to other users.

## 🔍 Your First Chat Session

Once you run `peerchat-cli start`, you'll see the interactive chat interface:

```
🚀 Xelvra P2P Messenger - Interactive Chat
Your Peer ID: 12D3KooWExample...
Your DID: did:xelvra:abc123...

Type /help for commands, or just start typing to send messages!
> 
```

### Essential Commands
- `/discover` - Find other Xelvra users on your network
- `/connect <peer_id>` - Connect to a specific user
- `/peers` - List your current connections
- `/status` - Check your node status
- `/help` - Show all available commands
- `/quit` - Exit the chat

### Discovering Other Users
```
> /discover
🔍 Starting peer discovery...
⏳ Scanning for 10 seconds...
..........
📊 Discovery completed
👥 Total discovered peers: 2
📋 Discovered peers:
  1. 12D3KooWExample1... (Alice)
  2. 12D3KooWExample2... (Bob)
💡 Use '/connect <peer_id>' to connect to a peer
```

### Connecting and Messaging
```
> /connect 12D3KooWExample1...
🔗 Attempting to connect to peer: 12D3KooWExample1...
✅ Successfully connected to peer: Alice

> Hello Alice! 👋
📤 Sending: Hello Alice! 👋
✅ Message sent to 1 peer(s): 'Hello Alice! 👋'

📨 Message from Alice: Hi there! Welcome to Xelvra! 🎉
```

## 🔧 Basic Configuration

### Configuration Files
Xelvra stores its configuration in `~/.xelvra/`:
- `config.yaml` - Main configuration
- `identity.key` - Your private key (keep secure!)
- `peerchat.log` - Application logs

### Important Settings
```yaml
# ~/.xelvra/config.yaml
network:
  listen_port: 0  # 0 = auto-select port
  discovery_port: 42424
  enable_mdns: true
  enable_udp_broadcast: true

logging:
  level: "info"
  file: "peerchat.log"
  max_size: 10  # MB
  max_backups: 3
```

## 🌐 Network Requirements

### Firewall Configuration
For optimal performance, allow these ports:
- **UDP 42424** - Peer discovery (recommended)
- **Dynamic TCP/UDP ports** - P2P connections (automatic)

### Network Types
Xelvra works on various networks:
- ✅ **Home WiFi** - Full functionality
- ✅ **Office Networks** - Usually works with discovery
- ✅ **Mobile Hotspots** - Full functionality
- ⚠️ **Public WiFi** - May have limitations
- ⚠️ **Corporate Networks** - May require firewall configuration

## 🔒 Security & Privacy

### Your Identity
- Your **Peer ID** is public and used for connections
- Your **DID** (Decentralized Identifier) is your unique identity
- Your **private key** is stored locally and never shared

### Message Security
- All messages are **end-to-end encrypted**
- **Forward secrecy** protects past messages
- **Metadata protection** hides communication patterns

### Best Practices
1. **Backup your identity**: Copy `~/.xelvra/` to a secure location
2. **Keep software updated**: Regular updates include security improvements
3. **Verify connections**: Only connect to trusted peers
4. **Monitor logs**: Check `peerchat.log` for unusual activity

## 🛠️ Troubleshooting

### Common Issues

#### "No peers found during discovery"
**Solutions:**
1. Check firewall settings (allow UDP 42424)
2. Ensure you're on the same network as other users
3. Try a different network (mobile hotspot)
4. Run `peerchat-cli doctor` for diagnostics

#### "Connection failed"
**Solutions:**
1. Verify the peer is still online
2. Check network connectivity
3. Try discovering again (addresses may change)
4. Ensure compatible versions

#### "Simulation mode detected"
**Solutions:**
1. Check internet connectivity
2. Verify firewall settings
3. Try different network interface
4. Run network diagnostics

### Getting Help
1. **Run diagnostics**: `peerchat-cli doctor`
2. **Check logs**: `~/.xelvra/peerchat.log`
3. **Debug mode**: `peerchat-cli listen` (shows all network activity)
4. **Community support**: [GitHub Discussions](https://github.com/Xelvra/peerchat/discussions)

## 📚 Next Steps

Now that you're up and running:

1. **[Read the User Manual](User-Manual)** - Learn about all features
2. **[Explore CLI Commands](CLI-Usage)** - Master the command-line interface
3. **[Understand P2P Networking](P2P-Networking)** - Learn how it works
4. **[Join the Community](https://github.com/Xelvra/peerchat/discussions)** - Connect with other users

## 🎉 Welcome to the Decentralized Future!

You're now part of a growing community of users who value privacy, security, and digital freedom. Welcome to **#XelvraFree**! 🌐

---

**Need more help?** Check out our [FAQ](FAQ) or [Troubleshooting Guide](Troubleshooting).
