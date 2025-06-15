# API Reference

Complete reference for Xelvra P2P Messenger APIs, including CLI commands, internal interfaces, and planned gRPC services.

## 📋 Table of Contents

- [CLI API](#cli-api)
- [Internal Go APIs](#internal-go-apis)
- [Configuration API](#configuration-api)
- [Future gRPC API](#future-grpc-api)
- [Protocol Specifications](#protocol-specifications)

## 🖥️ CLI API

The command-line interface provides the primary user API for Xelvra.

### Global Options

Available for all commands:

```bash
--config-dir PATH     # Configuration directory (default: ~/.xelvra)
--log-level LEVEL     # Logging level: debug, info, warn, error
--help               # Show help information
--version            # Show version information
```

### Core Commands

#### `init`
Initialize user identity and configuration.

```bash
peerchat-cli init [OPTIONS]
```

**Options:**
- `--config-dir PATH` - Custom configuration directory
- `--reset` - Reset existing configuration
- `--import FILE` - Import identity from file

**Examples:**
```bash
# Basic initialization
peerchat-cli init

# Initialize with custom directory
peerchat-cli init --config-dir ~/.xelvra-work

# Reset existing configuration
peerchat-cli init --reset
```

**Exit Codes:**
- `0` - Success
- `1` - Configuration error
- `2` - File system error

#### `start`
Start interactive chat mode.

```bash
peerchat-cli start [OPTIONS]
```

**Options:**
- `--port PORT` - Listening port (0 for auto-select)
- `--interface INTERFACE` - Network interface to bind
- `--no-discovery` - Disable peer discovery
- `--relay-only` - Use only relay connections

**Examples:**
```bash
# Start with default settings
peerchat-cli start

# Start on specific port and interface
peerchat-cli start --port 8080 --interface eth0

# Start without discovery
peerchat-cli start --no-discovery
```

#### `listen`
Start passive listening mode for debugging.

```bash
peerchat-cli listen [OPTIONS]
```

**Options:**
- `--verbose` - Show detailed debug information
- `--filter LEVEL` - Filter log messages by level
- `--output FILE` - Write logs to file

**Examples:**
```bash
# Basic listening mode
peerchat-cli listen

# Verbose mode with file output
peerchat-cli listen --verbose --output debug.log
```

#### `discover`
Discover peers on the network.

```bash
peerchat-cli discover [OPTIONS]
```

**Options:**
- `--timeout SECONDS` - Discovery timeout (default: 10)
- `--method METHOD` - Discovery method: mdns, udp, dht, all
- `--format FORMAT` - Output format: text, json, yaml

**Examples:**
```bash
# Basic discovery
peerchat-cli discover

# Extended discovery with JSON output
peerchat-cli discover --timeout 30 --format json

# Use specific discovery method
peerchat-cli discover --method mdns
```

**Output Format:**
```json
{
  "peers": [
    {
      "id": "12D3KooWExample...",
      "addresses": ["/ip4/192.168.1.100/tcp/4001"],
      "protocols": ["/xelvra/1.0.0"],
      "discovered_via": "mdns",
      "timestamp": "2024-01-15T10:30:00Z"
    }
  ],
  "total": 1,
  "duration": "10.5s"
}
```

#### `connect`
Connect to a specific peer.

```bash
peerchat-cli connect <PEER_ID> [OPTIONS]
```

**Arguments:**
- `PEER_ID` - Target peer ID or multiaddress

**Options:**
- `--timeout SECONDS` - Connection timeout (default: 30)
- `--retry COUNT` - Number of retry attempts
- `--force` - Force connection even if already connected

**Examples:**
```bash
# Connect to peer by ID
peerchat-cli connect 12D3KooWExample...

# Connect with custom timeout
peerchat-cli connect 12D3KooWExample... --timeout 60

# Connect using multiaddress
peerchat-cli connect /ip4/192.168.1.100/tcp/4001/p2p/12D3KooWExample...
```

#### `send`
Send a message to a peer.

```bash
peerchat-cli send <PEER_ID> <MESSAGE> [OPTIONS]
```

**Arguments:**
- `PEER_ID` - Target peer ID
- `MESSAGE` - Message text to send

**Options:**
- `--encrypt` - Force encryption (default: true)
- `--priority LEVEL` - Message priority: low, normal, high
- `--timeout SECONDS` - Send timeout

**Examples:**
```bash
# Send basic message
peerchat-cli send 12D3KooWExample... "Hello, world!"

# Send high-priority message
peerchat-cli send 12D3KooWExample... "Urgent message" --priority high
```

#### `send-file`
Send a file to a peer.

```bash
peerchat-cli send-file <PEER_ID> <FILE_PATH> [OPTIONS]
```

**Arguments:**
- `PEER_ID` - Target peer ID
- `FILE_PATH` - Path to file to send

**Options:**
- `--chunk-size SIZE` - Chunk size in bytes (default: 64KB)
- `--compress` - Compress file before sending
- `--verify` - Verify file integrity after transfer

**Examples:**
```bash
# Send file
peerchat-cli send-file 12D3KooWExample... document.pdf

# Send with compression and verification
peerchat-cli send-file 12D3KooWExample... large-file.zip --compress --verify
```

#### `status`
Display node status and connection information.

```bash
peerchat-cli status [OPTIONS]
```

**Options:**
- `--format FORMAT` - Output format: text, json, yaml
- `--watch` - Continuously update status
- `--interval SECONDS` - Update interval for watch mode

**Examples:**
```bash
# Basic status
peerchat-cli status

# JSON format with continuous updates
peerchat-cli status --format json --watch --interval 5
```

**Output Format:**
```json
{
  "node": {
    "id": "12D3KooWExample...",
    "addresses": ["/ip4/192.168.1.50/tcp/4001"],
    "protocols": ["/xelvra/1.0.0"],
    "uptime": "2h30m15s"
  },
  "connections": {
    "total": 3,
    "inbound": 1,
    "outbound": 2,
    "peers": [
      {
        "id": "12D3KooWPeer1...",
        "address": "/ip4/192.168.1.100/tcp/4001",
        "direction": "outbound",
        "duration": "1h45m30s"
      }
    ]
  },
  "network": {
    "nat_status": "public",
    "relay_connections": 0,
    "discovery_active": true
  }
}
```

#### `doctor`
Run system diagnostics and health checks.

```bash
peerchat-cli doctor [OPTIONS]
```

**Options:**
- `--fix` - Attempt to fix common issues
- `--detailed` - Show detailed diagnostic information
- `--test COMPONENT` - Test specific component: network, crypto, storage

**Examples:**
```bash
# Basic diagnostics
peerchat-cli doctor

# Detailed diagnostics with auto-fix
peerchat-cli doctor --detailed --fix

# Test specific component
peerchat-cli doctor --test network
```

#### `id`
Display identity information.

```bash
peerchat-cli id [OPTIONS]
```

**Options:**
- `--format FORMAT` - Output format: text, json, qr
- `--export FILE` - Export identity to file
- `--public-only` - Show only public information

**Examples:**
```bash
# Show identity
peerchat-cli id

# Export as QR code
peerchat-cli id --format qr

# Export to file
peerchat-cli id --export identity.json --public-only
```

### Interactive Commands

Available in interactive chat mode (`peerchat-cli start`):

#### Chat Commands
- `/help` - Show available commands
- `/peers` - List connected peers
- `/discover` - Discover new peers
- `/connect <peer_id>` - Connect to peer
- `/disconnect <peer_id>` - Disconnect from peer
- `/status` - Show node status
- `/clear` - Clear chat screen
- `/quit` or `/exit` - Exit chat mode

#### File Commands
- `/send-file <peer_id> <path>` - Send file to peer
- `/downloads` - Show download directory
- `/accept-file <id>` - Accept incoming file transfer
- `/reject-file <id>` - Reject incoming file transfer

#### Settings Commands
- `/set <key> <value>` - Change setting
- `/get <key>` - Show setting value
- `/settings` - Show all settings

## 🔧 Internal Go APIs

### P2P Node Interface

```go
package p2p

// Node represents a P2P node
type Node interface {
    // Start the node
    Start(ctx context.Context) error
    
    // Stop the node
    Stop() error
    
    // Connect to a peer
    Connect(ctx context.Context, peerID peer.ID) error
    
    // Disconnect from a peer
    Disconnect(peerID peer.ID) error
    
    // Send message to peer
    SendMessage(peerID peer.ID, msg []byte) error
    
    // Discover peers
    Discover(ctx context.Context) ([]peer.AddrInfo, error)
    
    // Get node status
    Status() NodeStatus
    
    // Get connected peers
    Peers() []peer.ID
}

// NodeStatus represents node status information
type NodeStatus struct {
    ID          peer.ID
    Addresses   []multiaddr.Multiaddr
    Connections int
    Uptime      time.Duration
    NATStatus   NATStatus
}
```

### Message Interface

```go
package message

// Message represents a P2P message
type Message interface {
    ID() string
    Sender() peer.ID
    Recipient() peer.ID
    Content() []byte
    Timestamp() time.Time
    Type() MessageType
    
    // Encryption methods
    Encrypt(key []byte) error
    Decrypt(key []byte) error
    IsEncrypted() bool
}

// MessageHandler handles incoming messages
type MessageHandler interface {
    HandleMessage(ctx context.Context, msg Message) error
    HandleFileTransfer(ctx context.Context, transfer FileTransfer) error
}
```

### Crypto Interface

```go
package crypto

// Identity represents a user identity
type Identity struct {
    DID        string
    PeerID     peer.ID
    PrivateKey crypto.PrivKey
    PublicKey  crypto.PubKey
    Created    time.Time
}

// Encryptor handles message encryption
type Encryptor interface {
    Encrypt(plaintext []byte, recipientKey crypto.PubKey) ([]byte, error)
    Decrypt(ciphertext []byte, senderKey crypto.PubKey) ([]byte, error)
    GenerateKeyPair() (crypto.PrivKey, crypto.PubKey, error)
}
```

## ⚙️ Configuration API

### Configuration Structure

```yaml
# ~/.xelvra/config.yaml
network:
  listen_port: 0              # 0 = auto-select
  discovery_port: 42424       # UDP discovery port
  enable_mdns: true           # mDNS discovery
  enable_udp_broadcast: true  # UDP broadcast discovery
  max_peers: 50              # Maximum connections
  connection_timeout: 30s     # Connection timeout
  nat_traversal: true        # Enable NAT traversal

user:
  display_name: ""           # Display name
  auto_accept_files: false   # Auto-accept file transfers
  download_directory: "downloads"

security:
  key_rotation_days: 60      # Key rotation interval
  max_message_size: 1048576  # 1MB message limit
  enable_forward_secrecy: true

logging:
  level: "info"              # debug, info, warn, error
  file: "peerchat.log"       # Log file name
  max_size: 10               # Max log file size (MB)
  max_backups: 3             # Number of backup files
  max_age: 30                # Max age in days
  compress: true             # Compress old logs
```

### Environment Variables

```bash
# Configuration
XELVRA_CONFIG_DIR="~/.xelvra"
XELVRA_LOG_LEVEL="info"
XELVRA_LISTEN_PORT="0"
XELVRA_DISCOVERY_PORT="42424"

# Network
XELVRA_MAX_PEERS="50"
XELVRA_CONNECTION_TIMEOUT="30s"
XELVRA_ENABLE_NAT_TRAVERSAL="true"

# Security
XELVRA_KEY_ROTATION_DAYS="60"
XELVRA_MAX_MESSAGE_SIZE="1048576"
```

## 🌐 Future gRPC API

Planned for Epoch 2 (API Service):

### Service Definition

```protobuf
syntax = "proto3";

package xelvra.v1;

service PeerChatService {
  // Node management
  rpc StartNode(StartNodeRequest) returns (StartNodeResponse);
  rpc StopNode(StopNodeRequest) returns (StopNodeResponse);
  rpc GetNodeStatus(GetNodeStatusRequest) returns (NodeStatus);
  
  // Peer management
  rpc DiscoverPeers(DiscoverPeersRequest) returns (stream PeerInfo);
  rpc ConnectToPeer(ConnectToPeerRequest) returns (ConnectToPeerResponse);
  rpc DisconnectFromPeer(DisconnectFromPeerRequest) returns (DisconnectFromPeerResponse);
  rpc ListPeers(ListPeersRequest) returns (ListPeersResponse);
  
  // Messaging
  rpc SendMessage(SendMessageRequest) returns (SendMessageResponse);
  rpc ReceiveMessages(ReceiveMessagesRequest) returns (stream Message);
  
  // File transfer
  rpc SendFile(stream SendFileRequest) returns (SendFileResponse);
  rpc ReceiveFile(ReceiveFileRequest) returns (stream ReceiveFileResponse);
  
  // Configuration
  rpc GetConfig(GetConfigRequest) returns (Config);
  rpc UpdateConfig(UpdateConfigRequest) returns (UpdateConfigResponse);
}
```

### Message Types

```protobuf
message NodeStatus {
  string id = 1;
  repeated string addresses = 2;
  int32 connection_count = 3;
  int64 uptime_seconds = 4;
  NATStatus nat_status = 5;
}

message PeerInfo {
  string id = 1;
  repeated string addresses = 2;
  repeated string protocols = 3;
  string discovered_via = 4;
  int64 timestamp = 5;
}

message Message {
  string id = 1;
  string sender_id = 2;
  string recipient_id = 3;
  bytes content = 4;
  int64 timestamp = 5;
  MessageType type = 6;
  bool encrypted = 7;
}
```

## 📡 Protocol Specifications

### Wire Protocol

Xelvra uses libp2p protocols for communication:

#### Protocol IDs
- `/xelvra/chat/1.0.0` - Chat messages
- `/xelvra/file/1.0.0` - File transfers
- `/xelvra/discovery/1.0.0` - Peer discovery
- `/xelvra/relay/1.0.0` - Relay connections

#### Message Format

```protobuf
message XelvraMessage {
  MessageHeader header = 1;
  oneof payload {
    ChatMessage chat = 2;
    FileChunk file_chunk = 3;
    DiscoveryMessage discovery = 4;
    StatusMessage status = 5;
  }
}

message MessageHeader {
  string message_id = 1;
  string sender_id = 2;
  string recipient_id = 3;
  int64 timestamp = 4;
  MessageType type = 5;
  bool encrypted = 6;
}
```

### Discovery Protocol

#### mDNS Discovery
- Service name: `_xelvra._tcp.local`
- TXT records contain peer information
- Automatic advertisement and discovery

#### UDP Broadcast Discovery
- Port: 42424 (configurable)
- Broadcast interval: 30 seconds
- Discovery message format: JSON

```json
{
  "version": "1.0.0",
  "peer_id": "12D3KooW...",
  "addresses": ["/ip4/192.168.1.100/tcp/4001"],
  "protocols": ["/xelvra/chat/1.0.0"],
  "timestamp": 1642248600
}
```

---

**For more technical details, see the [Developer Guide](Developer-Guide) and [Protocol Specifications](Protocol-Specifications).** 🔧
