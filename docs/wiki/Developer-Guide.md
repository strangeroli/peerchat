# Developer Guide

Comprehensive guide for developers who want to contribute to Xelvra P2P Messenger or understand its technical implementation.

## 📋 Table of Contents

- [Architecture Overview](#architecture-overview)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Core Components](#core-components)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Contributing](#contributing)
- [API Documentation](#api-documentation)

## 🏗️ Architecture Overview

Xelvra follows a modular, epoch-based development approach:

### Development Epochs

**Epoch 1: CLI Foundation** (Current)
- Core P2P networking with libp2p
- Command-line interface for testing and development
- Basic encryption and security features
- Local network discovery and NAT traversal

**Epoch 2: API Service** (Planned)
- gRPC API server for frontend communication
- Database layer with SQLite
- Event-driven architecture
- Monitoring and telemetry

**Epoch 3: GUI Application** (Planned)
- Cross-platform Flutter application
- Mobile-first design with desktop support
- Advanced UI/UX features
- Energy optimization for mobile devices

**Epoch 4: Advanced Features** (Future)
- Zero-knowledge proofs
- Quantum-resistant cryptography
- Voice and video communication
- Advanced mesh networking

### System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   GUI Client    │    │   CLI Client    │    │  Other Clients  │
│   (Flutter)     │    │     (Go)        │    │                 │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────┴───────────┐
                    │     API Service         │
                    │      (gRPC)            │
                    └─────────────┬───────────┘
                                 │
                    ┌─────────────┴───────────┐
                    │      P2P Core          │
                    │     (libp2p)           │
                    └─────────────────────────┘
```

## 🛠️ Development Setup

### Prerequisites

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Git** - Version control
- **Make** - Build automation (optional)
- **Docker** - For containerized development (optional)

### Environment Setup

1. **Clone the repository:**
```bash
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
```

2. **Install development tools:**
```bash
# Code formatting and linting
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Testing tools
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go install github.com/golang/mock/mockgen@latest

# Protocol buffer compiler (for future API development)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

3. **Download dependencies:**
```bash
go mod download
go mod verify
```

4. **Build the project:**
```bash
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go
```

5. **Run tests:**
```bash
go test ./...
```

### IDE Configuration

#### VS Code
Recommended extensions:
- Go (official Go extension)
- Go Test Explorer
- GitLens
- Markdown All in One

Configuration (`.vscode/settings.json`):
```json
{
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.formatTool": "goimports",
    "go.testFlags": ["-v", "-race"],
    "go.coverOnSave": true
}
```

#### GoLand/IntelliJ
- Enable Go modules support
- Configure golangci-lint as external tool
- Set up run configurations for CLI and tests

## 📁 Project Structure

```
peerchat/
├── cmd/                    # Application entry points
│   ├── peerchat-cli/      # CLI application
│   └── peerchat-api/      # API service (planned)
├── internal/              # Private application code
│   ├── p2p/              # P2P networking core
│   ├── crypto/           # Cryptography and security
│   ├── user/             # User identity management
│   ├── message/          # Message handling
│   ├── db/               # Database operations
│   ├── api/              # API handlers (planned)
│   └── util/             # Utility functions
├── pkg/                   # Public library code
│   └── proto/            # Protocol buffer definitions
├── tests/                 # Test files
│   ├── unit/             # Unit tests
│   ├── integration/      # Integration tests
│   └── e2e/              # End-to-end tests
├── scripts/              # Build and utility scripts
├── docs/                 # Documentation
├── .github/              # GitHub workflows and templates
└── deployments/          # Deployment configurations
```

### Key Directories

**`internal/`** - Core application logic, not importable by external packages
**`pkg/`** - Public APIs that can be imported by other projects
**`cmd/`** - Application entry points and main functions
**`tests/`** - All test files organized by type

## 🔧 Core Components

### P2P Networking (`internal/p2p/`)

The P2P core handles all networking functionality:

```go
// Core P2P node interface
type Node interface {
    Start(ctx context.Context) error
    Stop() error
    Connect(peerID peer.ID) error
    Disconnect(peerID peer.ID) error
    SendMessage(peerID peer.ID, msg []byte) error
    Discover() ([]peer.AddrInfo, error)
}
```

**Key files:**
- `node.go` - Main P2P node implementation
- `discovery.go` - Peer discovery mechanisms
- `transport.go` - Transport layer (QUIC/TCP)
- `nat.go` - NAT traversal and hole punching
- `relay.go` - Relay server functionality

### Cryptography (`internal/crypto/`)

Handles all cryptographic operations:

```go
// Encryption interface
type Encryptor interface {
    Encrypt(plaintext []byte, recipientKey []byte) ([]byte, error)
    Decrypt(ciphertext []byte, senderKey []byte) ([]byte, error)
    GenerateKeyPair() (PrivateKey, PublicKey, error)
}
```

**Key files:**
- `signal.go` - Signal Protocol implementation
- `keys.go` - Key management and rotation
- `identity.go` - Identity cryptography
- `x3dh.go` - X3DH key agreement protocol
- `double_ratchet.go` - Double Ratchet algorithm

### User Management (`internal/user/`)

Manages user identities and profiles:

```go
// User identity
type Identity struct {
    DID        string
    PeerID     peer.ID
    PrivateKey crypto.PrivKey
    PublicKey  crypto.PubKey
    Profile    Profile
}
```

### Message Handling (`internal/message/`)

Processes and routes messages:

```go
// Message interface
type Message interface {
    ID() string
    Sender() peer.ID
    Recipient() peer.ID
    Content() []byte
    Timestamp() time.Time
    Encrypt(key []byte) error
    Decrypt(key []byte) error
}
```

## 🔄 Development Workflow

### Git Workflow

We use a simplified Git flow:

1. **Main branch** - Stable releases
2. **Develop branch** - Integration branch for features
3. **Feature branches** - Individual features (`feature/feature-name`)
4. **Hotfix branches** - Critical fixes (`hotfix/fix-name`)

### Branch Naming

- `feature/add-mesh-networking`
- `bugfix/fix-connection-timeout`
- `hotfix/security-patch`
- `docs/update-api-reference`

### Commit Messages

Follow conventional commits:

```
type(scope): description

[optional body]

[optional footer]
```

Examples:
```
feat(p2p): add QUIC transport support
fix(crypto): resolve key rotation issue
docs(wiki): update installation guide
test(integration): add NAT traversal tests
```

### Pull Request Process

1. **Create feature branch** from develop
2. **Implement changes** with tests
3. **Update documentation** as needed
4. **Run full test suite**
5. **Submit pull request** with clear description
6. **Address review feedback**
7. **Merge after approval**

## 🧪 Testing

### Test Structure

```
tests/
├── unit/                  # Unit tests
│   ├── p2p_test.go
│   ├── crypto_test.go
│   └── message_test.go
├── integration/           # Integration tests
│   ├── node_integration_test.go
│   └── discovery_test.go
└── e2e/                  # End-to-end tests
    ├── chat_flow_test.go
    └── file_transfer_test.go
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test package
go test ./internal/p2p/

# Run specific test
go test -run TestNodeConnection ./internal/p2p/

# Run benchmarks
go test -bench=. ./...

# Run integration tests
go test -tags=integration ./tests/integration/
```

### Test Guidelines

1. **Unit Tests** - Test individual functions and methods
2. **Integration Tests** - Test component interactions
3. **End-to-End Tests** - Test complete user workflows
4. **Benchmarks** - Performance testing for critical paths

### Mock Generation

```bash
# Generate mocks for interfaces
mockgen -source=internal/p2p/node.go -destination=tests/mocks/node_mock.go
```

## 🤝 Contributing

### Code Style

Follow Go best practices:

```go
// Good: Clear function names and documentation
// ConnectToPeer establishes a connection to the specified peer
func (n *Node) ConnectToPeer(ctx context.Context, peerID peer.ID) error {
    if peerID == n.host.ID() {
        return ErrSelfConnection
    }
    
    // Implementation...
}

// Good: Proper error handling
if err := n.ConnectToPeer(ctx, peerID); err != nil {
    return fmt.Errorf("failed to connect to peer %s: %w", peerID, err)
}
```

### Documentation

- **Public APIs** - Document all exported functions and types
- **Complex Logic** - Add comments explaining the "why"
- **Examples** - Include usage examples in documentation
- **README Updates** - Keep README.md current with changes

### Security Considerations

- **Input Validation** - Validate all external inputs
- **Error Handling** - Don't leak sensitive information in errors
- **Cryptography** - Use established libraries and algorithms
- **Dependencies** - Keep dependencies updated and minimal

## 📚 API Documentation

### Current CLI API

The CLI provides the current public API:

```bash
# Core operations
peerchat-cli init
peerchat-cli start
peerchat-cli connect <peer_id>
peerchat-cli send <peer_id> <message>
peerchat-cli discover

# File operations
peerchat-cli send-file <peer_id> <file_path>

# Diagnostics
peerchat-cli status
peerchat-cli doctor
peerchat-cli id
```

### Future gRPC API (Epoch 2)

Planned gRPC service definition:

```protobuf
service PeerChatService {
    rpc Connect(ConnectRequest) returns (ConnectResponse);
    rpc SendMessage(SendMessageRequest) returns (SendMessageResponse);
    rpc DiscoverPeers(DiscoverRequest) returns (stream PeerInfo);
    rpc GetStatus(StatusRequest) returns (StatusResponse);
}
```

### Internal APIs

Key internal interfaces:

```go
// P2P Node interface
type Node interface {
    Start(context.Context) error
    Stop() error
    Connect(peer.ID) error
    SendMessage(peer.ID, []byte) error
}

// Message handler interface
type MessageHandler interface {
    HandleMessage(Message) error
    HandleFileTransfer(FileTransfer) error
}

// Discovery interface
type Discovery interface {
    Discover(context.Context) ([]peer.AddrInfo, error)
    Advertise(context.Context) error
}
```

## 🔍 Debugging

### Debug Logging

```bash
# Enable debug logging
export XELVRA_LOG_LEVEL=debug
peerchat-cli start

# Or use flag
peerchat-cli start --log-level debug
```

### Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

### Network Debugging

```bash
# Monitor network activity
peerchat-cli listen --verbose

# Test specific network functions
peerchat-cli doctor --test-nat
peerchat-cli doctor --test-discovery
```

## 📞 Getting Help

### Development Support

- **[GitHub Discussions](https://github.com/Xelvra/peerchat/discussions)** - Ask questions
- **[GitHub Issues](https://github.com/Xelvra/peerchat/issues)** - Report bugs
- **[Contributing Guide](https://github.com/Xelvra/peerchat/blob/main/CONTRIBUTING.md)** - Detailed contribution guidelines

### Resources

- **[Go Documentation](https://golang.org/doc/)**
- **[libp2p Documentation](https://docs.libp2p.io/)**
- **[Protocol Buffers](https://developers.google.com/protocol-buffers)**
- **[gRPC Documentation](https://grpc.io/docs/)**

---

**Ready to contribute?** Check out our [Contributing Guide](https://github.com/Xelvra/peerchat/blob/main/CONTRIBUTING.md) and join the **#XelvraFree** movement! 🚀
