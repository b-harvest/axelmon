package wallet

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"strings"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func getPrefix(addr string) string {
	i := strings.Index(addr, "1")
	if i == -1 {
		msg := fmt.Sprintf("Address converting error => %s", addr)
		panic(errors.New(msg))
	}
	return addr[:i]
}

func NewWallet(_ context.Context, addr string) (*Wallet, error) {
	prefix := getPrefix(addr)
	config := sdkTypes.NewConfig()
	config.SetBech32PrefixForAccount(prefix, prefix+"pub")
	config.SetBech32PrefixForValidator(prefix+"valoper", prefix+"valoperpub")
	config.SetBech32PrefixForConsensusNode(prefix+"valcons", prefix+"valvalconspub")
	config.Seal()

	address, err := sdkTypes.GetFromBech32(addr, prefix)
	if err != nil {
		return nil, err
	}
	hex := hex.EncodeToString(address)

	acc, err := sdkTypes.AccAddressFromHex(hex)
	if err != nil {
		return nil, err
	}
	val, err := sdkTypes.ValAddressFromHex(hex)
	if err != nil {
		return nil, err
	}
	cons, err := sdkTypes.ConsAddressFromHex(hex)
	if err != nil {
		return nil, err
	}

	wallet := &Wallet{
		config,
		acc,
		val,
		cons,
	}
	return wallet, nil
}

func (w *Wallet) PrintAcc() string {
	addrBytes := w.Acc.Bytes()

	bech32Addr, err := bech32.ConvertAndEncode(w.config.GetBech32AccountAddrPrefix(), addrBytes)
	if err != nil {
		panic(err)
	}
	return bech32Addr
}

func (w *Wallet) PrintValoper() string {
	addrBytes := w.Val.Bytes()
	bech32Addr, err := bech32.ConvertAndEncode(w.config.GetBech32ValidatorAddrPrefix(), addrBytes)
	if err != nil {
		panic(err)
	}
	return bech32Addr
}

func (w *Wallet) PrintCons() string {
	addrBytes := w.Cons.Bytes()
	bech32Addr, err := bech32.ConvertAndEncode(w.config.GetBech32ConsensusAddrPrefix(), addrBytes)
	if err != nil {
		panic(err)
	}
	return bech32Addr
}
