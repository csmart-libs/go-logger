package logger

import (
	"os"
	"strconv"
	"strings"
)

// ConfigFromEnv creates logger configuration from environment variables
func ConfigFromEnv() Config {
	config := DefaultConfig()

	// Get environment
	if env := os.Getenv("APP_ENV"); env != "" {
		config.Environment = strings.ToLower(env)
	}
	if env := os.Getenv("GIN_MODE"); env != "" {
		switch env {
		case "release":
			config.Environment = EnvProduction
		case "test":
			config.Environment = EnvTest
		default:
			config.Environment = EnvDevelopment
		}
	}

	// Get log level
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Level = strings.ToLower(level)
	}

	// Get encoding
	if encoding := os.Getenv("LOG_ENCODING"); encoding != "" {
		config.Encoding = strings.ToLower(encoding)
	}

	// Get output paths
	if outputs := os.Getenv("LOG_OUTPUT_PATHS"); outputs != "" {
		config.OutputPaths = strings.Split(outputs, ",")
	}

	// Get file options from environment
	if filename := os.Getenv("LOG_FILE"); filename != "" {
		config.FileOptions.Filename = filename
	}
	if maxSize := os.Getenv("LOG_FILE_MAX_SIZE"); maxSize != "" {
		if size, err := strconv.Atoi(maxSize); err == nil {
			config.FileOptions.MaxSize = size
		}
	}
	if maxAge := os.Getenv("LOG_FILE_MAX_AGE"); maxAge != "" {
		if age, err := strconv.Atoi(maxAge); err == nil {
			config.FileOptions.MaxAge = age
		}
	}
	if maxBackups := os.Getenv("LOG_FILE_MAX_BACKUPS"); maxBackups != "" {
		if backups, err := strconv.Atoi(maxBackups); err == nil {
			config.FileOptions.MaxBackups = backups
		}
	}
	if localTime := os.Getenv("LOG_FILE_LOCAL_TIME"); localTime != "" {
		config.FileOptions.LocalTime = strings.ToLower(localTime) == "true"
	}
	if compress := os.Getenv("LOG_FILE_COMPRESS"); compress != "" {
		config.FileOptions.Compress = strings.ToLower(compress) == "true"
	}
	if createDir := os.Getenv("LOG_FILE_CREATE_DIR"); createDir != "" {
		config.FileOptions.CreateDir = strings.ToLower(createDir) == "true"
	}
	if rotationMode := os.Getenv("LOG_FILE_ROTATION_MODE"); rotationMode != "" {
		config.FileOptions.RotationMode = RotationMode(strings.ToLower(rotationMode))
	}
	if timeInterval := os.Getenv("LOG_FILE_TIME_INTERVAL"); timeInterval != "" {
		config.FileOptions.TimeRotationInterval = TimeRotationInterval(strings.ToLower(timeInterval))
	}
	if timeFormat := os.Getenv("LOG_FILE_TIME_FORMAT"); timeFormat != "" {
		config.FileOptions.TimeRotationFormat = timeFormat
	}

	// Adjust config based on environment
	switch config.Environment {
	case EnvProduction:
		if config.Level == "" {
			config.Level = LevelInfo
		}
		config.Encoding = EncodingJSON
	case EnvStaging:
		if config.Level == "" {
			config.Level = LevelInfo
		}
		config.Encoding = EncodingJSON
	case EnvTest:
		if config.Level == "" {
			config.Level = LevelError
		}
		config.Encoding = EncodingConsole
	default: // development
		if config.Level == "" {
			config.Level = LevelDebug
		}
		config.Encoding = EncodingConsole
	}

	return config
}

// GetEffectiveConfig returns the effective configuration after applying defaults and validation
func GetEffectiveConfig() Config {
	config := ConfigFromEnv()
	config.Validate()
	return config
}
