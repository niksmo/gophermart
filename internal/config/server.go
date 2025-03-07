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

type ServerConfig struct {
	*net.TCPAddr
}

func NewServerConfig() ServerConfig {
	return ServerConfig{TCPAddr: getAddress()}
}

func (c *ServerConfig) Addr() string {
	return c.String()
}

func getAddress() *net.TCPAddr {
	flagValue := viper.GetString(addrFlag)
	envValue := viper.GetString(addrEnv)
	log := logger.Instance.With().Str("config", "server address").Logger()

	if envValue != "" {
		TCPAddr, err := net.ResolveTCPAddr("", envValue)
		envLog := log.With().Str("env", addrEnv).Str("value", envValue).Logger()
		if err == nil {
			envLog.Info().Msg("use env value")
			return TCPAddr
		}
		envLog.Warn().Err(err)
	}

	if flagValue != addrDefault {
		TCPAddr, err := net.ResolveTCPAddr("", flagValue)
		flagLog := log.With().
			Str("flag", addrFlagPrint).
			Str("value", flagValue).
			Logger()
		if err == nil {
			flagLog.Info().Msg("use flag value")
			return TCPAddr
		}
		flagLog.Warn().Err(err).Send()
	}

	TCPAddr, _ := net.ResolveTCPAddr("", addrDefault)
	log.Info().Str("flag", addrFlagPrint).Msg("use default value")

	return TCPAddr
}
