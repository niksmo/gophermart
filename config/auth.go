package config

import (
	"time"

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

	keyEnv       = "JWT_KEY"
	keyFlag      = "jkey"
	keyFlagShort = "k"
	keyUsage     = "jwt signing key"
	keyDefault   = "do_not_use_default_key_in_production"
	keyFlagPrint = "-" + keyFlagShort

	jwtLifetimeEvn       = "JWT_LIFETIME"
	jwtLifetimeFlag      = "jt"
	jwtLifetimeFlagShort = "t"
	jwtLifetimeUsage     = "jwt lifetime"
	jwtLifetimeDefault   = 5 * 24 * time.Hour // 5 days
	jwtLifetimeFlagPrint = "-" + jwtLifetimeFlagShort
)

type AuthConfig struct {
	cost        int
	key         string
	jwtLifetime time.Duration
}

func NewAuthConfig() AuthConfig {
	return AuthConfig{
		cost:        getCost(),
		key:         getKey(),
		jwtLifetime: getJWTLifetime(),
	}
}

func (c AuthConfig) Cost() int {
	return c.cost
}

func (c AuthConfig) Key() []byte {
	return []byte(c.key)
}

func (c AuthConfig) JWTLifetime() time.Duration {
	return c.jwtLifetime
}

func getCost() int {
	flagValue := viper.GetInt(costFlag)
	envValue := viper.GetInt(costEnv)
	log := logger.Instance.With().Str("config", "auth cost").Logger()

	if envValue > 0 {
		log.Info().
			Str("env", keyEnv).
			Int("value", envValue).
			Msg("use env value")
		return envValue
	}

	if flagValue != costDefault && flagValue <= 0 {
		log.Info().
			Str("flag", costFlagPrint).
			Int("value", flagValue).
			Msg("use flag value")
		return flagValue
	}

	log.Info().Str("flag", costFlagPrint).Msg("use default value")

	return costDefault
}

func getKey() string {
	flagValue := viper.GetString(keyFlag)
	envValue := viper.GetString(keyEnv)
	log := logger.Instance.With().Str("config", "auth key").Logger()

	if envValue != "" {
		log.Info().
			Str("env", keyEnv).
			Str("value", envValue).
			Msg("use env value")
		return envValue
	}

	if flagValue != keyDefault {
		log.Info().
			Str("flag", keyFlagPrint).
			Str("value", flagValue).
			Msg("use flag value")
		return flagValue
	}

	log.Info().Str("flag", keyFlagPrint).Msg("use default value")
	return keyDefault
}

func getJWTLifetime() time.Duration {
	flagValue := viper.GetDuration(jwtLifetimeFlag)
	envValue := viper.GetDuration(jwtLifetimeEvn)
	log := logger.Instance.With().Str("config", "auth jwtLifetime").Logger()

	if envValue > 0 {
		log.Info().
			Str("env", jwtLifetimeEvn).
			Dur("value", envValue).
			Msg("use env value")
		return envValue
	}

	if flagValue != jwtLifetimeDefault && flagValue > 0 {
		log.Info().
			Str("flag", jwtLifetimeFlagPrint).
			Dur("value", flagValue).
			Msg("use flag value")
		return flagValue
	}

	log.Info().Str("flag", jwtLifetimeFlagPrint).Msg("use default value")
	return jwtLifetimeDefault
}
