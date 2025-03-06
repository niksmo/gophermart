package config

import (
	"net"

	"github.com/niksmo/gophermart/internal/logger"
	"github.com/spf13/viper"
)

const (
	addrEnv       = "RUN_ADDRESS"
	addrFlag      = "addr"
	addrFlagShort = "a"
	addrUsage     = "run address"
	addrDefault   = "127.0.0.1:8080"
	addrFlagPrint = "-" + addrFlagShort
)

type AddressConfig struct {
	*net.TCPAddr
}

func NewAddressConfig() AddressConfig {
	flagValue := viper.GetString(addrFlag)
	envValue := viper.GetString(addrEnv)
	log := logger.Instance.With().Str("config", "runAddress").Logger()

	if envValue != "" {
		TCPAddr, err := net.ResolveTCPAddr("", envValue)
		envLog := log.With().Str("env", addrEnv).Str("value", envValue).Logger()
		if err == nil {
			envLog.Info().Msg("use env value")
			return AddressConfig{TCPAddr}
		}
		envLog.Warn().Err(err)
	}

	TCPAddr, err := net.ResolveTCPAddr("", flagValue)
	flagLog := log.With().
		Str("flag", addrFlagPrint).
		Str("value", flagValue).
		Logger()
	if err == nil {
		flagLog.Info().Msg("use flag value")
		return AddressConfig{TCPAddr}
	}
	flagLog.Warn().Err(err).Send()

	TCPAddr, _ = net.ResolveTCPAddr("", addrDefault)
	log.Info().Str("flag", addrFlagPrint).Msg("use default value")

	return AddressConfig{TCPAddr}
}

func (c *AddressConfig) Addr() string {
	return c.String()
}
