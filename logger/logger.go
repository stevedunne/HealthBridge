package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Level - the current logging level
	Level zapcore.Level
)

// Flags registers the command options for logging
func Flags() {
	zap.LevelFlag("log-level", zap.ErrorLevel, "set the logging level")
}

// NewLogger creates a new instance of the logger
func NewLogger() (*zap.Logger, error) {

	var config = getDevConfig()
	config.Level = zap.NewAtomicLevelAt(Level)

	logger, err := config.Build()

	return logger, err
}

func getProdConfig() *zap.Config {
	return &zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func getDevConfig() *zap.Config {
	return &zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}
