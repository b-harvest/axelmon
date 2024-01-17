package app

import (
	"context"
	"fmt"
	"strings"

	"bharvest.io/axelmon/client/api"
	"bharvest.io/axelmon/client/grpc"
	"bharvest.io/axelmon/log"
	"bharvest.io/axelmon/server"
	"bharvest.io/axelmon/tg"
)

func (c *Config) checkEVMVotes(ctx context.Context) error {
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

	result := make(map[string]server.VotesInfo)
	for _, chain := range chains {
		// If chain is included in except chains
		// then don't monitor that chain's EVM votes.
		if c.General.ExceptChains[strings.ToLower(chain.String())] {
			continue
		}

		votesInfo := server.VotesInfo{}

		resp, err := api.C.GetEVMVotes(chain.String(), c.EVMVote.CheckN, c.Wallet.Proxy.PrintAcc())
		if err != nil {
			return err
		}

		votesInfo.Missed = fmt.Sprintf("%d / %d", resp.MissCnt, c.EVMVote.CheckN)
		if resp.MissCnt >= c.EVMVote.MissCnt {
			votesInfo.Status = false

			msg := fmt.Sprintf("EVM votes status(%s): 🛑", chain)
			tg.SendMsg(msg)
			log.Info(msg)
		} else {
			votesInfo.Status = true

			msg := fmt.Sprintf("EVM votes status(%s): 🟢", chain)
			log.Info(msg)
		}
		result[chain.String()] = votesInfo
	}
	server.GlobalState.EVMVotes.Chain = result

	return nil
}
