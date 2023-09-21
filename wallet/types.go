package wallet

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

type Wallet struct {
	config *sdkTypes.Config
	Acc sdkTypes.AccAddress
	Val sdkTypes.ValAddress
	Cons sdkTypes.ConsAddress
}
