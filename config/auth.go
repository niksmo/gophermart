package config

import (
	"strconv"

	"github.com/niksmo/gophermart/pkg/logger"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

const (
	costEnv       = "BCRYPT_COST"
	costFlag      = "cost"
	costFlagShort = "c"
	costUsage     = "hashing cost"
	costDefault   = bcrypt.DefaultCost
	costFlagPrint = "-" + costFlagShort
)

type AuthConfig struct {
	cost int
}

func NewAuthConfig() AuthConfig {
	return AuthConfig{cost: getCost()}
}

func (c *AuthConfig) Cost() int {
	return c.cost
}

func getCost() int {
	flagValue := viper.GetInt(costFlag)
	envValue := viper.GetString(costEnv)
	log := logger.Instance.With().Str("config", "auth cost").Logger()

	if envValue != "" {
		cost, err := strconv.Atoi(envValue)
		envLog := log.With().Str("env", costEnv).Str("value", envValue).Logger()
		if err == nil {
			envLog.Info().Msg("use env value")
			return cost
		}
		envLog.Warn().Err(err)
	}

	if flagValue != costDefault {
		log.Info().
			Str("flag", costFlagPrint).
			Int("value", flagValue).
			Msg("use flag value")
		return flagValue
	}

	log.Info().Str("flag", costFlagPrint).Msg("use default value")

	return costDefault
}
