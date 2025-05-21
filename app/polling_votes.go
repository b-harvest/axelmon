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

	if c.PollingVote.CheckPeriodDays == 0 {
		c.PollingVote.CheckPeriodDays = 10
	}

	checkPeriod := time.Duration(c.PollingVote.CheckPeriodDays) * 24 * time.Hour

	for _, chain := range chains {
		// Skip excepted chains
		if c.General.ExceptChains[strings.ToLower(chain.String())] {
			continue
		}

		// Fetch polling vote data for this chain
		votesResults, err := api.C.GetPollingVotes(chain.String(), c.PollingVote.CheckN, pollingType, checkPeriod)
		if err != nil {
			return fmt.Errorf("error polling votes for chain %s: %w", chain, err)
		}

		// Aggregate by validator
		perValidator := make(map[string]*server.VotesInfo)
		for _, voteResult := range votesResults {
			v := perValidator[voteResult.Validator]
			if v == nil {
				v = &server.VotesInfo{}
				perValidator[voteResult.Validator] = v
			}
			v.MissedCnt += voteResult.MissCnt
			v.TotalCnt += voteResult.TotalVotes
		}

		// Determine status & set metrics
		for validator, info := range perValidator {
			info.Missed = fmt.Sprintf("%d / %d", info.MissedCnt, info.TotalCnt)

			// Expose metrics per validator per chain
			metrics.SetEVMVotesMissed(c.General.Network, validator, chain.String(), info.MissedCnt)
			metrics.SetEVMVotesSuccess(c.General.Network, validator, chain.String(), info.MissedCnt)

			missPercentage := float64(info.MissedCnt) / float64(info.TotalCnt) * 100
			if missPercentage > float64(c.PollingVote.MissPercentage) {
				info.Status = false
				c.alert(fmt.Sprintf("%s status(%s:%s)", pollingType, chain, validator), []string{}, false, false)
			} else {
				info.Status = true
				c.alert(fmt.Sprintf("%s status(%s:%s)", pollingType, chain, validator), []string{}, true, false)
			}

			result[chain.String()+"|"+validator] = *info
		}
	}

	server.GlobalState.EVMVotes.Chain = result
	return nil
}
