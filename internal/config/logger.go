package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger, _ = Cfg.Build()
var Cfg = zap.Config{
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
