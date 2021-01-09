package log

import (
	"time"
)

// FieldType loggable field types
type FieldType int

// Loggable field types
const (
	StringType FieldType = iota
	IntType
	DurationType
	ErrorType
)

// Field type
type Field struct {
	Key      string
	Type     FieldType
	String   string
	Int      int
	Duration time.Duration
	Error    error
}

// Logger interface
type Logger interface {
	Trace(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Sub(fields ...Field) Logger
}

// DefaultLogger is set to a Nop logger
//
// You may reassign this if you wish to avoid having to pass in With*Logger(yourLogger) options
// into many of the constructors to set a custom logger
var DefaultLogger Logger = NewNopLogger()

// String field value for strings
func String(key string, value string) Field {
	return Field{Type: StringType, Key: key, String: value}
}

// Int field value for ints
func Int(key string, value int) Field {
	return Field{Type: IntType, Key: key, Int: value}
}

// Duration field value to time.Durations
func Duration(key string, value time.Duration) Field {
	return Field{Type: DurationType, Key: key, Duration: value}
}

// Error field value for errors
func Error(err error) Field {
	return Field{Type: ErrorType, Key: "error", Error: err}
}

type nopLogger struct{}

// NewNopLogger returns a logger that logs to the void
func NewNopLogger() Logger {
	return nopLogger{}
}

func (l nopLogger) Trace(string, ...Field) {}
func (l nopLogger) Debug(string, ...Field) {}
func (l nopLogger) Info(string, ...Field)  {}
func (l nopLogger) Warn(string, ...Field)  {}
func (l nopLogger) Error(string, ...Field) {}
func (l nopLogger) Sub(...Field) Logger    { return l }
