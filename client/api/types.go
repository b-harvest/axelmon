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
		Chain      string
		MissCnt    int
		VoteInfos  []VoteInfo
		TotalVotes float64
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

type (
	SigningsRequest struct {
		Chain  string `json:"chain"`
		Size   int    `json:"size"`
	}

	Signing struct {
		CreatedAt int64  `json:"created_at"`
		ID        string `json:"id"`
		Signer     string `json:"signer"`
		Type      string `json:"type"`
		Sign      bool   `json:"sign"`
		Height    int    `json:"height"`
	}

	SigningsReturn struct {
		Chain      string
		MissCnt    int
		SigningInfos  []SigningInfo
		TotalSignings float64
	}
	SigningInfo struct {
		InitiatedTXHash string
		SessionID       float64
		
		// 0 => not signed
		// 1 => yes
		// 2 => no
		Sign byte
	}
)

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
