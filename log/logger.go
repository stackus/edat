package log

import (
	"time"
)

type FieldType int

const (
	StringType FieldType = iota
	IntType
	DurationType
	ErrorType
)

type Field struct {
	Key      string
	Type     FieldType
	String   string
	Int      int
	Duration time.Duration
	Error    error
}

type Logger interface {
	Trace(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Sub(fields ...Field) Logger
}

var DefaultLogger Logger = NewNopLogger()

func String(key string, value string) Field {
	return Field{Type: StringType, Key: key, String: value}
}

func Int(key string, value int) Field {
	return Field{Type: IntType, Key: key, Int: value}
}

func Duration(key string, value time.Duration) Field {
	return Field{Type: DurationType, Key: key, Duration: value}
}

func Error(err error) Field {
	return Field{Type: ErrorType, Key: "error", Error: err}
}

type nopLogger struct{}

func NewNopLogger() Logger {
	return nopLogger{}
}

func (l nopLogger) Trace(string, ...Field) {}
func (l nopLogger) Debug(string, ...Field) {}
func (l nopLogger) Info(string, ...Field)  {}
func (l nopLogger) Warn(string, ...Field)  {}
func (l nopLogger) Error(string, ...Field) {}
func (l nopLogger) Sub(...Field) Logger    { return l }
