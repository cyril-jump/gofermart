package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zap *zap.Logger
}

func New() *Logger {

	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.RFC3339TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, _ := cfg.Build()

	return &Logger{
		zap: logger,
	}
}

func (l *Logger) Close() {
	l.zap.Sync()
}

func (l *Logger) Info(msg string) {
	l.zap.Info(msg)
}

func (l *Logger) Warn(msg string, err error) {
	l.zap.Warn(msg, zap.Error(err))
}

func (l *Logger) Fatal(msg string, err error) {
	l.zap.Fatal(msg, zap.Error(err))
}
