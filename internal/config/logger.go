package config

import (
	"github.com/niksmo/gophermart/internal/logger"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

const (
	logLevelEnv       = "LOG_LVL"
	logLevelFlag      = "logLevel"
	logLevelFlagShort = "l"
	logLevelUsage     = "log level: debug, info, warn, error, fatal, panic"
	logLevelDefault   = "info"
	logLevelFlagPrint = "-" + logLevelFlagShort
)

type LoggerConfig struct {
	Level zerolog.Level
}

func NewLoggerConfig() LoggerConfig {
	flagValue := viper.GetString(logLevelFlag)
	envValue := viper.GetString(logLevelEnv)
	log := logger.Instance.With().Str("config", "logger").Logger()

	if envValue != "" {
		level, err := zerolog.ParseLevel(envValue)
		envLog := log.With().
			Str("env", logLevelEnv).
			Str("value", envValue).
			Logger()
		if err == nil {
			envLog.Info().Msg("use env value")
			return LoggerConfig{Level: level}
		}
		envLog.Warn().Err(err).Send()
	}

	level, err := zerolog.ParseLevel(flagValue)
	flagLog := log.With().
		Str("flag", logLevelFlagPrint).
		Str("value", flagValue).
		Logger()
	if err == nil {
		flagLog.Info().Msg("use flag value")
		return LoggerConfig{Level: level}
	}
	flagLog.Warn().Err(err).Send()

	defaultLevel, _ := zerolog.ParseLevel(logLevelDefault)
	log.Info().
		Str("flag", logLevelFlagPrint).
		Msg("use default value")

	return LoggerConfig{Level: defaultLevel}
}
