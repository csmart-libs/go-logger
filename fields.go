package logger

import (
	"go.uber.org/zap"
)

// Common field helpers

// String creates a string field
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

// Int creates an int field
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Int64 creates an int64 field
func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

// Uint creates a uint field
func Uint(key string, val uint) zap.Field {
	return zap.Uint(key, val)
}

// Uint32 creates a uint32 field
func Uint32(key string, val uint32) zap.Field {
	return zap.Uint32(key, val)
}

// Uint64 creates a uint64 field
func Uint64(key string, val uint64) zap.Field {
	return zap.Uint64(key, val)
}

// Float64 creates a float64 field
func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

// Bool creates a bool field
func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

// Any creates a field with any value
func Any(key string, val any) zap.Field {
	return zap.Any(key, val)
}

// Error creates an error field
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Duration creates a duration field
func Duration(key string, val any) zap.Field {
	return zap.Any(key, val)
}
