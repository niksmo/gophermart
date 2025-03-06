package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return logger
}

func SetLevel(level zerolog.Level) {
	zerolog.SetGlobalLevel(zerolog.Level(level))
}
