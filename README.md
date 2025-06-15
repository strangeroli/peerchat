# Messenger Xelvra: *#XelvraFree*

> 🚀 **Secure, decentralized P2P communication platform. Built on E2E encryption with AI-driven net prediction.**

**Messenger Xelvra** is a peer-to-peer (P2P) communication platform designed to restore privacy, security, and user control over digital communication. The project aims to create a secure, efficient, and decentralized platform that pushes the boundaries of P2P communication capabilities.

## 📢 Project Resources

[![GitHub Issues](https://img.shields.io/github/issues/Xelvra/peerchat)](https://github.com/Xelvra/peerchat/issues)
[![GitHub Wiki](https://img.shields.io/badge/GitHub-Wiki-blue)](https://github.com/Xelvra/peerchat/wiki)
[![GitHub Discussions](https://img.shields.io/github/discussions/Xelvra/peerchat)](https://github.com/Xelvra/peerchat/discussions)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

- 📖 **[Project Wiki](https://github.com/Xelvra/peerchat/wiki)** - Comprehensive documentation, tutorials, and guides
- 🐛 **[Issues](https://github.com/Xelvra/peerchat/issues)** - Bug reports, feature requests, and project tracking
- 💬 **[Discussions](https://github.com/Xelvra/peerchat/discussions)** - Community discussions, Q&A, and ideas
- 🔧 **[Releases](https://github.com/Xelvra/peerchat/releases)** - Download latest versions and release notes

## 📋 Table of Contents

- [About the Project](#i-vision)
- [Why Messenger Xelvra?](#why-messenger-xelvra)
- [Key Features](#key-features)
- [Development Epochs](#development-epochs)
- [Architecture](#ii-architecture)
- [Security](#iii-security)
- [Installation](#quick-start)
- [Usage](#usage)
- [Roadmap](#roadmap)
- [Troubleshooting](#troubleshooting)
- [Documentation](#-documentation)
- [Contributing](#ix-how-to-contribute)
- [License](#xi-licensing)

## 📚 Documentation

### 📖 Quick Access
| Document | Description | Location |
|----------|-------------|----------|
| [📖 User Guide](docs/USER_GUIDE.md) | Complete guide for end users | Repository |
| [🔧 Installation Guide](docs/INSTALLATION.md) | Platform-specific installation instructions | Repository |
| [👨‍💻 Developer Guide](docs/DEVELOPER_GUIDE.md) | Development setup and contribution guide | Repository |
| [📋 API Reference](docs/API_REFERENCE.md) | Complete API documentation | Repository |

### 🌐 GitHub Wiki
For comprehensive documentation, tutorials, and community-contributed guides, visit our **[GitHub Wiki](https://github.com/Xelvra/peerchat/wiki)**:

- **[Getting Started](https://github.com/Xelvra/peerchat/wiki/Getting-Started)** - Quick start guide for new users
- **[Installation Guide](https://github.com/Xelvra/peerchat/wiki/Installation)** - Detailed installation instructions for all platforms
- **[User Manual](https://github.com/Xelvra/peerchat/wiki/User-Manual)** - Complete user documentation
- **[Developer Documentation](https://github.com/Xelvra/peerchat/wiki/Developer-Guide)** - Technical documentation for contributors
- **[API Documentation](https://github.com/Xelvra/peerchat/wiki/API-Reference)** - Complete API reference
- **[Troubleshooting](https://github.com/Xelvra/peerchat/wiki/Troubleshooting)** - Common issues and solutions
- **[FAQ](https://github.com/Xelvra/peerchat/wiki/FAQ)** - Frequently asked questions

## Why Messenger Xelvra?

In today's digital landscape, centralized communication platforms have become the norm, but they come with significant drawbacks that threaten our fundamental rights to privacy and freedom of communication.

### The Problems We Solve

**🔒 Privacy Concerns**: Traditional messaging platforms often collect, analyze, and monetize personal data. Conversations, contacts, and communication patterns may become products in various business models.

**🏢 Centralized Architecture**: Messages typically pass through centralized servers, creating potential vulnerabilities to service interruptions, data breaches, and external pressures. Single points of failure can affect millions of users.

**🚫 Limited Communication Freedom**: Centralized platforms may face various pressures that could impact communication freedom, content availability, and user autonomy in different jurisdictions.

**💰 Data as Currency**: Private conversations often become part of business models focused on targeted advertising and behavioral analysis, where user privacy becomes a tradeable commodity.

### Our Solution: True Digital Freedom

Messenger Xelvra addresses these fundamental issues by providing:

- **Direct P2P Communication**: Your messages travel directly between devices without intermediaries
- **End-to-End Encryption**: Only you and your intended recipient can read your messages
- **Decentralized Architecture**: No single point of failure or control
- **User Data Ownership**: Your data stays on your devices, under your control
- **Censorship Resistance**: No central authority can block or monitor your communications
- **Open Source Transparency**: Every line of code is open for inspection and verification

This isn't just about better technology—it's about restoring the fundamental human right to private communication and digital freedom.

## Key Features

### 🔐 Security & Privacy
- **Signal Protocol Implementation**: Industry-standard end-to-end encryption with X3DH handshake and Double Ratchet
- **Forward Secrecy**: Automatic key rotation ensures past communications remain secure even if current keys are compromised
- **Metadata Protection**: Onion routing obfuscates communication patterns and network analysis
- **Zero-Knowledge Architecture**: No personal data stored on external servers

### 🌐 Decentralized Network
- **Hybrid P2P Model**: Direct peer connections with intelligent relay fallback
- **Multiple Discovery Methods**: Kademlia DHT, mDNS, UDP broadcast, and mesh networking
- **NAT Traversal**: Automatic hole-punching, STUN/TURN integration, and port-knocking for restrictive networks
- **Offline Capability**: Local mesh networking via Bluetooth LE and Wi-Fi Direct

### ⚡ Performance & Efficiency
- **QUIC Transport**: Ultra-low latency communication with TCP fallback
- **Resource Optimization**: <20MB memory usage, <1% CPU in idle mode
- **Energy Efficient**: <15mW power consumption on mobile devices
- **AI-Driven Routing**: Machine learning optimization for intelligent transport selection

### 🛠️ Developer Friendly
- **Modular Architecture**: Clean separation between CLI, API, and GUI components
- **gRPC API**: Modern, efficient communication between components
- **Cross-Platform**: Support for Linux, macOS, Windows, Android, and iOS
- **Comprehensive Testing**: Unit, integration, and chaos engineering tests

## Development Epochs

Messenger Xelvra follows a structured development approach divided into distinct epochs, each building upon the previous foundation:

### 🏗️ Epoch 1: CLI Foundation (Current)
**Status: ✅ Largely Complete**

The command-line interface serves as the foundation and testing ground for all core P2P functionality:

- ✅ **P2P Core**: libp2p integration with QUIC/TCP transports
- ✅ **Discovery Systems**: mDNS, UDP broadcast, and DHT peer discovery
- ✅ **NAT Traversal**: STUN integration with automatic public IP detection
- ✅ **File Transfer**: Secure P2P file sharing with chunking and progress tracking
- ✅ **CLI Commands**: Complete command set (init, start, connect, send, discover, doctor)
- ✅ **Logging & Diagnostics**: Comprehensive logging with rotation and network diagnostics
- 🔄 **In Progress**: Interactive chat UI, advanced encryption, and mesh networking

### 🔌 Epoch 2: API Service (Planned)
**Status: 📋 Planned**

Local gRPC API service to bridge P2P core with frontend applications:

- **gRPC Server**: High-performance API with event-driven architecture
- **Database Layer**: SQLite with WAL mode for persistent storage
- **Monitoring**: Prometheus metrics and OpenTelemetry tracing
- **Rate Limiting**: Protection against API abuse
- **Stream Processing**: Real-time message and event streaming

### 📱 Epoch 3: GUI Application (Planned)
**Status: 📋 Planned**

Cross-platform Flutter application with focus on mobile optimization:

- **Modern UI/UX**: Material Design with accessibility compliance (WCAG 2.1 AA)
- **Progressive Onboarding**: Visual P2P education and interactive demos
- **Energy Optimization**: <100mW active usage, intelligent sleep modes
- **Multi-Platform**: Android, iOS, Linux, macOS, Windows support
- **Advanced Features**: Group chats, file sharing, voice calls

### 🚀 Epoch 4: Advanced Features (Future)
**Status: 🔮 Future Vision**

Advanced capabilities and ecosystem expansion:

- **Zero-Knowledge Proofs**: Enhanced privacy with ZKP identity verification
- **Quantum Resistance**: Post-quantum cryptography integration
- **Voice & Video**: Real-time multimedia communication
- **Mesh Networks**: Advanced offline communication capabilities
- **Community Features**: Decentralized governance and Hash Token ecosystem

## 🚀 Quick Start

```bash
# Clone and build
git clone https://github.com/Xelvra/peerchat.git
cd peerchat
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go

# Initialize and start
./bin/peerchat-cli init
./bin/peerchat-cli start
```

For detailed instructions, see the [Installation Guide](docs/INSTALLATION.md).

## Usage

### Basic Commands

```bash
# Initialize your identity
./bin/peerchat-cli init

# Start the P2P node
./bin/peerchat-cli start

# Check network status
./bin/peerchat-cli status

# Discover peers
./bin/peerchat-cli discover

# Connect to a peer
./bin/peerchat-cli connect <peer_multiaddr>

# Send a message
./bin/peerchat-cli send <peer_multiaddr> "Hello, World!"

# Send a file
./bin/peerchat-cli send-file <peer_multiaddr> /path/to/file

# Interactive chat mode
./bin/peerchat-cli start

# Listen for messages (debugging)
./bin/peerchat-cli listen

# Network diagnostics
./bin/peerchat-cli doctor

# View help
./bin/peerchat-cli manual
```

For comprehensive usage instructions, see the [User Guide](docs/USER_GUIDE.md).

## Roadmap

### 🎯 Short Term (Next 3-6 months)
- **Complete CLI Implementation**: Finish interactive chat UI and advanced encryption features
- **Enhanced Security**: Implement full Signal Protocol with automatic key rotation
- **Mesh Networking**: Add Bluetooth LE and Wi-Fi Direct support for offline communication
- **Performance Optimization**: Achieve target metrics for latency and resource usage

### 🚀 Medium Term (6-12 months)
- **API Service**: Complete gRPC API implementation with monitoring and telemetry
- **GUI Development**: Begin Flutter application development with focus on mobile platforms
- **Advanced NAT Traversal**: Implement AI-driven transport prediction and port-knocking
- **Community Building**: Establish contributor guidelines and community governance

### 🌟 Long Term (1-2 years)
- **Cross-Platform GUI**: Complete multi-platform application with full feature parity
- **Voice & Video**: Real-time multimedia communication capabilities
- **Quantum Resistance**: Post-quantum cryptography integration
- **Ecosystem Expansion**: Hash Token system and decentralized governance (DAO)

### 🔮 Future Vision (2+ years)
- **Zero-Knowledge Features**: Advanced privacy with ZKP identity verification
- **IoT Integration**: Extend P2P communication to IoT devices and embedded systems
- **Protocol Standardization**: Work toward industry standardization of decentralized messaging
- **Global Adoption**: Scale to support millions of users in a truly decentralized network

## I. Vision

### A. The Problem

Current centralized communication platforms threaten user privacy by collecting, analyzing, and monetizing their data. User messages pass through servers beyond their control, making them vulnerable to surveillance and censorship.

### B. The Solution

Messenger Xelvra addresses this problem by providing a platform for direct, uncensored, and independent communication. The platform is designed to return control over data to users and restore trust in digital communication with emphasis on **extreme speed, minimal resource consumption, and top-tier security and robustness.**

## II. Architecture

### A. Principles

Xelvra Messenger architecture is built on principles of privacy protection, security, and decentralization. Key principles include:

1.  **P2P Communication:** Direct, encrypted communication between users without intermediaries.
2.  **Hybrid P2P Model:** Strategic use of direct P2P connections with priority and relay services as fallback to ensure functionality, sustainability, and trust building, without compromising user privacy. Implementation of **parallel transports (QUIC + pre-initialized TCP)** for latency minimization and resilience maximization, and **automated ICE framework with AI-driven prediction.**
3.  **User Experience:** Multi-platform support, efficiency, and intuitive design with **aggressive optimization for low resource consumption, including Progressive Onboarding (visual P2P explanation and interactive demo with local network simulator) and full Accessibility (WCAG 2.1 AA compliance, screen reader support).**

### **B. Technical Architecture**

Messenger Xelvra is a modular system consisting of three main components:

1.  **peerchat-cli (Go):** Command-line tool for development and testing of P2P logic and for running P2P node in the background as a system service.
2.  **peerchat-api (Go):** Local API service (gRPC) for communication with frontend applications, providing robust interface to P2P core and utilizing **event-driven architecture.**
3.  **peerchat_gui (Flutter):** Multi-platform graphical user interface, optimized for mobile devices with emphasis on energy efficiency.

**Detailed description of Go modules (peerchat/internal/):**

* **p2p/:** P2P network management (go-libp2p, **single Kademlia DHT with local in-memory LRU caching layer and adaptive polling**, NAT traversal with **aggressive hole-punching, embedded STUN/TURN, AI-driven prediction and port-knocking tactics**, connection management with recovery mechanisms and pre-warmed TCP connections, **QUIC transport with kernel-level/user-space UDP batching, hardware acceleration and dynamic window scaling (BBR+Cubic), TCP fallback, Onion routing for *all* metadata with multiple encryption layers**, **Bluetooth LE/Wi-Fi Direct as fallback for mesh networks with smart power management**).
* **crypto/:** Implementation of encryption protocols (Signal Protocol, X3DH, Double Ratchet), secure key management with **Memory Hardening (mlock(), canaries, memguard)**, **protection against Replay/DoS attacks (advanced rate-limiting, token buckets)**, and **resistance to timing attacks.**
* **user/:** User identity management (**DID format did:xelvra:\<hash\>, verification with Ed25519 signatures (ZKP planned for Epoch 4\)**), **implementation of peer discovery by DID, robust user blocking (with encrypted blacklist in DHT), Sybil Resistance (dynamic Proof-of-Work, automatic trust system, "Ghost" contact limitations).**
* **message/:** Message and file management (transfer, offline messages, pub/sub, **complex group management with roles and invitations**). Optimized large file transfer using chunking.
* **api/:** Implementation of gRPC server and API handlers with **robust error handling, input validation and rate limiting. Includes monitoring for Prometheus/Grafana and distributed tracing with OpenTelemetry.**
* **db/:** Database operations abstraction (**SQLite with WAL mode for high performance, low latency and corruption resistance, one encrypted userdata.db file per user with automatic WAL file checkpointing**).
* **util/:** Helper functions (logging, metrics, validation).

### **B.1 Protocol Specifications**

* **Message Framing:** All messages and data packets will be structured using Google Protobuf for efficient serialization and deserialization, ensuring compactness and transfer speed.
* **Handshake sequence:** Detailed flow diagram for X3DH and Double Ratchet protocol will be available in separate documentation, describing the exact sequence of key exchange and encrypted session establishment.
* **Onion routing for metadata:** Implementation of layered encryption for metadata (e.g., peer IP addresses in DHT queries, timestamps) inspired by Onion/Garlic routing principles, to make network graph analysis and determination of real source/destination of communication more difficult for external observers.
* **Mesh Networking Protocol (Example)**:

    ```proto
    // pkg/proto/mesh.proto
    syntax = "proto3";

    package xelvra.mesh;

    message MeshPacket {
      bytes sender_id = 1;       // Node identifier
      bytes message_id = 2;      // Unique message ID
      uint32 hop_limit = 3;      // TTL for flooding (prevents infinite loops)
      
      oneof payload {
        bytes raw_payload = 4;   // Encrypted onion-routed content
        // Add other payload types here if needed (e.g. debug types)
      }
    }
    ```

* **Transport layer for BLE:** Will use GATT profile with MTU (Maximum Transmission Unit) 512B for efficient data transfer.
* **Wi-Fi Direct:** Activation only with sufficient battery level (>50%) due to its higher consumption.

## **III. Security**

### **A. Philosophy**

Security is a key principle of Xelvra Messenger. The platform protects users from various threats, including passive eavesdropping, active attacks and censorship with emphasis on **proactive defense, minimization of exposed information and resistance to advanced threats.**

### **B. Cryptographic Core**

* **E2EE:** End-to-end encryption of messages between sender and recipient using Signal Protocol.
* **Cryptographic primitives:** Standardized and proven algorithms (AES-256, Curve25519, SHA-256/SHA-3, HKDF), **optimized using hardware acceleration (AES-NI).**
* **Secure key management:** Generation, storage and management of keys with emphasis on **protection in memory (mlock(), canaries, memguard) and on disk (encrypted SQLite files).**
* **Key Rotation (Zero-Touch):** Automatic rotation of long-term keys every **60 days** with user notification **48 hours** before expiration. During "grace period" of 72 hours, parallel encryption with both old and new keys is performed to ensure lossless transition. Maintaining key history for decrypting older messages. This process will be fully automated and transparent to users.
* **Data integrity:** Digital signatures and hashing for verification of origin and integrity of messages.
* **Zero-Knowledge Proof:** Implementation of ZKP mechanisms for **identity verification without revealing sensitive information (planned for Epoch 4, currently Ed25519 signatures).**

### **C. Metadata Protection**

Minimization of metadata and decentralized user identification. **Onion routing for obfuscation of *all* metadata with the goal of making network graph analysis more difficult.**

### **C.6 Forward Secrecy**

* **Key rotation:** Automatic key rotation for Double Ratchet algorithm every 100 sent messages or after 24 hours of inactivity, minimizing the amount of data encrypted with one key.
* **Automatic invalidation:** Keys will be automatically invalidated after 7 days of inactivity in conversation, ensuring that old sessions will not pose long-term risk.

### **D. Network Resilience**

Protection against **Sybil attacks (with dynamic Proof-of-Work for new DHT records, automatic trust system and limitation of new contacts for "Ghost" users)**, **DoS attacks (with advanced rate-limiting and multi-level connection management)** and **Replay attacks (with timestamps and sequence numbers).**

### **C.7 Protection against Advanced Threats**

* **Sybil Resistance:**
    * **Proof-of-Work Requirement:** For adding new records to the DHT (e.g., new user identities), dynamic Proof-of-Work will be required, whose difficulty will change based on network load to make DDoS attacks (PoW flooding requests) more difficult.
    * **Contact Limitation for "Ghost" Users:** Users in "Ghost" status will have a limited number of new contacts they can initiate per 24 hours (e.g., 3/day) to prevent spamming.
    * **Automatic Trust:** New users can communicate with 5 contacts/day without CAPTCHA. After verification (e.g., QR code from an existing and trusted contact), limits will disappear.
* **Quantum Resistance:**
    * **Hybrid Encryption:** For long-term protection, current (e.g., X25519) and post-quantum (e.g., Kyber768) algorithms will be combined for the handshake phase and establishing shared secrets.
    * **Migration Capability:** The architecture will be designed to allow future migration to purely post-quantum cryptographic schemes once they are standardized and proven.

### **E. Transparency**

Open-source code and independent security audits.

## **IV. Ecosystem**

### **A. Hash Tokens (HT)**

Internal virtual credits for rewarding contributions and ensuring sustainability. HT have no financial value outside the Xelvra Messenger ecosystem.

### **B. Path of Trust**

A system for building trust and reputation among users. User statuses: Ghost, User, Architect, Ambassador, God.

### **C. Community Governance**

Long-term goal: Decentralized network governance by the community (DAO).

## **V. Business Model**

Transparent and community-oriented funding. Crowdfunding and HT sales to finance further development.

## **VI. Quantifiable Goals and Energy Optimization**

### **A. Quantifiable Performance and Resource Goals**

* **P2P message latency (one way):**
    * < 50 ms for direct connections.
    * < 200 ms through relay.
    * **Maximum latency under load:** < 100ms at 100 messages/s.
* **Memory consumption (CLI/Backend idle):** < 20 MB (Go runtime).
* **Memory limit during active use:** < 50MB (Go runtime).
* **CPU consumption (CLI/Backend idle):** < 1%.
* **API call latency (internal):** < 10 ms.
* **API throughput:** > 1000 RPC/s for basic operations.

### **B. Energy Optimization**

Energy consumption optimization for Go backend and Flutter frontend with **explicit goals for mobile devices:**

* **Energy consumption (mobile, idle, background):** < 15 mW.
* **Energy consumption (mobile, active chat):** < 100 mW.
* **Energy demand (mobile):** < 5% battery/hour during active chatting.
* Implementation of intelligent sleep and wake strategies using platform-specific mechanisms (WorkManager for Android, Background Fetch/VOIP Push for iOS).

### **VI.B Implementation Strategy**

* **Operation Batching:** Grouping smaller network requests or database writes into larger batches to reduce overhead and optimize energy consumption. **Including kernel-level QUIC batching for Linux.**
* **Adaptive Polling:** Dynamically adjusting the frequency of heartbeats, DHT queries, and other periodic network activities based on battery status (e.g., L3 queries once every 10 min when battery <20%) and user activity to minimize unnecessary background activity.
* **Event-Driven Architecture:** Replacing polling mechanisms for communication between GUI and API and within the P2P network with push notifications (gRPC streams, WebSockets) to reduce CPU and battery load.
* **Hardware Acceleration:** Utilizing AES-NI instructions (if available on hardware) and other specific instruction sets for cryptographic operations to increase performance and reduce CPU consumption.
* **Battery-Aware GC:** Static setting of GOGC=30 + ballast alloc (e.g., 1GB dummy array) for memory stability and reduction of GC cycle frequency.
* **Deep Sleep Mode:** At low battery levels (<15%), deactivate DHT and switch to "mesh-only" mode (only mDNS/Bluetooth LE/Wi-Fi Direct) for minimal consumption.
    * **Resolution of Conflict Scenarios:**
        * **Incoming Call:** Using "light push" notifications (e.g., high-priority FCM with minimal payload) for local wake-up of the P2P node. Expected consumption: ~0.2 mW.
        * **Important Message:** A message stored in the local mesh network (via BLE/Wi-Fi Direct) or DHT will be notified only when the node wakes up from Deep Sleep mode (e.g., regular synchronization window). Expected consumption: ~0.1 mW.
        * **System Updates:** Synchronization of database/application updates in defined time windows (e.g., every 6 hours) during the night or when connected to a charger. Expected consumption: ~0.3 mW during synchronization.
        * **Periodic Ping (BLE beaconing):** To maintain minimal connectivity and facilitate node wake-up, even with WiFi/Bluetooth off.

## **VII. Deployment and Operations**

Application deployment, distribution, and maintenance strategy.

### **VII. Deployment Strategy**

* **Automated Builds:** Fully automated CI/CD pipeline that includes compilation, testing, and signing of all binaries and packages with digital signatures (using Sigstore/cosign).
* **Reproducible Builds:** Ensuring that binaries are reproducible from a given source code, allowing independent parties to verify integrity and absence of unauthorized changes.
* **Delta Updates (for GUI):** Integration of **bsdiff** for efficient delivery of small network updates (<100KB), minimizing the size of downloaded data (critical for mesh networks).
* **Offline Update Support:** The ability to deliver updates even in a local mesh network without internet access, increasing system robustness and user autonomy.

## **VIII. Testing and Quality Assurance**

Thorough testing to ensure reliability and security, including:

* Unit tests
* Integration tests
* E2E tests
* Performance and load tests against defined metrics.
* **Chaos Engineering:** Simulation of network and node failures (e.g., random node crashes, packet loss, delays, simulation of internet outages in Docker tests) to verify system resilience under unpredictable conditions.
    * **Integration into CI pipeline/docker-compose.yml:**

        ```yaml
        # Example for docker-compose.yml for network chaos
        services:
          network-chaos:
            image: nicholasjackson/chaos-http
            command: -target p2p-network -latency 100ms -jitter 50ms -loss 10%
            # Add networks and dependencies to target p2p-network
        ```

* **Fuzzing:** Testing the robustness of protocol parsing and inputs by generating random, potentially malicious data (especially **for QUIC handshake protocols and Protobuf messages**).
    * **Tools:** go-fuzz for Go modules.

        ```bash
        # Example of go-fuzz usage
        go-fuzz -bin=./message-fuzzer.zip -workdir=/fuzz
        ```

* **Penetration Testing:** Utilizing external tools and techniques for penetration testing (e.g., Nmap, Metasploit, OWASP ZAP – for testing exposed API, if relevant, otherwise for the network layer) to uncover vulnerabilities.
    * **QUIC Handshake Penetration Testing:** Test using [QUIC-Intruder](https://github.com/vanhauser-thc/thc-quic-intruder).
    * **Side-channel Attacks:** Verify resistance to side-channel attacks using [CacheScout](https://github.com/cachescout/cachescout) (to verify AES-NI implementation and other cryptographic operations).
    * **Timing Attack Resistance:** Analyze and add artificial, constant delays in cryptographic operations (e.g., key comparisons) to prevent timing attacks.
* **Real-world Network Testing:**
    * **Public WiFi with Captive Portals:** Testing connectivity and P2P functionality.
    * **Restrictive Firewalls:** Verifying the ability to traverse restrictive firewalls (ports 80/443).
    * **Mobile Networks with Frequent Handover:** Testing connection and P2P network resilience.
* **Energy Profiling:** Integrate into CI pipeline:

    ```bash
    # For Linux backend
    perf stat -e power/energy-pkg/ ./peerchat-cli test --duration 5m
    # For Android frontend
    adb shell dumpsys batterystats --enable full-wake-history
    adb bugreport > bugreport.txt # Analyze in Battery Historian
    ```

## **IX. How to Contribute**

Information on how to get involved in the development and community of Xelvra Messenger. For easier initial setup and development, the `peerchat-cli setup` command and a **Docker-based testing environment** will be available.

## **X. Code of Conduct: Building a Respectful Community**

Xelvra is a community built on trust, openness, and collaboration. To ensure a safe, welcoming, and inclusive environment for all, we have established this Code of Conduct. It applies to all project participants.

### **A. Our Values**

* Respect, inclusivity, openness, collaboration, safety.

### **B. Expected Behavior**

* Be welcoming and patient.
* Use welcoming and inclusive language.
* Be considerate, provide constructive criticism.
* Respect differing opinions.
* Respect privacy and security.
* Take responsibility for your mistakes.

### **C. Unacceptable Behavior**

* Harassment, discrimination, personal attacks, trolling, malicious code, publishing private information, coercion/threats.

### **D. Enforcement of the Code**

Cases of Code of Conduct violations are handled fairly and transparently.

## **XI. Licensing**

Messenger Xelvra is licensed under **GNU Affero General Public License v3.0 (AGPLv3)**.

## **XII. Glossary**

* **Kademlia DHT:** A distributed hash table used for peer discovery and data storage. In Xelvra Messenger, it will be used with a local in-memory LRU caching layer.
* **HT (Hash Token):** Internal virtual credits in the Xelvra Messenger ecosystem, used to reward users for contributions to the network (e.g., relaying messages, maintaining a DHT node) and to ensure sustainability. HT have no financial value outside the ecosystem.
* **Progressive Onboarding:** A user-friendly first-run process that gradually explains P2P concepts and guides the user through setup, including visual network simulations.
* **Zero-Touch Encryption:** Automatic management of cryptographic keys without the need for manual user intervention, including automatic rotation, "grace period," and notifications.
* **Kernel-Level QUIC Batching:** Optimization of data transfer in the QUIC protocol, where small packets are grouped and sent directly from the operating system kernel (e.g., using SO_ZEROCOPY and io_uring on Linux) to reduce overhead and improve throughput. Fallback to user-space batching for other OS.
* **Deep Sleep Mode:** An energy-saving mode for mobile applications where network activities are minimized and DHT is deactivated in favor of local mesh communication (mDNS, Bluetooth LE/Wi-Fi Direct) for maximum battery savings, with solutions for wake-up conflicts (light push, BLE beaconing).
* **AI-Driven Prediction / AI-Based Routing:** Use of lightweight machine learning models (e.g., ONNX Runtime) to dynamically evaluate network conditions and select the most efficient transport protocol or message path in real-time, with input validation and model sandboxing.
* **SQLite with WAL (Write-Ahead Logging):** A database mode that improves performance and crash resistance, minimizes fragmentation, and enables efficient checkpointing.
* **Port-Knocking:** A technique for opening ports on a firewall by sending a predefined sequence of packets to closed ports, which increases difficulty for attackers.

## Troubleshooting

### Common Issues and Solutions

#### 🔌 Connection Problems

**Can't connect to other peers?**
- Ensure the application is running and properly initialized
- Check that your firewall isn't blocking the application
- Try the automatic network diagnostics: run the doctor command for detailed analysis
- If behind a restrictive firewall, the application will automatically use relay servers

**Slow message delivery?**
- The application automatically selects the fastest available transport method
- Poor network quality may cause delays - try switching to a different network
- Check your internet connection stability
- The system will adapt to network conditions automatically

#### 📱 Performance Issues

**High battery drain on mobile?**
- Enable power-saving mode in the application settings
- The app automatically reduces background activity when battery is low
- Close unnecessary background applications
- Ensure you're using the latest version with energy optimizations

**Application running slowly?**
- Check available system memory and close other applications
- The application is designed to use minimal resources
- Restart the application if performance issues persist
- Check system logs for any error messages

#### 🔐 Security and Identity

**Lost access to your identity?**
- If you have a backup of your keys, you can restore your identity
- Without a backup, you'll need to create a new identity (previous message history will be lost)
- Always backup your identity keys in a secure location
- Consider using the built-in backup features when available

**Messages not being received?**
- Ensure both you and the sender are online and connected
- The application automatically manages encryption keys
- If there's been a long period of inactivity, keys may need to be re-synchronized
- Check that the sender hasn't been accidentally blocked

#### 🌐 Network and Discovery

**Can't find other users?**
- Ensure you're on the same network for local discovery
- Check that network discovery services are enabled
- For internet-wide discovery, ensure you have a stable internet connection
- The application will automatically try multiple discovery methods

**Offline communication not working?**
- Ensure Bluetooth and Wi-Fi are enabled for local mesh networking
- Check that other users are within range and have compatible devices
- Local mesh features may have limited range and capabilities
- Some features require internet connectivity

### Getting Help

If you continue to experience issues:

1. **Check the logs**: Application logs contain detailed diagnostic information
2. **Run diagnostics**: Use the built-in diagnostic tools for automated troubleshooting
3. **Consult documentation**: See the [GitHub Wiki](https://github.com/Xelvra/peerchat/wiki) for comprehensive guides
4. **Community support**:
   - 💬 [GitHub Discussions](https://github.com/Xelvra/peerchat/discussions) - Ask questions and get community help
   - 🐛 [GitHub Issues](https://github.com/Xelvra/peerchat/issues) - Report bugs and request features
   - 📖 [Wiki FAQ](https://github.com/Xelvra/peerchat/wiki/FAQ) - Check frequently asked questions
5. **Report bugs**: Use our [issue templates](https://github.com/Xelvra/peerchat/issues/new/choose) for detailed bug reports

For technical details and advanced troubleshooting, see the [Troubleshooting Guide](docs/TROUBLESHOOTING.md) or [Wiki Troubleshooting](https://github.com/Xelvra/peerchat/wiki/Troubleshooting).
