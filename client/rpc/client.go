package rpc

import (
	tendermintHttp "github.com/tendermint/tendermint/rpc/client/http"
)

func New(host string) (*Client, error) {
	result := &Client{
		host: host,
	}
	client, err := tendermintHttp.New(result.host, "/websocket")
	if err != nil {
		return nil, err
	}
	result.RPCClient = client

	return result, nil
}
