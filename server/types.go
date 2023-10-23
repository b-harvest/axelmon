package server

type VotesInfo struct {
	Status bool   `json:"status"`
	Missed string `json:"missed"`
}

type Response struct {
	Maintainers struct {
		Status      bool `json:"status"`
		Maintainer map[string]bool `json:"maintainer"`
	} `json:"maintainers"`

	Heartbeat struct {
		Status bool   `json:"status"`
		Missed string `json:"missed"`
	} `json:"heartbeat"`

	EVMVotes struct {
		Chain map[string]VotesInfo `json:"chain"`
	} `json:"EVMVotes"`
	
}
