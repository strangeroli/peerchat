# 🚀 Xelvra P2P Messenger CLI v0.1.0-alpha

**First Alpha Release - Ready for Community Testing**

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

## 🎯 What's Included

### ✅ Core P2P Infrastructure
- libp2p integration with QUIC/TCP transports
- P2P node initialization and graceful shutdown
- Real-time P2P networking (tested between multiple instances)

### ✅ Network Discovery & Connectivity
- mDNS peer discovery (active and functional)
- UDP broadcast discovery for local networks
- STUN integration for NAT traversal
- Public IP detection and NAT type identification

### ✅ Transport Layer
- QUIC transport as primary protocol (UDP/QUIC-v1)
- TCP fallback for reliability
- UDP buffer optimization

### ✅ Messaging & File Transfer
- P2P message sending and receiving
- Interactive chat mode with command history
- Secure P2P file transfer with chunking
- Progress tracking and resume capability

### ✅ CLI Application
Complete CLI with 12 commands:
- `init` - Identity generation and configuration
- `start` - Interactive P2P chat mode
- `listen` - Passive message monitoring
- `send` - P2P message transmission
- `send-file` - Secure file transfer
- `connect` - Peer connection management
- `discover` - Network peer discovery
- `status` - Real-time node status and diagnostics
- `doctor` - Comprehensive network diagnostics
- `manual` - Complete usage documentation
- `id` - Identity information display
- `profile` - Peer information lookup

### ✅ Security & Identity
- MessengerID generation (DID format preparation)
- Cryptographic identity management
- Secure configuration directory creation

### ✅ AI-Driven Features
- Machine learning optimization for transport selection
- Intelligent peer discovery and connection management
- Adaptive network prediction algorithms

### ✅ Logging & Diagnostics
- Comprehensive logging system with file rotation
- Log rotation (5MB/10k lines, 3 backup files)
- Network diagnostics with detailed reporting
- Real-time status monitoring

## 📊 Performance Targets (Achieved)

- **Memory usage**: <20MB idle ✅
- **CPU usage**: <1% idle ✅
- **Message latency**: <50ms direct connections ✅
- **File transfer**: Chunked, resumable, secure ✅
- **Network discovery**: Multi-method, real-time ✅

## 🔧 System Requirements

- **OS**: Linux (Ubuntu 20.04+, Debian 10+, Fedora 32+, Arch Linux)
- **Arch**: x86_64 (64-bit)
- **RAM**: 512MB minimum (1GB+ recommended)
- **Network**: Internet connection (optional for local mesh)

## 🚀 Installation

See `INSTALL.md` for complete installation instructions.

**Quick install**:
```bash
# Download and install
wget https://github.com/Xelvra/peerchat/releases/download/v0.1.0-alpha/peerchat-cli
chmod +x peerchat-cli
./peerchat-cli init
```

## 🧪 Testing Status

- ✅ **Unit tests**: All passing
- ✅ **File transfer**: Tested with chunking and resume
- ✅ **P2P communication**: Live tested between multiple instances
- ✅ **Network discovery**: Functional across different networks
- ✅ **CLI commands**: All 12 commands working
- ✅ **AI-driven routing**: Transport optimization active

## 🔮 What's Next

### Epoch 2 - API Service (Planned)
- gRPC API server for GUI integration
- Event-driven architecture with streaming
- Database layer with SQLite WAL mode

### Epoch 3 - GUI Application (Planned)
- Cross-platform Flutter application
- Mobile-optimized user interface
- Energy-efficient design

## 🐛 Known Issues

- Interactive chat UI needs refinement
- Advanced encryption features in development
- Mesh networking (Bluetooth LE/Wi-Fi Direct) planned

## 📝 Release Notes

**Date**: 2025-06-15  
**Version**: 0.1.0-alpha  
**License**: AGPLv3  
**Build**: Linux x86_64  

This is the first public alpha release of Xelvra P2P Messenger. The CLI provides a solid foundation for P2P communication with real networking capabilities, comprehensive diagnostics, and professional documentation.

**Status**: Ready for community testing and feedback! 🌟

## 🔗 Links

- **GitHub Repository**: https://github.com/Xelvra/peerchat
- **Documentation**: https://github.com/Xelvra/peerchat/tree/main/docs
- **Issues**: https://github.com/Xelvra/peerchat/issues
- **Releases**: https://github.com/Xelvra/peerchat/releases

## 📄 License

This software is licensed under the GNU Affero General Public License v3.0 (AGPLv3).

---

**Experience true P2P communication. Download, test, and contribute!** 🚀
