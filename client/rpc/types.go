package rpc

import (
	cometbftHttp "github.com/cometbft/cometbft/rpc/client/http"
)

type Client struct {
	RPCClient *cometbftHttp.HTTP
	host string
}
