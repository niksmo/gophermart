package config

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

const (
	logLevelEnv       = "LOG_LVL"
	logLevelFlag      = "logLevel"
	logLevelFlagShort = "l"
	logLevelUsage     = "log level: debug, info, warn, error, fatal, panic"
	logLevelDefault   = "info"
)

type LoggerConfig struct {
	Level zerolog.Level
}

func NewLoggerConfig(logger *zerolog.Logger) *LoggerConfig {
	flagValue := viper.GetString(addrFlag)
	envValue := viper.GetString(addrEnv)

	if envValue != "" {
		level, err := zerolog.ParseLevel(envValue)
		if err == nil {
			logger.Info().
				Str("env", logLevelEnv).
				Str("value", envValue).
				Msg("used")
			return &LoggerConfig{Level: level}
		}
		logger.Warn().
			Str("env", logLevelEnv).
			Str("value", envValue).
			Err(err).
			Send()
	}

	level, err := zerolog.ParseLevel(flagValue)
	if err == nil {
		return &LoggerConfig{Level: level}
	}
	logger.Warn().
		Str("flag", "-"+logLevelFlagShort).
		Str("value", flagValue).
		Err(err).
		Send()

	defaultLevel, _ := zerolog.ParseLevel(logLevelDefault)
	logger.Info().
		Str("flag", "-"+logLevelFlagShort).
		Msg("used default value")

	return &LoggerConfig{Level: defaultLevel}
}
