package logger

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// TimeRotatingWriter wraps lumberjack.Logger to support time-based rotation
type TimeRotatingWriter struct {
	*lumberjack.Logger
	options           FileOptions
	currentTimeFormat string
	lastRotationTime  time.Time
	mu                sync.Mutex
	baseFilename      string
}

// NewTimeRotatingWriter creates a new time-based rotating writer
func NewTimeRotatingWriter(options FileOptions) *TimeRotatingWriter {
	// Extract base filename and extension
	baseFilename := options.Filename

	// Set time format based on interval
	timeFormat := options.TimeRotationFormat
	if timeFormat == "" {
		switch options.TimeRotationInterval {
		case RotationHourly:
			timeFormat = "2006-01-02-15"
		case RotationDaily:
			timeFormat = "2006-01-02"
		case RotationWeekly:
			timeFormat = "2006-W01"
		case RotationMonthly:
			timeFormat = "2006-01"
		default:
			timeFormat = "2006-01-02"
		}
	}

	// Create initial filename with timestamp
	now := time.Now()
	if options.LocalTime {
		now = now.Local()
	} else {
		now = now.UTC()
	}

	timestampedFilename := generateTimestampedFilename(baseFilename, now, timeFormat)

	lj := &lumberjack.Logger{
		Filename:   timestampedFilename,
		MaxSize:    options.MaxSize,
		MaxAge:     options.MaxAge,
		MaxBackups: options.MaxBackups,
		LocalTime:  options.LocalTime,
		Compress:   options.Compress,
	}

	return &TimeRotatingWriter{
		Logger:            lj,
		options:           options,
		currentTimeFormat: timeFormat,
		lastRotationTime:  now,
		baseFilename:      baseFilename,
	}
}

// Write implements io.Writer interface with time-based rotation check
func (w *TimeRotatingWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	if w.options.LocalTime {
		now = now.Local()
	} else {
		now = now.UTC()
	}

	// Check if we need to rotate based on time
	if w.shouldRotateByTime(now) {
		if err := w.rotateByTime(now); err != nil {
			return 0, err
		}
	}

	return w.Logger.Write(p)
}

// shouldRotateByTime checks if rotation is needed based on time interval
func (w *TimeRotatingWriter) shouldRotateByTime(now time.Time) bool {
	switch w.options.TimeRotationInterval {
	case RotationHourly:
		return now.Hour() != w.lastRotationTime.Hour() ||
			now.Day() != w.lastRotationTime.Day() ||
			now.Month() != w.lastRotationTime.Month() ||
			now.Year() != w.lastRotationTime.Year()
	case RotationDaily:
		return now.Day() != w.lastRotationTime.Day() ||
			now.Month() != w.lastRotationTime.Month() ||
			now.Year() != w.lastRotationTime.Year()
	case RotationWeekly:
		_, thisWeek := now.ISOWeek()
		_, lastWeek := w.lastRotationTime.ISOWeek()
		return thisWeek != lastWeek || now.Year() != w.lastRotationTime.Year()
	case RotationMonthly:
		return now.Month() != w.lastRotationTime.Month() ||
			now.Year() != w.lastRotationTime.Year()
	}
	return false
}

// rotateByTime performs time-based rotation
func (w *TimeRotatingWriter) rotateByTime(now time.Time) error {
	// Close current file
	if err := w.Logger.Close(); err != nil {
		return err
	}

	// Generate new filename with current timestamp
	newFilename := generateTimestampedFilename(w.baseFilename, now, w.currentTimeFormat)

	// Update lumberjack logger with new filename
	w.Logger.Filename = newFilename
	w.lastRotationTime = now

	return nil
}

// generateTimestampedFilename creates a filename with timestamp
func generateTimestampedFilename(baseFilename string, t time.Time, timeFormat string) string {
	dir := filepath.Dir(baseFilename)
	filename := filepath.Base(baseFilename)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	timestamp := t.Format(timeFormat)
	timestampedName := fmt.Sprintf("%s-%s%s", nameWithoutExt, timestamp, ext)

	return filepath.Join(dir, timestampedName)
}
