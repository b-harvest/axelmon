package api

var C *Client

func Set(network string, axelarAPI string) {
	var axelarscan string
	if network == "mainnet" {
		axelarscan = "https://api.axelarscan.io:443"
	} else if network == "testnet" {
		axelarscan = "https://testnet.api.axelarscan.io:443"
	} else {
		axelarscan = network
	}

	C = &Client{
		axelar:     axelarAPI,
		axelarscan: axelarscan,
	}
}
