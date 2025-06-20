package mint

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/keeper"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper, ic types.InflationCalculationFn) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	// recalculate inflation rate
	//totalStakingSupply := k.StakingTokenSupply(ctx)
	//bondedRatio := k.BondedRatio(ctx)
	//minter.Inflation = ic(ctx, minter, params, bondedRatio)
	//minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalStakingSupply)
	k.SetMinter(ctx, minter)

	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params, ctx.BlockHeight())
	mintedCoins := sdk.NewCoins(mintedCoin)

	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	burnAmount := sdk.NewDecFromInt(mintedCoin.Amount).Mul(sdk.NewDecWithPrec(5, 2)).TruncateInt()
	burnedCoin := sdk.NewCoin(mintedCoin.Denom, burnAmount)
	burnedCoins := sdk.NewCoins(burnedCoin)
	err = k.BurnCoins(ctx, burnedCoins)
	if err != nil {
		panic(err)
	}
	remainingCoins := mintedCoins.Sub(burnedCoin)
	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, remainingCoins)
	if err != nil {
		panic(err)
	}

	if mintedCoin.Amount.IsInt64() {
		defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			//sdk.NewAttribute(types.AttributeKeyBondedRatio, bondedRatio.String()),
			//sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
			//sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBurnMinted,
			//sdk.NewAttribute(types.AttributeKeyBondedRatio, bondedRatio.String()),
			//sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
			//sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, burnedCoin.Amount.String()),
		),
	)
}
