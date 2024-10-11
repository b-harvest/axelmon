package grpc

import (
	"context"

	evmTypes "github.com/axelarnetwork/axelar-core/x/evm/types"
	"github.com/axelarnetwork/axelar-core/x/nexus/exported"
	nexusTypes "github.com/axelarnetwork/axelar-core/x/nexus/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	tx "github.com/cosmos/cosmos-sdk/types/tx"
)

// GetChainMaintainers fetch maintainers on specific chain.
// and this feature only apply to validator external chain support.(not amplifier)
func (c *Client) GetChainMaintainers(ctx context.Context, chain string) ([]sdkTypes.ValAddress, error) {
	resp, err := c.nexusQueryServiceClient.ChainMaintainers(
		ctx,
		&nexusTypes.ChainMaintainersRequest{
			Chain: chain,
		},
	)
	if err != nil {
		return nil, err
	}

	return resp.Maintainers, nil
}

func (c *Client) GetChains(ctx context.Context) ([]exported.ChainName, error) {
	resp, err := c.evmQueryServiceClient.Chains(
		ctx,
		&evmTypes.ChainsRequest{},
	)
	if err != nil {
		return nil, err
	}

	return resp.Chains, nil
}

func (c *Client) GetTxs(ctx context.Context, height int64) ([]*tx.Tx, error) {
	resp, err := c.txServiceClient.GetBlockWithTxs(
		ctx,
		&tx.GetBlockWithTxsRequest{
			Height: height,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp.Txs, nil
}
