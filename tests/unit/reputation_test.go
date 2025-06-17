package unit

import (
	"testing"

	"github.com/Xelvra/peerchat/internal/user"
	"github.com/sirupsen/logrus"
)

// TestReputationManagerCreation tests reputation manager creation
func TestReputationManagerCreation(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	rm := user.NewReputationManager(logger)
	if rm == nil {
		t.Fatal("Failed to create reputation manager")
	}
}

// TestUserReputationCreation tests user reputation creation
func TestUserReputationCreation(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	rm := user.NewReputationManager(logger)

	did := "did:xelvra:test123"
	rep := rm.CreateUserReputation(did)

	if rep == nil {
		t.Fatal("Failed to create user reputation")
	}

	if rep.DID != did {
		t.Errorf("Expected DID %s, got %s", did, rep.DID)
	}

	if rep.TrustLevel != user.TrustLevelGhost {
		t.Errorf("Expected trust level Ghost, got %s", rep.TrustLevel.String())
	}

	if rep.ReputationScore != 0 {
		t.Errorf("Expected reputation score 0, got %d", rep.ReputationScore)
	}

	if rep.ReliabilityScore != 1.0 {
		t.Errorf("Expected reliability score 1.0, got %f", rep.ReliabilityScore)
	}
}

// TestReputationRetrieval tests reputation retrieval
func TestReputationRetrieval(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	rm := user.NewReputationManager(logger)

	did := "did:xelvra:test123"
	originalRep := rm.CreateUserReputation(did)

	retrievedRep, err := rm.GetUserReputation(did)
	if err != nil {
		t.Fatalf("Failed to retrieve user reputation: %v", err)
	}

	if retrievedRep.DID != originalRep.DID {
		t.Errorf("Retrieved reputation DID mismatch")
	}

	if retrievedRep.TrustLevel != originalRep.TrustLevel {
		t.Errorf("Retrieved reputation trust level mismatch")
	}

	// Test non-existent user
	_, err = rm.GetUserReputation("did:xelvra:nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent user, got none")
	}
}

// TestActivityUpdates tests reputation activity updates
func TestActivityUpdates(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	rm := user.NewReputationManager(logger)

	did := "did:xelvra:test123"
	rm.CreateUserReputation(did)

	// Test message sent
	err := rm.UpdateActivity(did, "message_sent")
	if err != nil {
		t.Fatalf("Failed to update activity: %v", err)
	}

	rep, _ := rm.GetUserReputation(did)
	if rep.MessagesSent != 1 {
		t.Errorf("Expected 1 message sent, got %d", rep.MessagesSent)
	}
	if rep.ReputationScore != 1 {
		t.Errorf("Expected reputation score 1, got %d", rep.ReputationScore)
	}

	// Test file shared
	err = rm.UpdateActivity(did, "file_shared")
	if err != nil {
		t.Fatalf("Failed to update file sharing activity: %v", err)
	}

	rep, _ = rm.GetUserReputation(did)
	if rep.FilesShared != 1 {
		t.Errorf("Expected 1 file shared, got %d", rep.FilesShared)
	}
	if rep.ReputationScore != 6 { // 1 + 5
		t.Errorf("Expected reputation score 6, got %d", rep.ReputationScore)
	}

	// Test online hour
	err = rm.UpdateActivity(did, "online_hour")
	if err != nil {
		t.Fatalf("Failed to update online hour activity: %v", err)
	}

	rep, _ = rm.GetUserReputation(did)
	if rep.UptimeHours != 1.0 {
		t.Errorf("Expected 1.0 uptime hours, got %f", rep.UptimeHours)
	}
	if rep.ReputationScore != 8 { // 6 + 2
		t.Errorf("Expected reputation score 8, got %d", rep.ReputationScore)
	}
}

// TestTrustLevelRequirements tests trust level requirements
func TestTrustLevelRequirements(t *testing.T) {
	testCases := []struct {
		level                    user.TrustLevel
		expectedMinReputation    int64
		expectedMinUptime        float64
		expectedMinReliability   float64
		expectedMinVerifications int
	}{
		{user.TrustLevelUser, 100, 24, 0.8, 1},
		{user.TrustLevelArchitect, 1000, 168, 0.9, 3},
		{user.TrustLevelAmbassador, 10000, 720, 0.95, 5},
		{user.TrustLevelGod, 100000, 2160, 0.98, 10},
	}

	for _, tc := range testCases {
		t.Run(tc.level.String(), func(t *testing.T) {
			req := user.GetTrustLevelRequirements(tc.level)

			if req.MinReputationScore != tc.expectedMinReputation {
				t.Errorf("Expected min reputation %d, got %d", tc.expectedMinReputation, req.MinReputationScore)
			}

			if req.MinUptimeHours != tc.expectedMinUptime {
				t.Errorf("Expected min uptime %f, got %f", tc.expectedMinUptime, req.MinUptimeHours)
			}

			if req.MinReliabilityScore != tc.expectedMinReliability {
				t.Errorf("Expected min reliability %f, got %f", tc.expectedMinReliability, req.MinReliabilityScore)
			}

			if req.MinVerifications != tc.expectedMinVerifications {
				t.Errorf("Expected min verifications %d, got %d", tc.expectedMinVerifications, req.MinVerifications)
			}
		})
	}
}

// TestRateLimiting tests message rate limiting
func TestRateLimiting(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	rm := user.NewReputationManager(logger)

	did := "did:xelvra:test123"
	rm.CreateUserReputation(did)

	// Test initial rate limiting (Ghost level)
	canSend, reason := rm.CanSendMessage(did)
	if !canSend {
		t.Errorf("New user should be able to send first message, reason: %s", reason)
	}

	// Record message sent
	err := rm.RecordMessageSent(did)
	if err != nil {
		t.Fatalf("Failed to record message sent: %v", err)
	}

	// Test rate limiting (should be limited for Ghost level)
	canSend, reason = rm.CanSendMessage(did)
	if canSend {
		t.Error("Ghost user should be rate limited after sending message")
	}

	if reason == "" {
		t.Error("Rate limit reason should not be empty")
	}
}

// TestUserVerification tests user verification system
func TestUserVerification(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	rm := user.NewReputationManager(logger)

	verifierDID := "did:xelvra:verifier"
	targetDID := "did:xelvra:target"

	// Create users
	verifierRep := rm.CreateUserReputation(verifierDID)
	rm.CreateUserReputation(targetDID)

	// Promote verifier to User level manually for testing
	verifierRep.TrustLevel = user.TrustLevelUser
	verifierRep.ReputationScore = 1000

	// Test verification
	err := rm.VerifyUser(verifierDID, targetDID)
	if err != nil {
		t.Fatalf("Failed to verify user: %v", err)
	}

	// Check that target was verified
	targetRep, _ := rm.GetUserReputation(targetDID)
	if len(targetRep.VerifiedBy) != 1 {
		t.Errorf("Expected 1 verification, got %d", len(targetRep.VerifiedBy))
	}

	if targetRep.VerifiedBy[0] != verifierDID {
		t.Errorf("Expected verifier %s, got %s", verifierDID, targetRep.VerifiedBy[0])
	}

	// Check that verifier recorded the verification
	verifierRep, _ = rm.GetUserReputation(verifierDID)
	if len(verifierRep.Verified) != 1 {
		t.Errorf("Expected verifier to have 1 verification record, got %d", len(verifierRep.Verified))
	}

	// Test duplicate verification
	err = rm.VerifyUser(verifierDID, targetDID)
	if err == nil {
		t.Error("Expected error for duplicate verification, got none")
	}
}

// TestDailyMessageLimits tests daily message limits for different trust levels
func TestDailyMessageLimits(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	rm := user.NewReputationManager(logger)

	did := "did:xelvra:test123"
	rep := rm.CreateUserReputation(did)

	// Test Ghost level (5 messages/day)
	rep.DailyMessageCount = 5
	canSend, reason := rm.CanSendMessage(did)
	if canSend {
		t.Error("Ghost user should be limited at 5 messages/day")
	}
	if reason == "" {
		t.Error("Daily limit reason should not be empty")
	}

	// Test User level (100 messages/day)
	rep.TrustLevel = user.TrustLevelUser
	rep.DailyMessageCount = 100
	canSend, reason = rm.CanSendMessage(did)
	if canSend {
		t.Error("User should be limited at 100 messages/day")
	}

	// Test God level (unlimited)
	rep.TrustLevel = user.TrustLevelGod
	rep.DailyMessageCount = 10000
	canSend, _ = rm.CanSendMessage(did)
	if !canSend {
		t.Error("God level should have unlimited messages")
	}
}
