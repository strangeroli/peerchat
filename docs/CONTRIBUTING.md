# Contributing to Xelvra P2P Messenger

Thank you for your interest in contributing to Xelvra! This document provides guidelines and information for contributors.

## Table of Contents
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Guidelines](#contributing-guidelines)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing Requirements](#testing-requirements)
- [Documentation](#documentation)

## Code of Conduct

This project adheres to a Code of Conduct that all contributors are expected to follow. Please read the full text in the main README.md to understand what actions will and will not be tolerated.

### Our Values
- **Respect**: Treat all community members with respect
- **Inclusivity**: Welcome people of all backgrounds and experience levels
- **Openness**: Be open to feedback and different perspectives
- **Collaboration**: Work together towards common goals
- **Safety**: Maintain a safe environment for everyone

## Getting Started

### Prerequisites
- Go 1.21 or later
- Git
- Basic understanding of P2P networking concepts
- Familiarity with libp2p (helpful but not required)

### First Steps
1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/peerchat.git
   cd peerchat
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/Xelvra/peerchat.git
   ```
4. **Install dependencies**:
   ```bash
   go mod download
   ```
5. **Build and test**:
   ```bash
   go build -o bin/peerchat-cli cmd/peerchat-cli/main.go
   go test ./...
   ```

## Development Setup

### Environment Setup
```bash
# Install development tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Set up pre-commit hooks (optional but recommended)
cp scripts/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

### IDE Configuration
We recommend using VS Code with the Go extension. Configuration files are provided in `.vscode/`.

### Docker Development Environment
```bash
# Start development environment
docker-compose -f docker-compose.dev.yml up -d

# Run tests in container
docker-compose exec dev go test ./...
```

## Contributing Guidelines

### Types of Contributions

#### 🐛 Bug Reports
- Use the bug report template
- Include steps to reproduce
- Provide system information
- Include relevant logs

#### ✨ Feature Requests
- Use the feature request template
- Explain the use case
- Consider backwards compatibility
- Discuss implementation approach

#### 📝 Documentation
- Fix typos and improve clarity
- Add examples and tutorials
- Update API documentation
- Translate documentation

#### 🔧 Code Contributions
- Bug fixes
- New features
- Performance improvements
- Security enhancements

### Contribution Workflow

1. **Check existing issues** - Look for related issues or discussions
2. **Create an issue** - For significant changes, create an issue first
3. **Fork and branch** - Create a feature branch from main
4. **Develop** - Make your changes with tests
5. **Test** - Ensure all tests pass
6. **Document** - Update documentation as needed
7. **Submit PR** - Create a pull request with clear description

## Pull Request Process

### Before Submitting
- [ ] Code follows project style guidelines
- [ ] All tests pass locally
- [ ] Documentation is updated
- [ ] Commit messages follow conventions
- [ ] PR description is clear and complete

### PR Requirements
1. **Clear title** - Summarize the change in 50 characters or less
2. **Detailed description** - Explain what, why, and how
3. **Link issues** - Reference related issues with "Fixes #123"
4. **Test coverage** - Include tests for new functionality
5. **Documentation** - Update relevant documentation

### Review Process
1. **Automated checks** - CI/CD pipeline must pass
2. **Code review** - At least one maintainer review required
3. **Testing** - Manual testing for significant changes
4. **Approval** - Maintainer approval required for merge

### Merge Requirements
- All CI checks pass
- At least one approving review
- No unresolved conversations
- Up-to-date with main branch

## Coding Standards

### Go Style Guide
Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and these additional guidelines:

#### Code Formatting
```bash
# Format code
gofmt -w .
goimports -w .

# Lint code
golangci-lint run
```

#### Naming Conventions
- **Packages**: lowercase, single word when possible
- **Functions**: camelCase, exported functions start with uppercase
- **Variables**: camelCase, descriptive names
- **Constants**: CamelCase for exported, camelCase for internal

#### Error Handling
```go
// Good: Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to connect to peer %s: %w", peerID, err)
}

// Good: Check for specific error types
if errors.Is(err, context.DeadlineExceeded) {
    // Handle timeout
}
```

#### Logging
```go
// Use structured logging
logger.WithFields(logrus.Fields{
    "peer_id": peerID,
    "action":  "connect",
}).Info("Attempting peer connection")
```

### Security Guidelines
- **Input validation**: Validate all external inputs
- **Error messages**: Don't leak sensitive information
- **Cryptography**: Use established libraries and algorithms
- **Dependencies**: Keep dependencies up to date
- **Secrets**: Never commit secrets or keys

## Testing Requirements

### Test Types
1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test component interactions
3. **End-to-End Tests**: Test complete user workflows
4. **Performance Tests**: Benchmark critical paths

### Test Guidelines
```go
// Good: Table-driven tests
func TestPeerConnection(t *testing.T) {
    tests := []struct {
        name     string
        peerID   string
        expected error
    }{
        {"valid peer", "12D3KooW...", nil},
        {"invalid peer", "invalid", ErrInvalidPeerID},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ConnectToPeer(tt.peerID)
            assert.Equal(t, tt.expected, err)
        })
    }
}
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific tests
go test -run TestPeerConnection ./internal/p2p/

# Run benchmarks
go test -bench=. ./...
```

### Test Coverage
- Maintain >80% test coverage
- Focus on critical paths and edge cases
- Include error conditions and timeouts

## Documentation

### Documentation Standards
- **API docs**: Use Go doc comments for all public APIs
- **User docs**: Write clear, step-by-step instructions
- **Code comments**: Explain why, not what
- **Examples**: Include working code examples

### Documentation Structure
```
docs/
├── USER_GUIDE.md          # End-user documentation
├── DEVELOPER_GUIDE.md     # Developer setup and architecture
├── API_REFERENCE.md       # Complete API documentation
├── INSTALLATION.md        # Installation instructions
├── CONTRIBUTING.md        # This file
└── examples/              # Code examples and tutorials
```

### Writing Guidelines
- Use clear, concise language
- Include code examples
- Test all examples
- Keep documentation up to date

## Community

### Communication Channels
- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and ideas
- **Pull Requests**: Code review and collaboration

### Getting Help
- Check existing documentation
- Search GitHub issues
- Ask in GitHub Discussions
- Contact maintainers for security issues

### Recognition
Contributors are recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation

## License

By contributing to Xelvra, you agree that your contributions will be licensed under the GNU Affero General Public License v3.0 (AGPLv3).

---

Thank you for contributing to Xelvra! Together, we're building the future of decentralized communication.
