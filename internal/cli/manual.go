package cli

import (
	"fmt"
)

// RunManual handles the manual command
func RunManual(version string) {
	fmt.Printf(`
XELVRA P2P MESSENGER CLI MANUAL
===============================

NAME
    peerchat-cli - Xelvra P2P Messenger Command Line Interface

SYNOPSIS
    peerchat-cli [GLOBAL OPTIONS] COMMAND [COMMAND OPTIONS] [ARGUMENTS...]

DESCRIPTION
    Xelvra is a decentralized peer-to-peer messenger that operates without
    central servers. It uses libp2p for networking, Ed25519 for identity,
    and implements end-to-end encryption for secure communication.

    The CLI provides both interactive chat mode and standalone commands
    for managing your P2P node, discovering peers, and sending messages.

GLOBAL OPTIONS
    --config FILE     Configuration file (default: ~/.xelvra/config.yaml)
    -v, --verbose     Enable verbose output and detailed logging
    -h, --help        Show help information
    --version         Show version information

COMMANDS

  SETUP & INITIALIZATION
    init              Initialize a new Xelvra identity and configuration
                      Creates Ed25519 key pair and DID identifier
                      Sets up ~/.xelvra/ directory with configuration

                      Example:
                        peerchat-cli init

  INTERACTIVE CHAT
    start             Start interactive P2P chat mode with full features
                      Supports tab completion, command history, and real-time messaging
                      Use --daemon flag to run as background service

                      Examples:
                        peerchat-cli start
                        peerchat-cli start --daemon

  NODE MANAGEMENT
    status            Show detailed node status and network information
                      Displays peer connections, NAT info, and discovery status

                      Example:
                        peerchat-cli status

    listen            Start node in passive listening mode (debugging)
                      Shows all logs and network activity in real-time
                      No interactive input - use for monitoring and debugging

                      Example:
                        peerchat-cli listen

    stop              Stop running P2P node (not yet implemented)
                      Will terminate background daemon processes

                      Example:
                        peerchat-cli stop

  PEER DISCOVERY & CONNECTION
    discover          Discover peers on the local network
                      Uses mDNS, UDP broadcast, and DHT for peer discovery
                      Shows real-time discovery progress and results

                      Example:
                        peerchat-cli discover

  MESSAGING
    send              Send a message to a specific peer
                      Requires peer ID and message text as arguments
                      Currently requires running node for delivery

                      Example:
                        peerchat-cli send 12D3KooW... "Hello, World!"

  FILE TRANSFER
    send-file         Send a file to a peer (not yet implemented)
                      Will support chunked, resumable file transfers
                      Includes integrity verification and progress tracking

                      Example:
                        peerchat-cli send-file 12D3KooW... /path/to/file.txt

  IDENTITY & PROFILES
    id                Show your identity information
                      Displays DID, Peer ID, and network addresses

                      Example:
                        peerchat-cli id

    profile           Show profile information for a peer (not yet implemented)
                      Will display peer metadata and connection history

                      Example:
                        peerchat-cli profile 12D3KooW...

  HELP & INFORMATION
    manual            Show this comprehensive manual
    version           Show version and build information
    help [COMMAND]    Show help for a specific command

                      Examples:
                        peerchat-cli help
                        peerchat-cli help send

  DIAGNOSTICS & TROUBLESHOOTING
    doctor            Run comprehensive network diagnostics
                      Tests P2P connectivity, NAT traversal, and discovery
                      Provides troubleshooting suggestions for common issues

                      Example:
                        peerchat-cli doctor

    setup             Interactive setup wizard (not yet implemented)
                      Will guide through initial configuration and testing

                      Example:
                        peerchat-cli setup

INTERACTIVE CHAT COMMANDS
    When in interactive mode (peerchat-cli start), these commands are available:

    /help             Show available interactive commands
    /peers            List currently connected peers
    /discover         Discover new peers on the network
    /connect <id>     Connect to a specific peer (with tab completion)
    /disconnect <id>  Disconnect from a peer
    /status           Show current node status
    /clear            Clear the screen
    /quit, /exit      Exit interactive chat mode

    Regular messages (not starting with /) are sent to all connected peers.

KEYBOARD SHORTCUTS (Interactive Mode)
    Tab               Auto-complete commands and peer IDs
    ↑/↓ arrows        Navigate command history
    Ctrl+R            Search command history
    Ctrl+C            Exit chat mode
    Ctrl+L            Clear screen (same as /clear)
    Ctrl+A            Move cursor to beginning of line
    Ctrl+E            Move cursor to end of line

FILES AND DIRECTORIES
    ~/.xelvra/                    Main configuration directory
    ~/.xelvra/config.yaml         Node configuration file
    ~/.xelvra/identity.key        Private key file (Ed25519)
    ~/.xelvra/peerchat.log        Application log file (rotated)
    ~/.xelvra/chat_history        Interactive chat command history
    ~/.xelvra/offline_messages/   Stored offline messages
    ~/.xelvra/downloads/          Received files directory

CONFIGURATION
    The configuration file (~/.xelvra/config.yaml) contains:
    - Identity settings (DID, key paths)
    - Network configuration (ports, bootstrap peers)
    - Discovery settings (mDNS, DHT, UDP broadcast)
    - Logging configuration (level, rotation)

NETWORK PROTOCOLS
    - Transport: QUIC (primary), TCP (fallback)
    - Discovery: mDNS, UDP broadcast, DHT
    - Encryption: Ed25519 signatures, planned E2E encryption
    - NAT Traversal: STUN, UPnP, relay servers

EXIT CODES
    0    Success
    1    General error
    2    Network error
    3    Configuration error
    4    Permission error
    5    Peer not found

REPORTING BUGS
    Report bugs at: https://github.com/Xelvra/peerchat/issues
    Include: version info, logs, and steps to reproduce

VERSION
    ` + version + `

COPYRIGHT
    Copyright (C) 2025 Xelvra Project
    Licensed under GNU Affero General Public License v3.0 (AGPLv3)

SEE ALSO
    Project documentation: https://github.com/Xelvra/peerchat
    libp2p documentation: https://docs.libp2p.io/
`)
}
