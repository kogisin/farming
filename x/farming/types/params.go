package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyPrivatePlanCreationFee = []byte("PrivatePlanCreationFee")
	KeyEpochDays              = []byte("EpochDays")
	KeyFarmingFeeCollector    = []byte("FarmingFeeCollector")

	DefaultPrivatePlanCreationFee = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000)))
	DefaultEpochDays              = uint32(1)
	DefaultFarmingFeeCollector    = sdk.AccAddress(address.Module(ModuleName, []byte("FarmingFeeCollectorAcc"))).String()
	StakingReserveAcc             = sdk.AccAddress(address.Module(ModuleName, []byte("StakingReserveAcc")))
	RewardsReserveAcc             = sdk.AccAddress(address.Module(ModuleName, []byte("RewardsReserveAcc")))
)

var _ paramstypes.ParamSet = (*Params)(nil)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns the default farming module parameters.
func DefaultParams() Params {
	return Params{
		PrivatePlanCreationFee: DefaultPrivatePlanCreationFee,
		EpochDays:              DefaultEpochDays,
		FarmingFeeCollector:    DefaultFarmingFeeCollector,
	}
}

// ParamSetPairs implements paramstypes.ParamSet.
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyPrivatePlanCreationFee, &p.PrivatePlanCreationFee, validatePrivatePlanCreationFee),
		paramstypes.NewParamSetPair(KeyEpochDays, &p.EpochDays, validateEpochDays),
		paramstypes.NewParamSetPair(KeyFarmingFeeCollector, &p.FarmingFeeCollector, validateFarmingFeeCollector),
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate validates parameters.
func (p Params) Validate() error {
	for _, v := range []struct {
		value     interface{}
		validator func(interface{}) error
	}{
		{p.PrivatePlanCreationFee, validatePrivatePlanCreationFee},
		{p.EpochDays, validateEpochDays},
		{p.FarmingFeeCollector, validateFarmingFeeCollector},
	} {
		if err := v.validator(v.value); err != nil {
			return err
		}
	}
	return nil
}

func validatePrivatePlanCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := v.Validate(); err != nil {
		return err
	}

	return nil
}

func validateEpochDays(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("epoch days must be positive: %d", v)
	}

	return nil
}

func validateFarmingFeeCollector(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("farming fee collector address must not be empty")
	}

	_, err := sdk.AccAddressFromBech32(v)
	if err != nil {
		return fmt.Errorf("invalid account address: %v", v)
	}

	return nil
}
