# Xelvra P2P Messenger v0.3.0-alpha - Installation Guide

## 🚀 Quick Start

### Prerequisites
- Go 1.19 or later
- Git
- Network connectivity
- Linux, macOS, or Windows

### Installation Steps

1. **Clone the repository:**
```bash
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
```

2. **Build the application:**
```bash
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go
```

3. **Initialize your identity:**
```bash
./bin/peerchat-cli init
```

4. **Start the messenger:**
```bash
./bin/peerchat-cli start
```

## 🆕 New Features in v0.3.0-alpha

### Proof-of-Work Identity
Your identity now includes computational proof-of-work for Sybil resistance:
- Automatic PoW computation during `init`
- Configurable difficulty (default: 4 leading zeros)
- Enhanced security against fake identities

### Hierarchical Reputation System
Five trust levels with automatic progression:
- **Ghost**: New users (5 messages/day, 1 msg/minute)
- **User**: Verified users (100 messages/day, 1 msg/5sec)
- **Architect**: Contributors (500 messages/day, 1 msg/sec)
- **Ambassador**: Community leaders (1000 messages/day, 500ms)
- **God**: Core developers (unlimited)

### Energy Optimization
Battery-aware operations for mobile devices:
- Adaptive polling based on battery level
- Deep sleep mode at <15% battery
- Resource monitoring and optimization
- Performance targets: <20MB RAM, <1% CPU

### Enhanced Discovery
Hierarchical peer discovery system:
- Local-first: mDNS and UDP broadcast
- Global fallback: Kademlia DHT
- LRU caching for performance
- Smart local/remote detection

## 📋 CLI Commands

### Basic Commands
```bash
# Initialize identity with PoW
./bin/peerchat-cli init

# Start interactive chat
./bin/peerchat-cli start

# Check status (now includes energy info)
./bin/peerchat-cli status

# Discover peers
./bin/peerchat-cli discover

# Show your identity and reputation
./bin/peerchat-cli id

# Network diagnostics
./bin/peerchat-cli doctor

# View manual
./bin/peerchat-cli manual
```

### Interactive Chat Commands
```bash
/help          # Show available commands
/peers         # List connected peers
/discover      # Discover new peers
/connect <id>  # Connect to a peer
/status        # Show node status
/reputation    # Show reputation info
/energy        # Show energy profile
/clear         # Clear screen
/quit          # Exit chat
```

## 🔧 Configuration

### Energy Optimization
The system automatically optimizes based on conditions:
- **Full battery (>50%)**: Normal operation
- **Medium battery (20-50%)**: Reduced polling
- **Low battery (<20%)**: Conservative mode
- **Critical battery (<15%)**: Deep sleep mode

### Reputation System
Reputation grows through network participation:
- Sending messages: +1 point
- File sharing: +5 points
- Online time: +2 points/hour
- Peer verification: +50 points

### Discovery Settings
Hierarchical discovery prioritizes:
1. Local mDNS discovery (immediate)
2. UDP broadcast (local network)
3. DHT discovery (global network)
4. Bootstrap peers (fallback)

## 🛠️ Advanced Configuration

### Custom PoW Difficulty
For testing or high-security networks:
```bash
# Higher difficulty (more secure, slower)
export XELVRA_POW_DIFFICULTY=6

# Lower difficulty (faster, less secure)
export XELVRA_POW_DIFFICULTY=2
```

### Energy Thresholds
Customize energy optimization:
```bash
# Deep sleep threshold (default: 15%)
export XELVRA_DEEP_SLEEP_THRESHOLD=0.10

# DHT poll interval (default: 2 minutes)
export XELVRA_DHT_INTERVAL=300s
```

### Network Settings
Configure discovery behavior:
```bash
# Local cache size (default: 100 peers)
export XELVRA_CACHE_SIZE=200

# Bootstrap timeout (default: 10 seconds)
export XELVRA_BOOTSTRAP_TIMEOUT=15s
```

## 🔍 Troubleshooting

### Common Issues

**PoW computation takes too long:**
- Reduce difficulty: `export XELVRA_POW_DIFFICULTY=2`
- Check CPU usage during init
- Consider hardware limitations

**High energy consumption:**
- Check battery level reporting
- Verify deep sleep activation
- Monitor resource usage with `/energy`

**Peer discovery problems:**
- Check firewall settings
- Verify network connectivity
- Try manual peer connection
- Check NAT traversal status

**Reputation not advancing:**
- Ensure network activity
- Check message delivery success
- Verify peer connections
- Review trust level requirements

### Diagnostic Commands
```bash
# Comprehensive system check
./bin/peerchat-cli doctor

# Network status with energy info
./bin/peerchat-cli status

# View detailed logs
tail -f ~/.xelvra/peerchat.log
```

## 📊 Performance Monitoring

### Resource Targets
- **Memory**: <20MB idle usage
- **CPU**: <1% idle usage
- **Latency**: <50ms message delivery
- **Energy**: <15mW mobile consumption

### Monitoring Tools
```bash
# Check resource usage
./bin/peerchat-cli status

# Energy profile in interactive mode
/energy

# Network quality assessment
/status
```

## 🔗 Links

- **GitHub**: https://github.com/Xelvra/peerchat
- **Issues**: https://github.com/Xelvra/peerchat/issues
- **Documentation**: https://github.com/Xelvra/peerchat/tree/main/docs
- **License**: AGPLv3

## 💬 Support

- **GitHub Issues**: Report bugs and feature requests
- **GitHub Discussions**: Community Q&A and ideas
- **Documentation**: Comprehensive guides and API reference

---

**Note**: This is an alpha release. Features may change and bugs are expected. Please report issues on GitHub.
