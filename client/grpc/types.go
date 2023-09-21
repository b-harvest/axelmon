package grpc

import (
	evmTypes "github.com/axelarnetwork/axelar-core/x/evm/types"
	nexusTypes "github.com/axelarnetwork/axelar-core/x/nexus/types"
	tx "github.com/cosmos/cosmos-sdk/types/tx"
	"google.golang.org/grpc"
)

type Client struct {
	host string
	conn *grpc.ClientConn
	txServiceClient tx.ServiceClient
	nexusQueryServiceClient nexusTypes.QueryServiceClient
	evmQueryServiceClient evmTypes.QueryServiceClient
}
