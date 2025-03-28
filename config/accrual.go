package config

import (
	"net/url"

	"github.com/niksmo/gophermart/pkg/logger"
	"github.com/spf13/viper"
)

const (
	accrualEnv       = "ACCRUAL_SYSTEM_ADDRESS"
	accrualFlag      = "accrual"
	accrualFlagShort = "r"
	accrualUsage     = "accrual system address"
	accrualDefault   = "http://127.0.0.1:5050"
	accrualFlagPrint = "-" + accrualFlagShort

	ordersGetPath = "/api/orders"
)

type AccrualConfig struct {
	base          *url.URL
	ordersGetPath string
}

func NewAccrualConfig() AccrualConfig {
	var config AccrualConfig
	config.ordersGetPath = ordersGetPath

	flagValue := viper.GetString(accrualFlag)
	envValue := viper.GetString(accrualEnv)
	log := logger.Instance.With().Str("config", "accrual address").Logger()

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

	if flagValue != accrualDefault {
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
	}

	baseURL, _ := url.ParseRequestURI(accrualDefault)
	config.base = baseURL
	log.Info().Str("flag", accrualFlagPrint).Msg("use default value")

	return config
}

func (config *AccrualConfig) GetOrdersReqURL(order string) string {
	return config.base.JoinPath(config.ordersGetPath, order).String()
}
