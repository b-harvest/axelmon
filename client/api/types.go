package api

type Client struct {
	axelar     string
	axelarscan string
}

type (
	VotesRequest struct {
		Method string `json:"method"`
		Chain  string `json:"chain"`
		Size   int    `json:"size"`
	}

	Voter struct {
		Late      bool   `json:"late"`
		CreatedAt int64  `json:"created_at"`
		ID        string `json:"id"`
		Voter     string `json:"voter"`
		Type      string `json:"type"`
		Vote      bool   `json:"vote"`
		Confirmed bool   `json:"confirmed"`
		Height    int    `json:"height"`
	}

	VotesReturn struct {
		Chain     string
		MissCnt   int
		VoteInfos []VoteInfo
	}
	VoteInfo struct {
		InitiatedTXHash string
		PollID          string
		IsLate          bool

		// 0 => not voted
		// 1 => yes
		// 2 => no
		Vote byte
	}
)

type Proxy struct {
	Height string `json:"height"`
	Result struct {
		Address string `json:"address"`
		Status  string `json:"status"`
	} `json:"result"`
}
