package config

import (
	"net/url"

	"github.com/niksmo/gophermart/internal/logger"
	"github.com/spf13/viper"
)

const (
	accrualEnv       = "ACCRUAL_SYSTEM_ADDRESS"
	accrualFlag      = "accrual"
	accrualFlagShort = "r"
	accrualUsage     = "accrual system address"
	accrualDefault   = "http://127.0.0.1:5050"
	accrualFlagPrint = "-" + accrualFlagShort
)

type AccrualAddrConfig struct {
	base      *url.URL
	ordersGet string
}

func NewAccrualAddrConfig() AccrualAddrConfig {
	var config AccrualAddrConfig
	config.ordersGet = "/api/orders"

	flagValue := viper.GetString(accrualFlag)
	envValue := viper.GetString(accrualEnv)
	log := logger.Instance.With().Str("config", "accrualAddress").Logger()

	if envValue != "" {
		baseURL, err := url.ParseRequestURI(envValue)
		envLog := log.With().
			Str("env", accrualEnv).
			Str("value", envValue).
			Logger()
		if err == nil {
			config.base = baseURL
			envLog.Info().Msg("use env value")
			return config
		}
		envLog.Warn().Err(err).Send()
	}

	baseURL, err := url.ParseRequestURI(flagValue)
	flagLog := log.With().
		Str("flag", accrualFlagPrint).
		Str("value", flagValue).
		Logger()
	if err == nil {
		config.base = baseURL
		flagLog.Info().Msg("use flag value")
		return config
	}
	flagLog.Warn().Err(err).Send()

	baseURL, _ = url.ParseRequestURI(accrualDefault)
	config.base = baseURL
	log.Info().Str("flag", accrualFlagPrint).Msg("use default value")

	return config
}

func (config *AccrualAddrConfig) GetOrdersReqURL(order string) string {
	return config.base.JoinPath(config.ordersGet, order).String()
}
