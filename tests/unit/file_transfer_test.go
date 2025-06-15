package unit

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Xelvra/peerchat/internal/message"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileMetadataCreation(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, this is a test file for transfer!"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// Create file metadata
	metadata, err := message.CreateFileMetadata(testFile)
	require.NoError(t, err)
	assert.NotNil(t, metadata)

	// Verify metadata fields
	assert.Equal(t, "test.txt", metadata.Name)
	assert.Equal(t, int64(len(testContent)), metadata.Size)
	assert.NotEmpty(t, metadata.Hash)
	assert.Equal(t, "text/plain", metadata.MimeType)
	assert.NotEmpty(t, metadata.ID)
	assert.True(t, metadata.ChunkCount > 0)
	assert.Equal(t, message.FileChunkSize, metadata.ChunkSize)
}

func TestFileHashCalculation(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "hash_test.txt")
	testContent := "Test content for hash calculation"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// Calculate hash
	hash1, err := message.CalculateFileHash(testFile)
	require.NoError(t, err)
	assert.NotEmpty(t, hash1)

	// Calculate hash again - should be the same
	hash2, err := message.CalculateFileHash(testFile)
	require.NoError(t, err)
	assert.Equal(t, hash1, hash2)

	// Modify file and verify hash changes
	err = os.WriteFile(testFile, []byte(testContent+"modified"), 0644)
	require.NoError(t, err)

	hash3, err := message.CalculateFileHash(testFile)
	require.NoError(t, err)
	assert.NotEqual(t, hash1, hash3)
}

func TestMimeTypeDetection(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"test.txt", "text/plain"},
		{"document.pdf", "application/pdf"},
		{"image.jpg", "image/jpeg"},
		{"image.jpeg", "image/jpeg"},
		{"image.png", "image/png"},
		{"image.gif", "image/gif"},
		{"video.mp4", "video/mp4"},
		{"audio.mp3", "audio/mpeg"},
		{"archive.zip", "application/zip"},
		{"unknown.xyz", "application/octet-stream"},
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			tempDir := t.TempDir()
			testFile := filepath.Join(tempDir, test.filename)

			err := os.WriteFile(testFile, []byte("test content"), 0644)
			require.NoError(t, err)

			metadata, err := message.CreateFileMetadata(testFile)
			require.NoError(t, err)
			assert.Equal(t, test.expected, metadata.MimeType)
		})
	}
}

func TestFileTransferCreation(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	// Create test metadata
	metadata := message.FileMetadata{
		ID:         "test-transfer-123",
		Name:       "test.txt",
		Size:       1024,
		Hash:       "testhash",
		MimeType:   "text/plain",
		Timestamp:  time.Now(),
		ChunkCount: 1,
		ChunkSize:  message.FileChunkSize,
	}

	// Create a test peer ID
	peerID, err := peer.Decode("12D3KooWBhSxema2VqCGWW3dBkNQjzuUoTAozK9XP6y8JZtQZtjJ")
	require.NoError(t, err)

	// Test outgoing transfer
	transfer := message.NewFileTransfer("test-id", peerID, metadata, true, logger)
	assert.NotNil(t, transfer)
	assert.Equal(t, "test-id", transfer.ID)
	assert.Equal(t, peerID, transfer.PeerID)
	assert.Equal(t, metadata, transfer.Metadata)
	assert.True(t, transfer.IsOutgoing())
	assert.Equal(t, message.FileTransferPending, transfer.Status)
	assert.Equal(t, 0.0, transfer.Progress)
	assert.Equal(t, int64(1024), transfer.BytesTotal)

	// Test incoming transfer
	incomingTransfer := message.NewFileTransfer("test-id-2", peerID, metadata, false, logger)
	assert.False(t, incomingTransfer.IsOutgoing())
}

func TestFileTransferManager(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	manager := message.NewFileTransferManager(logger)
	assert.NotNil(t, manager)

	// Test empty manager
	transfers := manager.ListTransfers()
	assert.Empty(t, transfers)

	// Test getting non-existent transfer
	transfer, exists := manager.GetTransfer("non-existent")
	assert.Nil(t, transfer)
	assert.False(t, exists)
}

func TestFileTransferStatus(t *testing.T) {
	tests := []struct {
		status   message.FileTransferStatus
		expected string
	}{
		{message.FileTransferPending, "pending"},
		{message.FileTransferActive, "active"},
		{message.FileTransferCompleted, "completed"},
		{message.FileTransferFailed, "failed"},
		{message.FileTransferCancelled, "cancelled"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			assert.Equal(t, test.expected, test.status.String())
		})
	}
}

func TestFileTransferProgressUpdate(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	metadata := message.FileMetadata{
		ID:   "test-progress",
		Size: 1000,
	}

	peerID, err := peer.Decode("12D3KooWBhSxema2VqCGWW3dBkNQjzuUoTAozK9XP6y8JZtQZtjJ")
	require.NoError(t, err)

	// Test outgoing transfer progress
	outgoing := message.NewFileTransfer("test-out", peerID, metadata, true, logger)

	// Initial progress should be 0
	assert.Equal(t, 0.0, outgoing.Progress)

	// Simulate sending 500 bytes
	outgoing.BytesSent = 500
	outgoing.UpdateProgress()
	assert.Equal(t, 0.5, outgoing.Progress)

	// Simulate sending all bytes
	outgoing.BytesSent = 1000
	outgoing.UpdateProgress()
	assert.Equal(t, 1.0, outgoing.Progress)

	// Test incoming transfer progress
	incoming := message.NewFileTransfer("test-in", peerID, metadata, false, logger)

	// Simulate receiving 250 bytes
	incoming.BytesReceived = 250
	incoming.UpdateProgress()
	assert.Equal(t, 0.25, incoming.Progress)

	// Simulate receiving all bytes
	incoming.BytesReceived = 1000
	incoming.UpdateProgress()
	assert.Equal(t, 1.0, incoming.Progress)
}

func TestFileTransferConstants(t *testing.T) {
	// Verify file transfer constants are reasonable
	assert.Equal(t, 32*1024, message.FileChunkSize)                        // 32KB chunks
	assert.Equal(t, 1024, message.FileHeaderSize)                          // 1KB header limit
	assert.Equal(t, uint32(0x58454C56), uint32(message.FileTransferMagic)) // "XELV" magic
}
