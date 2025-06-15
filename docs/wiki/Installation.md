# Installation Guide

This comprehensive guide covers installation of Xelvra P2P Messenger on all supported platforms.

## 🚀 Quick Installation

### Pre-built Binaries (Recommended)

Download the latest release for your platform:

#### Linux
```bash
# Download and install
curl -L https://github.com/Xelvra/peerchat/releases/latest/download/peerchat-cli-linux -o peerchat-cli
chmod +x peerchat-cli
sudo mv peerchat-cli /usr/local/bin/

# Verify installation
peerchat-cli version
```

#### macOS
```bash
# Intel Macs
curl -L https://github.com/Xelvra/peerchat/releases/latest/download/peerchat-cli-darwin-amd64 -o peerchat-cli

# Apple Silicon Macs
curl -L https://github.com/Xelvra/peerchat/releases/latest/download/peerchat-cli-darwin-arm64 -o peerchat-cli

# Make executable and install
chmod +x peerchat-cli
sudo mv peerchat-cli /usr/local/bin/

# Verify installation
peerchat-cli version
```

#### Windows
```powershell
# Download using PowerShell
Invoke-WebRequest -Uri "https://github.com/Xelvra/peerchat/releases/latest/download/peerchat-cli-windows.exe" -OutFile "peerchat-cli.exe"

# Add to PATH (optional)
# Move peerchat-cli.exe to a directory in your PATH

# Verify installation
.\peerchat-cli.exe version
```

## 🔨 Build from Source

### Prerequisites

#### All Platforms
- **Go 1.21 or later** - [Download Go](https://golang.org/dl/)
- **Git** - [Download Git](https://git-scm.com/downloads)
- **Network connectivity** for downloading dependencies

#### Platform-Specific Requirements

**Linux:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install build-essential git

# CentOS/RHEL/Fedora
sudo dnf install gcc git make
# or
sudo yum install gcc git make
```

**macOS:**
```bash
# Install Xcode Command Line Tools
xcode-select --install

# Or install via Homebrew
brew install git go
```

**Windows:**
- Install [Git for Windows](https://git-scm.com/download/win)
- Install [Go for Windows](https://golang.org/dl/)
- Install [Build Tools for Visual Studio](https://visualstudio.microsoft.com/downloads/#build-tools-for-visual-studio-2022) (optional, for CGO)

### Build Steps

1. **Clone the repository:**
```bash
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
```

2. **Download dependencies:**
```bash
go mod download
go mod verify
```

3. **Build the CLI:**
```bash
# Linux/macOS
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go

# Windows
go build -o bin/peerchat-cli.exe cmd/peerchat-cli/main.go
```

4. **Verify the build:**
```bash
# Linux/macOS
./bin/peerchat-cli version

# Windows
.\bin\peerchat-cli.exe version
```

5. **Install system-wide (optional):**
```bash
# Linux/macOS
sudo cp bin/peerchat-cli /usr/local/bin/

# Windows - Add bin/ directory to your PATH
```

## 📦 Package Managers

### Linux Package Managers

#### Snap (Universal Linux)
```bash
# Install from Snap Store (when available)
sudo snap install xelvra-peerchat
```

#### Flatpak
```bash
# Install from Flathub (when available)
flatpak install flathub org.xelvra.PeerChat
```

#### Arch Linux (AUR)
```bash
# Install from AUR (when available)
yay -S xelvra-peerchat
# or
paru -S xelvra-peerchat
```

### macOS Package Managers

#### Homebrew
```bash
# Install via Homebrew (when available)
brew install xelvra/tap/peerchat
```

#### MacPorts
```bash
# Install via MacPorts (when available)
sudo port install peerchat
```

### Windows Package Managers

#### Chocolatey
```powershell
# Install via Chocolatey (when available)
choco install xelvra-peerchat
```

#### Scoop
```powershell
# Install via Scoop (when available)
scoop bucket add xelvra https://github.com/Xelvra/scoop-bucket
scoop install peerchat
```

## 🐳 Docker Installation

### Using Docker Hub
```bash
# Pull the latest image
docker pull xelvra/peerchat:latest

# Run interactively
docker run -it --rm \
  -v ~/.xelvra:/root/.xelvra \
  -p 42424:42424/udp \
  xelvra/peerchat:latest

# Run specific command
docker run --rm \
  -v ~/.xelvra:/root/.xelvra \
  xelvra/peerchat:latest doctor
```

### Build Docker Image
```bash
# Clone repository
git clone https://github.com/Xelvra/peerchat.git
cd peerchat

# Build image
docker build -t xelvra/peerchat .

# Run
docker run -it --rm \
  -v ~/.xelvra:/root/.xelvra \
  -p 42424:42424/udp \
  xelvra/peerchat
```

## ⚙️ Post-Installation Setup

### 1. Initialize Your Identity
```bash
peerchat-cli init
```

This creates:
- `~/.xelvra/config.yaml` - Configuration file
- `~/.xelvra/identity.key` - Your private key (keep secure!)
- `~/.xelvra/peerchat.log` - Log file

### 2. Test Your Installation
```bash
peerchat-cli doctor
```

This checks:
- Network connectivity
- Firewall configuration
- Port availability
- System requirements

### 3. Configure Firewall

#### Linux (UFW)
```bash
# Allow peer discovery
sudo ufw allow 42424/udp

# Allow dynamic P2P ports (optional)
sudo ufw allow out 1024:65535/tcp
sudo ufw allow out 1024:65535/udp
```

#### Linux (iptables)
```bash
# Allow peer discovery
sudo iptables -A INPUT -p udp --dport 42424 -j ACCEPT

# Allow outgoing P2P connections
sudo iptables -A OUTPUT -p tcp --dport 1024:65535 -j ACCEPT
sudo iptables -A OUTPUT -p udp --dport 1024:65535 -j ACCEPT
```

#### macOS
```bash
# macOS Firewall usually allows outgoing connections by default
# If you have strict firewall rules, allow UDP 42424
```

#### Windows Firewall
```powershell
# Allow UDP 42424 for peer discovery
New-NetFirewallRule -DisplayName "Xelvra Peer Discovery" -Direction Inbound -Protocol UDP -LocalPort 42424 -Action Allow

# Allow outgoing P2P connections (usually allowed by default)
New-NetFirewallRule -DisplayName "Xelvra P2P Outbound" -Direction Outbound -Protocol TCP -RemotePort 1024-65535 -Action Allow
```

## 🔧 Configuration

### Basic Configuration
Edit `~/.xelvra/config.yaml`:

```yaml
# Network settings
network:
  listen_port: 0          # 0 = auto-select
  discovery_port: 42424   # UDP discovery port
  enable_mdns: true       # Local network discovery
  enable_udp_broadcast: true
  max_peers: 50          # Maximum concurrent connections

# Logging settings
logging:
  level: "info"          # debug, info, warn, error
  file: "peerchat.log"
  max_size: 10           # MB
  max_backups: 3
  max_age: 30            # days

# Security settings
security:
  key_rotation_days: 60  # Automatic key rotation
  max_message_size: 1048576  # 1MB
```

### Advanced Configuration
```yaml
# Performance tuning
performance:
  connection_timeout: 30s
  read_timeout: 10s
  write_timeout: 10s
  max_concurrent_streams: 100

# Development settings (for testing)
development:
  enable_simulation: false
  log_network_events: false
  disable_encryption: false  # Never use in production!
```

## 🧪 Verification

### Test Basic Functionality
```bash
# Check version
peerchat-cli version

# Run diagnostics
peerchat-cli doctor

# Test identity
peerchat-cli id

# Start interactive mode (Ctrl+C to exit)
peerchat-cli start
```

### Test Network Discovery
```bash
# In one terminal
peerchat-cli listen

# In another terminal
peerchat-cli discover
```

### Test File Transfer
```bash
# Create test file
echo "Hello, Xelvra!" > test.txt

# Send to discovered peer
peerchat-cli send-file <peer_id> test.txt
```

## 🚨 Troubleshooting Installation

### Common Issues

#### "Command not found"
**Solution:** Add the binary to your PATH or use the full path:
```bash
# Linux/macOS
export PATH=$PATH:/path/to/peerchat-cli
# or
/full/path/to/peerchat-cli version
```

#### "Permission denied"
**Solution:** Make the binary executable:
```bash
chmod +x peerchat-cli
```

#### "Go version too old"
**Solution:** Update Go to version 1.21 or later:
```bash
# Check current version
go version

# Update Go (download from https://golang.org/dl/)
```

#### Build fails with "missing dependencies"
**Solution:** Ensure you have build tools installed:
```bash
# Linux
sudo apt install build-essential

# macOS
xcode-select --install
```

#### "Network unreachable" during build
**Solution:** Check internet connection and proxy settings:
```bash
# Set proxy if needed
export GOPROXY=https://proxy.golang.org,direct
export GOSUMDB=sum.golang.org
```

### Getting Help

If you encounter issues:

1. **Check the logs:** `~/.xelvra/peerchat.log`
2. **Run diagnostics:** `peerchat-cli doctor`
3. **Search [GitHub Issues](https://github.com/Xelvra/peerchat/issues)**
4. **Ask in [Discussions](https://github.com/Xelvra/peerchat/discussions)**
5. **Create a [Bug Report](https://github.com/Xelvra/peerchat/issues/new/choose)**

## 🎉 Next Steps

Once installed:

1. **[Get Started](Getting-Started)** - Quick start guide
2. **[Read the User Manual](User-Manual)** - Complete documentation
3. **[Join the Community](https://github.com/Xelvra/peerchat/discussions)** - Connect with other users

Welcome to **#XelvraFree**! 🌐
