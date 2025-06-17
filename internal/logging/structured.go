package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// LogLevel represents log levels
type LogLevel string

const (
	LevelTrace LogLevel = "trace"
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
	LevelFatal LogLevel = "fatal"
	LevelPanic LogLevel = "panic"
)

// LogFormat represents log output formats
type LogFormat string

const (
	FormatJSON LogFormat = "json"
	FormatText LogFormat = "text"
)

// Config represents logging configuration
type Config struct {
	Level      LogLevel  `yaml:"level" json:"level"`
	Format     LogFormat `yaml:"format" json:"format"`
	Output     string    `yaml:"output" json:"output"`         // file path or "stdout", "stderr"
	MaxSize    int       `yaml:"max_size" json:"max_size"`     // MB
	MaxBackups int       `yaml:"max_backups" json:"max_backups"`
	MaxAge     int       `yaml:"max_age" json:"max_age"`       // days
	Compress   bool      `yaml:"compress" json:"compress"`
	
	// Component-specific log levels
	Components map[string]LogLevel `yaml:"components" json:"components"`
}

// DefaultConfig returns default logging configuration
func DefaultConfig() *Config {
	return &Config{
		Level:      LevelInfo,
		Format:     FormatText,
		Output:     "stdout",
		MaxSize:    10,  // 10MB
		MaxBackups: 5,
		MaxAge:     30,  // 30 days
		Compress:   true,
		Components: make(map[string]LogLevel),
	}
}

// StructuredLogger provides advanced logging capabilities
type StructuredLogger struct {
	logger     *logrus.Logger
	config     *Config
	writers    map[string]io.Writer
	components map[string]*logrus.Entry
}

// NewStructuredLogger creates a new structured logger
func NewStructuredLogger(config *Config) (*StructuredLogger, error) {
	if config == nil {
		config = DefaultConfig()
	}
	
	logger := logrus.New()
	
	// Set log level
	level, err := logrus.ParseLevel(string(config.Level))
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %v", err)
	}
	logger.SetLevel(level)
	
	// Set formatter
	switch config.Format {
	case FormatJSON:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	case FormatText:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	default:
		return nil, fmt.Errorf("invalid log format: %s", config.Format)
	}
	
	sl := &StructuredLogger{
		logger:     logger,
		config:     config,
		writers:    make(map[string]io.Writer),
		components: make(map[string]*logrus.Entry),
	}
	
	// Set output
	if err := sl.setOutput(config.Output); err != nil {
		return nil, fmt.Errorf("failed to set log output: %v", err)
	}
	
	return sl, nil
}

// setOutput configures the log output
func (sl *StructuredLogger) setOutput(output string) error {
	switch strings.ToLower(output) {
	case "stdout":
		sl.logger.SetOutput(os.Stdout)
	case "stderr":
		sl.logger.SetOutput(os.Stderr)
	default:
		// File output with rotation
		if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %v", err)
		}
		
		writer := &lumberjack.Logger{
			Filename:   output,
			MaxSize:    sl.config.MaxSize,
			MaxBackups: sl.config.MaxBackups,
			MaxAge:     sl.config.MaxAge,
			Compress:   sl.config.Compress,
		}
		
		sl.logger.SetOutput(writer)
		sl.writers["main"] = writer
	}
	
	return nil
}

// GetLogger returns the underlying logrus logger
func (sl *StructuredLogger) GetLogger() *logrus.Logger {
	return sl.logger
}

// GetComponentLogger returns a logger for a specific component
func (sl *StructuredLogger) GetComponentLogger(component string) *logrus.Entry {
	if entry, exists := sl.components[component]; exists {
		return entry
	}
	
	entry := sl.logger.WithField("component", component)
	
	// Apply component-specific log level if configured
	if level, exists := sl.config.Components[component]; exists {
		if logLevel, err := logrus.ParseLevel(string(level)); err == nil {
			// Create a new logger instance for this component with specific level
			componentLogger := logrus.New()
			componentLogger.SetLevel(logLevel)
			componentLogger.SetFormatter(sl.logger.Formatter)
			componentLogger.SetOutput(sl.logger.Out)
			
			entry = componentLogger.WithField("component", component)
		}
	}
	
	sl.components[component] = entry
	return entry
}

// WithFields creates a new logger entry with additional fields
func (sl *StructuredLogger) WithFields(fields map[string]interface{}) *logrus.Entry {
	return sl.logger.WithFields(logrus.Fields(fields))
}

// WithField creates a new logger entry with an additional field
func (sl *StructuredLogger) WithField(key string, value interface{}) *logrus.Entry {
	return sl.logger.WithField(key, value)
}

// Trace logs a trace message
func (sl *StructuredLogger) Trace(args ...interface{}) {
	sl.logger.Trace(args...)
}

// Tracef logs a formatted trace message
func (sl *StructuredLogger) Tracef(format string, args ...interface{}) {
	sl.logger.Tracef(format, args...)
}

// Debug logs a debug message
func (sl *StructuredLogger) Debug(args ...interface{}) {
	sl.logger.Debug(args...)
}

// Debugf logs a formatted debug message
func (sl *StructuredLogger) Debugf(format string, args ...interface{}) {
	sl.logger.Debugf(format, args...)
}

// Info logs an info message
func (sl *StructuredLogger) Info(args ...interface{}) {
	sl.logger.Info(args...)
}

// Infof logs a formatted info message
func (sl *StructuredLogger) Infof(format string, args ...interface{}) {
	sl.logger.Infof(format, args...)
}

// Warn logs a warning message
func (sl *StructuredLogger) Warn(args ...interface{}) {
	sl.logger.Warn(args...)
}

// Warnf logs a formatted warning message
func (sl *StructuredLogger) Warnf(format string, args ...interface{}) {
	sl.logger.Warnf(format, args...)
}

// Error logs an error message
func (sl *StructuredLogger) Error(args ...interface{}) {
	sl.logger.Error(args...)
}

// Errorf logs a formatted error message
func (sl *StructuredLogger) Errorf(format string, args ...interface{}) {
	sl.logger.Errorf(format, args...)
}

// Fatal logs a fatal message and exits
func (sl *StructuredLogger) Fatal(args ...interface{}) {
	sl.logger.Fatal(args...)
}

// Fatalf logs a formatted fatal message and exits
func (sl *StructuredLogger) Fatalf(format string, args ...interface{}) {
	sl.logger.Fatalf(format, args...)
}

// Panic logs a panic message and panics
func (sl *StructuredLogger) Panic(args ...interface{}) {
	sl.logger.Panic(args...)
}

// Panicf logs a formatted panic message and panics
func (sl *StructuredLogger) Panicf(format string, args ...interface{}) {
	sl.logger.Panicf(format, args...)
}

// SetLevel sets the log level
func (sl *StructuredLogger) SetLevel(level LogLevel) error {
	logLevel, err := logrus.ParseLevel(string(level))
	if err != nil {
		return fmt.Errorf("invalid log level: %v", err)
	}
	
	sl.logger.SetLevel(logLevel)
	sl.config.Level = level
	return nil
}

// SetComponentLevel sets the log level for a specific component
func (sl *StructuredLogger) SetComponentLevel(component string, level LogLevel) {
	sl.config.Components[component] = level
	
	// Update existing component logger if it exists
	if _, exists := sl.components[component]; exists {
		delete(sl.components, component) // Force recreation with new level
	}
}

// GetLevel returns the current log level
func (sl *StructuredLogger) GetLevel() LogLevel {
	return LogLevel(sl.logger.GetLevel().String())
}

// Close closes any file writers
func (sl *StructuredLogger) Close() error {
	for _, writer := range sl.writers {
		if closer, ok := writer.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Rotate rotates log files (if using file output)
func (sl *StructuredLogger) Rotate() error {
	for _, writer := range sl.writers {
		if rotator, ok := writer.(*lumberjack.Logger); ok {
			return rotator.Rotate()
		}
	}
	return nil
}

// GetStats returns logging statistics
func (sl *StructuredLogger) GetStats() map[string]interface{} {
	stats := map[string]interface{}{
		"level":      sl.config.Level,
		"format":     sl.config.Format,
		"output":     sl.config.Output,
		"components": len(sl.components),
	}
	
	componentLevels := make(map[string]LogLevel)
	for component, level := range sl.config.Components {
		componentLevels[component] = level
	}
	stats["component_levels"] = componentLevels
	
	return stats
}
