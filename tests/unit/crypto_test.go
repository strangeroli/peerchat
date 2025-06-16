package unit

import (
	"testing"

	"github.com/Xelvra/peerchat/internal/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateKeyPair(t *testing.T) {
	keyPair, err := crypto.GenerateKeyPair()
	require.NoError(t, err)
	require.NotNil(t, keyPair)

	assert.Len(t, keyPair.PrivateKey, crypto.PrivateKeySize)
	assert.Len(t, keyPair.PublicKey, crypto.PublicKeySize)
	assert.NotEqual(t, make([]byte, crypto.PrivateKeySize), keyPair.PrivateKey)
	assert.NotEqual(t, make([]byte, crypto.PublicKeySize), keyPair.PublicKey)
}

func TestKeyPairDestroy(t *testing.T) {
	keyPair, err := crypto.GenerateKeyPair()
	require.NoError(t, err)

	// Store original private key for comparison
	originalPrivateKey := make([]byte, len(keyPair.PrivateKey))
	copy(originalPrivateKey, keyPair.PrivateKey)

	// Destroy the key pair
	keyPair.Destroy()

	// Verify private key is zeroed out
	assert.Equal(t, make([]byte, crypto.PrivateKeySize), keyPair.PrivateKey)
	assert.NotEqual(t, originalPrivateKey, keyPair.PrivateKey)
}

func TestSignalCrypto(t *testing.T) {
	sc, err := crypto.NewSignalCrypto()
	require.NoError(t, err)
	require.NotNil(t, sc)

	identityKey := sc.GetIdentityKey()
	assert.Len(t, identityKey, crypto.PublicKeySize)
	assert.NotEqual(t, make([]byte, crypto.PublicKeySize), identityKey)
}

func TestEncryptDecryptMessage(t *testing.T) {
	sc, err := crypto.NewSignalCrypto()
	require.NoError(t, err)

	// Test message
	plaintext := []byte("Hello, Xelvra P2P!")
	chainKey := make([]byte, crypto.SharedKeySize)
	copy(chainKey, "test-chain-key-32-bytes-long!!")

	// Encrypt message
	ciphertext, err := sc.EncryptMessage(plaintext, chainKey)
	require.NoError(t, err)
	assert.NotEqual(t, plaintext, ciphertext)
	assert.Greater(t, len(ciphertext), len(plaintext))

	// Decrypt message
	decrypted, err := sc.DecryptMessage(ciphertext, chainKey)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestReplayAttackProtection(t *testing.T) {
	sc, err := crypto.NewSignalCrypto()
	require.NoError(t, err)

	plaintext := []byte("Test message")
	chainKey := make([]byte, crypto.SharedKeySize)
	copy(chainKey, "test-chain-key-32-bytes-long!!")

	// Encrypt message
	ciphertext, err := sc.EncryptMessage(plaintext, chainKey)
	require.NoError(t, err)

	// First decryption should succeed
	decrypted1, err := sc.DecryptMessage(ciphertext, chainKey)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted1)

	// Second decryption with same ciphertext should fail (replay attack)
	_, err = sc.DecryptMessage(ciphertext, chainKey)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "replay attack detected")
}

func TestSignalCryptoDestroy(t *testing.T) {
	sc, err := crypto.NewSignalCrypto()
	require.NoError(t, err)

	// Destroy should not panic
	assert.NotPanics(t, func() {
		sc.Destroy()
	})

	// Multiple destroys should not panic
	assert.NotPanics(t, func() {
		sc.Destroy()
	})
}

func TestEncryptDecryptWithInvalidChainKey(t *testing.T) {
	sc, err := crypto.NewSignalCrypto()
	require.NoError(t, err)

	plaintext := []byte("Test message")
	chainKey := make([]byte, crypto.SharedKeySize)
	copy(chainKey, "test-chain-key-32-bytes-long!!")

	// Encrypt with valid chain key
	ciphertext, err := sc.EncryptMessage(plaintext, chainKey)
	require.NoError(t, err)

	// Try to decrypt with different chain key
	wrongChainKey := make([]byte, crypto.SharedKeySize)
	copy(wrongChainKey, "wrong-chain-key-32-bytes-long!")

	_, err = sc.DecryptMessage(ciphertext, wrongChainKey)
	assert.Error(t, err)
}

func TestDecryptInvalidCiphertext(t *testing.T) {
	sc, err := crypto.NewSignalCrypto()
	require.NoError(t, err)

	chainKey := make([]byte, crypto.SharedKeySize)
	copy(chainKey, "test-chain-key-32-bytes-long!!")

	// Test with too short ciphertext
	shortCiphertext := []byte("short")
	_, err = sc.DecryptMessage(shortCiphertext, chainKey)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ciphertext too short")

	// Test with corrupted ciphertext
	corruptedCiphertext := make([]byte, crypto.NonceSize+crypto.TagSize+10)
	_, err = sc.DecryptMessage(corruptedCiphertext, chainKey)
	assert.Error(t, err)
}
