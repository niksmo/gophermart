package config

import (
	"errors"
	"net/url"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

const (
	dbURIEnv       = "DATABASE_URI"
	dbURIFlag      = "dsn"
	dbURIFlagShort = "d"
	dbURIDefault   = ""
	dbURIUsage     = "example: " +
		"'postgres://user:pwd@127.0.0.1:5432/db_name/sslmode=disable'"
	dbURIFlagPrint = "-" + dbURIFlagShort
)

type DatabaseCofig struct {
	URI string
}

func NewDatabaseConfig(logger zerolog.Logger) (config DatabaseCofig) {
	flagValue := viper.GetString(dbURIFlag)
	envValue := viper.GetString(dbURIEnv)
	configLogger := logger.With().Str("config", "database").Logger()

	if envValue != "" {
		DSN, err := parseDSN(envValue)
		if err == nil {
			config = DatabaseCofig{URI: DSN}
			return
		}
		configLogger.Warn().
			Str("env", dbURIEnv).
			Str("value", envValue).
			Err(err).
			Send()
	}

	if flagValue != "" {
		DSN, err := parseDSN(flagValue)
		if err == nil {
			config = DatabaseCofig{URI: DSN}
			return
		}
		configLogger.Warn().
			Str("flag", dbURIFlagPrint).
			Str("value", flagValue).
			Err(err).
			Send()
	}

	configLogger.Fatal().Msg("URI not set")
	return
}

func parseDSN(value string) (string, error) {
	URL, err := url.ParseRequestURI(value)
	if err != nil {
		return "", err
	}
	if URL.Scheme == "" || URL.User == nil {
		return "", errors.New("invalid database DSN")
	}
	return URL.String(), nil
}
