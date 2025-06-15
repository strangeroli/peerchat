# Xelvra P2P Messenger Development Guide

## Project Structure

```
xelvra/
├── bin/                    # Compiled binaries (gitignored)
├── cmd/                    # Application entry points
│   ├── peerchat-cli/      # CLI application
│   └── peerchat-api/      # API server (planned)
├── dist/                   # Distribution packages (gitignored)
├── docs/                   # Documentation
├── internal/               # Internal packages
│   ├── crypto/            # Cryptographic operations
│   ├── db/                # Database operations
│   ├── message/           # Message handling
│   ├── p2p/               # P2P networking
│   ├── user/              # User identity management
│   └── util/              # Utility functions
├── pkg/                    # Public packages
├── scripts/                # Build and utility scripts
├── tests/                  # Test files
└── peerchat_gui/          # Flutter GUI (planned)
```

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional)

### Clone and Setup

```bash
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
go mod tidy
```

### Build

```bash
# Build for current platform
./scripts/build.sh

# Cross-platform build
./scripts/build-cross.sh
```

### Testing

```bash
# Run all tests
./scripts/test.sh

# Run specific tests
go test -v ./tests/...
go test -v ./internal/crypto/...

# Run benchmarks
go test -bench=. -benchmem ./tests/...
```

## Architecture Overview

### Core Components

1. **P2P Networking** (`internal/p2p/`)
   - libp2p-based networking
   - QUIC and TCP transports
   - Peer discovery and management

2. **Cryptography** (`internal/crypto/`)
   - Signal Protocol implementation
   - X3DH key agreement
   - Double Ratchet encryption

3. **User Identity** (`internal/user/`)
   - DID-based identity system
   - Ed25519 signatures
   - Trust level management

4. **Message Handling** (`internal/message/`)
   - Message routing and processing
   - Protocol handlers
   - Encryption/decryption

5. **Database** (`internal/db/`)
   - SQLite with WAL mode
   - Message persistence
   - User data storage

### Design Principles

- **Performance First**: Target <50ms latency, <20MB memory, <1% CPU
- **Security by Design**: End-to-end encryption, memory protection
- **Decentralization**: No central servers or authorities
- **Modularity**: Clean separation of concerns
- **Testability**: Comprehensive test coverage

## Coding Standards

### Go Style

Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html).

### Key Guidelines

1. **Error Handling**: Always handle errors explicitly
2. **Context**: Use context.Context for cancellation and timeouts
3. **Logging**: Use structured logging with logrus
4. **Testing**: Write tests for all public functions
5. **Documentation**: Document all exported functions and types

### Example Code Style

```go
// Package crypto provides cryptographic operations for Xelvra messenger.
package crypto

import (
    "context"
    "fmt"
    
    "github.com/sirupsen/logrus"
)

// SignalCrypto provides Signal Protocol cryptographic operations.
type SignalCrypto struct {
    logger *logrus.Logger
}

// NewSignalCrypto creates a new Signal Protocol crypto instance.
func NewSignalCrypto(logger *logrus.Logger) (*SignalCrypto, error) {
    if logger == nil {
        return nil, fmt.Errorf("logger cannot be nil")
    }
    
    return &SignalCrypto{
        logger: logger,
    }, nil
}
```

## Testing Guidelines

### Test Structure

- Unit tests in the same package as the code
- Integration tests in `tests/` directory
- Benchmark tests for performance-critical code

### Test Naming

```go
func TestFunctionName(t *testing.T)           // Unit test
func TestFunctionName_ErrorCase(t *testing.T) // Error case
func BenchmarkFunctionName(b *testing.B)      // Benchmark
```

### Test Coverage

Aim for >80% test coverage. Check coverage with:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Performance Guidelines

### Memory Management

- Use object pools for frequently allocated objects
- Avoid unnecessary allocations in hot paths
- Use `sync.Pool` for reusable buffers

### Concurrency

- Use goroutines for I/O operations
- Protect shared state with mutexes
- Use channels for communication between goroutines

### Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.

# Memory profiling
go test -memprofile=mem.prof -bench=.

# Analyze profiles
go tool pprof cpu.prof
go tool pprof mem.prof
```

## Security Considerations

### Cryptographic Operations

- Use constant-time operations for sensitive data
- Properly handle key material
- Implement secure random number generation

### Memory Protection

- Zero out sensitive data after use
- Use memory protection libraries where possible
- Avoid logging sensitive information

### Input Validation

- Validate all external inputs
- Use safe parsing functions
- Implement rate limiting

## Contributing

### Workflow

1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Run the test suite
5. Submit a pull request

### Commit Messages

Use conventional commit format:

```
feat: add QUIC transport support
fix: resolve memory leak in message handler
docs: update CLI usage guide
test: add benchmarks for crypto operations
```

### Pull Request Guidelines

- Include tests for new functionality
- Update documentation as needed
- Ensure all tests pass
- Follow the coding standards

## Build System

### Scripts

- `scripts/build.sh`: Build for current platform
- `scripts/build-cross.sh`: Cross-platform build
- `scripts/test.sh`: Run comprehensive tests

### Build Flags

```bash
# Development build
go build -o bin/peerchat-cli ./cmd/peerchat-cli

# Production build with optimizations
go build -ldflags="-s -w" -o bin/peerchat-cli ./cmd/peerchat-cli

# Cross-compilation
GOOS=windows GOARCH=amd64 go build -o bin/peerchat-cli.exe ./cmd/peerchat-cli
```

## Debugging

### Debug Builds

```bash
go build -gcflags="all=-N -l" -o bin/peerchat-cli-debug ./cmd/peerchat-cli
```

### Using Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the CLI
dlv debug ./cmd/peerchat-cli -- start --verbose
```

### Logging

Enable debug logging:

```bash
peerchat-cli start --verbose
```

## Release Process

### Version Tagging

```bash
git tag -a v0.1.0 -m "Release version 0.1.0"
git push origin v0.1.0
```

### Release Builds

```bash
./scripts/build-cross.sh
```

This creates distribution packages in the `dist/` directory.

## Dependencies

### Core Dependencies

- `github.com/libp2p/go-libp2p`: P2P networking
- `github.com/spf13/cobra`: CLI framework
- `github.com/sirupsen/logrus`: Structured logging
- `github.com/mattn/go-sqlite3`: SQLite database

### Development Dependencies

- Testing: Go standard library
- Benchmarking: Go standard library
- Profiling: `go tool pprof`

## Troubleshooting

### Common Build Issues

**CGO errors with SQLite:**
```bash
# Install SQLite development headers
sudo apt-get install libsqlite3-dev  # Ubuntu/Debian
sudo yum install sqlite-devel        # CentOS/RHEL
```

**libp2p compilation issues:**
```bash
# Update Go to latest version
go mod tidy
go clean -modcache
```

### Performance Issues

- Use profiling tools to identify bottlenecks
- Check for goroutine leaks
- Monitor memory usage patterns
- Optimize hot paths identified by benchmarks
