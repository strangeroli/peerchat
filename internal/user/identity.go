package user

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
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

	// Proof-of-Work constants
	DefaultPOWDifficulty = 4                // Number of leading zeros required in hash
	MaxPOWIterations     = 1000000          // Maximum iterations to prevent infinite loops
	POWTimeout           = 30 * time.Second // Maximum time for PoW computation
)

// ProofOfWork represents a proof-of-work solution for DID creation
type ProofOfWork struct {
	Nonce      uint64    // Nonce value that solves the PoW
	Difficulty int       // Required difficulty (number of leading zeros)
	Hash       []byte    // Resulting hash that meets difficulty requirement
	ComputedAt time.Time // When the PoW was computed
}

// MessengerID represents a decentralized identity in Xelvra format
type MessengerID struct {
	DID         string             // did:xelvra:<hash>
	PublicKey   ed25519.PublicKey  // Ed25519 public key
	PrivateKey  ed25519.PrivateKey // Ed25519 private key (TODO: add memory protection)
	PeerID      peer.ID            // libp2p peer ID
	CreatedAt   time.Time          // Creation timestamp
	ProofOfWork *ProofOfWork       // PoW solution for Sybil resistance
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
	MessengerID       *MessengerID
	DisplayName       string
	TrustLevel        TrustLevel
	Reputation        int64
	LastSeen          time.Time
	IsBlocked         bool
	ContactsSince     time.Time
	MessagesSent      int64     // Track messages sent for rate limiting
	LastMessageTime   time.Time // Last message timestamp for rate limiting
	DailyMessageCount int       // Messages sent today
	LastDayReset      time.Time // When daily counter was last reset
}

// GenerateMessengerID creates a new MessengerID with cryptographic identity and PoW
func GenerateMessengerID() (*MessengerID, error) {
	return GenerateMessengerIDWithDifficulty(DefaultPOWDifficulty)
}

// GenerateMessengerIDWithDifficulty creates a new MessengerID with specified PoW difficulty
func GenerateMessengerIDWithDifficulty(difficulty int) (*MessengerID, error) {
	// Generate Ed25519 key pair for identity
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ed25519 key pair: %w", err)
	}

	// Create libp2p peer ID from Ed25519 public key
	libp2pPrivKey, err := crypto.UnmarshalEd25519PrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create libp2p private key: %w", err)
	}

	peerID, err := peer.IDFromPrivateKey(libp2pPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer ID: %w", err)
	}

	// Compute Proof-of-Work for the public key
	pow, err := computeProofOfWork(publicKey, difficulty)
	if err != nil {
		return nil, fmt.Errorf("failed to compute proof-of-work: %w", err)
	}

	// Generate DID from public key hash (including PoW validation)
	did := generateDIDWithPOW(publicKey, pow)

	return &MessengerID{
		DID:         did,
		PublicKey:   publicKey,
		PrivateKey:  privateKey,
		PeerID:      peerID,
		CreatedAt:   time.Now(),
		ProofOfWork: pow,
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

// generateDID creates a DID from a public key (legacy function)
func generateDID(publicKey ed25519.PublicKey) string {
	// Hash the public key
	hash := sha256.Sum256(publicKey)

	// Encode as base58 for readability
	encoded := base58.Encode(hash[:])

	return DIDPrefix + encoded
}

// generateDIDWithPOW creates a DID from a public key with PoW validation
func generateDIDWithPOW(publicKey ed25519.PublicKey, pow *ProofOfWork) string {
	// Create combined data: publicKey + nonce + difficulty
	data := make([]byte, len(publicKey)+8+4)
	copy(data, publicKey)
	binary.LittleEndian.PutUint64(data[len(publicKey):], pow.Nonce)
	binary.LittleEndian.PutUint32(data[len(publicKey)+8:], uint32(pow.Difficulty))

	// Hash the combined data
	hash := sha256.Sum256(data)

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

// CanSendMessageNow checks if user can send a message right now (rate limiting)
func (up *UserProfile) CanSendMessageNow() bool {
	// Reset daily counter if it's a new day
	if time.Since(up.LastDayReset) > 24*time.Hour {
		up.DailyMessageCount = 0
		up.LastDayReset = time.Now()
	}

	// Check daily limit
	limit := up.GetDailyMessageLimit()
	if limit > 0 && up.DailyMessageCount >= limit {
		return false
	}

	// Check rate limiting based on trust level
	switch up.TrustLevel {
	case TrustLevelGhost:
		// Ghost users: max 1 message per minute
		return time.Since(up.LastMessageTime) >= time.Minute
	case TrustLevelUser:
		// Regular users: max 1 message per 5 seconds
		return time.Since(up.LastMessageTime) >= 5*time.Second
	default:
		// Higher trust levels: no rate limiting
		return true
	}
}

// RecordMessageSent records that a message was sent
func (up *UserProfile) RecordMessageSent() {
	up.MessagesSent++
	up.LastMessageTime = time.Now()
	up.DailyMessageCount++
}

// computeProofOfWork computes a proof-of-work for the given public key
func computeProofOfWork(publicKey ed25519.PublicKey, difficulty int) (*ProofOfWork, error) {
	if difficulty <= 0 || difficulty > 32 {
		return nil, fmt.Errorf("invalid difficulty: %d (must be 1-32)", difficulty)
	}

	startTime := time.Now()
	var nonce uint64

	// Create the target: difficulty leading zeros
	target := make([]byte, 32)
	leadingZeros := difficulty / 8
	remainingBits := difficulty % 8

	// Set the target with leading zeros
	for i := 0; i < leadingZeros; i++ {
		target[i] = 0x00
	}
	if remainingBits > 0 {
		target[leadingZeros] = byte(0xFF >> remainingBits)
	} else if leadingZeros < 32 {
		target[leadingZeros] = 0xFF
	}

	for nonce < MaxPOWIterations {
		// Check timeout
		if time.Since(startTime) > POWTimeout {
			return nil, fmt.Errorf("proof-of-work computation timed out after %v", POWTimeout)
		}

		// Create data to hash: publicKey + nonce + difficulty
		data := make([]byte, len(publicKey)+8+4)
		copy(data, publicKey)
		binary.LittleEndian.PutUint64(data[len(publicKey):], nonce)
		binary.LittleEndian.PutUint32(data[len(publicKey)+8:], uint32(difficulty))

		// Compute hash
		hash := sha256.Sum256(data)

		// Check if hash meets difficulty requirement
		if meetsTarget(hash[:], target, difficulty) {
			return &ProofOfWork{
				Nonce:      nonce,
				Difficulty: difficulty,
				Hash:       hash[:],
				ComputedAt: time.Now(),
			}, nil
		}

		nonce++
	}

	return nil, fmt.Errorf("failed to find proof-of-work solution within %d iterations", MaxPOWIterations)
}

// meetsTarget checks if a hash meets the difficulty target
func meetsTarget(hash, target []byte, difficulty int) bool {
	leadingZeros := difficulty / 8
	remainingBits := difficulty % 8

	// Check full zero bytes
	for i := 0; i < leadingZeros; i++ {
		if hash[i] != 0x00 {
			return false
		}
	}

	// Check remaining bits if any
	if remainingBits > 0 && leadingZeros < len(hash) {
		mask := byte(0xFF << (8 - remainingBits))
		return (hash[leadingZeros] & mask) == 0x00
	}

	return true
}

// ValidateProofOfWork validates a proof-of-work solution
func ValidateProofOfWork(publicKey ed25519.PublicKey, pow *ProofOfWork) bool {
	if pow == nil {
		return false
	}

	// Recreate the data that was hashed
	data := make([]byte, len(publicKey)+8+4)
	copy(data, publicKey)
	binary.LittleEndian.PutUint64(data[len(publicKey):], pow.Nonce)
	binary.LittleEndian.PutUint32(data[len(publicKey)+8:], uint32(pow.Difficulty))

	// Compute hash
	hash := sha256.Sum256(data)

	// Verify the hash matches the stored hash
	if len(hash) != len(pow.Hash) {
		return false
	}
	for i := range hash {
		if hash[i] != pow.Hash[i] {
			return false
		}
	}

	// Create target for difficulty check
	target := make([]byte, 32)
	leadingZeros := pow.Difficulty / 8
	remainingBits := pow.Difficulty % 8

	for i := 0; i < leadingZeros; i++ {
		target[i] = 0x00
	}
	if remainingBits > 0 {
		target[leadingZeros] = byte(0xFF >> remainingBits)
	} else if leadingZeros < 32 {
		target[leadingZeros] = 0xFF
	}

	// Check if hash meets the difficulty requirement
	return meetsTarget(hash[:], target, pow.Difficulty)
}

// IsOnline checks if user was seen recently (within 5 minutes)
func (up *UserProfile) IsOnline() bool {
	return time.Since(up.LastSeen) < 5*time.Minute
}
