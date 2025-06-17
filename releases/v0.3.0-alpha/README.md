# Xelvra P2P Messenger v0.3.0-alpha

**Release Date:** December 17, 2024  
**Version:** v0.3.0-alpha  
**Build:** Linux x86_64

## 📦 Release Contents

```
peerchat-cli          # Main executable binary (Linux x86_64)
SHA256SUMS           # SHA256 checksums for verification
INSTALL.md           # Installation and setup guide
RELEASE_NOTES.md     # Detailed release notes
STATUS.md            # Project status and progress
README.md            # This file
```

## 🔐 Security Verification

**SHA256 Checksum:**
```bash
sha256sum -c SHA256SUMS
```

**Expected Output:**
```
peerchat-cli: OK
```

## 🚀 Quick Start

### 1. Download and Verify
```bash
# Download the release
wget https://github.com/Xelvra/peerchat/releases/download/v0.3.0-alpha/peerchat-cli

# Verify integrity
sha256sum -c SHA256SUMS

# Make executable
chmod +x peerchat-cli
```

### 2. Initialize Your Identity
```bash
# Create your Proof-of-Work protected identity
./peerchat-cli init
```

### 3. Start Messaging
```bash
# Start interactive chat
./peerchat-cli start

# Or listen for messages
./peerchat-cli listen
```

## 🆕 What's New in v0.3.0-alpha

### 🔐 Proof-of-Work Identity System
- **Sybil-Resistant DID**: Computational proof required for identity creation
- **Configurable Difficulty**: Adjustable security level (default: 4 leading zeros)
- **Automatic Validation**: All identities verified on network entry

### 🏆 Hierarchical Reputation System
- **5 Trust Levels**: Ghost → User → Architect → Ambassador → God
- **Merit-Based Progression**: Automatic advancement through network contribution
- **Smart Rate Limiting**: Message limits based on trust level
- **Peer Verification**: Community-driven trust building

### 🌐 Advanced Peer Discovery
- **Hierarchical Protocol**: Local-first (IPv6, mDNS) → Global (DHT) → Relay fallback
- **LRU Caching**: Intelligent local peer caching (100 peers)
- **Smart Routing**: Automatic local vs. remote peer detection
- **NAT Traversal**: Hole punching with relay server creation

### ⚡ Energy Optimization
- **Adaptive Polling**: Battery-aware network operations
- **Deep Sleep Mode**: Ultra-low power at <15% battery
- **Resource Monitoring**: Real-time CPU/memory tracking
- **Performance Targets**: <20MB RAM, <1% CPU, <50ms latency

## 📋 System Requirements

### Minimum Requirements
- **OS**: Linux x86_64 (Ubuntu 18.04+, CentOS 7+, Debian 9+)
- **RAM**: 32MB available memory
- **CPU**: Any x86_64 processor
- **Network**: Internet connection for P2P operations
- **Storage**: 10MB free disk space

### Recommended Requirements
- **OS**: Linux x86_64 (Ubuntu 20.04+, CentOS 8+, Debian 11+)
- **RAM**: 64MB available memory
- **CPU**: Multi-core x86_64 processor
- **Network**: Broadband internet connection
- **Storage**: 50MB free disk space

## 🔧 Configuration

### Environment Variables
```bash
# Proof-of-Work difficulty (1-32, default: 4)
export XELVRA_POW_DIFFICULTY=4

# Deep sleep battery threshold (0.0-1.0, default: 0.15)
export XELVRA_DEEP_SLEEP_THRESHOLD=0.15

# Local peer cache size (default: 100)
export XELVRA_CACHE_SIZE=100

# DHT polling interval (default: 2m)
export XELVRA_DHT_INTERVAL=2m

# Bootstrap timeout (default: 10s)
export XELVRA_BOOTSTRAP_TIMEOUT=10s
```

### Configuration Files
```bash
~/.xelvra/config.yaml     # Main configuration
~/.xelvra/identity.json   # Your identity and keys
~/.xelvra/peers.db        # Known peers database
~/.xelvra/messages.db     # Message history
~/.xelvra/peerchat.log    # Application logs
```

## 🎯 Trust Level Progression

| Level | Name | Daily Messages | Rate Limit | Requirements |
|-------|------|----------------|------------|--------------|
| 0 | Ghost | 5 | 1/minute | New user |
| 1 | User | 100 | 1/5sec | 100 reputation, 24h uptime |
| 2 | Architect | 500 | 1/sec | 1000 reputation, 1 week uptime |
| 3 | Ambassador | 1000 | 500ms | 10000 reputation, 1 month uptime |
| 4 | God | Unlimited | None | 100000 reputation, 3 months uptime |

## 🌐 Network Discovery Protocol

### Discovery Sequence
1. **IPv6 Local**: Direct IPv6 link-local discovery
2. **mDNS**: Multicast DNS local network discovery
3. **UDP Broadcast**: Network broadcast discovery
4. **DHT Global**: Kademlia distributed hash table
5. **Hole Punching**: NAT traversal attempts
6. **Relay Server**: Fallback relay creation

### Connection Priority
1. **Direct P2P**: Fastest, most efficient
2. **NAT Traversal**: Hole-punched connections
3. **Relay Server**: Fallback for restrictive networks

## 📊 Performance Metrics

### Resource Usage (Idle)
- **Memory**: <20MB RSS
- **CPU**: <1% utilization
- **Network**: <1KB/s background traffic
- **Energy**: <15mW on mobile devices

### Network Performance
- **Latency**: <50ms P2P message delivery
- **Discovery**: Local peers <1s, Global peers <30s
- **Throughput**: QUIC primary, TCP fallback
- **Reliability**: 99%+ message delivery rate

## 🔍 Troubleshooting

### Common Issues

**Identity creation takes too long:**
```bash
# Reduce PoW difficulty for testing
export XELVRA_POW_DIFFICULTY=2
./peerchat-cli init
```

**No peers discovered:**
```bash
# Check network connectivity
./peerchat-cli doctor

# Manual peer connection
./peerchat-cli connect <peer-id>
```

**High resource usage:**
```bash
# Check energy profile
./peerchat-cli status

# Enable deep sleep mode
export XELVRA_DEEP_SLEEP_THRESHOLD=0.5
```

### Diagnostic Commands
```bash
./peerchat-cli status    # System status with energy info
./peerchat-cli doctor    # Network diagnostics
./peerchat-cli discover  # Manual peer discovery
./peerchat-cli manual    # Built-in manual
```

## 📞 Support

- **GitHub Issues**: https://github.com/Xelvra/peerchat/issues
- **Documentation**: https://github.com/Xelvra/peerchat/tree/main/docs
- **Discussions**: https://github.com/Xelvra/peerchat/discussions
- **Wiki**: https://github.com/Xelvra/peerchat/wiki

## 📄 License

This software is licensed under the GNU Affero General Public License v3.0 (AGPLv3).
See the LICENSE file for full license text.

## 🔗 Links

- **Source Code**: https://github.com/Xelvra/peerchat
- **Releases**: https://github.com/Xelvra/peerchat/releases
- **Documentation**: https://github.com/Xelvra/peerchat/tree/main/docs
- **Project Website**: https://xelvra.github.io/peerchat

---

**Warning**: This is an alpha release intended for testing and development.
Not recommended for production use.
