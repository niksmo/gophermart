package config

import (
	"errors"

	"github.com/niksmo/gophermart/pkg/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Init() {
	pflag.ErrHelp = errors.New("gophermart: help requested")
	pflag.StringP(addrFlag, addrFlagShort, addrDefault, addrUsage)
	pflag.StringP(dbURIFlag, dbURIFlagShort, dbURIDefault, dbURIUsage)
	pflag.StringP(
		accrualFlag, accrualFlagShort, accrualDefault, accrualUsage,
	)
	pflag.StringP(
		logLevelFlag, logLevelFlagShort, logLevelDefault, logLevelUsage,
	)
	pflag.IntP(costFlag, costFlagShort, costDefault, costUsage)
	pflag.StringP(keyFlag, keyFlagShort, keyDefault, keyUsage)
	pflag.DurationP(
		jwtLifetimeFlag, jwtLifetimeFlagShort, jwtLifetimeDefault, jwtLifetimeUsage,
	)
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		logger.Instance.Fatal().Err(err).Caller().Send()
	}

	envVars := []string{
		addrEnv,
		dbURIEnv,
		accrualEnv,
		logLevelEnv,
		costEnv,
		keyEnv,
		jwtLifetimeEvn,
	}

	for _, env := range envVars {
		err = viper.BindEnv(env)
		if err != nil {
			logger.Instance.Fatal().Err(err).Caller().Send()
		}
	}
}
