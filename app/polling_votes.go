package app

import (
	"context"
	"fmt"
	"github.com/axelarnetwork/axelar-core/x/nexus/exported"
	"strings"
	"time"

	"bharvest.io/axelmon/client/api"
	"bharvest.io/axelmon/client/grpc"
	"bharvest.io/axelmon/metrics"
	"bharvest.io/axelmon/server"
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

		if c.PollingVote.CheckPeriodDays == 0 {
			c.PollingVote.CheckPeriodDays = 10
		}
		resp, err := api.C.GetPollingVotes(chain.String(), c.PollingVote.CheckN, c.Wallet.Proxy.PrintAcc(), pollingType,
			time.Duration(c.PollingVote.CheckPeriodDays)*time.Hour*24)
		if err != nil {
			return err
		}

		votesInfo.Missed = fmt.Sprintf("%d / %d", resp.MissCnt, int(resp.TotalVotes))
		metrics.SetEVMVotesMissed(chain.String(), resp.MissCnt)
		metrics.SetEVMVotesSuccess(chain.String(), int(resp.TotalVotes)-resp.MissCnt)

		if (float64(resp.MissCnt)/resp.TotalVotes)*100 > float64(c.PollingVote.MissPercentage) {
			votesInfo.Status = false

			msg := fmt.Sprintf("%s status(%s)", pollingType, chain)
			c.alert(msg, []string{}, false, false)
		} else {
			votesInfo.Status = true

			msg := fmt.Sprintf("%s status(%s)", pollingType, chain)
			c.alert(msg, []string{}, true, false)
		}

		result[chain.String()] = votesInfo
	}
	server.GlobalState.EVMVotes.Chain = result

	return nil
}
