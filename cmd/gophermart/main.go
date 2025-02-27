package main

import (
	"log"

	"github.com/niksmo/gophermart/internal/config"
)

func main() {
	config.Init()
	addrConfig := config.NewAddressConfig()
	accrualConfig := config.NewAccrualAddrConfig()
	log.Println("addr:", addrConfig)
	log.Println("accrual addr:", accrualConfig.GetOrdersReqURL("54321"))
}
