# Xelvra v0.4.0-alpha - Event-Driven Architecture Release

**Release Date:** June 17, 2025  
**Version:** 0.4.0-alpha  
**Codename:** Event Storm  

## 🚀 What's New in v0.4.0-alpha

This release introduces a complete **event-driven architecture** overhaul, making Xelvra more responsive, scalable, and maintainable. The new architecture provides real-time event processing, advanced logging capabilities, and a foundation for future GUI integration.

### ✨ Major Features

#### 🎯 Event-Driven Architecture
- **Centralized Event Bus System** - High-performance event processing with worker goroutines
- **Event Emitter Components** - Standardized event emission for all P2P operations
- **Advanced Callback Manager** - Priority-based callbacks with timeout, retry, and debounce support
- **Asynchronous Event Processing** - Non-blocking event handling for improved responsiveness

#### 📊 Advanced Logging System
- **Structured Logging** - JSON and text format support with configurable levels
- **Log Rotation** - Automatic log file rotation with size and age limits
- **Component-Specific Logging** - Individual log levels per component
- **Performance Optimized** - Minimal overhead logging with buffering

#### 🔧 Technical Improvements
- **Fixed CLI Version Display** - Dynamic version display instead of hardcoded values
- **Enhanced Build System** - Improved error handling for missing components
- **Comprehensive Unit Tests** - Full test coverage for event system components
- **golangci-lint Integration** - Code quality enforcement

## 📋 Event Types Supported

The new event system supports comprehensive event types:

### P2P Events
- `peer.connected` - Peer connection established
- `peer.disconnected` - Peer connection lost
- `peer.discovered` - New peer discovered via DHT/mDNS

### Message Events
- `message.received` - Incoming message received
- `message.sent` - Message successfully sent
- `message.failed` - Message delivery failed

### File Transfer Events
- `file.transfer.started` - File transfer initiated
- `file.transfer.progress` - Transfer progress update
- `file.transfer.completed` - Transfer completed successfully
- `file.transfer.failed` - Transfer failed

### Node Events
- `node.started` - P2P node started
- `node.stopped` - P2P node stopped
- `node.error` - Node error occurred

### Network Events
- `network.connected` - Network connection established
- `network.disconnected` - Network connection lost
- `network.error` - Network error occurred

## 🏗️ Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Layer     │    │   Event Bus     │    │  P2P Layer      │
│                 │    │                 │    │                 │
│ • Commands      │◄──►│ • Event Queue   │◄──►│ • DHT           │
│ • Interactive   │    │ • Workers       │    │ • NAT Traversal │
│ • Daemon        │    │ • Subscriptions │    │ • Transport     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Event Emitters  │    │ Callback Mgr    │    │ Structured Log  │
│                 │    │                 │    │                 │
│ • P2P Events    │    │ • Priorities    │    │ • Rotation      │
│ • Message Events│    │ • Timeouts      │    │ • Levels        │
│ • File Events   │    │ • Retries       │    │ • Components    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🛠️ Installation & Usage

### Prerequisites
- Go 1.21 or higher
- Linux/macOS/Windows support

### Build from Source
```bash
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
git checkout v0.4.0-alpha
./scripts/build.sh
```

### Quick Start
```bash
# Initialize your identity
./bin/peerchat-cli init

# Start interactive chat mode
./bin/peerchat-cli start

# Listen for incoming connections
./bin/peerchat-cli listen

# Discover peers on network
./bin/peerchat-cli discover
```

## 📊 Performance Improvements

- **Event Processing**: Up to 10,000 events/second with configurable worker pools
- **Memory Usage**: Optimized event queuing with configurable buffer sizes
- **Log Performance**: Structured logging with minimal allocation overhead
- **Startup Time**: Improved CLI startup with dynamic version loading

## 🧪 Testing

This release includes comprehensive testing:

```bash
# Run all tests
go test ./tests/unit/ -v

# Run event system tests specifically
go test ./tests/unit/ -run TestEvent -v
go test ./tests/unit/ -run TestCallback -v

# Run with coverage
go test ./tests/unit/ -cover
```

## 🔄 Migration from v0.3.0-alpha

The event system is backward compatible. Existing functionality continues to work while new event-driven features are available for integration.

### Breaking Changes
- None - this is a feature-additive release

### Deprecated Features
- Direct polling mechanisms (replaced by event-driven alternatives)

## 🐛 Known Issues

- Prometheus metrics integration not yet implemented (planned for v0.4.1-alpha)
- GUI integration pending (planned for Epoch 2)

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](../../docs/CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the AGPLv3 License - see the [LICENSE](../../LICENSE) file for details.

## 🔗 Links

- **GitHub Repository**: https://github.com/Xelvra/peerchat
- **Documentation**: https://github.com/Xelvra/peerchat/wiki
- **Issues**: https://github.com/Xelvra/peerchat/issues
- **Releases**: https://github.com/Xelvra/peerchat/releases

## 📈 Roadmap

- **v0.4.1-alpha**: Prometheus metrics integration
- **v0.5.0-alpha**: Security enhancements and onion routing
- **Epoch 2**: API server implementation
- **Epoch 3**: GUI client development

---

**Built with ❤️ by the Xelvra team**  
*#XelvraFree - Decentralized messaging for everyone*
