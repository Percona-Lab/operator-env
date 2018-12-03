package logger

import (
	"github.com/rs/zerolog"
	"os"
)

// NewLogger return new zerolog.Logger instance with configuration based on the config.
func NewLogger(verbose bool) (logger *zerolog.Logger) {
	switch verbose {
	case true:
		devLogger := zerolog.New(zerolog.ConsoleWriter{
			NoColor: false,
			Out:     os.Stdout,
		}).Level(zerolog.DebugLevel).With().Timestamp().Logger()
		logger = &devLogger
	default:
		prodLogger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Timestamp().Logger()
		logger = &prodLogger
	}
	return logger
}
