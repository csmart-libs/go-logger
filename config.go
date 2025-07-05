package logger

import (
	"os"
)

// RotationMode defines how log files should be rotated
type RotationMode string

const (
	// RotationModeSize rotates based on file size only
	RotationModeSize RotationMode = "size"
	// RotationModeTime rotates based on time only
	RotationModeTime RotationMode = "time"
	// RotationModeBoth rotates based on both size and time (whichever comes first)
	RotationModeBoth RotationMode = "both"
)

// TimeRotationInterval defines the time interval for rotation
type TimeRotationInterval string

const (
	// RotationHourly rotates every hour
	RotationHourly TimeRotationInterval = "hourly"
	// RotationDaily rotates every day
	RotationDaily TimeRotationInterval = "daily"
	// RotationWeekly rotates every week
	RotationWeekly TimeRotationInterval = "weekly"
	// RotationMonthly rotates every month
	RotationMonthly TimeRotationInterval = "monthly"
)

// Environment constants
const (
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"
	EnvTest        = "test"
)

// Log level constants
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
	LevelPanic = "panic"
)

// Encoding constants
const (
	EncodingJSON    = "json"
	EncodingConsole = "console"
)

// FileOptions holds file-specific logging options
type FileOptions struct {
	// Filename is the file to write logs to. If empty, logs will only go to stdout
	Filename string `json:"filename" yaml:"filename"`

	// MaxSize is the maximum size in megabytes of the log file before it gets rotated
	MaxSize int `json:"max_size" yaml:"max_size"`

	// MaxAge is the maximum number of days to retain old log files
	MaxAge int `json:"max_age" yaml:"max_age"`

	// MaxBackups is the maximum number of old log files to retain
	MaxBackups int `json:"max_backups" yaml:"max_backups"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time. Default is UTC time.
	LocalTime bool `json:"local_time" yaml:"local_time"`

	// Compress determines if the rotated log files should be compressed using gzip
	Compress bool `json:"compress" yaml:"compress"`

	// FileMode is the file mode to use when creating log files
	FileMode os.FileMode `json:"file_mode" yaml:"file_mode"`

	// CreateDir determines if the directory should be created if it doesn't exist
	CreateDir bool `json:"create_dir" yaml:"create_dir"`

	// RotationMode determines how files should be rotated (size, time, or both)
	RotationMode RotationMode `json:"rotation_mode" yaml:"rotation_mode"`

	// TimeRotationInterval determines the time interval for rotation when using time-based rotation
	TimeRotationInterval TimeRotationInterval `json:"time_rotation_interval" yaml:"time_rotation_interval"`

	// TimeRotationFormat is the format string for time-based file naming
	// Default formats:
	// - Hourly: "2006-01-02-15"
	// - Daily: "2006-01-02"
	// - Weekly: "2006-W01"
	// - Monthly: "2006-01"
	TimeRotationFormat string `json:"time_rotation_format" yaml:"time_rotation_format"`
}

// Config holds logger configuration
type Config struct {
	Level       string      `json:"level" yaml:"level"`
	Environment string      `json:"environment" yaml:"environment"`
	OutputPaths []string    `json:"output_paths" yaml:"output_paths"`
	Encoding    string      `json:"encoding" yaml:"encoding"`
	FileOptions FileOptions `json:"file_options" yaml:"file_options"`
}

// DefaultFileOptions returns default file options
func DefaultFileOptions() FileOptions {
	return FileOptions{
		Filename:             "",               // Empty means no file output
		MaxSize:              100,              // 100 MB
		MaxAge:               30,               // 30 days
		MaxBackups:           10,               // Keep 10 backup files
		LocalTime:            false,            // Use UTC time
		Compress:             true,             // Compress old files
		FileMode:             0644,             // Standard file permissions
		CreateDir:            true,             // Create directory if not exists
		RotationMode:         RotationModeSize, // Default to size-based rotation
		TimeRotationInterval: RotationDaily,    // Default to daily rotation
		TimeRotationFormat:   "",               // Will be set based on interval
	}
}

// DefaultConfig returns default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:       "info",
		Environment: "development",
		OutputPaths: []string{"stdout"},
		Encoding:    "console",
		FileOptions: DefaultFileOptions(),
	}
}

// DevelopmentConfig returns configuration optimized for development
func DevelopmentConfig() Config {
	config := DefaultConfig()
	config.Level = LevelDebug
	config.Environment = EnvDevelopment
	config.OutputPaths = []string{"stdout"}
	config.Encoding = EncodingConsole
	return config
}

// ProductionConfig returns configuration optimized for production
func ProductionConfig() Config {
	config := DefaultConfig()
	config.Level = LevelInfo
	config.Environment = EnvProduction
	config.OutputPaths = []string{"stdout"}
	config.Encoding = EncodingJSON
	return config
}

// ProductionConfigWithFile returns production configuration with file output
func ProductionConfigWithFile(filename string) Config {
	config := ProductionConfig()
	config.FileOptions.Filename = filename
	config.OutputPaths = []string{"file"} // Only file output for production
	return config
}

// TestConfig returns configuration optimized for testing
func TestConfig() Config {
	config := DefaultConfig()
	config.Level = LevelError
	config.Environment = EnvTest
	config.OutputPaths = []string{"stdout"}
	config.Encoding = EncodingConsole
	config.FileOptions.Filename = "" // No file output for tests
	return config
}
