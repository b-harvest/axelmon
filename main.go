package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"bharvest.io/axelmon/app"
	"bharvest.io/axelmon/client/api"
	"bharvest.io/axelmon/log"
	"bharvest.io/axelmon/server"
	"bharvest.io/axelmon/tg"
	"bharvest.io/axelmon/wallet"
	"github.com/pelletier/go-toml/v2"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx := context.Background()

	cfgPath := flag.String("config", "", "Config file")
	flag.Parse()
	if *cfgPath == "" {
		panic("Error: Please input config file path with -config flag.")
	}

	f, err := os.ReadFile(*cfgPath)
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

	cfg.Wallet.Validator, err = wallet.NewWallet(ctx, cfg.General.ValidatorAcc)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	api.Set(cfg.General.Network, cfg.General.API)
	log.Info(fmt.Sprintf("Start Axelmon (for %s)", cfg.General.Network))
	proxyAcc, err := api.C.GetProxyByVal(cfg.Wallet.Validator.PrintValoper())
	if err != nil {
		log.Error(err)
		panic(err)
	}

	if proxyAcc != "" {
		cfg.Wallet.Proxy, err = wallet.NewWallet(ctx, proxyAcc)
		if err != nil {
			log.Error(err)
			panic(err)
		}
	} else {
		log.Warn("Cannot fetch proxy acc. it may occur errors while retrieving voting infos.")
		cfg.Wallet.Proxy = cfg.Wallet.Validator
		if err != nil {
			log.Error(err)
			panic(err)
		}
	}

	tgTitle := fmt.Sprintf("🤖 Axelmon for %s 🤖", cfg.General.Network)
	tg.SetTg(cfg.Tg.Enable, tgTitle, cfg.Tg.Token, cfg.Tg.ChatID)

	cfg.General.ExceptChains = map[string]bool{}
	exceptChains := strings.Split(strings.ReplaceAll(cfg.General.ExceptChainsString, " ", ""), ",")
	for _, exceptChain := range exceptChains {
		cfg.General.ExceptChains[strings.ToLower(exceptChain)] = true
	}

	go server.Run(cfg.General.ListenPort)
	for {
		app.Run(ctx, &cfg)
		time.Sleep(time.Duration(cfg.General.Period) * time.Minute)
	}
}
