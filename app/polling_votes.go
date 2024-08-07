package app

import (
	"context"
	"fmt"
	"github.com/axelarnetwork/axelar-core/x/nexus/exported"
	"strings"

	"bharvest.io/axelmon/client/api"
	"bharvest.io/axelmon/client/grpc"
	"bharvest.io/axelmon/log"
	"bharvest.io/axelmon/metrics"
	"bharvest.io/axelmon/server"
	"bharvest.io/axelmon/tg"
	"github.com/prometheus/client_golang/prometheus"
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

	return c.checkPollingVotes(ctx, api.EVM_POLLING_TYPE, chains)
}

func (c *Config) checkVMVotes(ctx context.Context) error {
	chains, err := api.C.GetVerifierSupportedChains(c.Wallet.Proxy.PrintAcc())
	if err != nil {
		return err
	}
	return c.checkPollingVotes(ctx, api.VM_POLLING_TYPE, chains)
}

func (c *Config) checkPollingVotes(ctx context.Context, pollingType api.PollingType, chains []exported.ChainName) error {

	result := make(map[string]server.VotesInfo)
	for _, chain := range chains {
		// If chain is included in except chains
		// then don't monitor that chain's EVM votes.
		if c.General.ExceptChains[strings.ToLower(chain.String())] {
			continue
		}

		votesInfo := server.VotesInfo{}

		resp, err := api.C.GetPollingVotes(chain.String(), c.PollingVote.CheckN, c.Wallet.Proxy.PrintAcc(), pollingType)
		if err != nil {
			return err
		}

		votesInfo.Missed = fmt.Sprintf("%d / %d", resp.MissCnt, c.PollingVote.CheckN)
		metrics.EVMVotesCounter.With(prometheus.Labels{"network_name": chain.String(), "status": "missed"}).Add(float64(resp.MissCnt))
		// check if the total number of votes is higher than the number of votes checked
		if resp.TotalVotes < float64(c.PollingVote.CheckN) {
			metrics.EVMVotesCounter.With(prometheus.Labels{"network_name": chain.String(), "status": "success"}).Add(float64(resp.TotalVotes - float64(resp.MissCnt)))
		} else {
			metrics.EVMVotesCounter.With(prometheus.Labels{"network_name": chain.String(), "status": "success"}).Add(float64(c.PollingVote.CheckN - resp.MissCnt))
		}

		if resp.MissCnt >= c.PollingVote.MissCnt {
			votesInfo.Status = false

			msg := fmt.Sprintf("%s status(%s): 🛑", pollingType, chain)
			tg.SendMsg(msg)
			log.Info(msg)
		} else {
			votesInfo.Status = true

			msg := fmt.Sprintf("%s status(%s): 🟢", pollingType, chain)
			log.Info(msg)
		}
		result[chain.String()] = votesInfo
	}
	server.GlobalState.EVMVotes.Chain = result

	return nil
}
