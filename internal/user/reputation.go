package user

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ReputationManager manages the hierarchical reputation system
// Implements the Ghost -> User -> Architect -> Ambassador -> God system from tmp/Proof-of-Work.md
type ReputationManager struct {
	logger *logrus.Logger
	mu     sync.RWMutex

	// User reputation tracking
	userReputations map[string]*UserReputation // DID -> UserReputation

	// Trust network (who trusts whom)
	trustNetwork map[string][]string // DID -> list of trusted DIDs

	// Verification network (who verified whom)
	verificationNetwork map[string][]string // DID -> list of verified DIDs
}

// UserReputation represents a user's reputation in the network
type UserReputation struct {
	DID              string
	TrustLevel       TrustLevel
	ReputationScore  int64
	MessagesSent     int64
	MessagesReceived int64
	FilesShared      int64
	UptimeHours      float64
	LastActivity     time.Time
	CreatedAt        time.Time

	// Behavioral metrics
	ReliabilityScore    float64 // 0.0 - 1.0 (message delivery success rate)
	ResponsivenessScore float64 // 0.0 - 1.0 (response time quality)
	HelpfulnessScore    float64 // 0.0 - 1.0 (community contribution)

	// Trust relationships
	TrustedBy  []string // DIDs that trust this user
	Trusts     []string // DIDs this user trusts
	VerifiedBy []string // DIDs that verified this user
	Verified   []string // DIDs this user verified

	// Rate limiting
	DailyMessageCount int
	LastMessageTime   time.Time
	LastDayReset      time.Time

	// Violations and penalties
	ViolationCount int
	LastViolation  time.Time
	PenaltyUntil   time.Time
}

// TrustLevelRequirements defines requirements for each trust level
type TrustLevelRequirements struct {
	MinReputationScore     int64
	MinUptimeHours         float64
	MinReliabilityScore    float64
	MinVerifications       int
	MinTimeInPreviousLevel time.Duration
}

// GetTrustLevelRequirements returns requirements for each trust level
func GetTrustLevelRequirements(level TrustLevel) TrustLevelRequirements {
	switch level {
	case TrustLevelUser:
		return TrustLevelRequirements{
			MinReputationScore:     100,
			MinUptimeHours:         24,                 // 24 hours of activity
			MinReliabilityScore:    0.8,                // 80% message delivery success
			MinVerifications:       1,                  // At least 1 verification from higher level
			MinTimeInPreviousLevel: 7 * 24 * time.Hour, // 1 week as Ghost
		}
	case TrustLevelArchitect:
		return TrustLevelRequirements{
			MinReputationScore:     1000,
			MinUptimeHours:         168,                 // 1 week of activity
			MinReliabilityScore:    0.9,                 // 90% reliability
			MinVerifications:       3,                   // 3 verifications from Architect+ level
			MinTimeInPreviousLevel: 30 * 24 * time.Hour, // 1 month as User
		}
	case TrustLevelAmbassador:
		return TrustLevelRequirements{
			MinReputationScore:     10000,
			MinUptimeHours:         720,                 // 1 month of activity
			MinReliabilityScore:    0.95,                // 95% reliability
			MinVerifications:       5,                   // 5 verifications from Ambassador+ level
			MinTimeInPreviousLevel: 90 * 24 * time.Hour, // 3 months as Architect
		}
	case TrustLevelGod:
		return TrustLevelRequirements{
			MinReputationScore:     100000,
			MinUptimeHours:         2160,                 // 3 months of activity
			MinReliabilityScore:    0.98,                 // 98% reliability
			MinVerifications:       10,                   // 10 verifications from God level
			MinTimeInPreviousLevel: 180 * 24 * time.Hour, // 6 months as Ambassador
		}
	default:
		return TrustLevelRequirements{} // Ghost level has no requirements
	}
}

// NewReputationManager creates a new reputation manager
func NewReputationManager(logger *logrus.Logger) *ReputationManager {
	return &ReputationManager{
		logger:              logger,
		userReputations:     make(map[string]*UserReputation),
		trustNetwork:        make(map[string][]string),
		verificationNetwork: make(map[string][]string),
	}
}

// GetUserReputation returns reputation for a user DID
func (rm *ReputationManager) GetUserReputation(did string) (*UserReputation, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	if rep, exists := rm.userReputations[did]; exists {
		// Return a copy to avoid race conditions
		repCopy := *rep
		return &repCopy, nil
	}

	return nil, fmt.Errorf("user reputation not found for DID: %s", did)
}

// CreateUserReputation creates initial reputation for a new user
func (rm *ReputationManager) CreateUserReputation(did string) *UserReputation {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rep := &UserReputation{
		DID:                 did,
		TrustLevel:          TrustLevelGhost,
		ReputationScore:     0,
		CreatedAt:           time.Now(),
		LastActivity:        time.Now(),
		ReliabilityScore:    1.0, // Start with perfect score
		ResponsivenessScore: 1.0,
		HelpfulnessScore:    0.5, // Neutral start
		TrustedBy:           make([]string, 0),
		Trusts:              make([]string, 0),
		VerifiedBy:          make([]string, 0),
		Verified:            make([]string, 0),
		LastDayReset:        time.Now(),
	}

	rm.userReputations[did] = rep

	rm.logger.WithFields(logrus.Fields{
		"did":         did,
		"trust_level": rep.TrustLevel.String(),
	}).Info("Created new user reputation")

	return rep
}

// UpdateActivity updates user activity metrics
func (rm *ReputationManager) UpdateActivity(did string, activityType string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rep, exists := rm.userReputations[did]
	if !exists {
		return fmt.Errorf("user reputation not found for DID: %s", did)
	}

	rep.LastActivity = time.Now()

	switch activityType {
	case "message_sent":
		rep.MessagesSent++
		rep.ReputationScore += 1
	case "message_received":
		rep.MessagesReceived++
	case "file_shared":
		rep.FilesShared++
		rep.ReputationScore += 5
	case "online_hour":
		rep.UptimeHours += 1.0
		rep.ReputationScore += 2
	}

	// Check for trust level promotion
	rm.checkTrustLevelPromotion(rep)

	return nil
}

// checkTrustLevelPromotion checks if user qualifies for trust level promotion
func (rm *ReputationManager) checkTrustLevelPromotion(rep *UserReputation) {
	currentLevel := rep.TrustLevel
	nextLevel := currentLevel + 1

	// Can't promote beyond God level
	if nextLevel > TrustLevelGod {
		return
	}

	requirements := GetTrustLevelRequirements(nextLevel)
	timeInCurrentLevel := time.Since(rep.CreatedAt)

	// Check if user meets all requirements
	if rep.ReputationScore >= requirements.MinReputationScore &&
		rep.UptimeHours >= requirements.MinUptimeHours &&
		rep.ReliabilityScore >= requirements.MinReliabilityScore &&
		len(rep.VerifiedBy) >= requirements.MinVerifications &&
		timeInCurrentLevel >= requirements.MinTimeInPreviousLevel {

		oldLevel := rep.TrustLevel
		rep.TrustLevel = nextLevel

		rm.logger.WithFields(logrus.Fields{
			"did":        rep.DID,
			"old_level":  oldLevel.String(),
			"new_level":  rep.TrustLevel.String(),
			"reputation": rep.ReputationScore,
		}).Info("User promoted to higher trust level")
	}
}

// VerifyUser allows higher-level users to verify lower-level users
func (rm *ReputationManager) VerifyUser(verifierDID, targetDID string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	verifier, exists := rm.userReputations[verifierDID]
	if !exists {
		return fmt.Errorf("verifier not found: %s", verifierDID)
	}

	target, exists := rm.userReputations[targetDID]
	if !exists {
		return fmt.Errorf("target user not found: %s", targetDID)
	}

	// Only higher or equal level users can verify
	if verifier.TrustLevel < target.TrustLevel {
		return fmt.Errorf("insufficient trust level to verify user")
	}

	// Check if already verified
	for _, did := range target.VerifiedBy {
		if did == verifierDID {
			return fmt.Errorf("user already verified by this verifier")
		}
	}

	// Add verification
	target.VerifiedBy = append(target.VerifiedBy, verifierDID)
	verifier.Verified = append(verifier.Verified, targetDID)

	// Bonus reputation for verification
	target.ReputationScore += 50
	verifier.ReputationScore += 10 // Verifier gets small bonus

	rm.logger.WithFields(logrus.Fields{
		"verifier":       verifierDID,
		"target":         targetDID,
		"verifier_level": verifier.TrustLevel.String(),
		"target_level":   target.TrustLevel.String(),
	}).Info("User verification completed")

	// Check for promotion after verification
	rm.checkTrustLevelPromotion(target)

	return nil
}

// CanSendMessage checks if user can send a message based on trust level and rate limits
func (rm *ReputationManager) CanSendMessage(did string) (bool, string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rep, exists := rm.userReputations[did]
	if !exists {
		return false, "User reputation not found"
	}

	// Check if user is under penalty
	if time.Now().Before(rep.PenaltyUntil) {
		return false, fmt.Sprintf("User under penalty until %s", rep.PenaltyUntil.Format("15:04:05"))
	}

	// Reset daily counter if it's a new day
	if time.Since(rep.LastDayReset) > 24*time.Hour {
		rep.DailyMessageCount = 0
		rep.LastDayReset = time.Now()
	}

	// Check daily limit based on trust level
	dailyLimit := getDailyMessageLimit(rep.TrustLevel)
	if dailyLimit > 0 && rep.DailyMessageCount >= dailyLimit {
		return false, fmt.Sprintf("Daily message limit reached (%d/%d)", rep.DailyMessageCount, dailyLimit)
	}

	// Check rate limiting based on trust level
	minInterval := getMinMessageInterval(rep.TrustLevel)
	if time.Since(rep.LastMessageTime) < minInterval {
		return false, fmt.Sprintf("Rate limit: wait %v", minInterval-time.Since(rep.LastMessageTime))
	}

	return true, ""
}

// RecordMessageSent records that a message was sent
func (rm *ReputationManager) RecordMessageSent(did string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rep, exists := rm.userReputations[did]
	if !exists {
		return fmt.Errorf("user reputation not found for DID: %s", did)
	}

	rep.LastMessageTime = time.Now()
	rep.DailyMessageCount++

	return nil
}

// getDailyMessageLimit returns daily message limit for trust level
func getDailyMessageLimit(level TrustLevel) int {
	switch level {
	case TrustLevelGhost:
		return 5
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

// getMinMessageInterval returns minimum interval between messages
func getMinMessageInterval(level TrustLevel) time.Duration {
	switch level {
	case TrustLevelGhost:
		return 1 * time.Minute
	case TrustLevelUser:
		return 5 * time.Second
	case TrustLevelArchitect:
		return 1 * time.Second
	case TrustLevelAmbassador:
		return 500 * time.Millisecond
	case TrustLevelGod:
		return 0 // No rate limiting
	default:
		return 1 * time.Minute
	}
}
