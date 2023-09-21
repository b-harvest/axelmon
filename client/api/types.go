package api

type VotesRequest struct {
	Method string 	`json:"method"`
	Chain  string 	`json:"chain"`
	Size   int 		`json:"size"`
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

type VoteInfo struct {
	InitiatedTXHash string
	PollID string
	IsLate bool

	// 0 => not voted
	// 1 => yes
	// 2 => no
	Vote byte
}

type VotesResponse struct {
	Chain string
	MissCnt byte
	VoteInfos []VoteInfo
}

type Proxy struct {
	Height string `json:"height"`
	Result struct {
		Address string `json:"address"`
		Status  string `json:"status"`
	} `json:"result"`
}
