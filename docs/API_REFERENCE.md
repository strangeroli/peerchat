# Xelvra P2P Messenger - API Reference

## Table of Contents
- [CLI Commands](#cli-commands)
- [P2P Wrapper API](#p2p-wrapper-api)
- [Message Manager API](#message-manager-api)
- [Discovery Manager API](#discovery-manager-api)
- [Identity Manager API](#identity-manager-api)
- [Configuration API](#configuration-api)

## CLI Commands

### Node Management

#### `peerchat-cli start`
Start interactive P2P chat mode.

**Usage:**
```bash
peerchat-cli start [--verbose]
```

**Options:**
- `--verbose, -v`: Enable verbose logging

**Interactive Commands:**
- `/help` - Show available commands
- `/peers` - List connected peers
- `/discover` - Discover peers in network
- `/connect <peer_id>` - Connect to specific peer
- `/status` - Show node status
- `/quit` - Exit chat

#### `peerchat-cli listen`
Start passive listening mode for debugging.

**Usage:**
```bash
peerchat-cli listen [--verbose]
```

**Description:**
Displays all logs and network activity on console. No interaction available.

#### `peerchat-cli status`
Show detailed node status and network information.

**Usage:**
```bash
peerchat-cli status
```

**Output:**
- Peer ID and DID
- Listen addresses
- Connected peers count
- Discovery status
- Network information

### Communication

#### `peerchat-cli send`
Send message to specific peer (CLI mode).

**Usage:**
```bash
peerchat-cli send <peer_id> <message>
```

**Arguments:**
- `peer_id`: Target peer ID or multiaddr
- `message`: Text message to send

**Example:**
```bash
peerchat-cli send 12D3KooW... "Hello, World!"
```

#### `peerchat-cli send-file`
Send file to specific peer.

**Usage:**
```bash
peerchat-cli send-file <peer_id> <file_path>
```

**Arguments:**
- `peer_id`: Target peer ID
- `file_path`: Path to file to send

### Discovery

#### `peerchat-cli discover`
Discover peers in the network.

**Usage:**
```bash
peerchat-cli discover [--timeout=10s]
```

**Options:**
- `--timeout`: Discovery timeout duration

### Identity

#### `peerchat-cli init`
Initialize new Xelvra identity and configuration.

**Usage:**
```bash
peerchat-cli init [--force]
```

**Options:**
- `--force`: Overwrite existing identity

#### `peerchat-cli id`
Show your identity information.

**Usage:**
```bash
peerchat-cli id
```

**Output:**
- DID (Decentralized Identifier)
- Peer ID
- Public key
- Listen addresses

### Utilities

#### `peerchat-cli doctor`
Diagnose network and configuration issues.

**Usage:**
```bash
peerchat-cli doctor
```

**Checks:**
- System requirements
- Network connectivity
- P2P node creation
- Firewall configuration

#### `peerchat-cli version`
Show version information.

**Usage:**
```bash
peerchat-cli version
```

#### `peerchat-cli manual`
Show detailed usage manual.

**Usage:**
```bash
peerchat-cli manual
```

## P2P Wrapper API

### Constructor

#### `NewP2PWrapper(ctx context.Context, useSimulation bool) *P2PWrapper`
Create new P2P wrapper instance.

**Parameters:**
- `ctx`: Context for cancellation
- `useSimulation`: Force simulation mode if true

**Returns:**
- `*P2PWrapper`: New wrapper instance

### Methods

#### `Start() error`
Start the P2P node.

**Returns:**
- `error`: Error if startup fails

#### `Stop()`
Stop the P2P node and cleanup resources.

#### `GetNodeInfo() *NodeInfo`
Get current node information.

**Returns:**
- `*NodeInfo`: Node status and addresses

#### `IsUsingSimulation() bool`
Check if running in simulation mode.

**Returns:**
- `bool`: True if using simulation

#### `GetDiscoveredPeers() []string`
Get list of discovered peer IDs.

**Returns:**
- `[]string`: Array of peer ID strings

#### `GetConnectedPeers() []string`
Get list of currently connected peer IDs.

**Returns:**
- `[]string`: Array of connected peer ID strings

#### `ConnectToPeer(peerID string) bool`
Attempt to connect to specific peer.

**Parameters:**
- `peerID`: Target peer ID string

**Returns:**
- `bool`: True if connection successful

#### `SendMessage(peerID, message string) error`
Send message to specific peer.

**Parameters:**
- `peerID`: Target peer ID
- `message`: Message text

**Returns:**
- `error`: Error if sending fails

#### `SendMessageToMultiplePeers(message string, peerIDs []string) bool`
Send message to multiple peers.

**Parameters:**
- `message`: Message text
- `peerIDs`: Array of target peer IDs

**Returns:**
- `bool`: True if all sends successful

## Message Manager API

### Types

#### `MessageType`
```go
type MessageType int

const (
    MessageTypeText MessageType = iota
    MessageTypeFile
    MessageTypeImage
    MessageTypeAudio
    MessageTypeVideo
)
```

#### `Message`
```go
type Message struct {
    ID        string      `json:"id"`
    Type      MessageType `json:"type"`
    From      string      `json:"from"`
    To        string      `json:"to"`
    Content   []byte      `json:"content"`
    Timestamp time.Time   `json:"timestamp"`
    Signature []byte      `json:"signature"`
}
```

### Methods

#### `NewMessageManager(host host.Host, identity *user.Identity) *MessageManager`
Create new message manager.

**Parameters:**
- `host`: libp2p host instance
- `identity`: User identity for signing

**Returns:**
- `*MessageManager`: New message manager

#### `SendTextMessage(peerID peer.ID, text string) error`
Send text message to peer.

**Parameters:**
- `peerID`: Target peer ID
- `text`: Message text

**Returns:**
- `error`: Error if sending fails

#### `SetMessageHandler(handler func(*Message))`
Set handler for incoming messages.

**Parameters:**
- `handler`: Function to handle incoming messages

## Discovery Manager API

### Methods

#### `NewDiscoveryManager(host host.Host, logger *logrus.Logger) *DiscoveryManager`
Create new discovery manager.

**Parameters:**
- `host`: libp2p host instance
- `logger`: Logger instance

**Returns:**
- `*DiscoveryManager`: New discovery manager

#### `Start() error`
Start peer discovery services.

**Returns:**
- `error`: Error if startup fails

#### `Stop() error`
Stop peer discovery services.

**Returns:**
- `error`: Error if shutdown fails

#### `GetDiscoveredPeers() []peer.ID`
Get list of discovered peers.

**Returns:**
- `[]peer.ID`: Array of discovered peer IDs

#### `GetPeerAddresses(peerID peer.ID) []multiaddr.Multiaddr`
Get addresses for specific peer.

**Parameters:**
- `peerID`: Target peer ID

**Returns:**
- `[]multiaddr.Multiaddr`: Array of peer addresses

#### `GetStatus() *DiscoveryStatus`
Get current discovery status.

**Returns:**
- `*DiscoveryStatus`: Discovery status information

## Identity Manager API

### Types

#### `Identity`
```go
type Identity struct {
    DID        string
    PrivateKey crypto.PrivKey
    PublicKey  crypto.PubKey
}
```

### Methods

#### `NewIdentity() (*Identity, error)`
Create new cryptographic identity.

**Returns:**
- `*Identity`: New identity instance
- `error`: Error if creation fails

#### `LoadIdentity(configDir string) (*Identity, error)`
Load existing identity from configuration.

**Parameters:**
- `configDir`: Configuration directory path

**Returns:**
- `*Identity`: Loaded identity
- `error`: Error if loading fails

#### `SaveIdentity(configDir string) error`
Save identity to configuration directory.

**Parameters:**
- `configDir`: Configuration directory path

**Returns:**
- `error`: Error if saving fails

#### `GetDID() string`
Get Decentralized Identifier.

**Returns:**
- `string`: DID string

#### `Sign(data []byte) ([]byte, error)`
Sign data with private key.

**Parameters:**
- `data`: Data to sign

**Returns:**
- `[]byte`: Signature bytes
- `error`: Error if signing fails

#### `Verify(data, signature []byte, publicKey crypto.PubKey) bool`
Verify signature against data.

**Parameters:**
- `data`: Original data
- `signature`: Signature to verify
- `publicKey`: Public key for verification

**Returns:**
- `bool`: True if signature valid

## Configuration API

### Types

#### `NodeConfig`
```go
type NodeConfig struct {
    ListenAddrs    []string
    BootstrapPeers []string
    LogLevel       logrus.Level
    Logger         *logrus.Logger
    DataDir        string
}
```

### Methods

#### `DefaultNodeConfig() *NodeConfig`
Get default node configuration.

**Returns:**
- `*NodeConfig`: Default configuration

#### `LoadConfig(configPath string) (*NodeConfig, error)`
Load configuration from file.

**Parameters:**
- `configPath`: Path to configuration file

**Returns:**
- `*NodeConfig`: Loaded configuration
- `error`: Error if loading fails

#### `SaveConfig(config *NodeConfig, configPath string) error`
Save configuration to file.

**Parameters:**
- `config`: Configuration to save
- `configPath`: Target file path

**Returns:**
- `error`: Error if saving fails

---

For usage examples, see the [User Guide](USER_GUIDE.md).
For development information, see the [Developer Guide](DEVELOPER_GUIDE.md).
