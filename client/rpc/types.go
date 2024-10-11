package rpc

import (
	tendermintHttp "github.com/tendermint/tendermint/rpc/client/http"
)

type Client struct {
	RPCClient *tendermintHttp.HTTP
	host      string
}
