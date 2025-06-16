# Xelvra P2P Messenger CLI - Installation Guide

## 🚀 Version 0.2.0-alpha

This is the second alpha release of Xelvra P2P Messenger CLI - featuring enhanced interactive experience, improved code quality, and comprehensive documentation.

## 📋 System Requirements

- **Operating System**: Linux (Ubuntu 20.04+, Debian 10+, Fedora 32+, Arch Linux)
- **Architecture**: x86_64 (64-bit)
- **Memory**: Minimum 512MB RAM (recommended 1GB+)
- **Network**: Internet connection for P2P discovery (optional for local mesh)
- **Terminal**: Modern terminal with readline support recommended
- **Permissions**: User-level permissions (no root required)

## 📦 Installation Methods

### Method 1: Direct Binary Installation (Recommended)

1. **Download the binary**:
   ```bash
   # Download from GitHub releases
   wget https://github.com/Xelvra/peerchat/releases/download/v0.2.0-alpha/peerchat-cli
   
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

3. **Start the enhanced interactive chat**:
   ```bash
   ./peerchat-cli start
   ```
   This starts the interactive chat mode with tab completion and command history.

## 📖 Basic Usage

### Essential Commands

```bash
# Show help
./peerchat-cli --help

# Initialize identity (first time only)
./peerchat-cli init

# Start enhanced interactive chat
./peerchat-cli start

# Check node status with detailed diagnostics
./peerchat-cli status

# Discover peers with network information
./peerchat-cli discover

# Connect to a peer
./peerchat-cli connect <peer_multiaddr>

# Send a message
./peerchat-cli send <peer_multiaddr> "Hello, World!"

# Send a file
./peerchat-cli send-file <peer_multiaddr> /path/to/file

# Listen for messages (passive mode)
./peerchat-cli listen

# Enhanced network diagnostics
./peerchat-cli doctor

# Show your identity
./peerchat-cli id

# View complete manual (updated)
./peerchat-cli manual
```

### Enhanced Interactive Chat Commands (NEW in v0.2.0)

When in `start` mode, you can use these commands with **tab completion**:

```
/help          - Show available commands
/peers         - List connected peers
/discover      - Find nearby peers with diagnostics
/connect <id>  - Connect to a specific peer (tab completion)
/disconnect <id> - Disconnect from peer (NEW)
/status        - Show detailed node status
/clear         - Clear screen (NEW)
/quit, /exit   - Exit chat mode
```

### New Keyboard Shortcuts (v0.2.0)

- **Tab** - Auto-complete commands and peer IDs
- **↑/↓ arrows** - Navigate command history
- **Ctrl+C** - Exit chat mode
- **Ctrl+L** - Clear screen (alternative to /clear)
- **Ctrl+A** - Move to beginning of line
- **Ctrl+E** - Move to end of line
- **Ctrl+R** - Search command history

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
# Run enhanced network diagnostics
./peerchat-cli doctor

# Check if ports are blocked
sudo netstat -tulpn | grep 4001
```

**Interactive features not working**:
- Ensure you have a modern terminal
- Try different terminal emulator if issues persist
- Tab completion requires readline support

**High memory usage**:
- Normal usage: <20MB
- If higher, restart the application

### Log Files

Logs are stored in `~/.xelvra/peerchat.log` with automatic rotation:
```bash
# View recent logs
tail -f ~/.xelvra/peerchat.log

# View all logs
cat ~/.xelvra/peerchat.log

# Check rotated logs
ls ~/.xelvra/peerchat.log.*
```

## 🔐 Security Notes

- Your identity keys are stored in `~/.xelvra/`
- **Backup your identity**: Copy the entire `~/.xelvra/` directory
- Never share your private keys
- The application uses end-to-end encryption by default
- Enhanced error handling for security operations

## 📊 Performance Metrics

**Achieved Performance** (v0.2.0-alpha):
- Memory usage: <20MB (idle)
- CPU usage: <1% (idle)
- Message latency: <50ms (direct connections)
- Interactive responsiveness: <10ms (input handling)
- File transfer: Chunked with resume capability
- AI-driven routing: Active transport optimization

## 🆕 What's New in v0.2.0-alpha

### Enhanced Interactive Experience
- ✅ Tab completion for commands and peer IDs
- ✅ Command history with arrow key navigation
- ✅ Keyboard shortcuts (Ctrl+C, Ctrl+L, Ctrl+A, Ctrl+E, Ctrl+R)
- ✅ Screen clearing with `/clear` command
- ✅ Peer disconnection with `/disconnect` command
- ✅ History search with Ctrl+R

### Improved Code Quality
- ✅ Fixed all compilation and GitHub Actions errors
- ✅ Full gofmt compliance for professional standards
- ✅ Enhanced error handling and graceful failures
- ✅ Comprehensive linting integration
- ✅ Memory safety improvements

### Enhanced Documentation
- ✅ Updated CLI manual with interactive features
- ✅ Comprehensive keyboard shortcuts guide
- ✅ Context-aware help system
- ✅ Professional documentation standards

### Network Improvements
- ✅ Real NAT detection with type display
- ✅ Public IP address visibility
- ✅ Enhanced peer discovery diagnostics
- ✅ Network quality indicators

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

**Ready to experience enhanced P2P communication? Start with `./peerchat-cli init`!** 🚀
