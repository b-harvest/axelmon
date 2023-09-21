package main

import (
	"context"
	"os"
	"runtime"

	"bharvest.io/axelmon/app"
	"bharvest.io/axelmon/client/api"
	"bharvest.io/axelmon/log"
	"bharvest.io/axelmon/tg"
	"bharvest.io/axelmon/wallet"
	"github.com/pelletier/go-toml/v2"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx := context.Background()

	f, err := os.ReadFile("config.toml")
	if err != nil {
		log.Error(err)
		panic(err)
	}
	cfg := app.Config{}
	err = toml.Unmarshal(f, &cfg)
	if err != nil {
		log.Error(err)
		panic(err)
	}


	cfg.Wallet.Validator, err = wallet.NewWallet(ctx, cfg.Wallet.ValidatorAcc)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	proxyAcc, err := api.GetProxyByVal(cfg.Wallet.Validator.PrintValoper())
	if err != nil {
		log.Error(err)
		panic(err)
	}
	cfg.Wallet.Proxy, err = wallet.NewWallet(ctx, proxyAcc)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	tg.SetTg(cfg.Tg.Enable,cfg.Tg.Token, cfg.Tg.ChatID)

	err = app.Run(ctx, &cfg)
	if err != nil {
		log.Error(err)
		return
	}
}
