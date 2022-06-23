package logger

import "go.uber.org/zap"

type Logger struct {
	Zap *zap.Logger
}

func New() *Logger {
	logger, _ := zap.NewProduction()
	return &Logger{
		Zap: logger,
	}
}

func (l *Logger) Close() {
	l.Zap.Sync()
}
