package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/axelarnetwork/axelar-core/x/nexus/exported"
	"io"
	"net/http"
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
	return nil, errors.New("didnt' found any verifier matched with your Acc -> " + proxyAcc)
}

func (c *Client) GetPollingVotes(chain string, size int, proxyAcc string, pollingType PollingType) (*VotesReturn, error) {
	// VotesResponse MissCnt is byte type.
	// Therefore, the maximum number of evm votes should be
	// less than 256
	if size > 255 {
		return nil, errors.New("exceed maximum # evm votes")
	}

	reqBytes, err := json.Marshal(VotesRequest{
		string(pollingType),
		chain,
		size,
	})
	if err != nil {
		return nil, err
	}
	reqBody := bytes.NewBuffer(reqBytes)

	url := fmt.Sprintf("%s/validator/%s", c.axelarscan, pollingType)
	req, err := http.NewRequest("POST", url, reqBody)
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

	var data []map[string]any
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return nil, err
	}

	result := VotesReturn{}
	result.Chain = chain
	result.MissCnt = 0
	if total, ok := res["total"].(float64); ok {
		result.TotalVotes = total
	} else {
		result.TotalVotes = 0
	}
	result.VoteInfos = make([]VoteInfo, len(data))
	for i, d := range data {
		if d["initiated_txhash"] != nil {
			result.VoteInfos[i].InitiatedTXHash = d["initiated_txhash"].(string)
		}
		if d["id"] != nil {
			result.VoteInfos[i].PollID = d["id"].(string)
		}

		voter := d[proxyAcc]
		if voter != nil {
			voteInfoBytes, err := json.Marshal(voter)
			if err != nil {
				return nil, err
			}

			voteInfo := Voter{}
			err = json.Unmarshal(voteInfoBytes, &voteInfo)
			if err != nil {
				return nil, err
			}

			result.VoteInfos[i].IsLate = voteInfo.Late

			if voteInfo.Vote {
				// vote => yes
				result.VoteInfos[i].Vote = 1
			} else {
				// vote => no
				result.VoteInfos[i].Vote = 2
			}
		} else {
			// vote => not voted
			result.VoteInfos[i].Vote = 0
		}

		if result.VoteInfos[i].IsLate || result.VoteInfos[i].Vote != 1 {
			result.MissCnt++
		}
	}

	return &result, nil
}
