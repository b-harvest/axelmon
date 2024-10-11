package app

import (
	"bharvest.io/axelmon/wallet"
	"context"
	"sync"
)

type Config struct {
	General struct {
		Network              string `toml:"network"`
		Period               uint   `toml:"period"`
		ExceptChainsString   string `toml:"except_chains"`
		ExceptChains         map[string]bool
		ValidatorAcc         string      `toml:"validator_acc"`
		RPC                  string      `toml:"rpc"`
		API                  string      `toml:"api"`
		GRPC                 string      `toml:"grpc"`
		GRPCSecureConnection bool        `toml:"grpc_secure_connection"`
		ListenPort           int         `toml:"listen_port"`
		TargetSvcs           []TargetSvc `toml:"target_svcs"`
	} `toml:"general"`
	Wallet struct {
		Validator *wallet.Wallet
		Proxy     *wallet.Wallet
	} `toml:"wallet"`
	Alerts struct {
		Tg struct {
			Enabled  bool     `toml:"enable"`
			Token    string   `toml:"token"`
			ChatID   string   `toml:"chat_id"`
			Mentions []string `toml:"mentions"`
		} `toml:"telegram"`
		Slack struct {
			Enabled  bool     `toml:"enable"`
			Webhook  string   `toml:"webhook"`
			Mentions []string `toml:"mentions"`
		} `toml:"slack"`
	} `toml:"alerts"`

	Heartbeat struct {
		CheckN  int `toml:"check_n"`
		MissCnt int `toml:"miss_cnt"`
	} `toml:"heartbeat"`
	PollingVote struct {
		CheckN          int `toml:"check_n"`
		MissPercentage  int `toml:"miss_percentage"`
		CheckPeriodDays int `toml:"check_period_days"`
	} `toml:"external_chain_vote"`

	Ctx       context.Context
	Cancel    context.CancelFunc
	alertChan chan *alertMsg
	alertMux  sync.RWMutex
}

type TargetSvc string

const (
	MaintainerTargetSvc TargetSvc = "maintainer"
	HeartbeatTargetSvc  TargetSvc = "heartbeat"
	EVMVoteTargetSvc    TargetSvc = "evm"
	VMVoteTargetSvc     TargetSvc = "vm"
)
