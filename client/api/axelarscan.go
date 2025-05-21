package api

import (
	"bharvest.io/axelmon/log"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/axelarnetwork/axelar-core/x/nexus/exported"
	"io"
	"net/http"
	"time"
)

type PollingType string

const (
	EVM_POLLING_TYPE PollingType = "searchPolls"
	VM_POLLING_TYPE  PollingType = "searchVMPolls"
)

func (c *Client) GetVerifierSupportedChains(proxyAcc string) ([]exported.ChainName, error) {
	url := fmt.Sprintf("%s/validator/getVerifiers", c.axelarscan)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("Failed to get proper data from axelarscan")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res map[string]any
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(res["data"])
	if err != nil {
		return nil, err
	}

	var data []VerifierAccount
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return nil, err
	}

	for _, d := range data {
		if d.Address == proxyAcc {
			var result []exported.ChainName
			for _, chainName := range d.SupportedChains {
				result = append(result, exported.ChainName(chainName))
			}
			return result, nil
		}
	}
	return nil, errors.New("didn't found any verifier matched with your Acc -> " + proxyAcc)
}

func (c *Client) GetPollingVotes(chain string, size int, pollingType PollingType, checkPeriod time.Duration) ([]VotesReturn, error) {
	if size > 255 {
		return nil, errors.New("exceed maximum # of votes (max 255)")
	}

	reqBytes, err := json.Marshal(VotesRequest{
		string(pollingType),
		chain,
		size,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/validator/%s", c.axelarscan, pollingType)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response from axelarscan: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var res struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	now := time.Now()
	var voteResults []VotesReturn

	for _, entry := range res.Data {
		createdAtRaw, ok := entry["created_at"].(map[string]any)
		if !ok || createdAtRaw["ms"] == nil {
			continue
		}
		createdMs := int64(createdAtRaw["ms"].(float64))
		createdTime := time.Unix(createdMs/1000, 0)
		if createdTime.Before(now.Add(-checkPeriod)) {
			log.Debug("skipping old vote")
			continue
		}

		for key, value := range entry {
			// skip known meta fields
			switch key {
			case "id", "created_at", "key", "initiated_txhash":
				continue
			}

			valueMap, ok := value.(map[string]any)
			if !ok {
				continue
			}

			voteRaw, err := json.Marshal(valueMap)
			if err != nil {
				log.Errorf("marshal vote for %s failed: %v", key, err)
				continue
			}

			var vote Voter
			if err := json.Unmarshal(voteRaw, &vote); err != nil {
				log.Errorf("unmarshal vote for %s failed: %v", key, err)
				continue
			}

			isLate := false
			if vote.Late != nil {
				isLate = *vote.Late
			}
			v := VoteInfo{
				Validator:       key,
				PollID:          safeString(entry["id"]),
				InitiatedTXHash: safeString(entry["initiated_txhash"]),
				IsLate:          isLate,
			}

			switch {
			case !vote.Vote:
				v.Vote = 2
			case vote.Vote:
				v.Vote = 1
			default:
				v.Vote = 0
			}

			// append to result
			appendVote(&voteResults, chain, v)
		}
	}

	return voteResults, nil
}

func safeString(val any) string {
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

func appendVote(results *[]VotesReturn, chain string, vote VoteInfo) {
	for i := range *results {
		if (*results)[i].Chain == chain && (*results)[i].Validator == vote.Validator {
			r := &(*results)[i]
			r.VoteInfos = append(r.VoteInfos, vote)
			r.TotalVotes++
			if vote.IsLate || vote.Vote != 1 {
				r.MissCnt++
			}
			return
		}
	}

	// if not found, create new
	newResult := VotesReturn{
		Chain:      chain,
		Validator:  vote.Validator,
		VoteInfos:  []VoteInfo{vote},
		MissCnt:    boolToInt(vote.IsLate || vote.Vote != 1),
		TotalVotes: 1,
	}
	*results = append(*results, newResult)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
