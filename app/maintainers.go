package app

import (
	"context"
	"fmt"

	"bharvest.io/axelmon/client/grpc"
	"bharvest.io/axelmon/log"
	"bharvest.io/axelmon/server"
	"bharvest.io/axelmon/tg"
)

func (c *Config) checkMaintainers(ctx context.Context) error {
	client := grpc.New(c.General.GRPC)
	err := client.Connect(ctx, c.General.GRPCSecureConnection)
	defer client.Terminate(ctx)
	if err != nil {
		return err
	}

	chains, err := client.GetChains(ctx)
	if err != nil {
		return err
	}

	result := make(map[string]bool)
	for _, chain := range chains {
		maintainers, err := client.GetChainMaintainers(ctx, chain.String())
		if err != nil {
			return err
		}
		for _, acc := range maintainers {
			if acc.Equals(c.Wallet.Validator.Cons) {
				result[chain.String()] = true
			}
		}
	}

	check := true
	msg := "Maintainer list: "
	for k, v := range result {
		msg += fmt.Sprintf("(%s: %v) ", k, v)
		if v == false {
			m := fmt.Sprint("Maintainer status(): 🛑", k)
			tg.SendMsg(m)
			check = false
		}
	}

	server.GlobalState.Maintainers.Maintainer = result
	if check {
		server.GlobalState.Maintainers.Status = true

		log.Info("Maintainer status: 🟢")
	} else {
		server.GlobalState.Maintainers.Status = false

		log.Info("Maintainer status: 🛑")
	}

	return nil
}
