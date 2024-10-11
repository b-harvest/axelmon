package main

import (
	"bharvest.io/axelmon/app"
	"bharvest.io/axelmon/client/api"
	"bharvest.io/axelmon/log"
	"bharvest.io/axelmon/server"
	"bharvest.io/axelmon/wallet"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx, cancel := context.WithCancel(context.Background())

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

	if cfg.PollingVote.MissPercentage < 0 || cfg.PollingVote.MissPercentage > 100 {
		panic(errors.New("MissPercentage must be between 0 to 100"))
	} else if cfg.PollingVote.MissPercentage == 0 {
		log.Debug("MissPercentage seems like zero. it'll alert if there is any failed record")
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

	cfg.General.ExceptChains = map[string]bool{}
	exceptChains := strings.Split(strings.ReplaceAll(cfg.General.ExceptChainsString, " ", ""), ",")
	for _, exceptChain := range exceptChains {
		cfg.General.ExceptChains[strings.ToLower(exceptChain)] = true
	}

	cfg.Ctx = ctx

	go server.Run(cfg.General.ListenPort)

	quitting := make(chan os.Signal, 1)
	signal.Notify(quitting, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	ticker := time.NewTicker(time.Duration(cfg.General.Period) * time.Minute)
	go app.Run(ctx, &cfg)

	for {
		select {
		case <-ctx.Done():
			app.SaveOnExit(server.STATE_FILE_PATH)
			return
		case <-quitting:
			cancel()
			app.SaveOnExit(server.STATE_FILE_PATH)
			return
		case <-ticker.C:
			go app.Run(ctx, &cfg)
		}
	}

}
