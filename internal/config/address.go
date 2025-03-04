package config

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/viper"
)

const (
	addrEnv       = "RUN_ADDRESS"
	addrFlag      = "addr"
	addrFlagShort = "a"
	addrUsage     = "run address"
	addrDefault   = "127.0.0.1:8080"
)

type AddressConfig struct {
	*net.TCPAddr
}

func NewAddressConfig() *AddressConfig {
	flagValue := viper.GetString(addrFlag)
	envValue := viper.GetString(addrEnv)
	errParsePrefix := "parse run address"

	if envValue != "" {
		TCPAddr, err := net.ResolveTCPAddr("", envValue)
		if err == nil {
			return &AddressConfig{TCPAddr}
		}
		log.Println(fmt.Errorf(errParsePrefix+" env value err:%w", err))
	}

	TCPAddr, err := net.ResolveTCPAddr("", flagValue)
	if err == nil {
		return &AddressConfig{TCPAddr}
	}
	log.Println(fmt.Errorf(errParsePrefix+" flag value err: %w", err))

	TCPAddr, _ = net.ResolveTCPAddr("", addrDefault)
	log.Println("use default run address:", addrDefault)

	return &AddressConfig{TCPAddr}
}
