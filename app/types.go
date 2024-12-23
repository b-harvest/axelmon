package app

import (
	"context"
	"sync"
	"time"

	"bharvest.io/axelmon/wallet"
)

type Duration time.Duration

func (d *Duration) UnmarshalText(b []byte) error {
	x, err := time.ParseDuration(string(b))
	if err != nil {
		return err
	}
	*d = Duration(x)
	return nil
}

type Config struct {
	General struct {
		Network              string    `toml:"network"`
		Period               *Duration `toml:"period"`
		ExceptChainsString   string    `toml:"except_chains"`
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
		ResendDuration *Duration `toml:"resend_duration"`
		Tg             struct {
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
	PollingSigning struct {
		CheckN          int `toml:"check_n"`
		MissPercentage  int `toml:"miss_percentage"`
		CheckPeriodDays int `toml:"check_period_days"`
	} `toml:"external_chain_signing"`

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
	VMSigningTargetSvc  TargetSvc = "vmSigning"
)
