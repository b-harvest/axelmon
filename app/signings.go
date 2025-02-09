package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/axelarnetwork/axelar-core/x/nexus/exported"

	"bharvest.io/axelmon/client/api"
	"bharvest.io/axelmon/metrics"
	"bharvest.io/axelmon/server"
	"github.com/prometheus/client_golang/prometheus"
)

func (c *Config) checkVMSignings(ctx context.Context) error {
	chains, err := api.C.GetVerifierSupportedChains(c.Wallet.Proxy.PrintAcc())
	if err != nil {
		return err
	}
	return c.getSignings(ctx, chains)
}

func (c *Config) getSignings(ctx context.Context, chains []exported.ChainName) error {

	result := make(map[string]server.VotesInfo)
	for _, chain := range chains {
		// If chain is included in except chains
		// then don't monitor that chain's VM signings.
		if c.General.ExceptChains[strings.ToLower(chain.String())] {
			continue
		}

		votesInfo := server.VotesInfo{}

		if c.PollingSigning.CheckPeriodDays == 0 {
			c.PollingSigning.CheckPeriodDays = 10
		}
		resp, err := api.C.GetPollingSignings(chain.String(), c.PollingSigning.CheckN, c.Wallet.Proxy.PrintAcc(),
			time.Duration(c.PollingSigning.CheckPeriodDays)*time.Hour*24)
		if err != nil {
			return err
		}

		votesInfo.Missed = fmt.Sprintf("%d / %d", resp.MissCnt, int(resp.TotalSignings))
		metrics.VMSigningsCounter.With(prometheus.Labels{"network_name": chain.String(), "status": "missed"}).Add(float64(resp.MissCnt))
		// check if the total number of signings is higher than the number of signings checked
		if resp.TotalSignings < float64(c.PollingVote.CheckN) {
			metrics.VMSigningsCounter.With(prometheus.Labels{"network_name": chain.String(), "status": "success"}).Add(float64(int(resp.TotalSignings) - resp.MissCnt))
		} else {
			metrics.VMSigningsCounter.With(prometheus.Labels{"network_name": chain.String(), "status": "success"}).Add(resp.TotalSignings - float64(resp.MissCnt))
		}

		if (float64(resp.MissCnt)/resp.TotalSignings)*100 > float64(c.PollingVote.MissPercentage) {
			votesInfo.Status = false

			msg := fmt.Sprintf("status(%s)", chain)
			c.alert(msg, []string{}, false, false)
		} else {
			votesInfo.Status = true

			msg := fmt.Sprintf("status(%s)", chain)
			c.alert(msg, []string{}, true, false)
		}

		result[chain.String()] = votesInfo
	}
	server.GlobalState.VMSignings.Chain = result

	return nil
}
