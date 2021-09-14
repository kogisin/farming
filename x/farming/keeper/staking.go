package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/farming/x/farming/types"
)

// GetStaking returns a specific staking identified by id.
func (k Keeper) GetStaking(ctx sdk.Context, stakingCoinDenom string, farmerAcc sdk.AccAddress) (staking types.Staking, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetStakingKey(stakingCoinDenom, farmerAcc))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshal(bz, &staking)
	found = true
	return
}

// SetStaking implements Staking.
func (k Keeper) SetStaking(ctx sdk.Context, stakingCoinDenom string, farmerAcc sdk.AccAddress, staking types.Staking) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&staking)
	store.Set(types.GetStakingKey(stakingCoinDenom, farmerAcc), bz)
	store.Set(types.GetStakingIndexKey(farmerAcc, stakingCoinDenom), []byte{})
}

func (k Keeper) DeleteStaking(ctx sdk.Context, stakingCoinDenom string, farmerAcc sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetStakingKey(stakingCoinDenom, farmerAcc))
	store.Delete(types.GetStakingIndexKey(farmerAcc, stakingCoinDenom))
}

func (k Keeper) IterateStakings(ctx sdk.Context, cb func(stakingCoinDenom string, farmerAcc sdk.AccAddress, staking types.Staking) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.StakingKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var staking types.Staking
		k.cdc.MustUnmarshal(iter.Value(), &staking)
		stakingCoinDenom, farmerAcc := types.ParseStakingKey(iter.Key())
		if cb(stakingCoinDenom, farmerAcc, staking) {
			break
		}
	}
}

func (k Keeper) IterateStakingsByFarmer(ctx sdk.Context, farmerAcc sdk.AccAddress, cb func(stakingCoinDenom string, staking types.Staking) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetStakingsByFarmerPrefix(farmerAcc))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		farmerAcc, stakingCoinDenom := types.ParseStakingIndexKey(iter.Key())
		staking, _ := k.GetStaking(ctx, stakingCoinDenom, farmerAcc)
		if cb(stakingCoinDenom, staking) {
			break
		}
	}
}

func (k Keeper) GetAllStakedCoinsByFarmer(ctx sdk.Context, farmerAcc sdk.AccAddress) sdk.Coins {
	stakedCoins := sdk.NewCoins()
	k.IterateStakingsByFarmer(ctx, farmerAcc, func(stakingCoinDenom string, staking types.Staking) (stop bool) {
		stakedCoins = stakedCoins.Add(sdk.NewCoin(stakingCoinDenom, staking.Amount))
		return false
	})
	return stakedCoins
}

func (k Keeper) GetQueuedStaking(ctx sdk.Context, stakingCoinDenom string, farmerAcc sdk.AccAddress) (queuedStaking types.QueuedStaking, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetQueuedStakingKey(stakingCoinDenom, farmerAcc))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshal(bz, &queuedStaking)
	found = true
	return
}

func (k Keeper) GetAllQueuedStakedCoinsByFarmer(ctx sdk.Context, farmerAcc sdk.AccAddress) sdk.Coins {
	stakedCoins := sdk.NewCoins()
	k.IterateQueuedStakingsByFarmer(ctx, farmerAcc, func(stakingCoinDenom string, queuedStaking types.QueuedStaking) (stop bool) {
		stakedCoins = stakedCoins.Add(sdk.NewCoin(stakingCoinDenom, queuedStaking.Amount))
		return false
	})
	return stakedCoins
}

func (k Keeper) SetQueuedStaking(ctx sdk.Context, stakingCoinDenom string, farmerAcc sdk.AccAddress, queuedStaking types.QueuedStaking) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&queuedStaking)
	store.Set(types.GetQueuedStakingKey(stakingCoinDenom, farmerAcc), bz)
	store.Set(types.GetQueuedStakingIndexKey(farmerAcc, stakingCoinDenom), []byte{})
}

func (k Keeper) DeleteQueuedStaking(ctx sdk.Context, stakingCoinDenom string, farmerAcc sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetQueuedStakingKey(stakingCoinDenom, farmerAcc))
	store.Delete(types.GetQueuedStakingIndexKey(farmerAcc, stakingCoinDenom))
}

func (k Keeper) IterateQueuedStakings(ctx sdk.Context, cb func(stakingCoinDenom string, farmerAcc sdk.AccAddress, queuedStaking types.QueuedStaking) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.QueuedStakingKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var queuedStaking types.QueuedStaking
		k.cdc.MustUnmarshal(iter.Value(), &queuedStaking)
		stakingCoinDenom, farmerAcc := types.ParseQueuedStakingKey(iter.Key())
		if cb(stakingCoinDenom, farmerAcc, queuedStaking) {
			break
		}
	}
}

func (k Keeper) IterateQueuedStakingsByFarmer(ctx sdk.Context, farmerAcc sdk.AccAddress, cb func(stakingCoinDenom string, queuedStaking types.QueuedStaking) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetQueuedStakingByFarmerPrefix(farmerAcc))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		farmerAcc, stakingCoinDenom := types.ParseQueuedStakingIndexKey(iter.Key())
		queuedStaking, _ := k.GetQueuedStaking(ctx, stakingCoinDenom, farmerAcc)
		if cb(stakingCoinDenom, queuedStaking) {
			break
		}
	}
}

func (k Keeper) GetTotalStaking(ctx sdk.Context, stakingCoinDenom string) (totalStaking types.TotalStaking, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTotalStakingKey(stakingCoinDenom))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshal(bz, &totalStaking)
	found = true
	return
}

func (k Keeper) SetTotalStaking(ctx sdk.Context, stakingCoinDenom string, totalStaking types.TotalStaking) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&totalStaking)
	store.Set(types.GetTotalStakingKey(stakingCoinDenom), bz)
}

// ReserveStakingCoins sends staking coins to the staking reserve account.
func (k Keeper) ReserveStakingCoins(ctx sdk.Context, farmerAcc sdk.AccAddress, stakingCoins sdk.Coins) error {
	if err := k.bankKeeper.SendCoins(ctx, farmerAcc, k.GetStakingReservePoolAcc(ctx), stakingCoins); err != nil {
		return err
	}
	return nil
}

// ReleaseStakingCoins sends staking coins back to the farmer.
func (k Keeper) ReleaseStakingCoins(ctx sdk.Context, farmerAcc sdk.AccAddress, unstakingCoins sdk.Coins) error {
	if err := k.bankKeeper.SendCoins(ctx, k.GetStakingReservePoolAcc(ctx), farmerAcc, unstakingCoins); err != nil {
		return err
	}
	return nil
}

// Stake stores staking coins to queued coins, and it will be processed in the next epoch.
func (k Keeper) Stake(ctx sdk.Context, farmerAcc sdk.AccAddress, amount sdk.Coins) error {
	if err := k.ReserveStakingCoins(ctx, farmerAcc, amount); err != nil {
		return err
	}

	for _, coin := range amount {
		queuedStaking, found := k.GetQueuedStaking(ctx, coin.Denom, farmerAcc)
		if !found {
			queuedStaking.Amount = sdk.ZeroInt()
		}
		queuedStaking.Amount = queuedStaking.Amount.Add(coin.Amount)
		k.SetQueuedStaking(ctx, coin.Denom, farmerAcc, queuedStaking)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeStake,
			sdk.NewAttribute(types.AttributeKeyFarmer, farmerAcc.String()),
			sdk.NewAttribute(types.AttributeKeyStakingCoins, amount.String()),
		),
	})

	return nil
}

// Unstake unstakes an amount of staking coins from the staking reserve account.
func (k Keeper) Unstake(ctx sdk.Context, farmerAcc sdk.AccAddress, amount sdk.Coins) error {
	// TODO: send coins at once, not in every WithdrawRewards

	for _, coin := range amount {
		staking, found := k.GetStaking(ctx, coin.Denom, farmerAcc)
		if !found {
			staking.Amount = sdk.ZeroInt()
		}

		queuedStaking, found := k.GetQueuedStaking(ctx, coin.Denom, farmerAcc)
		if !found {
			queuedStaking.Amount = sdk.ZeroInt()
		}

		availableAmt := staking.Amount.Add(queuedStaking.Amount)
		if availableAmt.LT(coin.Amount) {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInsufficientFunds, "%s%s is smaller than %s%s", availableAmt, coin.Denom, coin.Amount, coin.Denom)
		}

		if staking.Amount.IsPositive() {
			if _, err := k.WithdrawRewards(ctx, farmerAcc, coin.Denom); err != nil {
				return err
			}
		}

		removedFromStaking := sdk.ZeroInt()

		queuedStaking.Amount = queuedStaking.Amount.Sub(coin.Amount)
		if queuedStaking.Amount.IsNegative() {
			staking.Amount = staking.Amount.Add(queuedStaking.Amount)
			removedFromStaking = queuedStaking.Amount.Neg()
			queuedStaking.Amount = sdk.ZeroInt()
			if staking.Amount.IsPositive() {
				currentEpoch := k.GetCurrentEpoch(ctx, coin.Denom)
				staking.StartingEpoch = currentEpoch
				k.SetStaking(ctx, coin.Denom, farmerAcc, staking)
			} else {
				k.DeleteStaking(ctx, coin.Denom, farmerAcc)
			}
		}

		if queuedStaking.Amount.IsPositive() {
			k.SetQueuedStaking(ctx, coin.Denom, farmerAcc, queuedStaking)
		} else {
			k.DeleteQueuedStaking(ctx, coin.Denom, farmerAcc)
		}

		totalStaking, found := k.GetTotalStaking(ctx, coin.Denom)
		if !found {
			totalStaking.Amount = sdk.ZeroInt()
		}
		totalStaking.Amount = totalStaking.Amount.Sub(removedFromStaking)
		k.SetTotalStaking(ctx, coin.Denom, totalStaking)
	}

	if err := k.ReleaseStakingCoins(ctx, farmerAcc, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnstake,
			sdk.NewAttribute(types.AttributeKeyFarmer, farmerAcc.String()),
			sdk.NewAttribute(types.AttributeKeyUnstakingCoins, amount.String()),
		),
	})

	return nil
}

// ProcessQueuedCoins moves queued coins into staked coins.
func (k Keeper) ProcessQueuedCoins(ctx sdk.Context) {
	k.IterateQueuedStakings(ctx, func(stakingCoinDenom string, farmerAcc sdk.AccAddress, queuedStaking types.QueuedStaking) (stop bool) {
		staking, found := k.GetStaking(ctx, stakingCoinDenom, farmerAcc)
		if found {
			if _, err := k.WithdrawRewards(ctx, farmerAcc, stakingCoinDenom); err != nil {
				panic(err)
			}
		} else {
			staking.Amount = sdk.ZeroInt()
		}

		k.DeleteQueuedStaking(ctx, stakingCoinDenom, farmerAcc)
		k.SetStaking(ctx, stakingCoinDenom, farmerAcc, types.Staking{
			Amount:        staking.Amount.Add(queuedStaking.Amount),
			StartingEpoch: k.GetCurrentEpoch(ctx, stakingCoinDenom),
		})

		totalStaking, found := k.GetTotalStaking(ctx, stakingCoinDenom)
		if !found {
			totalStaking.Amount = sdk.ZeroInt()
		}
		k.SetTotalStaking(ctx, stakingCoinDenom, types.TotalStaking{
			Amount: totalStaking.Amount.Add(queuedStaking.Amount),
		})

		return false
	})
}
