package user

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/mr-tron/base58"
)

const (
	// DID format: did:xelvra:<hash>
	DIDPrefix = "did:xelvra:"

	// Identity key sizes
	Ed25519PrivateKeySize = 64
	Ed25519PublicKeySize  = 32
	Ed25519SignatureSize  = 64
)

// MessengerID represents a decentralized identity in Xelvra format
type MessengerID struct {
	DID        string             // did:xelvra:<hash>
	PublicKey  ed25519.PublicKey  // Ed25519 public key
	PrivateKey ed25519.PrivateKey // Ed25519 private key (TODO: add memory protection)
	PeerID     peer.ID            // libp2p peer ID
	CreatedAt  time.Time          // Creation timestamp
}

// TrustLevel represents the trust level of a user in the network
type TrustLevel int

const (
	TrustLevelGhost      TrustLevel = iota // 0 - New user, limited privileges
	TrustLevelUser                         // 1 - Basic verified user
	TrustLevelArchitect                    // 2 - Contributor to the network
	TrustLevelAmbassador                   // 3 - Community leader
	TrustLevelGod                          // 4 - Core developer/admin
)

// String returns the string representation of TrustLevel
func (tl TrustLevel) String() string {
	switch tl {
	case TrustLevelGhost:
		return "Ghost"
	case TrustLevelUser:
		return "User"
	case TrustLevelArchitect:
		return "Architect"
	case TrustLevelAmbassador:
		return "Ambassador"
	case TrustLevelGod:
		return "God"
	default:
		return "Unknown"
	}
}

// UserProfile represents a user's profile information
type UserProfile struct {
	MessengerID   *MessengerID
	DisplayName   string
	TrustLevel    TrustLevel
	Reputation    int64
	LastSeen      time.Time
	IsBlocked     bool
	ContactsSince time.Time
}

// GenerateMessengerID creates a new MessengerID with cryptographic identity
func GenerateMessengerID() (*MessengerID, error) {
	// Generate Ed25519 key pair for identity
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ed25519 key pair: %w", err)
	}

	// Generate DID from public key hash
	did := generateDID(publicKey)

	// Create libp2p peer ID from Ed25519 public key
	libp2pPrivKey, err := crypto.UnmarshalEd25519PrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create libp2p private key: %w", err)
	}

	peerID, err := peer.IDFromPrivateKey(libp2pPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer ID: %w", err)
	}

	return &MessengerID{
		DID:        did,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		PeerID:     peerID,
		CreatedAt:  time.Now(),
	}, nil
}

// Sign creates an Ed25519 signature for the given data
func (mid *MessengerID) Sign(data []byte) ([]byte, error) {
	if mid.PrivateKey == nil {
		return nil, fmt.Errorf("private key not available")
	}

	// Create signature using Ed25519
	signature := ed25519.Sign(mid.PrivateKey, data)
	return signature, nil
}

// Verify verifies an Ed25519 signature
func (mid *MessengerID) Verify(data, signature []byte) bool {
	return ed25519.Verify(mid.PublicKey, data, signature)
}

// VerifySignature verifies a signature from another MessengerID
func VerifySignature(publicKey ed25519.PublicKey, data, signature []byte) bool {
	return ed25519.Verify(publicKey, data, signature)
}

// GetDID returns the DID string
func (mid *MessengerID) GetDID() string {
	return mid.DID
}

// GetPublicKeyHex returns the public key as hex string
func (mid *MessengerID) GetPublicKeyHex() string {
	return hex.EncodeToString(mid.PublicKey)
}

// GetPeerID returns the libp2p peer ID
func (mid *MessengerID) GetPeerID() peer.ID {
	return mid.PeerID
}

// Destroy securely destroys the MessengerID
func (mid *MessengerID) Destroy() {
	if mid.PrivateKey != nil {
		// TODO: Securely zero out the private key memory
		for i := range mid.PrivateKey {
			mid.PrivateKey[i] = 0
		}
	}
}

// generateDID creates a DID from a public key
func generateDID(publicKey ed25519.PublicKey) string {
	// Hash the public key
	hash := sha256.Sum256(publicKey)

	// Encode as base58 for readability
	encoded := base58.Encode(hash[:])

	return DIDPrefix + encoded
}

// ParseDID parses a DID string and extracts the hash
func ParseDID(did string) ([]byte, error) {
	if len(did) < len(DIDPrefix) {
		return nil, fmt.Errorf("invalid DID format")
	}

	if did[:len(DIDPrefix)] != DIDPrefix {
		return nil, fmt.Errorf("invalid DID prefix")
	}

	hashStr := did[len(DIDPrefix):]
	hash, err := base58.Decode(hashStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode DID hash: %w", err)
	}

	return hash, nil
}

// ValidateDID validates a DID format
func ValidateDID(did string) bool {
	_, err := ParseDID(did)
	return err == nil
}

// CreateUserProfile creates a new user profile with default settings
func CreateUserProfile(messengerID *MessengerID, displayName string) *UserProfile {
	return &UserProfile{
		MessengerID:   messengerID,
		DisplayName:   displayName,
		TrustLevel:    TrustLevelGhost, // Start as Ghost
		Reputation:    0,
		LastSeen:      time.Now(),
		IsBlocked:     false,
		ContactsSince: time.Now(),
	}
}

// CanSendMessage checks if user can send messages based on trust level
func (up *UserProfile) CanSendMessage() bool {
	// Ghost users have limited messaging capabilities
	if up.TrustLevel == TrustLevelGhost {
		// Check if user has been active for at least 24 hours
		return time.Since(up.ContactsSince) > 24*time.Hour
	}
	return true
}

// GetDailyMessageLimit returns the daily message limit based on trust level
func (up *UserProfile) GetDailyMessageLimit() int {
	switch up.TrustLevel {
	case TrustLevelGhost:
		return 5 // Limited for new users
	case TrustLevelUser:
		return 100
	case TrustLevelArchitect:
		return 500
	case TrustLevelAmbassador:
		return 1000
	case TrustLevelGod:
		return -1 // Unlimited
	default:
		return 5
	}
}

// CanCreateGroup checks if user can create groups
func (up *UserProfile) CanCreateGroup() bool {
	return up.TrustLevel >= TrustLevelUser
}

// UpdateLastSeen updates the last seen timestamp
func (up *UserProfile) UpdateLastSeen() {
	up.LastSeen = time.Now()
}

// IsOnline checks if user was seen recently (within 5 minutes)
func (up *UserProfile) IsOnline() bool {
	return time.Since(up.LastSeen) < 5*time.Minute
}
