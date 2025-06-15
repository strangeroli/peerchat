# Xelvra P2P Messenger CLI Usage Guide

## Overview

The Xelvra CLI (`peerchat-cli`) is a command-line interface for the Xelvra P2P messaging network. It provides a secure, decentralized way to communicate without relying on centralized servers.

## Installation

### From Source

```bash
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
./scripts/build.sh
```

The binary will be available in the `bin/` directory.

### Pre-built Binaries

Download the latest release from the GitHub releases page and extract the binary for your platform.

## Quick Start

### 1. Initialize Your Identity

Before using Xelvra, you need to create a decentralized identity:

```bash
./bin/peerchat-cli init
```

This command:
- Generates a new Ed25519 key pair
- Creates a DID (Decentralized Identifier) in the format `did:xelvra:<hash>`
- Sets up the local database
- Creates a configuration file at `~/.xelvra/config.yaml`

### 2. Start the P2P Node

```bash
./bin/peerchat-cli start
```

This command starts the P2P node and begins listening for connections. The node will:
- Connect to the P2P network
- Listen for incoming messages
- Maintain peer connections
- Provide real-time status updates

### 3. Check Node Status

```bash
./bin/peerchat-cli status
```

Shows current node information including:
- Node identity (DID and Peer ID)
- Connected peers
- Network addresses
- Performance metrics

## Commands Reference

### `init`

Initialize a new Xelvra identity and configuration.

```bash
peerchat-cli init [flags]
```

**Options:**
- `--config string`: Custom config file path
- `-v, --verbose`: Enable verbose output

**Example:**
```bash
peerchat-cli init --verbose
```

### `start`

Start the P2P node and begin networking.

```bash
peerchat-cli start [flags]
```

**Options:**
- `--config string`: Custom config file path
- `-v, --verbose`: Enable verbose output

**Example:**
```bash
peerchat-cli start --verbose
```

### `status`

Display current node status and statistics.

```bash
peerchat-cli status [flags]
```

**Options:**
- `--config string`: Custom config file path

**Example:**
```bash
peerchat-cli status
```

### `version`

Show version information.

```bash
peerchat-cli version
```

**Example Output:**
```
Xelvra P2P Messenger CLI v0.1.0-alpha
Built with Go and libp2p
https://github.com/Xelvra/peerchat
```

### `help`

Show help information for any command.

```bash
peerchat-cli help [command]
```

**Examples:**
```bash
peerchat-cli help
peerchat-cli help init
peerchat-cli help start
```

## Configuration

The configuration file is located at `~/.xelvra/config.yaml` by default. You can specify a custom location using the `--config` flag.

### Configuration Structure

```yaml
# Identity configuration
identity:
  did: "did:xelvra:5B5CDn5SvTvYHnyuShbAoLGxRzrcGQthUNYHz61TjCei"
  peer_id: "12D3KooWEhwTXBCkpm61HyS25wjiE4zwf5s6Bwq7efxddqZXkAMd"

# Network configuration
network:
  listen_addrs:
    - "/ip4/0.0.0.0/tcp/0"
    - "/ip4/0.0.0.0/udp/0/quic-v1"
  bootstrap_peers: []
  enable_quic: true
  enable_tcp: true

# Logging configuration
logging:
  level: "info"
  format: "json"

# Database configuration
database:
  path: "/home/user/.xelvra/userdata.db"
```

## Performance Targets

Xelvra is designed with aggressive performance targets:

- **Latency**: <50ms for direct P2P connections
- **Memory Usage**: <20MB when idle
- **CPU Usage**: <1% when idle

## Security Features

- **End-to-End Encryption**: All messages are encrypted using Signal Protocol
- **Decentralized Identity**: No central authority controls your identity
- **Memory Protection**: Sensitive data is protected in memory (planned)
- **Forward Secrecy**: Messages cannot be decrypted even if keys are compromised

## Trust System

Xelvra implements a 5-level trust system:

1. **Ghost** (Level 0): New users with limited privileges
2. **User** (Level 1): Basic verified users
3. **Architect** (Level 2): Contributors to the network
4. **Ambassador** (Level 3): Community leaders
5. **God** (Level 4): Core developers/administrators

## Troubleshooting

### Common Issues

**Node won't start:**
- Check if ports are available
- Verify configuration file syntax
- Check firewall settings

**Can't connect to peers:**
- Ensure internet connectivity
- Check NAT/firewall configuration
- Verify bootstrap peers are reachable

**High memory usage:**
- Check for memory leaks in logs
- Restart the node periodically
- Monitor peer connections

### Debug Mode

Enable verbose logging for troubleshooting:

```bash
peerchat-cli start --verbose
```

### Log Files

Logs are written to:
- Console output (when using `--verbose`)
- System journal (on systemd systems)

## Advanced Usage

### Custom Bootstrap Peers

Add custom bootstrap peers to your configuration:

```yaml
network:
  bootstrap_peers:
    - "/ip4/192.168.1.100/tcp/4001/p2p/12D3KooW..."
    - "/dns4/bootstrap.xelvra.com/tcp/4001/p2p/12D3KooW..."
```

### Network Interfaces

Specify custom listen addresses:

```yaml
network:
  listen_addrs:
    - "/ip4/0.0.0.0/tcp/4001"
    - "/ip6/::/tcp/4001"
```

## Support

For support and questions:
- GitHub Issues: https://github.com/Xelvra/peerchat/issues
- Documentation: https://github.com/Xelvra/peerchat/docs
- Community: [Coming Soon]
