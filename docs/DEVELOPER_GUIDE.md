# Xelvra P2P Messenger - Developer Guide

## Table of Contents
- [Architecture Overview](#architecture-overview)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Core Components](#core-components)
- [Building and Testing](#building-and-testing)
- [Contributing](#contributing)
- [API Reference](#api-reference)

## Architecture Overview

Xelvra is built on a modular P2P architecture using libp2p as the networking foundation:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Layer     │    │  Web Interface  │    │  Mobile App     │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│                 │    │                 │    │                 │
│  cmd/           │    │  web/           │    │  mobile/        │
│  peerchat-cli/  │    │  (future)       │    │  (future)       │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Core P2P      │
                    │   Engine        │
                    ├─────────────────┤
                    │                 │
                    │  internal/      │
                    │  ├─ p2p/        │
                    │  ├─ message/    │
                    │  ├─ user/       │
                    │  └─ crypto/     │
                    │                 │
                    └─────────────────┘
                             │
                    ┌─────────────────┐
                    │   libp2p        │
                    │   Foundation    │
                    ├─────────────────┤
                    │                 │
                    │  • Host         │
                    │  • Transport    │
                    │  • Discovery    │
                    │  • Security     │
                    │  • Routing      │
                    │                 │
                    └─────────────────┘
```

### Key Design Principles
- **Modularity**: Clear separation between CLI, core engine, and networking
- **Testability**: Comprehensive unit and integration tests
- **Extensibility**: Plugin architecture for future features
- **Security**: Defense in depth with multiple security layers
- **Performance**: Optimized for low latency and resource usage

## Development Setup

### Prerequisites
```bash
# Install Go 1.21+
go version

# Install development tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/securecodewarrior/sast-scan@latest
```

### Clone and Setup
```bash
git clone https://github.com/Xelvra/peerchat.git
cd peerchat

# Install dependencies
go mod download

# Run tests
go test ./...

# Build CLI
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go
```

### Development Workflow
```bash
# Format code
goimports -w .

# Lint code
golangci-lint run

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o bin/peerchat-cli-linux cmd/peerchat-cli/main.go
GOOS=darwin GOARCH=amd64 go build -o bin/peerchat-cli-darwin cmd/peerchat-cli/main.go
GOOS=windows GOARCH=amd64 go build -o bin/peerchat-cli-windows.exe cmd/peerchat-cli/main.go
```

## Project Structure

```
peerchat/
├── cmd/                    # Command-line applications
│   └── peerchat-cli/      # Main CLI application
│       └── main.go        # CLI entry point and commands
├── internal/              # Private application code
│   ├── p2p/              # P2P networking layer
│   │   ├── node.go       # Main P2P node implementation
│   │   ├── wrapper.go    # P2P wrapper with fallback
│   │   ├── discovery.go  # Peer discovery (mDNS, UDP)
│   │   └── config.go     # P2P configuration
│   ├── message/          # Message handling
│   │   ├── manager.go    # Message routing and delivery
│   │   ├── types.go      # Message type definitions
│   │   └── crypto.go     # Message encryption
│   ├── user/             # User identity management
│   │   ├── identity.go   # DID and key management
│   │   └── profile.go    # User profile handling
│   └── crypto/           # Cryptographic utilities
│       ├── keys.go       # Key generation and management
│       └── encryption.go # Encryption/decryption
├── tests/                # Integration and end-to-end tests
├── scripts/              # Build and deployment scripts
├── docs/                 # Documentation
├── bin/                  # Compiled binaries (gitignored)
├── dist/                 # Distribution packages (gitignored)
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── LICENSE              # AGPLv3 license
└── README.md            # Project overview
```

## Core Components

### P2P Node (`internal/p2p/node.go`)
The main P2P node implementation that:
- Manages libp2p host and networking
- Handles peer connections and discovery
- Routes messages between peers
- Maintains node status and metrics

### P2P Wrapper (`internal/p2p/wrapper.go`)
Provides a high-level interface with:
- Automatic fallback to simulation mode
- Simplified API for CLI applications
- Error handling and logging
- Configuration management

### Discovery Manager (`internal/p2p/discovery.go`)
Implements hierarchical peer discovery using:
- **Local Discovery**: mDNS and UDP broadcast for immediate local network peers
- **Global Discovery**: Kademlia DHT for distributed peer finding
- **LRU Caching**: Intelligent local peer caching for performance optimization
- **Bootstrap**: IPFS-compatible bootstrap peers for DHT initialization

### Message Manager (`internal/message/manager.go`)
Handles message routing with:
- End-to-end encryption
- Message queuing and delivery
- Peer-to-peer routing
- Message type handling

### Identity Manager (`internal/user/identity.go`)
Manages user identity:
- DID (Decentralized Identifier) generation
- Cryptographic key management
- Identity verification
- Profile management

## Building and Testing

### Unit Tests
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...

# Run specific package tests
go test ./internal/p2p/
```

### Integration Tests
```bash
# Run integration tests
go test -tags=integration ./tests/

# Run end-to-end tests
go test -tags=e2e ./tests/
```

### Performance Testing
```bash
# Run benchmarks
go test -bench=. ./...

# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=. ./internal/p2p/
go tool pprof cpu.prof

# Profile memory usage
go test -memprofile=mem.prof -bench=. ./internal/p2p/
go tool pprof mem.prof
```

### Code Quality
```bash
# Static analysis
golangci-lint run

# Security scanning
gosec ./...

# Dependency checking
go mod verify
go list -m -u all
```

## Contributing

### Code Style
- Follow standard Go conventions
- Use `gofmt` and `goimports` for formatting
- Write comprehensive tests for new features
- Document public APIs with Go doc comments
- Use meaningful variable and function names

### Commit Guidelines
```
type(scope): description

[optional body]

[optional footer]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Example:
```
feat(p2p): add UDP broadcast discovery

Implement UDP broadcast-based peer discovery for local networks.
This complements mDNS discovery and works on networks where
multicast is disabled.

Closes #123
```

### Pull Request Process
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes with tests
4. Run the full test suite: `go test ./...`
5. Commit your changes: `git commit -m 'feat: add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Development Guidelines
- **Security First**: All code must be security-reviewed
- **Test Coverage**: Maintain >80% test coverage
- **Documentation**: Update docs for user-facing changes
- **Backwards Compatibility**: Don't break existing APIs
- **Performance**: Profile performance-critical code

## API Reference

### P2P Wrapper API
```go
// Create new P2P wrapper
wrapper := p2p.NewP2PWrapper(ctx, useSimulation)

// Start the node
err := wrapper.Start()

// Get node information
nodeInfo := wrapper.GetNodeInfo()

// Discover peers
peers := wrapper.GetDiscoveredPeers()

// Connect to peer
success := wrapper.ConnectToPeer(peerID)

// Send message
err := wrapper.SendMessage(peerID, message)

// Stop the node
wrapper.Stop()
```

### Message API
```go
// Send text message
err := messageManager.SendTextMessage(peerID, text)

// Send file
err := messageManager.SendFile(peerID, filePath)

// Handle incoming messages
messageManager.SetMessageHandler(func(msg *message.Message) {
    // Process message
})
```

### Identity API
```go
// Create new identity
identity, err := user.NewIdentity()

// Get DID
did := identity.GetDID()

// Sign data
signature, err := identity.Sign(data)

// Verify signature
valid := identity.Verify(data, signature, publicKey)
```

---

For user documentation, see the [User Guide](USER_GUIDE.md).
For API documentation, see the [API Reference](API_REFERENCE.md).
