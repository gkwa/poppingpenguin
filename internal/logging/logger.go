package logging

import (
	"fmt"
	"io"
	"os"
)

// LogLevel defines the severity of a log message
type LogLevel int

const (
	// ErrorLevel logs only errors
	ErrorLevel LogLevel = iota
	// WarningLevel logs warnings and errors
	WarningLevel
	// InfoLevel logs info, warnings, and errors
	InfoLevel
	// DebugLevel logs debug, info, warnings, and errors
	DebugLevel
)

// Logger defines the interface for logging
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Error(format string, args ...interface{})
}

// ConsoleLogger implements Logger for console output
type ConsoleLogger struct {
	level  LogLevel
	writer io.Writer
}

// NewLogger creates a new logger with the specified verbosity level
func NewLogger(verbosity int) Logger {
	level := ErrorLevel

	switch verbosity {
	case 1:
		level = WarningLevel
	case 2:
		level = InfoLevel
	case 3:
		level = DebugLevel
	}

	return &ConsoleLogger{
		level:  level,
		writer: os.Stderr,
	}
}

// Debug logs debug messages
func (l *ConsoleLogger) Debug(format string, args ...interface{}) {
	if l.level >= DebugLevel {
		fmt.Fprintf(l.writer, "DEBUG: "+format+"\n", args...)
	}
}

// Info logs informational messages
func (l *ConsoleLogger) Info(format string, args ...interface{}) {
	if l.level >= InfoLevel {
		fmt.Fprintf(l.writer, "INFO: "+format+"\n", args...)
	}
}

// Warning logs warning messages
func (l *ConsoleLogger) Warning(format string, args ...interface{}) {
	if l.level >= WarningLevel {
		fmt.Fprintf(l.writer, "WARNING: "+format+"\n", args...)
	}
}

// Error logs error messages
func (l *ConsoleLogger) Error(format string, args ...interface{}) {
	if l.level >= ErrorLevel {
		fmt.Fprintf(l.writer, "ERROR: "+format+"\n", args...)
	}
}
