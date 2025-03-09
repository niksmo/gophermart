package config

import (
	"errors"
	"net/url"

	"github.com/niksmo/gophermart/pkg/logger"
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
	dsn string
}

func NewDatabaseConfig() (config DatabaseCofig) {
	flagValue := viper.GetString(dbURIFlag)
	envValue := viper.GetString(dbURIEnv)
	log := logger.Instance.With().Str("config", "database").Logger()

	if envValue != "" {
		DSN, err := parseDSN(envValue)
		envLog := log.With().
			Str("env", dbURIEnv).
			Str("value", envValue).
			Logger()
		if err == nil {
			config = DatabaseCofig{dsn: DSN}
			envLog.Info().Msg("use env value")
			return
		}
		envLog.Warn().Err(err).Send()
	}

	if flagValue != "" {
		DSN, err := parseDSN(flagValue)
		flagLog := log.With().
			Str("flag", dbURIFlagPrint).
			Str("value", flagValue).
			Logger()
		if err == nil {
			config = DatabaseCofig{dsn: DSN}
			flagLog.Info().Msg("use flag value")
			return
		}
		flagLog.Warn().Err(err).Send()
	}

	log.Fatal().Msg("URI not set")
	return
}

func (config DatabaseCofig) URI() string {
	return config.dsn
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
