package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/tendermint/farming/x/farming/types"
)

// Simulation parameter constants.
const (
	PrivatePlanCreationFee = "private_plan_creation_fee"
	EpochDays              = "epoch_days"
	FarmingFeeCollector    = "farming_fee_collector"
)

// GenPrivatePlanCreationFee return randomized private plan creation fee.
func GenPrivatePlanCreationFee(r *rand.Rand) sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(simulation.RandIntBetween(r, 0, 100_000_000))))
}

// GenEpochDays return default epoch days.
func GenEpochDays(r *rand.Rand) uint32 {
	return uint32(simulation.RandIntBetween(r, int(types.DefaultEpochDays), 10))
}

// GenFarmingFeeCollector returns default farming fee collector.
func GenFarmingFeeCollector(r *rand.Rand) string {
	return types.DefaultFarmingFeeCollector
}

// RandomizedGenState generates a random GenesisState for farming.
func RandomizedGenState(simState *module.SimulationState) {
	var privatePlanCreationFee sdk.Coins
	simState.AppParams.GetOrGenerate(
		simState.Cdc, PrivatePlanCreationFee, &privatePlanCreationFee, simState.Rand,
		func(r *rand.Rand) { privatePlanCreationFee = GenPrivatePlanCreationFee(r) },
	)

	var epochDays uint32
	simState.AppParams.GetOrGenerate(
		simState.Cdc, EpochDays, &epochDays, simState.Rand,
		func(r *rand.Rand) { epochDays = GenEpochDays(r) },
	)

	var feeCollector string
	simState.AppParams.GetOrGenerate(
		simState.Cdc, FarmingFeeCollector, &feeCollector, simState.Rand,
		func(r *rand.Rand) { feeCollector = GenFarmingFeeCollector(r) },
	)

	farmingGenesis := types.GenesisState{
		Params: types.Params{
			PrivatePlanCreationFee: privatePlanCreationFee,
			EpochDays:              epochDays,
			FarmingFeeCollector:    feeCollector,
		},
	}

	bz, _ := json.MarshalIndent(&farmingGenesis, "", " ")
	fmt.Printf("Selected randomly generated farming parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&farmingGenesis)
}
