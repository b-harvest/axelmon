package api

type Client struct {
	axelar     string
	axelarscan string
}

type VotesRequest struct {
	Method string `json:"method"`
	Chain  string `json:"chain"`
	Size   int    `json:"size"`
}

type Voter struct {
	Late      bool   `json:"late"`
	CreatedAt int64  `json:"created_at"`
	ID        string `json:"id"`
	Voter     string `json:"voter"`
	Type      string `json:"type"`
	Vote      bool   `json:"vote"`
	Confirmed bool   `json:"confirmed"`
	Height    int    `json:"height"`
}

type VotesReturn struct {
	Chain      string     `json:"chain"`
	Validator  string     `json:"validator"`
	VoteInfos  []VoteInfo `json:"votes"`
	MissCnt    int        `json:"miss_count"`
	TotalVotes int        `json:"total_votes"`
}

type VoteInfo struct {
	PollID          string `json:"poll_id"`
	InitiatedTXHash string `json:"initiated_tx_hash"`
	Vote            int    `json:"vote"` // 0=not voted, 1=yes, 2=no
	IsLate          bool   `json:"is_late"`
	Validator       string `json:"validator"`
}

type Proxy struct {
	Height string `json:"height"`
	Result struct {
		Address string `json:"address"`
		Status  string `json:"status"`
	} `json:"result"`
}

type VerifierAccount struct {
	Address         string   `json:"address"`
	SupportedChains []string `json:"supportedChains"`
}
