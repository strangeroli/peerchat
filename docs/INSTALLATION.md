# Xelvra P2P Messenger - Installation Guide

> 📖 **For the most comprehensive and up-to-date installation instructions, visit our [GitHub Wiki Installation Guide](https://github.com/Xelvra/peerchat/wiki/Installation)**

## Table of Contents
- [System Requirements](#system-requirements)
- [Quick Start](#quick-start)
- [Installation Methods](#installation-methods)
- [Platform-Specific Instructions](#platform-specific-instructions)
- [Configuration](#configuration)
- [Verification](#verification)
- [Troubleshooting](#troubleshooting)

## System Requirements

### Minimum Requirements
- **Operating System**: Linux, macOS, or Windows
- **Architecture**: x86_64 (amd64) or ARM64
- **Memory**: 512 MB RAM
- **Storage**: 100 MB free space
- **Network**: Internet connection for initial setup

### Recommended Requirements
- **Memory**: 2 GB RAM or more
- **Storage**: 1 GB free space
- **Network**: Stable broadband connection
- **Firewall**: UDP port 42424 open for peer discovery

### Software Dependencies
- **Go**: Version 1.21 or later (for building from source)
- **Git**: For cloning the repository
- **C Compiler**: GCC or Clang (for CGO dependencies)

## Quick Start

### 1. Download and Build
```bash
# Clone the repository
git clone https://github.com/Xelvra/peerchat.git
cd peerchat

# Build the CLI
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go

# Make it executable (Linux/macOS)
chmod +x bin/peerchat-cli
```

### 2. Initialize
```bash
# Initialize your identity
./bin/peerchat-cli init

# Test your setup
./bin/peerchat-cli doctor
```

### 3. Start Chatting
```bash
# Start interactive chat
./bin/peerchat-cli start
```

## Installation Methods

### Method 1: Build from Source (Recommended)

#### Prerequisites
```bash
# Install Go 1.21+
# Linux (Ubuntu/Debian)
sudo apt update
sudo apt install golang-go git build-essential

# macOS (with Homebrew)
brew install go git

# Windows (with Chocolatey)
choco install golang git mingw
```

#### Build Steps
```bash
# Clone repository
git clone https://github.com/Xelvra/peerchat.git
cd peerchat

# Download dependencies
go mod download

# Build CLI
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go

# Optional: Install to system PATH
sudo cp bin/peerchat-cli /usr/local/bin/
```

### Method 2: Pre-built Binaries (Future)

Pre-built binaries will be available for download from GitHub Releases:

```bash
# Linux x86_64
wget https://github.com/Xelvra/peerchat/releases/latest/download/peerchat-cli-linux-amd64
chmod +x peerchat-cli-linux-amd64
sudo mv peerchat-cli-linux-amd64 /usr/local/bin/peerchat-cli

# macOS x86_64
wget https://github.com/Xelvra/peerchat/releases/latest/download/peerchat-cli-darwin-amd64
chmod +x peerchat-cli-darwin-amd64
sudo mv peerchat-cli-darwin-amd64 /usr/local/bin/peerchat-cli

# Windows x86_64
# Download peerchat-cli-windows-amd64.exe from releases page
```

### Method 3: Package Managers (Future)

Package manager support is planned:

```bash
# Homebrew (macOS/Linux)
brew install xelvra/tap/peerchat-cli

# Snap (Linux)
sudo snap install peerchat-cli

# Chocolatey (Windows)
choco install peerchat-cli

# APT (Ubuntu/Debian)
sudo apt install peerchat-cli

# RPM (Fedora/RHEL)
sudo dnf install peerchat-cli
```

## Platform-Specific Instructions

### Linux

#### Ubuntu/Debian
```bash
# Install dependencies
sudo apt update
sudo apt install golang-go git build-essential

# Clone and build
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go

# Install system-wide
sudo cp bin/peerchat-cli /usr/local/bin/
```

#### Fedora/RHEL/CentOS
```bash
# Install dependencies
sudo dnf install golang git gcc

# Clone and build
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go

# Install system-wide
sudo cp bin/peerchat-cli /usr/local/bin/
```

#### Arch Linux
```bash
# Install dependencies
sudo pacman -S go git gcc

# Clone and build
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go

# Install system-wide
sudo cp bin/peerchat-cli /usr/local/bin/
```

### macOS

#### Using Homebrew
```bash
# Install dependencies
brew install go git

# Clone and build
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go

# Install system-wide
sudo cp bin/peerchat-cli /usr/local/bin/
```

#### Using MacPorts
```bash
# Install dependencies
sudo port install go git

# Clone and build
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go

# Install system-wide
sudo cp bin/peerchat-cli /usr/local/bin/
```

### Windows

#### Using Git Bash
```bash
# Install Go and Git from official websites
# Then in Git Bash:

git clone https://github.com/Xelvra/peerchat.git
cd peerchat
go build -o bin/peerchat-cli.exe cmd/peerchat-cli/main.go

# Add to PATH or copy to desired location
```

#### Using PowerShell
```powershell
# Install dependencies with Chocolatey
choco install golang git mingw

# Clone and build
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
go build -o bin/peerchat-cli.exe cmd/peerchat-cli/main.go

# Add to PATH
$env:PATH += ";$(pwd)\bin"
```

## Configuration

### Initial Setup
```bash
# Initialize identity and configuration
peerchat-cli init

# This creates ~/.xelvra/ with:
# ├── config.yaml      # Node configuration
# ├── identity.json    # Your cryptographic identity
# └── peerchat.log     # Application logs
```

### Configuration File
Edit `~/.xelvra/config.yaml`:

```yaml
# Network configuration
listen_addrs:
  - "/ip4/0.0.0.0/tcp/0"
  - "/ip4/0.0.0.0/udp/0/quic-v1"

# Bootstrap peers (for internet-wide discovery)
bootstrap_peers: []

# Logging configuration
log_level: "info"
log_file: "~/.xelvra/peerchat.log"

# Discovery configuration
discovery:
  mdns_enabled: true
  udp_broadcast_enabled: true
  dht_enabled: false
```

### Environment Variables
```bash
# Override configuration directory
export XELVRA_CONFIG_DIR="/custom/path"

# Override log level
export XELVRA_LOG_LEVEL="debug"

# Override listen addresses
export XELVRA_LISTEN_ADDRS="/ip4/0.0.0.0/tcp/4001"
```

## Verification

### Test Installation
```bash
# Check version
peerchat-cli version

# Run diagnostics
peerchat-cli doctor

# Test basic functionality
peerchat-cli id
```

### Network Test
```bash
# Start in one terminal
peerchat-cli start

# In another terminal, discover peers
peerchat-cli discover

# Check status
peerchat-cli status
```

### Expected Output
```
$ peerchat-cli version
Xelvra P2P Messenger CLI v0.1.0-alpha

$ peerchat-cli doctor
🩺 Network Diagnostics
======================
✅ System checks:
  - OS: Linux
  - Go version: 1.21+

✅ Network connectivity:
  - Internet: Available
  - DNS: Functional

🔧 P2P node checks:
  - Node creation: ✅ Success
  - Peer ID: 12D3KooW...
  - Listen addresses: 2 configured

✅ All diagnostics passed
🎉 Your Xelvra node is ready for P2P communication!
```

## Troubleshooting

### Common Issues

#### Build Errors
**Problem**: `go build` fails with dependency errors.

**Solution**:
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Try building again
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go
```

#### Permission Denied
**Problem**: Cannot execute binary or access configuration.

**Solution**:
```bash
# Make binary executable
chmod +x bin/peerchat-cli

# Fix configuration permissions
chmod 700 ~/.xelvra
chmod 600 ~/.xelvra/*
```

#### Network Issues
**Problem**: Cannot discover peers or connect.

**Solution**:
```bash
# Check firewall
sudo ufw allow 42424/udp  # Linux
# Or configure Windows Firewall / macOS Firewall

# Test network connectivity
peerchat-cli doctor

# Try different network
# Use mobile hotspot to test
```

#### CGO Errors
**Problem**: Build fails with CGO-related errors.

**Solution**:
```bash
# Install C compiler
# Ubuntu/Debian
sudo apt install build-essential

# macOS
xcode-select --install

# Windows
# Install MinGW or Visual Studio Build Tools
```

### Getting Help

1. **Run diagnostics**: `peerchat-cli doctor`
2. **Check logs**: `~/.xelvra/peerchat.log`
3. **GitHub Issues**: https://github.com/Xelvra/peerchat/issues
4. **Documentation**: See [User Guide](USER_GUIDE.md)

---

For usage instructions, see the [User Guide](USER_GUIDE.md).
For development setup, see the [Developer Guide](DEVELOPER_GUIDE.md).
