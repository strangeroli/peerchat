package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Xelvra/peerchat/internal/message"
	"github.com/Xelvra/peerchat/internal/user"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"
)

const (
	// Database configuration
	DatabaseName = "userdata.db"
	WALMode      = "WAL"

	// Performance settings
	CheckpointInterval = 1000 // transactions
	CheckpointSizeMB   = 50   // MB

	// Encryption settings
	EncryptionKeySize   = 32    // AES-256 key size
	EncryptionNonceSize = 12    // GCM nonce size
	PBKDF2Iterations    = 100000 // PBKDF2 iterations for key derivation
)

// SQLiteDB represents the SQLite database with WAL mode and encryption
type SQLiteDB struct {
	db     *sql.DB
	logger *logrus.Logger
	dbPath string

	// Encryption
	encryptionKey []byte
	mutex         sync.RWMutex

	// Transaction counters for WAL checkpointing
	transactionCount int64
	lastCheckpoint   time.Time
}

// NewSQLiteDB creates a new SQLite database with optimized settings and encryption
func NewSQLiteDB(dataDir string, password string, logger *logrus.Logger) (*SQLiteDB, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Derive encryption key from password
	encryptionKey := deriveEncryptionKey(password)

	dbPath := filepath.Join(dataDir, DatabaseName)
	
	// Open database with WAL mode and optimizations
	dsn := fmt.Sprintf("%s?_journal_mode=%s&_synchronous=NORMAL&_cache_size=10000&_temp_store=memory", 
		dbPath, WALMode)
	
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	// Configure connection pool
	db.SetMaxOpenConns(1)  // SQLite works best with single connection
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)
	
	sqliteDB := &SQLiteDB{
		db:            db,
		logger:        logger,
		dbPath:        dbPath,
		encryptionKey: encryptionKey,
		lastCheckpoint: time.Now(),
	}
	
	// Initialize database schema
	if err := sqliteDB.initSchema(); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			sqliteDB.logger.WithError(closeErr).Error("Failed to close database after schema init error")
		}
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}
	
	// Start WAL checkpoint routine
	go sqliteDB.walCheckpointRoutine()

	logger.WithField("path", dbPath).Info("SQLite database initialized with WAL mode and encryption")
	return sqliteDB, nil
}

// deriveEncryptionKey derives an encryption key from password using PBKDF2
func deriveEncryptionKey(password string) []byte {
	// Use a fixed salt for simplicity - in production, store salt separately
	salt := []byte("xelvra_messenger_salt_2024")
	return pbkdf2.Key([]byte(password), salt, PBKDF2Iterations, EncryptionKeySize, sha256.New)
}

// Close closes the database connection
func (db *SQLiteDB) Close() error {
	if db.db != nil {
		// Perform final checkpoint
		if err := db.checkpoint(); err != nil {
			db.logger.WithError(err).Warn("Failed to perform final checkpoint")
		}

		db.logger.Info("Database closed successfully")
		return db.db.Close()
	}
	return nil
}

// encrypt encrypts data using AES-GCM
func (db *SQLiteDB) encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(db.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, EncryptionNonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// decrypt decrypts data using AES-GCM
func (db *SQLiteDB) decrypt(data []byte) ([]byte, error) {
	if len(data) < EncryptionNonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	block, err := aes.NewCipher(db.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := data[:EncryptionNonceSize]
	ciphertext := data[EncryptionNonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// walCheckpointRoutine runs periodic WAL checkpoints
func (db *SQLiteDB) walCheckpointRoutine() {
	ticker := time.NewTicker(5 * time.Minute) // Check every 5 minutes
	defer ticker.Stop()

	for range ticker.C {
		db.checkpointWALIfNeeded()
	}
}

// checkpointWALIfNeeded performs WAL checkpoint if needed
func (db *SQLiteDB) checkpointWALIfNeeded() {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	shouldCheckpoint := false

	// Check transaction count
	if db.transactionCount >= CheckpointInterval {
		shouldCheckpoint = true
		db.logger.WithField("transaction_count", db.transactionCount).Debug("WAL checkpoint triggered by transaction count")
	}

	// Check WAL file size
	if walSize, err := db.getWALSize(); err == nil && walSize > CheckpointSizeMB*1024*1024 {
		shouldCheckpoint = true
		db.logger.WithField("wal_size_mb", walSize/(1024*1024)).Debug("WAL checkpoint triggered by file size")
	}

	if shouldCheckpoint {
		if err := db.checkpoint(); err != nil {
			db.logger.WithError(err).Error("Failed to perform WAL checkpoint")
		} else {
			db.transactionCount = 0
			db.lastCheckpoint = time.Now()
			db.logger.Debug("WAL checkpoint completed successfully")
		}
	}
}

// getWALSize returns the size of the WAL file in bytes
func (db *SQLiteDB) getWALSize() (int64, error) {
	walPath := db.dbPath + "-wal"
	if info, err := os.Stat(walPath); err == nil {
		return info.Size(), nil
	}
	return 0, nil
}

// initSchema initializes the database schema
func (db *SQLiteDB) initSchema() error {
	schema := `
	-- Users table for storing user identities and profiles
	CREATE TABLE IF NOT EXISTS users (
		did TEXT PRIMARY KEY,
		public_key TEXT NOT NULL,
		display_name TEXT,
		trust_level INTEGER DEFAULT 0,
		reputation INTEGER DEFAULT 0,
		last_seen DATETIME,
		is_blocked BOOLEAN DEFAULT FALSE,
		contacts_since DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	-- Messages table for storing message history
	CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		type INTEGER NOT NULL,
		from_did TEXT NOT NULL,
		to_did TEXT,
		group_id TEXT,
		content BLOB,
		metadata TEXT, -- JSON
		timestamp DATETIME NOT NULL,
		signature BLOB,
		is_encrypted BOOLEAN DEFAULT FALSE,
		is_read BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (from_did) REFERENCES users(did),
		FOREIGN KEY (to_did) REFERENCES users(did)
	);
	
	-- Groups table for group chat management
	CREATE TABLE IF NOT EXISTS groups (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		creator_did TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (creator_did) REFERENCES users(did)
	);
	
	-- Group members table
	CREATE TABLE IF NOT EXISTS group_members (
		group_id TEXT NOT NULL,
		user_did TEXT NOT NULL,
		role TEXT DEFAULT 'member', -- member, admin, owner
		joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (group_id, user_did),
		FOREIGN KEY (group_id) REFERENCES groups(id),
		FOREIGN KEY (user_did) REFERENCES users(did)
	);
	
	-- Contacts table for managing user contacts
	CREATE TABLE IF NOT EXISTS contacts (
		owner_did TEXT NOT NULL,
		contact_did TEXT NOT NULL,
		display_name TEXT,
		is_blocked BOOLEAN DEFAULT FALSE,
		added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (owner_did, contact_did),
		FOREIGN KEY (owner_did) REFERENCES users(did),
		FOREIGN KEY (contact_did) REFERENCES users(did)
	);
	
	-- Files table for file transfer tracking
	CREATE TABLE IF NOT EXISTS files (
		id TEXT PRIMARY KEY,
		message_id TEXT,
		filename TEXT NOT NULL,
		file_size INTEGER NOT NULL,
		file_hash TEXT NOT NULL,
		mime_type TEXT,
		local_path TEXT,
		upload_progress REAL DEFAULT 0.0,
		download_progress REAL DEFAULT 0.0,
		status TEXT DEFAULT 'pending', -- pending, transferring, completed, failed
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (message_id) REFERENCES messages(id)
	);

	-- User settings table for encrypted configuration storage
	CREATE TABLE IF NOT EXISTS user_settings (
		key TEXT PRIMARY KEY,
		value BLOB NOT NULL, -- Encrypted value
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- File transfers table for tracking file transfer sessions
	CREATE TABLE IF NOT EXISTS file_transfers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		transfer_id TEXT UNIQUE NOT NULL,
		peer_id TEXT NOT NULL,
		file_name TEXT NOT NULL,
		file_size INTEGER NOT NULL,
		file_hash TEXT NOT NULL,
		status INTEGER NOT NULL, -- FileTransferStatus
		direction INTEGER NOT NULL, -- 0=outgoing, 1=incoming
		progress REAL DEFAULT 0.0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		completed_at DATETIME
	);
	
	-- Create indexes for better performance
	CREATE INDEX IF NOT EXISTS idx_messages_from_did ON messages(from_did);
	CREATE INDEX IF NOT EXISTS idx_messages_to_did ON messages(to_did);
	CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages(timestamp);
	CREATE INDEX IF NOT EXISTS idx_messages_group_id ON messages(group_id);
	CREATE INDEX IF NOT EXISTS idx_users_last_seen ON users(last_seen);
	CREATE INDEX IF NOT EXISTS idx_contacts_owner_did ON contacts(owner_did);
	CREATE INDEX IF NOT EXISTS idx_files_message_id ON files(message_id);
	CREATE INDEX IF NOT EXISTS idx_file_transfers_peer_id ON file_transfers(peer_id);
	CREATE INDEX IF NOT EXISTS idx_file_transfers_status ON file_transfers(status);
	
	-- Create triggers for updating timestamps
	CREATE TRIGGER IF NOT EXISTS update_users_timestamp 
		AFTER UPDATE ON users
		BEGIN
			UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE did = NEW.did;
		END;
	
	CREATE TRIGGER IF NOT EXISTS update_groups_timestamp 
		AFTER UPDATE ON groups
		BEGIN
			UPDATE groups SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;
	`
	
	_, err := db.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	
	db.logger.Info("Database schema initialized successfully")
	return nil
}

// SaveUser saves or updates a user profile
func (db *SQLiteDB) SaveUser(profile *user.UserProfile) error {
	query := `
		INSERT OR REPLACE INTO users 
		(did, public_key, display_name, trust_level, reputation, last_seen, is_blocked, contacts_since)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := db.db.Exec(query,
		profile.MessengerID.GetDID(),
		profile.MessengerID.GetPublicKeyHex(),
		profile.DisplayName,
		int(profile.TrustLevel),
		profile.Reputation,
		profile.LastSeen,
		profile.IsBlocked,
		profile.ContactsSince,
	)
	
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	
	db.incrementTransactionCount()
	return nil
}

// LoadUser loads a user profile by DID
func (db *SQLiteDB) LoadUser(did string) (*user.UserProfile, error) {
	query := `
		SELECT did, public_key, display_name, trust_level, reputation, 
		       last_seen, is_blocked, contacts_since
		FROM users WHERE did = ?
	`
	
	row := db.db.QueryRow(query, did)
	
	var profile user.UserProfile
	var publicKeyHex string
	var trustLevel int
	
	err := row.Scan(
		&did,
		&publicKeyHex,
		&profile.DisplayName,
		&trustLevel,
		&profile.Reputation,
		&profile.LastSeen,
		&profile.IsBlocked,
		&profile.ContactsSince,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %s", did)
		}
		return nil, fmt.Errorf("failed to load user: %w", err)
	}
	
	profile.TrustLevel = user.TrustLevel(trustLevel)
	
	// TODO: Reconstruct MessengerID from stored data
	// This would require storing and loading the full identity
	
	return &profile, nil
}

// SaveMessage saves a message to the database with encryption
func (db *SQLiteDB) SaveMessage(msg *message.Message) error {
	query := `
		INSERT INTO messages
		(id, type, from_did, to_did, group_id, content, metadata, timestamp, signature, is_encrypted)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Encrypt sensitive content
	var encryptedContent []byte
	var err error
	if len(msg.Content) > 0 {
		encryptedContent, err = db.encrypt(msg.Content)
		if err != nil {
			return fmt.Errorf("failed to encrypt message content: %w", err)
		}
	}

	var metadataJSON string
	if msg.Metadata != nil {
		// TODO: Serialize metadata to JSON
		metadataJSON = "{}"
	}

	_, err = db.db.Exec(query,
		msg.ID,
		int(msg.Type),
		msg.From,
		msg.To,
		msg.GroupID,
		encryptedContent,
		metadataJSON,
		msg.Timestamp,
		msg.Signature,
		msg.IsEncrypted,
	)

	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	db.incrementTransactionCount()
	return nil
}

// LoadMessages loads messages for a conversation with decryption
func (db *SQLiteDB) LoadMessages(fromDID, toDID string, limit int) ([]*message.Message, error) {
	query := `
		SELECT id, type, from_did, to_did, group_id, content, metadata,
		       timestamp, signature, is_encrypted
		FROM messages
		WHERE (from_did = ? AND to_did = ?) OR (from_did = ? AND to_did = ?)
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := db.db.Query(query, fromDID, toDID, toDID, fromDID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			db.logger.WithError(err).Error("Failed to close rows")
		}
	}()

	var messages []*message.Message

	for rows.Next() {
		var msg message.Message
		var msgType int
		var metadataJSON string
		var encryptedContent []byte

		err := rows.Scan(
			&msg.ID,
			&msgType,
			&msg.From,
			&msg.To,
			&msg.GroupID,
			&encryptedContent,
			&metadataJSON,
			&msg.Timestamp,
			&msg.Signature,
			&msg.IsEncrypted,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		// Decrypt content if present
		if len(encryptedContent) > 0 {
			decryptedContent, err := db.decrypt(encryptedContent)
			if err != nil {
				db.logger.WithError(err).Warn("Failed to decrypt message content")
				// Continue with encrypted content rather than failing
				msg.Content = encryptedContent
			} else {
				msg.Content = decryptedContent
			}
		}

		msg.Type = message.MessageType(msgType)
		// TODO: Deserialize metadata from JSON

		messages = append(messages, &msg)
	}

	return messages, nil
}

// SaveSetting saves an encrypted setting to the database
func (db *SQLiteDB) SaveSetting(key string, value []byte) error {
	encryptedValue, err := db.encrypt(value)
	if err != nil {
		return fmt.Errorf("failed to encrypt setting value: %w", err)
	}

	query := `INSERT OR REPLACE INTO user_settings (key, value) VALUES (?, ?)`
	_, err = db.db.Exec(query, key, encryptedValue)
	if err != nil {
		return fmt.Errorf("failed to save setting: %w", err)
	}

	db.incrementTransactionCount()
	return nil
}

// LoadSetting loads and decrypts a setting from the database
func (db *SQLiteDB) LoadSetting(key string) ([]byte, error) {
	query := `SELECT value FROM user_settings WHERE key = ?`

	var encryptedValue []byte
	err := db.db.QueryRow(query, key).Scan(&encryptedValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("setting not found: %s", key)
		}
		return nil, fmt.Errorf("failed to load setting: %w", err)
	}

	value, err := db.decrypt(encryptedValue)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt setting: %w", err)
	}

	return value, nil
}

// SaveFileTransfer saves file transfer information
func (db *SQLiteDB) SaveFileTransfer(transfer *message.FileTransfer) error {
	query := `
		INSERT OR REPLACE INTO file_transfers
		(transfer_id, peer_id, file_name, file_size, file_hash, status, direction, progress, completed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var completedAt *time.Time
	if transfer.Status == message.FileTransferCompleted {
		completedAt = &transfer.EndTime
	}

	direction := 0 // outgoing
	if !transfer.IsOutgoing() {
		direction = 1 // incoming
	}

	_, err := db.db.Exec(query,
		transfer.ID,
		transfer.PeerID.String(),
		transfer.Metadata.Name,
		transfer.Metadata.Size,
		transfer.Metadata.Hash,
		int(transfer.Status),
		direction,
		transfer.Progress,
		completedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save file transfer: %w", err)
	}

	db.incrementTransactionCount()
	return nil
}

// LoadFileTransfers loads file transfer history for a peer
func (db *SQLiteDB) LoadFileTransfers(peerID string, limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT transfer_id, peer_id, file_name, file_size, file_hash,
		       status, direction, progress, created_at, completed_at
		FROM file_transfers
		WHERE peer_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := db.db.Query(query, peerID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query file transfers: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			db.logger.WithError(err).Error("Failed to close rows")
		}
	}()

	var transfers []map[string]interface{}

	for rows.Next() {
		transfer := make(map[string]interface{})
		var transferID, peerIDStr, fileName, fileHash string
		var fileSize, status, direction int64
		var progress float64
		var createdAt time.Time
		var completedAt *time.Time

		err := rows.Scan(
			&transferID, &peerIDStr, &fileName, &fileSize, &fileHash,
			&status, &direction, &progress, &createdAt, &completedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan file transfer: %w", err)
		}

		transfer["transfer_id"] = transferID
		transfer["peer_id"] = peerIDStr
		transfer["file_name"] = fileName
		transfer["file_size"] = fileSize
		transfer["file_hash"] = fileHash
		transfer["status"] = status
		transfer["direction"] = direction
		transfer["progress"] = progress
		transfer["created_at"] = createdAt
		transfer["completed_at"] = completedAt

		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

// incrementTransactionCount increments the transaction counter and performs checkpoint if needed
func (db *SQLiteDB) incrementTransactionCount() {
	db.transactionCount++

	if db.transactionCount%CheckpointInterval == 0 {
		if err := db.checkpoint(); err != nil {
			db.logger.WithError(err).Warn("Failed to perform WAL checkpoint")
		}
	}
}

// checkpoint performs a WAL checkpoint to optimize database performance
func (db *SQLiteDB) checkpoint() error {
	_, err := db.db.Exec("PRAGMA wal_checkpoint(TRUNCATE)")
	if err != nil {
		return fmt.Errorf("failed to checkpoint WAL: %w", err)
	}
	
	db.logger.Debug("WAL checkpoint completed")
	return nil
}

// GetStats returns database statistics
func (db *SQLiteDB) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})
	
	// Get database size
	if info, err := os.Stat(db.dbPath); err == nil {
		stats["db_size_bytes"] = info.Size()
	}
	
	// Get WAL file size
	walPath := db.dbPath + "-wal"
	if info, err := os.Stat(walPath); err == nil {
		stats["wal_size_bytes"] = info.Size()
	}
	
	stats["transaction_count"] = db.transactionCount
	
	return stats
}
