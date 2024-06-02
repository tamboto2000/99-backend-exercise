package logger

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogHandler struct {
	l *zap.Logger
}

// NewDefaultHandler created the default log handler implementation
// for Logger with configurable log level by lvl
func NewDefaultHandler(lvl string, w io.Writer) (LogHandler, error) {
	config := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = TimeKey
	encoderConfig.StacktraceKey = zapcore.OmitKey // to hide stacktrace info
	encoderConfig.LevelKey = LevelKey
	encoderConfig.CallerKey = CallerKey
	encoderConfig.MessageKey = MessageKey
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig = encoderConfig

	if lvl != "" {
		logLevel, err := zapcore.ParseLevel(lvl)
		if err != nil {
			panic(err)
		}
		config.Level.SetLevel(logLevel)
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	wsync := zapcore.AddSync(w)
	core := zapcore.NewCore(encoder, wsync, config.Level)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

	zl := zapLogHandler{l: logger}

	return &zl, nil
}

// Log implements LogHandler
func (zl *zapLogHandler) Log(lvl string, msg string, attrs ...Attr) {
	zl.l.Log(logLvlToZapLvl(lvl), msg, attrsToZapFields(attrs...)...)
}

func logLvlToZapLvl(lvl string) zapcore.Level {
	switch lvl {
	case DebugLevel:
		return zapcore.DebugLevel

	case InfoLevel:
		return zapcore.InfoLevel

	case WarnLevel:
		return zapcore.WarnLevel

	case ErrorLevel:
		return zapcore.ErrorLevel

	case FatalLevel:
		return zapcore.FatalLevel
	}

	return zapcore.InvalidLevel
}

func attrsToZapFields(attrs ...Attr) []zap.Field {
	var fields []zap.Field
	for _, a := range attrs {
		fields = append(fields, zap.Any(a.Key, a.Val))
	}

	return fields
}
