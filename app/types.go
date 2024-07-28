package app

import "bharvest.io/axelmon/wallet"

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
	Tg struct {
		Enable bool   `toml:"enable"`
		Token  string `toml:"token"`
		ChatID string `toml:"chat_id"`
	} `toml:"tg"`

	Heartbeat struct {
		CheckN  int `toml:"check_n"`
		MissCnt int `toml:"miss_cnt"`
	} `toml:"heartbeat"`
	PollingVote struct {
		CheckN  int `toml:"check_n"`
		MissCnt int `toml:"miss_cnt"`
	} `toml:"evm_vote"`
}

type TargetSvc string

const (
	MaintainerTargetSvc TargetSvc = "maintainer"
	HeartbeatTargetSvc  TargetSvc = "heartbeat"
	EVMVoteTargetSvc    TargetSvc = "evm"
	VMVoteTargetSvc     TargetSvc = "vm"
)
