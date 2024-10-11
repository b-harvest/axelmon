package server

import (
	"sync"
	"time"
)

type VotesInfo struct {
	Status bool   `json:"status"`
	Missed string `json:"missed"`
}

type Response struct {
	Maintainers struct {
		Status     bool            `json:"status"`
		Maintainer map[string]bool `json:"maintainer"`
	} `json:"maintainers"`

	Heartbeat struct {
		Status bool   `json:"status"`
		Missed string `json:"missed"`
	} `json:"heartbeat"`

	EVMVotes struct {
		Chain map[string]VotesInfo `json:"chain"`
	} `json:"externalChainVotes"`

	Alerts struct {
		SentTgAlarms  map[string]time.Time `json:"sent_tg_alarms"`
		SentSlkAlarms map[string]time.Time `json:"sent_slk_alarms"`
		AllAlarms     map[string]time.Time `json:"sent_all_alarms"`
		NotifyMux     sync.RWMutex
	} `json:"alerts"`
}
