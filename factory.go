package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Global logger instance
var globalLogger Logger

// Initialize initializes the global logger with the given configuration
func Initialize(config Config) error {
	logger, err := NewLogger(config)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// NewLogger creates a new logger instance with the given configuration
func NewLogger(config Config) (Logger, error) {
	// Parse log level
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// Create encoder config based on environment
	var encoderConfig zapcore.EncoderConfig
	if config.Environment == "production" {
		encoderConfig = zap.NewProductionEncoderConfig()
		config.Encoding = "json"
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Configure time encoding
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create encoder
	var encoder zapcore.Encoder
	if config.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Create writer syncer
	var writeSyncer zapcore.WriteSyncer

	// Check if we need file output
	if config.FileOptions.Filename != "" {
		// Create directory if needed
		if config.FileOptions.CreateDir {
			dir := filepath.Dir(config.FileOptions.Filename)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, err
			}
		}

		var fileWriter io.Writer

		// Choose writer based on rotation mode
		switch config.FileOptions.RotationMode {
		case RotationModeTime, RotationModeBoth:
			// Use time-based rotating writer
			fileWriter = NewTimeRotatingWriter(config.FileOptions)
		default:
			// Use size-based rotating writer (lumberjack)
			fileWriter = &lumberjack.Logger{
				Filename:   config.FileOptions.Filename,
				MaxSize:    config.FileOptions.MaxSize,
				MaxAge:     config.FileOptions.MaxAge,
				MaxBackups: config.FileOptions.MaxBackups,
				LocalTime:  config.FileOptions.LocalTime,
				Compress:   config.FileOptions.Compress,
			}
		}

		// Combine stdout and file output if needed
		if len(config.OutputPaths) > 0 && config.OutputPaths[0] != "stdout" {
			// Only file output
			writeSyncer = zapcore.AddSync(fileWriter)
		} else {
			// Both stdout and file output
			writeSyncer = zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(os.Stdout),
				zapcore.AddSync(fileWriter),
			)
		}
	} else {
		// Only stdout output
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// Create core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// Create logger
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &ZapLogger{logger: zapLogger}, nil
}

// GetLogger returns the global logger instance
func GetLogger() Logger {
	if globalLogger == nil {
		// Initialize with default config if not initialized
		config := DefaultConfig()
		if env := os.Getenv("APP_ENV"); env != "" {
			config.Environment = env
		}
		if level := os.Getenv("LOG_LEVEL"); level != "" {
			config.Level = strings.ToLower(level)
		}
		Initialize(config)
	}
	return globalLogger
}

// Global logger functions

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

// Panic logs a panic message and panics
func Panic(msg string, fields ...zap.Field) {
	GetLogger().Panic(msg, fields...)
}

// With creates a child logger with additional fields
func With(fields ...zap.Field) Logger {
	return GetLogger().With(fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return GetLogger().Sync()
}
