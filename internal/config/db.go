package config

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/spf13/viper"
)

const (
	dbURIEnv       = "DATABASE_URI"
	dbURIFlag      = "dsn"
	dbURIFlagShort = "d"
	dbURIDefault   = ""
	dbURIUsage     = "example: " +
		"'postgres://user:pwd@127.0.0.1:5432/db_name/sslmode=disable'"
)

type DatabaseCofig struct {
	URI string
}

func NewDatabaseConfig() *DatabaseCofig {
	flagValue := viper.GetString(dbURIFlag)
	envValue := viper.GetString(dbURIEnv)

	if envValue != "" {
		DSN, err := parseDSN(envValue)
		if err != nil {
			panic(fmt.Errorf("%w; %s", err, dbURIUsage))
		}
		return &DatabaseCofig{URI: DSN}
	}

	if flagValue != "" {
		DSN, err := parseDSN(flagValue)
		if err != nil {
			panic(fmt.Errorf("%w; %s", err, dbURIUsage))
		}
		return &DatabaseCofig{URI: DSN}
	}

	panic("Database URI not set; " + dbURIUsage)
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
