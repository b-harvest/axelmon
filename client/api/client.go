package api

var C *Client

func Set(network string, axelarAPI string) {
	var axelarscan string
	if network == "mainnet" {
		axelarscan = "https://api.axelarscan.io:443"
	} else if network == "testnet" {
		axelarscan = "https://testnet.api.axelarscan.io:443"
	} else {
		panic("You must input mainnet or testnet to network field in your config.")
	}

	C = &Client{
		axelar: axelarAPI,
		axelarscan: axelarscan,
	}
}
