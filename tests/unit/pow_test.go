package unit

import (
	"fmt"
	"testing"
	"time"

	"github.com/Xelvra/peerchat/internal/user"
)

// TestProofOfWorkGeneration tests PoW identity generation
func TestProofOfWorkGeneration(t *testing.T) {
	// Test with low difficulty for speed
	identity, err := user.GenerateMessengerIDWithDifficulty(2)
	if err != nil {
		t.Fatalf("Failed to generate PoW identity: %v", err)
	}

	if identity == nil {
		t.Fatal("Generated identity is nil")
	}

	if identity.ProofOfWork == nil {
		t.Fatal("Identity missing proof-of-work")
	}

	if identity.ProofOfWork.Difficulty != 2 {
		t.Errorf("Expected difficulty 2, got %d", identity.ProofOfWork.Difficulty)
	}

	if len(identity.ProofOfWork.Hash) != 32 {
		t.Errorf("Expected hash length 32, got %d", len(identity.ProofOfWork.Hash))
	}

	if identity.DID == "" {
		t.Error("DID is empty")
	}

	if len(identity.PublicKey) == 0 {
		t.Error("Public key is empty")
	}

	if len(identity.PrivateKey) == 0 {
		t.Error("Private key is empty")
	}
}

// TestProofOfWorkValidation tests PoW validation
func TestProofOfWorkValidation(t *testing.T) {
	identity, err := user.GenerateMessengerIDWithDifficulty(2)
	if err != nil {
		t.Fatalf("Failed to generate PoW identity: %v", err)
	}

	// Test valid PoW
	if !user.ValidateProofOfWork(identity.PublicKey, identity.ProofOfWork) {
		t.Error("Valid proof-of-work failed validation")
	}

	// Test invalid PoW (modify nonce)
	invalidPow := *identity.ProofOfWork
	invalidPow.Nonce = identity.ProofOfWork.Nonce + 1
	if user.ValidateProofOfWork(identity.PublicKey, &invalidPow) {
		t.Error("Invalid proof-of-work passed validation")
	}

	// Test nil PoW
	if user.ValidateProofOfWork(identity.PublicKey, nil) {
		t.Error("Nil proof-of-work passed validation")
	}
}

// TestProofOfWorkDifficulty tests different difficulty levels
func TestProofOfWorkDifficulty(t *testing.T) {
	difficulties := []int{1, 2, 3}

	for _, difficulty := range difficulties {
		t.Run(fmt.Sprintf("Difficulty_%d", difficulty), func(t *testing.T) {
			start := time.Now()
			identity, err := user.GenerateMessengerIDWithDifficulty(difficulty)
			duration := time.Since(start)

			if err != nil {
				t.Fatalf("Failed to generate PoW with difficulty %d: %v", difficulty, err)
			}

			if identity.ProofOfWork.Difficulty != difficulty {
				t.Errorf("Expected difficulty %d, got %d", difficulty, identity.ProofOfWork.Difficulty)
			}

			// Higher difficulty should generally take longer (though not guaranteed)
			t.Logf("Difficulty %d took %v", difficulty, duration)

			// Validate the PoW
			if !user.ValidateProofOfWork(identity.PublicKey, identity.ProofOfWork) {
				t.Errorf("Generated PoW with difficulty %d failed validation", difficulty)
			}
		})
	}
}

// TestProofOfWorkTimeout tests PoW timeout handling
func TestProofOfWorkTimeout(t *testing.T) {
	// This test would require modifying the timeout constant or using dependency injection
	// For now, we'll test with a reasonable difficulty that should complete quickly
	identity, err := user.GenerateMessengerIDWithDifficulty(4)
	if err != nil {
		t.Fatalf("Failed to generate PoW identity: %v", err)
	}

	if identity.ProofOfWork.ComputedAt.IsZero() {
		t.Error("PoW computation time not recorded")
	}

	// Check that computation time is recent
	if time.Since(identity.ProofOfWork.ComputedAt) > time.Minute {
		t.Error("PoW computation time seems too old")
	}
}

// TestInvalidDifficulty tests invalid difficulty values
func TestInvalidDifficulty(t *testing.T) {
	invalidDifficulties := []int{0, -1, 33, 100}

	for _, difficulty := range invalidDifficulties {
		t.Run(fmt.Sprintf("Invalid_%d", difficulty), func(t *testing.T) {
			_, err := user.GenerateMessengerIDWithDifficulty(difficulty)
			if err == nil {
				t.Errorf("Expected error for invalid difficulty %d, but got none", difficulty)
			}
		})
	}
}

// BenchmarkProofOfWork benchmarks PoW generation
func BenchmarkProofOfWork(b *testing.B) {
	difficulties := []int{1, 2, 3}

	for _, difficulty := range difficulties {
		b.Run(fmt.Sprintf("Difficulty_%d", difficulty), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := user.GenerateMessengerIDWithDifficulty(difficulty)
				if err != nil {
					b.Fatalf("Failed to generate PoW: %v", err)
				}
			}
		})
	}
}
