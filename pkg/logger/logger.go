package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Instance zerolog.Logger

func Init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	zerolog.MessageFieldName = "msg"
	Instance = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func SetLevel(level zerolog.Level) {
	zerolog.SetGlobalLevel(zerolog.Level(level))
}
