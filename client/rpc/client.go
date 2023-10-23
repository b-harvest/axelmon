package rpc

import (
	"context"

	"bharvest.io/axelmon/log"
	cometbftHttp "github.com/cometbft/cometbft/rpc/client/http"
)

func New(host string) (*Client, error) {
	result := &Client{
		host: host,
	}
	client, err := cometbftHttp.New(result.host, "/websocket")
	if err != nil {
		return nil, err
	}
	result.RPCClient = client

	return result, nil
}

func (c *Client) Connect(ctx context.Context) error {
	// For websocket connection
	err := c.RPCClient.Start()
	if err != nil {
		return err
	}

	log.Info("RPC connected")
	return nil
}

func (c *Client) Terminate(_ context.Context) error {
	// For websocket connection
	err := c.RPCClient.Stop()
	if err != nil {
		return err
	}

	log.Info("RPC connection terminated")
	return nil
}
