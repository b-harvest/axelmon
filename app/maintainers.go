package app

import (
	"context"
	"fmt"
	"strings"

	"bharvest.io/axelmon/client/grpc"
	"bharvest.io/axelmon/log"
	"bharvest.io/axelmon/metrics"
	"bharvest.io/axelmon/server"
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
		// If chain is included in except chains
		// then don't monitor that chain's maintainers.
		if c.General.ExceptChains[strings.ToLower(chain.String())] {
			result[chain.String()] = true
			continue
		}

		maintainers, err := client.GetChainMaintainers(ctx, chain.String())
		if err != nil {
			return err
		}
		for _, acc := range maintainers {
			if acc.Equals(c.Wallet.Validator.Cons) {
				result[chain.String()] = true
				break
			}
		}

		if !result[chain.String()] {
			result[chain.String()] = false
		}

		var maintainerInNetwork int
		if result[chain.String()] {
			maintainerInNetwork = 1
		} else {
			maintainerInNetwork = 0
		}
		metrics.SetMaintainersStatus(chain.String(), maintainerInNetwork)
	}

	check := true
	prefix := "Maintainers status: "
	msg := prefix
	var alerts []string
	for k, v := range result {
		msg += fmt.Sprintf("(%s: %v) ", k, v)
		if v == false {
			check = false
			alerts = append(alerts, k)
		}
	}
	log.Debug(msg)

	server.GlobalState.Maintainers.Maintainer = result
	if check {
		server.GlobalState.Maintainers.Status = true

		c.alert(prefix, alerts, true, false)

		log.Info(msg)
	} else {
		server.GlobalState.Maintainers.Status = false

		c.alert(prefix, alerts, false, false)
		log.Info(msg)
	}

	return nil
}
