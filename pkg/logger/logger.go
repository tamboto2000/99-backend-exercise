// Package logger provides simple APIs for logging.
// It already comes with its default log handler that you can directly use,
// or create your own handler implementation.
package logger

import "os"

// Attributes key for default
// log attributes
const (
	TimeKey    = "ts"
	LevelKey   = "level"
	CallerKey  = "caller"
	MessageKey = "msg"
)

// Log levels
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

type Attr struct {
	Key string
	Val any
}

// Any simplify creation of Attr
func Any(key string, val any) Attr {
	return Attr{
		Key: key,
		Val: val,
	}
}

// LogHandler provide API for implementin log handler.
// You can create your own handler by implementing this
// interface
type LogHandler interface {
	Log(lvl string, msg string, attrs ...Attr)
}

var gLog = func() *Logger {
	h, _ := NewDefaultHandler(DebugLevel, os.Stderr)
	return NewLogger(h)
}()

// SetDefault set default global logger
func SetDefault(logger *Logger) {
	gLog = logger
}

// GetDefault get a copy of default global logger
func GetDefault() *Logger {
	gLogC := *gLog
	return &gLogC
}

// Logger handles logging with multiple level.
// It comes with logging methods ranging from
// level info to fatal.
type Logger struct {
	h      LogHandler
	wattrs []Attr
}

func NewLogger(h LogHandler) *Logger {
	return &Logger{h: h}
}

// WithAttrs set attrs to wattrs, wattrs will be
// included as log attributes everytime you call
// any logging methods. This is a convenient way
// to convey a context of the log that being produced,
func (l *Logger) WithAttrs(attrs ...Attr) {
	l.wattrs = attrs
}

// AppendAttrs will only append to existing attributes
// in wattrs
func (l *Logger) AppendAttrs(attrs ...Attr) {
	l.wattrs = append(l.wattrs, attrs...)
}

// Debug logs a message at DebugLevel.
// The message includes any attrs passed at the log site, as well as wattrs set by WithAttrs or AppendAttrs.
func (l *Logger) Debug(msg string, attrs ...Attr) {
	attrs = append(l.wattrs, attrs...)
	l.h.Log(DebugLevel, msg, attrs...)
}

// Info logs a message at InfoLevel.
// The message includes any attrs passed at the log site, as well as wattrs set by WithAttrs or AppendAttrs.
func (l *Logger) Info(msg string, attrs ...Attr) {
	attrs = append(l.wattrs, attrs...)
	l.h.Log(InfoLevel, msg, attrs...)
}

// Warn logs a message at WarnLevel.
// The message includes any attrs passed at the log site, as well as wattrs set by WithAttrs or AppendAttrs.
func (l *Logger) Warn(msg string, attrs ...Attr) {
	attrs = append(l.wattrs, attrs...)
	l.h.Log(WarnLevel, msg, attrs...)
}

// Error logs a message at ErrorLevel.
// The message includes any attrs passed at the log site, as well as wattrs set by WithAttrs or AppendAttrs.
func (l *Logger) Error(msg string, attrs ...Attr) {
	attrs = append(l.wattrs, attrs...)
	l.h.Log(ErrorLevel, msg, attrs...)
}

// Fatal logs a message at FatalLevel.
// The message includes any attrs passed at the log site, as well as wattrs set by WithAttrs or AppendAttrs.
func (l *Logger) Fatal(msg string, attrs ...Attr) {
	attrs = append(l.wattrs, attrs...)
	l.h.Log(FatalLevel, msg, attrs...)
}

// Debug logs a message at DebugLevel.
// This function is using default global logger.
// The message includes any attrs passed at the log site, as well as wattrs set by WithAttrs or AppendAttrs.
func Debug(msg string, attrs ...Attr) {
	attrs = append(gLog.wattrs, attrs...)
	gLog.h.Log(DebugLevel, msg, attrs...)
}

// Info logs a message at InfoLevel.
// This function is using default global logger.
// The message includes any attrs passed at the log site, as well as wattrs set by WithAttrs or AppendAttrs.
func Info(msg string, attrs ...Attr) {
	attrs = append(gLog.wattrs, attrs...)
	gLog.h.Log(InfoLevel, msg, attrs...)
}

// Warn logs a message at WarnLevel.
// This function is using default global logger.
// The message includes any attrs passed at the log site, as well as wattrs set by WithAttrs or AppendAttrs.
func Warn(msg string, attrs ...Attr) {
	attrs = append(gLog.wattrs, attrs...)
	gLog.h.Log(WarnLevel, msg, attrs...)
}

// Error logs a message at ErrorLevel.
// This function is using default global logger.
// The message includes any attrs passed at the log site, as well as wattrs set by WithAttrs or AppendAttrs.
func Error(msg string, attrs ...Attr) {
	attrs = append(gLog.wattrs, attrs...)
	gLog.h.Log(ErrorLevel, msg, attrs...)
}

// Fatal logs a message at FatalLevel.
// This function is using default global logger.
// The message includes any attrs passed at the log site, as well as wattrs set by WithAttrs or AppendAttrs.
func Fatal(msg string, attrs ...Attr) {
	attrs = append(gLog.wattrs, attrs...)
	gLog.h.Log(FatalLevel, msg, attrs...)
}
