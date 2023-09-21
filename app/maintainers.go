package app

import (
	"context"
	"fmt"

	"bharvest.io/axelmon/client/grpc"
	"bharvest.io/axelmon/log"
	"bharvest.io/axelmon/tg"
)

func (c *Config) checkMaintainers(ctx context.Context) (error) {
	client := grpc.New(c.Node.GRPC)
	err := client.Connect(ctx, c.Node.GRPCSecureConnection)
	defer client.Terminate(ctx)
	if err != nil {
		log.Error(err)
		return err
	}

	chains, err := client.GetChains(ctx)
	if err != nil {
		log.Error(err)
		return err
	}

	result := make(map[string]bool)
	for _, chain:= range chains {
		maintainers, err := client.GetChainMaintainers(ctx, chain.String())
		if err != nil {
			log.Error(err)
			return err
		}
		for _, acc := range maintainers {
			if acc.Equals(c.Wallet.Validator.Cons) {
				result[chain.String()] = true
			}
		}
	}

	check := true
	for k, v := range result {
		if v == false {
			m := fmt.Sprint("🛑 Axelar Maintainer를 확인해주세요.%0AMaintainer: ", k)
			tg.SendMsg(m)
			check = false
		}
	}

	fmt.Println("=================== Maintainer ===================")
	if check {
		fmt.Println("Maintainer: 🟢")
	} else {
		fmt.Println("Status: 🛑")
	}
	fmt.Println(result)

	return nil
}
