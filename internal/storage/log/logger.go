package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(level string) *zap.Logger {
	zLevel := getLevel(level)

	const (
		timestamp  = "@timestamp"
		severity   = "level"
		loggerName = "logger"
		caller     = "caller"
		message    = "message"
		stacktrace = "stacktrace"
	)

	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				TimeKey:        timestamp,
				LevelKey:       severity,
				NameKey:        loggerName,
				CallerKey:      caller,
				MessageKey:     message,
				StacktraceKey:  stacktrace,
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}),
			zapcore.AddSync(os.Stdout),
			zLevel,
		),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	return logger
}

func getLevel(level string) zap.AtomicLevel {
	zLevel := zapcore.InfoLevel

	switch level {
	case "info":
		zLevel = zapcore.InfoLevel
	case "warn":
		zLevel = zapcore.WarnLevel
	case "error":
		zLevel = zapcore.ErrorLevel
	}

	return zap.NewAtomicLevelAt(zLevel)
}
