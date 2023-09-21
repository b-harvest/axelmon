package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"bharvest.io/axelmon/client/grpc"
	"bharvest.io/axelmon/client/rpc"
	"bharvest.io/axelmon/tg"
	rewardTypes "github.com/axelarnetwork/axelar-core/x/reward/types"
	tssTypes "github.com/axelarnetwork/axelar-core/x/tss/types"
)

func (c *Config) checkHeartbeat(ctx context.Context) error {
	clientGRPC := grpc.New(c.Node.GRPC)
	err := clientGRPC.Connect(ctx, c.Node.GRPCSecureConnection)
	defer clientGRPC.Terminate(ctx)
	if err != nil {
		return err
	}

	heartbeatHeight, err := c.findHeartBeatHeight(ctx)
	if err != nil {
		return err
	}

	n_heart := 3
	cnt := 0
	fmt.Println("=================== Heartbeat ===================")
	fmt.Println("Broadcaster:", c.Wallet.Proxy.PrintAcc())
	for i:=0; i<n_heart; i++ {
		isFound, err := c.findHeartbeat(ctx, clientGRPC, heartbeatHeight, 5)
		if err != nil {
			return err
		}
		heartbeatHeight -= 50

		if isFound {
			cnt++
		}
	}

	if cnt == n_heart-2 {
		// # 2 heartbeats missing
		tg.SendMsg("🛑 연속적인 heartbeat missing 발생%0A노드를 확인해주세요.")
	} else {
		fmt.Println("Heartbeat: 🟢")
	}

	return nil
}

func (c *Config) findHeartbeat(ctx context.Context, clientGRPC *grpc.Client, heartbeatHeight int64, tryCnt int) (bool, error) {
	for j:= 0; j<tryCnt; j++ {
		txs, err := clientGRPC.GetTxs(ctx, heartbeatHeight)
		if err != nil {
			return false, err
		}
		for _, tx := range txs {
			for _, msg := range tx.Body.Messages {
				if msg.TypeUrl == "/axelar.reward.v1beta1.RefundMsgRequest" {
					refundMsg := rewardTypes.RefundMsgRequest{}
					err = refundMsg.Unmarshal(msg.Value)
					if err != nil {
						return false, err
					}
					if refundMsg.InnerMessage.TypeUrl == "/axelar.tss.v1beta1.HeartBeatRequest" {
						heartbeat := tssTypes.HeartBeatRequest{}
						err = heartbeat.Unmarshal(refundMsg.InnerMessage.Value)
						if err != nil {
							return false, err
						}
						if heartbeat.Sender.Equals(c.Wallet.Proxy.Acc) && len(heartbeat.KeyIDs) >= 1 {

							hash := sha256.New()
							hash.Write([]byte(fmt.Sprint(heartbeat.KeyIDs)))
							md := hash.Sum(nil)
							fmt.Println("hash: ", hex.EncodeToString(md))

							fmt.Println("# signed:", len(heartbeat.KeyIDs))
							fmt.Println("==============")
							return true, nil
						}
					}
				}
			}
		}
		heartbeatHeight++
	}

	return false, errors.New(fmt.Sprintf("Didn't heartbeat signal on height: %d", heartbeatHeight))
}

func (c *Config) findHeartBeatHeight(ctx context.Context) (int64, error) {
	client, err := rpc.New(c.Node.RPC)
	if err != nil {
		return 0, err
	}

	height, err := client.GetLatestHeight(ctx)
	if err != nil {
		return 0, err
	}

	var heartbeatHeight int64
	if height % 50 != 0 {
		heartbeatHeight = height - (height % 50) + 1
	} else {
		heartbeatHeight = heartbeatHeight-50 + 1
	}

	return heartbeatHeight, nil
}
