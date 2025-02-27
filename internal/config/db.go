package config

import (
	"github.com/spf13/viper"
)

const (
	dbURIEnv       = "DATABASE_URI"
	dbURIFlag      = "dsn"
	dbURIFlagShort = "d"
	dbURIUsage     = "Database DSN"
	dbURIDefault   = ""
)

type DatabaseCofig struct{}

func NewDatabaseConfig() *DatabaseCofig {
	flagValue := viper.GetString(dbURIFlag)
	envValue := viper.GetString(dbURIEnv)

	if envValue != "" {

	}

	if flagValue != "" {

	}

	return nil
}
