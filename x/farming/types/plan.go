package types

import (
	"strconv"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	MaxNameLength int = 140
)

var (
	_ PlanI = (*FixedAmountPlan)(nil)
	_ PlanI = (*RatioPlan)(nil)
)

// NewBasePlan creates a new BasePlan object
//nolint:interfacer
func NewBasePlan(id uint64, name string, typ PlanType, farmingPoolAddr, terminationAddr string, coinWeights sdk.DecCoins, startTime, endTime time.Time) *BasePlan {
	basePlan := &BasePlan{
		Id:                 id,
		Name:               name,
		Type:               typ,
		FarmingPoolAddress: farmingPoolAddr,
		RewardPoolAddress:  GenerateRewardPoolAcc(PlanUniqueKey(id, typ, farmingPoolAddr)).String(),
		TerminationAddress: terminationAddr,
		StakingCoinWeights: coinWeights,
		StartTime:          startTime,
		EndTime:            endTime,
	}
	return basePlan
}

func (plan BasePlan) GetId() uint64 { //nolint:golint
	return plan.Id
}

func (plan *BasePlan) SetId(id uint64) error { //nolint:golint
	plan.Id = id
	return nil
}

func (plan BasePlan) GetType() PlanType {
	return plan.Type
}

func (plan *BasePlan) SetType(typ PlanType) error {
	plan.Type = typ
	return nil
}

func (plan BasePlan) GetFarmingPoolAddress() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(plan.FarmingPoolAddress)
	return addr
}

func (plan *BasePlan) SetFarmingPoolAddress(addr sdk.AccAddress) error {
	plan.FarmingPoolAddress = addr.String()
	return nil
}

func (plan BasePlan) GetRewardPoolAddress() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(plan.RewardPoolAddress)
	return addr
}

func (plan *BasePlan) SetRewardPoolAddress(addr sdk.AccAddress) error {
	plan.RewardPoolAddress = addr.String()
	return nil
}

func (plan BasePlan) GetTerminationAddress() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(plan.TerminationAddress)
	return addr
}

func (plan *BasePlan) SetTerminationAddress(addr sdk.AccAddress) error {
	plan.TerminationAddress = addr.String()
	return nil
}

func (plan BasePlan) GetStakingCoinWeights() sdk.DecCoins {
	return plan.StakingCoinWeights
}

func (plan *BasePlan) SetStakingCoinWeights(coinWeights sdk.DecCoins) error {
	plan.StakingCoinWeights = coinWeights
	return nil
}

func (plan BasePlan) GetStartTime() time.Time {
	return plan.StartTime
}

func (plan *BasePlan) SetStartTime(t time.Time) error {
	plan.StartTime = t
	return nil
}

func (plan BasePlan) GetEndTime() time.Time {
	return plan.EndTime
}

func (plan *BasePlan) SetEndTime(t time.Time) error {
	plan.EndTime = t
	return nil
}

// Validate checks for errors on the Plan fields
func (plan BasePlan) Validate() error {
	if plan.Type != PlanTypePrivate && plan.Type != PlanTypePublic {
		return sdkerrors.Wrapf(ErrInvalidPlanType, "unknown plan type: %s", plan.Type)
	}
	if _, err := sdk.AccAddressFromBech32(plan.FarmingPoolAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid farming pool address %q: %v", plan.FarmingPoolAddress, err)
	}
	if _, err := sdk.AccAddressFromBech32(plan.RewardPoolAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid reward pool address %q: %v", plan.RewardPoolAddress, err)
	}
	if _, err := sdk.AccAddressFromBech32(plan.TerminationAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid termination address %q: %v", plan.TerminationAddress, err)
	}
	if len(plan.Name) > MaxNameLength {
		return sdkerrors.Wrapf(ErrInvalidNameLength, "plan name cannot be longer than max length of %d", MaxNameLength)
	}
	if err := plan.StakingCoinWeights.Validate(); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid staking coin weights: %v", err)
	}
	if !plan.EndTime.After(plan.StartTime) {
		return sdkerrors.Wrapf(ErrInvalidPlanEndTime, "end time %s must be greater than start time %s", plan.EndTime, plan.StartTime)
	}
	return nil
}

func (plan BasePlan) String() string {
	out, _ := plan.MarshalYAML()
	return out.(string)
}

// MarshalYAML returns the YAML representation of an Plan.
func (plan BasePlan) MarshalYAML() (interface{}, error) {
	bz, err := codec.MarshalYAML(codec.NewProtoCodec(codectypes.NewInterfaceRegistry()), &plan)
	if err != nil {
		return nil, err
	}
	return string(bz), err
}

func NewFixedAmountPlan(basePlan *BasePlan, epochAmount sdk.Coins) *FixedAmountPlan {
	return &FixedAmountPlan{
		BasePlan:    basePlan,
		EpochAmount: epochAmount,
	}
}

func NewRatioPlan(basePlan *BasePlan, epochRatio sdk.Dec) *RatioPlan {
	return &RatioPlan{
		BasePlan:   basePlan,
		EpochRatio: epochRatio,
	}
}

// PlanUniqueKey returns unique name of the plan consists of given Id, Type and FarmingPoolAddress.
func PlanUniqueKey(id uint64, typ PlanType, farmingPoolAddr string) string {
	poolNameObjects := make([]string, 3)
	poolNameObjects[0] = strconv.FormatUint(id, 10)
	poolNameObjects[1] = strconv.FormatInt(int64(typ), 10)
	poolNameObjects[2] = farmingPoolAddr
	return strings.Join(poolNameObjects, "/")
}

// GenerateRewardPoolAcc returns deterministically generated reward pool account for the given plan name
func GenerateRewardPoolAcc(name string) sdk.AccAddress {
	return address.Module(ModuleName, []byte(strings.Join([]string{RewardPoolAccKeyPrefix, name}, "/")))
}

type PlanI interface {
	proto.Message

	GetId() uint64
	SetId(uint64) error

	GetType() PlanType
	SetType(PlanType) error

	GetFarmingPoolAddress() sdk.AccAddress
	SetFarmingPoolAddress(sdk.AccAddress) error

	GetRewardPoolAddress() sdk.AccAddress
	SetRewardPoolAddress(sdk.AccAddress) error

	GetTerminationAddress() sdk.AccAddress
	SetTerminationAddress(sdk.AccAddress) error

	GetStakingCoinWeights() sdk.DecCoins
	SetStakingCoinWeights(sdk.DecCoins) error

	GetStartTime() time.Time
	SetStartTime(time.Time) error

	GetEndTime() time.Time
	SetEndTime(time.Time) error

	String() string
}

// ValidatePlanName validates if the plan name exists in store.
func ValidatePlanName(plans []PlanI, name string) error {
	for _, plan := range plans {
		switch p := plan.(type) {
		case *FixedAmountPlan:
			if p.Name == name {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "plan name '%s' already exists", name)
			}
		case *RatioPlan:
			if p.Name == name {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "plan name '%s' already exists", name)
			}
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized plan type: %T", p)
		}
	}
	return nil
}

// ValidateOverflowEpochRatio validates that a single account cannot create ratio plans that exceeds
// a total epoch ratio of 1 (100%) because it will overflow the account's balance
func ValidateOverflowEpochRatio(plans []PlanI, farmingPoolAddr string, epochRatio sdk.Dec) error {
	farmerEpochRatio := make(map[string]sdk.Dec)

	for _, plan := range plans {
		farmingPoolAddr := plan.GetFarmingPoolAddress().String()

		if plan, ok := plan.(*RatioPlan); ok {
			farmerEpochRatio[farmingPoolAddr] = plan.EpochRatio
		}
	}

	if epochRatio, ok := farmerEpochRatio[farmingPoolAddr]; ok {
		epochRatio = epochRatio.Add(epochRatio)

		if epochRatio.GT(sdk.NewDec(1)) {
			return sdkerrors.Wrapf(ErrInvalidPlanEpochRatio, "total epoch ratio must be lower than 1")
		}
	}

	return nil
}
