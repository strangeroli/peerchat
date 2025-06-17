# Messenger Xelvra: *#XelvraFree*

[![GitHub Issues](https://img.shields.io/github/issues/Xelvra/peerchat)](https://github.com/Xelvra/peerchat/issues)
[![GitHub Wiki](https://img.shields.io/badge/GitHub-Wiki-blue)](https://github.com/Xelvra/peerchat/wiki)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)

> 🚀 **Secure, decentralized P2P messenger. Built on E2E encryption with AI-driven net prediction.**

**Messenger Xelvra** is a peer-to-peer (P2P) communication platform designed to restore privacy, security, and user control over digital communication. The project aims to create a secure, efficient, and decentralized platform that pushes the boundaries of P2P communication capabilities.

## 📋 Table of Contents

- [Why Xelvra?](#-why-xelvra)
- [Key Features](#-key-features)
- [Quick Start](#-quick-start)
- [Usage](#-usage)
- [Development Status](#-development-status)
- [Roadmap](#-roadmap)
- [Documentation](#-documentation)
- [Contributing](#-contributing)
- [License](#-license)

## 🌟 Why Xelvra?

In an era where digital privacy is under constant threat, Xelvra offers a fundamentally different approach to messaging.

### The Problem with Traditional Messengers

- **🔍 Data Collection**: Your conversations become products for advertising
- **🏢 Central Control**: Single points of failure and censorship
- **🚫 Limited Freedom**: Platform policies can restrict your communication
- **💰 Privacy as Currency**: Your personal data is monetized without your consent

### The Xelvra Solution

- **🔒 True Privacy**: End-to-end encryption with no data collection
- **🌐 Decentralized**: Direct peer-to-peer communication
- **🛡️ Censorship Resistant**: No central authority can block you
- **📱 Your Data**: Everything stays on your devices
- **🔓 Open Source**: Transparent, auditable code

## 🚀 Key Features

### 🔐 Security First
- **End-to-End Encryption**: Signal Protocol with X3DH handshake and Double Ratchet
- **Proof-of-Work Identity**: Sybil-resistant identity creation
- **Hierarchical Trust**: Ghost → User → Architect → Ambassador → God reputation system
- **Forward Secrecy**: Automatic key rotation protects past conversations
- **Metadata Protection**: Onion routing hides communication patterns

### 🌐 True Decentralization
- **6-Phase Discovery**: IPv6 → mDNS → UDP → DHT → Hole Punching → Relay
- **Smart NAT Traversal**: Multiple strategies for restrictive networks
- **Local-First**: Prioritizes direct connections, falls back gracefully
- **Intelligent Caching**: LRU cache for optimal peer management

### ⚡ High Performance
- **QUIC Transport**: Ultra-low latency with TCP fallback
- **Resource Efficient**: <20MB memory, <1% CPU idle
- **Energy Optimized**: Battery-aware operations, deep sleep mode
- **AI-Driven**: Machine learning for optimal routing

### 🛠️ Developer Ready
- **Modular Design**: Clean CLI, API, and GUI separation
- **gRPC API**: Modern, efficient inter-component communication
- **Cross-Platform**: Linux, macOS, Windows, Android, iOS
- **Comprehensive Testing**: Unit, integration, and chaos engineering

## 🚀 Quick Start

```bash
# Clone and build
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go

# Initialize and start
./bin/peerchat-cli init
./bin/peerchat-cli start
```

For detailed instructions, see our [Installation Guide](docs/INSTALLATION.md).

## 📱 Usage

### Basic Commands

```bash
# Initialize your identity
./bin/peerchat-cli init

# Start interactive chat
./bin/peerchat-cli start

# Check network status
./bin/peerchat-cli status

# Discover peers
./bin/peerchat-cli discover

# Network diagnostics
./bin/peerchat-cli doctor
```

### Interactive Chat Commands

```bash
/help          # Show available commands
/peers         # List connected peers
/discover      # Discover new peers
/connect <id>  # Connect to a peer
/status        # Show node status
/quit          # Exit chat
```

For detailed usage, see our [User Guide](docs/USER_GUIDE.md).

## 🏗️ Development Status

### ✅ Epoch 1: CLI Foundation (COMPLETED)
- **P2P Core**: libp2p integration with QUIC/TCP transports
- **6-Phase Discovery**: IPv6 → mDNS → UDP → DHT → Hole punching → Relay
- **Security**: Signal Protocol with E2E encryption and memory protection
- **Advanced Features**: Proof-of-Work DID, hierarchical reputation, energy optimization
- **CLI Interface**: Complete command set with interactive chat
- **System Integration**: Daemon mode, systemd support, comprehensive logging
- **Event-Driven Architecture**: v0.4.0-alpha introduces centralized event bus system
- **Advanced Logging**: Structured logging with rotation and component-specific levels

### 🔌 Epoch 2: API Service (PLANNED)
- **gRPC Server**: High-performance API with event-driven architecture
- **Database Layer**: SQLite with WAL mode for persistent storage
- **Monitoring**: Prometheus metrics and OpenTelemetry tracing
- **Stream Processing**: Real-time message and event streaming

### 📱 Epoch 3: GUI Application (PLANNED)
- **Cross-Platform**: Flutter app for Android, iOS, Linux, macOS, Windows
- **Modern UI/UX**: Material Design with accessibility compliance
- **Energy Optimization**: <100mW active usage, intelligent sleep modes
- **Advanced Features**: Group chats, file sharing, voice calls

### 🚀 Epoch 4: Advanced Features (FUTURE)
- **Zero-Knowledge Proofs**: Enhanced privacy with ZKP identity verification
- **Quantum Resistance**: Post-quantum cryptography integration
- **Voice & Video**: Real-time multimedia communication
- **Mesh Networks**: Advanced offline communication capabilities

## 🗺️ Roadmap

### 🎯 Short Term (3-6 months)
- Enhanced security features and automatic key rotation
- Mesh networking with Bluetooth LE and Wi-Fi Direct
- Performance optimization to meet target metrics

### 🚀 Medium Term (6-12 months)
- Complete gRPC API implementation
- Begin Flutter GUI development
- Advanced NAT traversal with AI-driven prediction

### 🌟 Long Term (1-2 years)
- Cross-platform GUI with full feature parity
- Voice & video communication
- Quantum-resistant cryptography
- Ecosystem expansion with governance features

## 📚 Documentation

### 📖 Quick Access
| Document | Description |
|----------|-------------|
| [📖 User Guide](docs/USER_GUIDE.md) | Complete guide for end users |
| [🔧 Installation Guide](docs/INSTALLATION.md) | Platform-specific installation instructions |
| [👨‍💻 Developer Guide](docs/DEVELOPER_GUIDE.md) | Development setup and contribution guide |
| [📋 API Reference](docs/API_REFERENCE.md) | Complete API documentation |

### 🌐 GitHub Resources
- **[📖 Wiki](https://github.com/Xelvra/peerchat/wiki)** - Comprehensive documentation and tutorials
- **[🐛 Issues](https://github.com/Xelvra/peerchat/issues)** - Bug reports and feature requests
- **[💬 Discussions](https://github.com/Xelvra/peerchat/discussions)** - Community Q&A and ideas
- **[🔧 Releases](https://github.com/Xelvra/peerchat/releases)** - Download latest versions

## 🤝 Contributing

We welcome contributions from developers, security researchers, and privacy advocates!

### How to Get Started
1. **Fork the repository** and clone it locally
2. **Read our [Developer Guide](docs/DEVELOPER_GUIDE.md)** for setup instructions
3. **Check [open issues](https://github.com/Xelvra/peerchat/issues)** for tasks to work on
4. **Join [discussions](https://github.com/Xelvra/peerchat/discussions)** to connect with the community

### Ways to Contribute
- 🐛 **Bug Reports**: Help us identify and fix issues
- 💡 **Feature Requests**: Suggest new functionality
- 🔧 **Code Contributions**: Submit pull requests for improvements
- 📚 **Documentation**: Improve guides and tutorials
- 🔍 **Security Audits**: Help us maintain security standards
- 🌍 **Translations**: Make Xelvra accessible worldwide

### Development Environment
```bash
# Quick setup for contributors
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
go mod download
./scripts/build.sh
./scripts/test.sh
```

## 📄 License

Messenger Xelvra is licensed under the **GNU Affero General Public License v3.0 (AGPLv3)**.

This ensures that:
- ✅ **Free to use** for personal and commercial purposes
- ✅ **Open source** - all code is transparent and auditable
- ✅ **Copyleft protection** - modifications must remain open source
- ✅ **Network copyleft** - even SaaS deployments must share source code

See the [LICENSE](LICENSE) file for full details.

---

## 🚀 Ready to Experience True Digital Freedom?

**Download Xelvra today and join the decentralized communication revolution!**

- 📥 **[Latest Release](https://github.com/Xelvra/peerchat/releases/latest)** - Download the newest version
- 📖 **[Quick Start Guide](https://github.com/Xelvra/peerchat/wiki/Getting-Started)** - Get up and running in minutes
- 💬 **[Join the Community](https://github.com/Xelvra/peerchat/discussions)** - Connect with other users and developers

*Your conversations, your control. Experience the future of private communication with Messenger Xelvra.*


