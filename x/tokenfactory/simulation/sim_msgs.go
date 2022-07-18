package simulation

import (
	"errors"

	legacysimulationtype "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/osmosis-labs/osmosis/v10/osmoutils"
	simulation "github.com/osmosis-labs/osmosis/v10/simulation/types"
	"github.com/osmosis-labs/osmosis/v10/x/tokenfactory/keeper"
	"github.com/osmosis-labs/osmosis/v10/x/tokenfactory/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RandomMsgCreateDenom creates a random tokenfactory denom that is no greater than 44 alphanumeric characters
func RandomMsgCreateDenom(k keeper.Keeper, sim *simulation.SimCtx, ctx sdk.Context) (*types.MsgCreateDenom, error) {
	return &types.MsgCreateDenom{
		Sender:   sim.RandomSimAccount().Address.String(),
		Subdenom: sim.RandStringOfLength(types.MaxSubdenomLength),
	}, nil
}

func RandomMsgMintDenom(k keeper.Keeper, sim *simulation.SimCtx, ctx sdk.Context) (*types.MsgMint, error) {
	acc, senderExists := sim.RandomSimAccountWithConstraint(accountCreatedTokenFactoryDenom(k, ctx))
	if !senderExists {
		return nil, errors.New("no addr has created a tokenfactory coin")
	}
	// Pick denom
	store := k.GetCreatorPrefixStore(ctx, acc.Address.String())
	denoms := osmoutils.GatherAllKeysFromStore(store)
	denom := simulation.RandSelect(sim, denoms...)

	// TODO: Replace with an improved rand exponential coin
	mintAmount := sim.RandPositiveInt(sdk.NewIntFromUint64(1000_000000))
	return &types.MsgMint{
		Sender: acc.Address.String(),
		Amount: sdk.NewCoin(denom, mintAmount),
	}, nil
}

// TODO: We are going to need to index the owner of an account as well, rather than creator
// to simulate admin changes
func accountCreatedTokenFactoryDenom(k keeper.Keeper, ctx sdk.Context) simulation.SimAccountConstraint {
	return func(acc legacysimulationtype.Account) bool {
		store := k.GetCreatorPrefixStore(ctx, acc.Address.String())
		iterator := store.Iterator(nil, nil)
		defer iterator.Close()
		return iterator.Valid()
	}
}
