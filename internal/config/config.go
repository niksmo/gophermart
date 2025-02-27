package config

import (
	"errors"
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Init() {
	pflag.ErrHelp = errors.New("gophermart: help requested")
	pflag.StringP(addrFlag, addrFlagShort, addrDefault, addrUsage)
	pflag.StringP(dbURIFlag, dbURIFlagShort, dbURIDefault, dbURIUsage)
	pflag.StringP(accrualFlag, accrualFlagShort, accrualDefault, accrualUsage)
	pflag.Parse()

	var err error
	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Panic(err)
	}

	for _, env := range []string{addrEnv, dbURIEnv, accrualEnv} {
		err = viper.BindEnv(env)
		if err != nil {
			log.Panic(err)
		}
	}
}
