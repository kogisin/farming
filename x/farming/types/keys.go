package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName is the name of the farming module
	ModuleName = "farming"

	// RouterKey is the message router key for the farming module
	RouterKey = ModuleName

	// StoreKey is the default store key for the farming module
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the farming module
	QuerierRoute = ModuleName
)

var (
	GlobalPlanIdKey        = []byte("globalPlanId")
	GlobalLastEpochTimeKey = []byte("globalLastEpochTime")

	PlanKeyPrefix               = []byte{0x11}
	PlansByFarmerIndexKeyPrefix = []byte{0x12}

	StakingKeyPrefix            = []byte{0x21}
	StakingIndexKeyPrefix       = []byte{0x22}
	QueuedStakingKeyPrefix      = []byte{0x23}
	QueuedStakingIndexKeyPrefix = []byte{0x24}
	TotalStakingKeyPrefix       = []byte{0x25}

	HistoricalRewardsKeyPrefix = []byte{0x31}
	CurrentEpochKeyPrefix      = []byte{0x32}
)

// GetPlanKey returns kv indexing key of the plan
func GetPlanKey(planID uint64) []byte {
	return append(PlanKeyPrefix, sdk.Uint64ToBigEndian(planID)...)
}

// GetPlansByFarmerIndexKey returns kv indexing key of the plan indexed by reserve account
func GetPlansByFarmerIndexKey(farmerAcc sdk.AccAddress) []byte {
	return append(PlansByFarmerIndexKeyPrefix, address.MustLengthPrefix(farmerAcc.Bytes())...)
}

// GetPlanByFarmerAddrIndexKey returns kv indexing key of the plan indexed by reserve account
func GetPlanByFarmerAddrIndexKey(farmerAcc sdk.AccAddress, planID uint64) []byte {
	return append(append(PlansByFarmerIndexKeyPrefix, address.MustLengthPrefix(farmerAcc.Bytes())...), sdk.Uint64ToBigEndian(planID)...)
}

// GetStakingKey returns a key for staking of corresponding the id
func GetStakingKey(stakingCoinDenom string, farmerAcc sdk.AccAddress) []byte {
	return append(append(StakingKeyPrefix, LengthPrefixString(stakingCoinDenom)...), farmerAcc...)
}

func GetStakingIndexKey(farmerAcc sdk.AccAddress, stakingCoinDenom string) []byte {
	return append(append(StakingIndexKeyPrefix, address.MustLengthPrefix(farmerAcc)...), []byte(stakingCoinDenom)...)
}

func GetStakingsByFarmerPrefix(farmerAcc sdk.AccAddress) []byte {
	return append(StakingIndexKeyPrefix, address.MustLengthPrefix(farmerAcc)...)
}

func GetQueuedStakingKey(stakingCoinDenom string, farmerAcc sdk.AccAddress) []byte {
	return append(append(QueuedStakingKeyPrefix, LengthPrefixString(stakingCoinDenom)...), farmerAcc...)
}

func GetQueuedStakingIndexKey(farmerAcc sdk.AccAddress, stakingCoinDenom string) []byte {
	return append(append(QueuedStakingIndexKeyPrefix, address.MustLengthPrefix(farmerAcc)...), []byte(stakingCoinDenom)...)
}

func GetQueuedStakingByFarmerPrefix(farmerAcc sdk.AccAddress) []byte {
	return append(QueuedStakingIndexKeyPrefix, address.MustLengthPrefix(farmerAcc)...)
}

func GetTotalStakingKey(stakingCoinDenom string) []byte {
	return append(TotalStakingKeyPrefix, []byte(stakingCoinDenom)...)
}

func GetHistoricalRewardsKey(stakingCoinDenom string, epoch uint64) []byte {
	return append(append(HistoricalRewardsKeyPrefix, LengthPrefixString(stakingCoinDenom)...), sdk.Uint64ToBigEndian(epoch)...)
}

func GetCurrentEpochKey(stakingCoinDenom string) []byte {
	return append(CurrentEpochKeyPrefix, []byte(stakingCoinDenom)...)
}

func ParseStakingKey(key []byte) (stakingCoinDenom string, farmerAcc sdk.AccAddress) {
	if !bytes.HasPrefix(key, StakingKeyPrefix) {
		panic("key does not have proper prefix")
	}
	denomLen := key[1]
	stakingCoinDenom = string(key[2 : 2+denomLen])
	farmerAcc = key[2+denomLen:]
	return
}

func ParseStakingIndexKey(key []byte) (farmerAcc sdk.AccAddress, stakingCoinDenom string) {
	if !bytes.HasPrefix(key, StakingIndexKeyPrefix) {
		panic("key does not have proper prefix")
	}
	addrLen := key[1]
	farmerAcc = key[2 : 2+addrLen]
	stakingCoinDenom = string(key[2+addrLen:])
	return
}

func ParseQueuedStakingKey(key []byte) (stakingCoinDenom string, farmerAcc sdk.AccAddress) {
	if !bytes.HasPrefix(key, QueuedStakingKeyPrefix) {
		panic("key does not have proper prefix")
	}
	denomLen := key[1]
	stakingCoinDenom = string(key[2 : 2+denomLen])
	farmerAcc = key[2+denomLen:]
	return
}

func ParseQueuedStakingIndexKey(key []byte) (farmerAcc sdk.AccAddress, stakingCoinDenom string) {
	if !bytes.HasPrefix(key, QueuedStakingIndexKeyPrefix) {
		panic("key does not have proper prefix")
	}
	addrLen := key[1]
	farmerAcc = key[2 : 2+addrLen]
	stakingCoinDenom = string(key[2+addrLen:])
	return
}

func ParseHistoricalRewardsKey(key []byte) (stakingCoinDenom string, epoch uint64) {
	if !bytes.HasPrefix(key, HistoricalRewardsKeyPrefix) {
		panic("key does not have proper prefix")
	}
	denomLen := key[1]
	stakingCoinDenom = string(key[2 : 2+denomLen])
	epoch = sdk.BigEndianToUint64(key[2+denomLen:])
	return
}

func ParseCurrentEpochKey(key []byte) (stakingCoinDenom string) {
	if !bytes.HasPrefix(key, CurrentEpochKeyPrefix) {
		panic("key does not have proper prefix")
	}
	stakingCoinDenom = string(key[1:])
	return
}

// LengthPrefixString is LengthPrefix for string.
func LengthPrefixString(s string) []byte {
	bz := []byte(s)
	bzLen := len(bz)
	if bzLen == 0 {
		return bz
	}
	return append([]byte{byte(bzLen)}, bz...)
}
