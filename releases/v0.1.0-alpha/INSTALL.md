# Xelvra P2P Messenger CLI - Installation Guide

## 🚀 Version 0.1.0-alpha

This is the first alpha release of Xelvra P2P Messenger CLI - a secure, decentralized peer-to-peer communication platform.

## 📋 System Requirements

- **Operating System**: Linux (Ubuntu 20.04+, Debian 10+, Fedora 32+, Arch Linux)
- **Architecture**: x86_64 (64-bit)
- **Memory**: Minimum 512MB RAM (recommended 1GB+)
- **Network**: Internet connection for P2P discovery (optional for local mesh)
- **Permissions**: User-level permissions (no root required)

## 📦 Installation Methods

### Method 1: Direct Binary Installation (Recommended)

1. **Download the binary**:
   ```bash
   # Download from GitHub releases
   wget https://github.com/Xelvra/peerchat/releases/download/v0.1.0-alpha/peerchat-cli
   
   # Or if you have this file locally
   # Copy peerchat-cli to your desired location
   ```

2. **Make it executable**:
   ```bash
   chmod +x peerchat-cli
   ```

3. **Install to system PATH** (optional):
   ```bash
   # Install for current user
   mkdir -p ~/.local/bin
   mv peerchat-cli ~/.local/bin/
   
   # Add to PATH if not already (add to ~/.bashrc or ~/.zshrc)
   echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
   source ~/.bashrc
   
   # OR install system-wide (requires sudo)
   sudo mv peerchat-cli /usr/local/bin/
   ```

### Method 2: Build from Source

1. **Install Go** (version 1.19+):
   ```bash
   # Ubuntu/Debian
   sudo apt update && sudo apt install golang-go
   
   # Fedora
   sudo dnf install golang
   
   # Arch Linux
   sudo pacman -S go
   ```

2. **Clone and build**:
   ```bash
   git clone https://github.com/Xelvra/peerchat.git
   cd peerchat
   go build -o peerchat-cli cmd/peerchat-cli/main.go
   ```

## 🔧 First Setup

1. **Initialize your identity**:
   ```bash
   ./peerchat-cli init
   ```
   This creates your cryptographic identity and configuration directory at `~/.xelvra/`

2. **Test network connectivity**:
   ```bash
   ./peerchat-cli doctor
   ```
   This diagnoses your network setup and provides recommendations.

3. **Start the P2P node**:
   ```bash
   ./peerchat-cli start
   ```
   This starts the interactive chat mode.

## 📖 Basic Usage

### Essential Commands

```bash
# Show help
./peerchat-cli --help

# Initialize identity (first time only)
./peerchat-cli init

# Start interactive chat
./peerchat-cli start

# Check node status
./peerchat-cli status

# Discover peers
./peerchat-cli discover

# Connect to a peer
./peerchat-cli connect <peer_multiaddr>

# Send a message
./peerchat-cli send <peer_multiaddr> "Hello, World!"

# Send a file
./peerchat-cli send-file <peer_multiaddr> /path/to/file

# Listen for messages (passive mode)
./peerchat-cli listen

# Network diagnostics
./peerchat-cli doctor

# Show your identity
./peerchat-cli id

# View complete manual
./peerchat-cli manual
```

### Interactive Chat Commands

When in `start` mode, you can use these commands:

```
/help       - Show available commands
/peers      - List connected peers
/discover   - Find nearby peers
/connect    - Connect to a specific peer
/status     - Show node status
/quit       - Exit chat mode
```

## 🔥 Firewall Configuration

If you're behind a firewall, you may need to open ports:

```bash
# Ubuntu/Debian (ufw)
sudo ufw allow 4001/tcp
sudo ufw allow 4001/udp

# Fedora/RHEL (firewalld)
sudo firewall-cmd --permanent --add-port=4001/tcp
sudo firewall-cmd --permanent --add-port=4001/udp
sudo firewall-cmd --reload

# Manual iptables
sudo iptables -A INPUT -p tcp --dport 4001 -j ACCEPT
sudo iptables -A INPUT -p udp --dport 4001 -j ACCEPT
```

## 🐛 Troubleshooting

### Common Issues

**"Permission denied" error**:
```bash
chmod +x peerchat-cli
```

**"Command not found"**:
- Ensure the binary is in your PATH
- Use `./peerchat-cli` if running from current directory

**Network connectivity issues**:
```bash
# Run network diagnostics
./peerchat-cli doctor

# Check if ports are blocked
sudo netstat -tulpn | grep 4001
```

**High memory usage**:
- Normal usage: <20MB
- If higher, restart the application

### Log Files

Logs are stored in `~/.xelvra/peerchat.log`:
```bash
# View recent logs
tail -f ~/.xelvra/peerchat.log

# View all logs
cat ~/.xelvra/peerchat.log
```

## 🔐 Security Notes

- Your identity keys are stored in `~/.xelvra/`
- **Backup your identity**: Copy the entire `~/.xelvra/` directory
- Never share your private keys
- The application uses end-to-end encryption by default

## 📊 Performance Metrics

**Achieved Performance** (v0.1.0-alpha):
- Memory usage: <20MB (idle)
- CPU usage: <1% (idle)
- Message latency: <50ms (direct connections)
- File transfer: Chunked with resume capability
- AI-driven routing: Active transport optimization

## 🆕 What's New in v0.1.0-alpha

- ✅ Complete P2P networking with libp2p
- ✅ QUIC/TCP transport protocols
- ✅ mDNS and UDP broadcast discovery
- ✅ STUN integration for NAT traversal
- ✅ Secure file transfer with chunking
- ✅ Interactive CLI with 12 commands
- ✅ Real-time messaging and chat
- ✅ AI-driven routing and transport optimization
- ✅ Comprehensive logging and diagnostics
- ✅ Cross-platform build system

## 🔗 Links

- **GitHub**: https://github.com/Xelvra/peerchat
- **Issues**: https://github.com/Xelvra/peerchat/issues
- **Documentation**: https://github.com/Xelvra/peerchat/tree/main/docs
- **License**: AGPLv3

## 💬 Support

- **GitHub Issues**: Report bugs and feature requests
- **GitHub Discussions**: Community support and questions
- **Documentation**: Complete guides in the `docs/` directory

---

**Ready to experience true P2P communication? Start with `./peerchat-cli init`!** 🚀
