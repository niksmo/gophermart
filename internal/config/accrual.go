package config

import (
	"fmt"
	"log"
	"net/url"

	"github.com/spf13/viper"
)

const (
	accrualEnv       = "ACCRUAL_SYSTEM_ADDRESS"
	accrualFlag      = "accrual"
	accrualFlagShort = "r"
	accrualUsage     = "accrual system address"
	accrualDefault   = "http://127.0.0.1:5050"
)

type AccrualAddrConfig struct {
	base      *url.URL
	ordersGet string
}

func NewAccrualAddrConfig() *AccrualAddrConfig {
	var config AccrualAddrConfig
	config.ordersGet = "/api/orders"

	flagValue := viper.GetString(accrualFlag)
	envValue := viper.GetString(accrualEnv)

	errParsePrefix := "parse accrual system address"

	if envValue != "" {
		baseURL, err := url.ParseRequestURI(envValue)
		if err == nil {
			config.base = baseURL
			return &config
		}
		log.Println(fmt.Errorf(
			errParsePrefix+" env value err: %w", err,
		))
	}

	baseURL, err := url.ParseRequestURI(flagValue)
	if err == nil {
		config.base = baseURL
		return &config
	}
	log.Println(fmt.Errorf(
		errParsePrefix+" flag value err: %w", err,
	))

	baseURL, _ = url.ParseRequestURI(accrualDefault)
	config.base = baseURL
	log.Println("use default accrual system address:", addrDefault)

	return &config
}

func (config *AccrualAddrConfig) GetOrdersReqURL(order string) string {
	return config.base.JoinPath(config.ordersGet, order).String()
}
