package app

import "bharvest.io/axelmon/wallet"

type Config struct {
	Node struct {
		GRPC string `toml:"grpc"`
		GRPCSecureConnection bool `toml:"grpc_secure_connection"`
		RPC  string `toml:"rpc"`
	} `toml:"node"`
	Wallet struct {
		ValidatorAcc   string `toml:"validator_acc"`
		Validator *wallet.Wallet
		Proxy *wallet.Wallet
	} `toml:"wallet"`
	Tg struct {
		Enable bool   `toml:"enable"`
		Token  string `toml:"token"`
		ChatID string `toml:"chat_id"`
	} `toml:"tg"`
}
