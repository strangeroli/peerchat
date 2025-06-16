# 🚀 Xelvra P2P Messenger CLI v0.2.0-alpha

**Second Alpha Release - Enhanced Interactive Experience**

## 📦 Release Contents

- `peerchat-cli` - Linux x86_64 binary (40MB)
- `INSTALL.md` - Complete installation and setup guide
- `README.md` - This file
- `SHA256SUMS` - Checksums for verification

## ⚡ Quick Start

```bash
# Make executable
chmod +x peerchat-cli

# Initialize your identity
./peerchat-cli init

# Test network connectivity
./peerchat-cli doctor

# Start interactive chat
./peerchat-cli start
```

## 🆕 What's New in v0.2.0-alpha

### ✨ Enhanced Interactive Chat Experience
- **Tab Completion**: Auto-complete commands and peer IDs
- **Command History**: Navigate with ↑/↓ arrow keys
- **Keyboard Shortcuts**: Full readline support with Ctrl+C, Ctrl+L, Ctrl+A, Ctrl+E
- **Screen Clearing**: New `/clear` command for clean interface
- **Peer Management**: New `/disconnect <peer_id>` command
- **History Search**: Ctrl+R for command history search

### 🔧 Improved Code Quality & Reliability
- **Fixed Compilation Issues**: Resolved all GitHub Actions errors
- **Code Formatting**: Full gofmt compliance for professional standards
- **Enhanced Error Handling**: Better error messages and graceful failures
- **Linting Integration**: Comprehensive code quality checks
- **Memory Safety**: Improved resource management and cleanup

### 📚 Comprehensive Documentation Updates
- **Enhanced CLI Manual**: Complete interactive features documentation
- **Keyboard Shortcuts Guide**: Detailed shortcuts and commands reference
- **Updated Help System**: Context-aware help in interactive mode
- **Professional Documentation**: GitHub-ready documentation standards

### 🌐 Network Diagnostics Improvements
- **Real NAT Detection**: Shows actual NAT type (port_restricted, etc.)
- **Public IP Display**: Visible external IP addresses
- **Enhanced Discovery**: Better peer discovery with multiple methods
- **Network Quality Indicators**: Real-time connection quality assessment

## 🎯 Complete Feature Set

### ✅ Core P2P Infrastructure
- libp2p integration with QUIC/TCP transports
- P2P node initialization and graceful shutdown
- Real-time P2P networking (tested between multiple instances)
- Multi-instance support on same machine with different ports

### ✅ Network Discovery & Connectivity
- mDNS peer discovery (active and functional)
- UDP broadcast discovery for local networks
- STUN integration for NAT traversal
- Public IP detection and NAT type identification
- Real-time network diagnostics with detailed reporting

### ✅ Transport Layer
- QUIC transport as primary protocol (UDP/QUIC-v1)
- TCP fallback for reliability
- UDP buffer optimization
- Automatic transport selection based on network conditions

### ✅ Interactive Chat System
- **Advanced Input Handling**: Full readline integration
- **Tab Completion**: Commands and peer IDs
- **Command History**: Persistent history with search
- **Keyboard Shortcuts**: Professional terminal experience
- **Real-time Messaging**: Live peer-to-peer communication
- **Screen Management**: Clear screen and interface control

### ✅ CLI Application
Complete CLI with 12 commands + interactive features:
- `init` - Identity generation and configuration
- `start` - Enhanced interactive P2P chat mode
- `listen` - Passive message monitoring
- `send` - P2P message transmission
- `send-file` - Secure file transfer
- `connect` - Peer connection management
- `discover` - Network peer discovery with diagnostics
- `status` - Real-time node status and diagnostics
- `doctor` - Comprehensive network diagnostics
- `manual` - Complete usage documentation (updated)
- `id` - Identity information display
- `profile` - Peer information lookup

### ✅ Interactive Commands (in chat mode)
- `/help` - Show available commands
- `/peers` - List connected peers
- `/discover` - Discover peers in network
- `/connect <id>` - Connect to specific peer
- `/disconnect <id>` - Disconnect from peer (NEW)
- `/status` - Show node status
- `/clear` - Clear screen (NEW)
- `/quit`, `/exit` - Exit chat mode

### ✅ Security & Identity
- MessengerID generation (DID format preparation)
- Cryptographic identity management
- Secure configuration directory creation
- Enhanced error handling for security operations

### ✅ AI-Driven Features
- Machine learning optimization for transport selection
- Intelligent peer discovery and connection management
- Adaptive network prediction algorithms
- Smart routing based on network conditions

### ✅ Logging & Diagnostics
- Comprehensive logging system with file rotation
- Log rotation (5MB/10k lines, 3 backup files)
- Network diagnostics with detailed reporting
- Real-time status monitoring
- Enhanced error logging and debugging

## 📊 Performance Targets (Achieved)

- **Memory usage**: <20MB idle ✅
- **CPU usage**: <1% idle ✅
- **Message latency**: <50ms direct connections ✅
- **File transfer**: Chunked, resumable, secure ✅
- **Network discovery**: Multi-method, real-time ✅
- **Interactive responsiveness**: <10ms input handling ✅

## 🔧 System Requirements

- **OS**: Linux (Ubuntu 20.04+, Debian 10+, Fedora 32+, Arch Linux)
- **Arch**: x86_64 (64-bit)
- **RAM**: 512MB minimum (1GB+ recommended)
- **Network**: Internet connection (optional for local mesh)
- **Terminal**: Modern terminal with readline support recommended

## 🚀 Installation

See `INSTALL.md` for complete installation instructions.

**Quick install**:
```bash
# Download and install
wget https://github.com/Xelvra/peerchat/releases/download/v0.2.0-alpha/peerchat-cli
chmod +x peerchat-cli
./peerchat-cli init
```

## 🧪 Testing Status

- ✅ **Unit tests**: All passing
- ✅ **File transfer**: Tested with chunking and resume
- ✅ **P2P communication**: Live tested between multiple instances
- ✅ **Network discovery**: Functional across different networks
- ✅ **CLI commands**: All 12 commands working
- ✅ **Interactive features**: Tab completion, history, shortcuts tested
- ✅ **Code quality**: gofmt, go vet, compilation checks passing
- ✅ **NAT traversal**: Tested with port_restricted NAT
- ✅ **Multi-instance**: Multiple nodes on same machine verified

## 🔮 What's Next

### Epoch 2 - API Service (In Progress)
- gRPC API server for GUI integration
- Event-driven architecture with streaming
- Database layer with SQLite WAL mode

### Epoch 3 - GUI Application (Planned)
- Cross-platform Flutter application
- Mobile-optimized user interface
- Energy-efficient design

## 🐛 Known Issues

- DHT discovery inactive in local testing (normal behavior)
- Advanced encryption features in development
- Mesh networking (Bluetooth LE/Wi-Fi Direct) planned

## 📝 Release Notes

**Date**: 2025-06-16  
**Version**: 0.2.0-alpha  
**License**: AGPLv3  
**Build**: Linux x86_64  

This second alpha release focuses on user experience improvements, code quality, and interactive features. The CLI now provides a professional-grade interactive experience with comprehensive documentation and enhanced reliability.

**Major improvements**:
- Enhanced interactive chat with tab completion and command history
- Fixed all compilation and formatting issues
- Comprehensive documentation updates
- Improved network diagnostics and NAT detection
- Professional code quality standards

**Status**: Ready for extended community testing! 🌟

## 🔗 Links

- **GitHub Repository**: https://github.com/Xelvra/peerchat
- **Documentation**: https://github.com/Xelvra/peerchat/tree/main/docs
- **Issues**: https://github.com/Xelvra/peerchat/issues
- **Releases**: https://github.com/Xelvra/peerchat/releases

## 📄 License

This software is licensed under the GNU Affero General Public License v3.0 (AGPLv3).

---

**Experience enhanced P2P communication with professional interactive features!** 🚀
