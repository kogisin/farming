package simulation_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	farmingapp "github.com/tendermint/farming/app"
	"github.com/tendermint/farming/app/params"
	"github.com/tendermint/farming/x/farming/simulation"
	"github.com/tendermint/farming/x/farming/types"
)

// TestWeightedOperations tests the weights of the operations.
func TestWeightedOperations(t *testing.T) {
	app, ctx := createTestApp(false)

	ctx.WithChainID("test-chain")

	cdc := app.AppCodec()
	appParams := make(simtypes.AppParams)

	weightedOps := simulation.WeightedOperations(appParams, cdc, app.AccountKeeper, app.BankKeeper, app.FarmingKeeper)

	s := rand.NewSource(1)
	r := rand.New(s)
	accs := getTestingAccounts(t, r, app, ctx, 1)

	expected := []struct {
		weight     int
		opMsgRoute string
		opMsgName  string
	}{
		{params.DefaultWeightMsgCreateFixedAmountPlan, types.ModuleName, types.TypeMsgCreateFixedAmountPlan},
		{params.DefaultWeightMsgCreateRatioPlan, types.ModuleName, types.TypeMsgCreateRatioPlan},
		{params.DefaultWeightMsgStake, types.ModuleName, types.TypeMsgStake},
		{params.DefaultWeightMsgUnstake, types.ModuleName, types.TypeMsgUnstake},
		{params.DefaultWeightMsgHarvest, types.ModuleName, types.TypeMsgHarvest},
	}

	for i, w := range weightedOps {
		operationMsg, _, _ := w.Op()(r, app.BaseApp, ctx, accs, ctx.ChainID())
		// the following checks are very much dependent from the ordering of the output given
		// by WeightedOperations. if the ordering in WeightedOperations changes some tests
		// will fail
		require.Equal(t, expected[i].weight, w.Weight(), "weight should be the same")
		require.Equal(t, expected[i].opMsgRoute, operationMsg.Route, "route should be the same")
		require.Equal(t, expected[i].opMsgName, operationMsg.Name, "operation Msg name should be the same")
	}
}

// TestSimulateMsgCreateFixedAmountPlan tests the normal scenario of a valid message of type TypeMsgCreateFixedAmountPlan.
// Abnormal scenarios, where the message are created by an errors are not tested here.
func TestSimulateMsgCreateFixedAmountPlan(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup a single account
	s := rand.NewSource(1)
	r := rand.New(s)

	accounts := getTestingAccounts(t, r, app, ctx, 1)

	// setup randomly generated private plan creation fees
	feeCoins := simulation.GenPrivatePlanCreationFee(r)
	params := app.FarmingKeeper.GetParams(ctx)
	params.PrivatePlanCreationFee = feeCoins
	app.FarmingKeeper.SetParams(ctx, params)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// execute operation
	op := simulation.SimulateMsgCreateFixedAmountPlan(app.AccountKeeper, app.BankKeeper, app.FarmingKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgCreateFixedAmountPlan
	err = app.AppCodec().UnmarshalJSON(operationMsg.Msg, &msg)
	require.NoError(t, err)

	require.True(t, operationMsg.OK)
	require.Equal(t, types.TypeMsgCreateFixedAmountPlan, msg.Type())
	require.Equal(t, "simulation", msg.Name)
	require.Equal(t, "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj3", msg.Creator)
	require.Equal(t, "1.000000000000000000stake", msg.StakingCoinWeights.String())
	require.Equal(t, "656203300pool93E069B333B5ECEBFE24C6E1437E814003248E0DD7FF8B9F82119F4587449BA5", msg.EpochAmount.String())
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgCreateRatioPlan tests the normal scenario of a valid message of type TypeMsgCreateRatioPlan.
// Abnormal scenarios, where the message are created by an errors are not tested here.
func TestSimulateMsgCreateRatioPlan(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup a single account
	s := rand.NewSource(1)
	r := rand.New(s)

	accounts := getTestingAccounts(t, r, app, ctx, 1)

	// setup randomly generated private plan creation fees
	feeCoins := simulation.GenPrivatePlanCreationFee(r)
	params := app.FarmingKeeper.GetParams(ctx)
	params.PrivatePlanCreationFee = feeCoins
	app.FarmingKeeper.SetParams(ctx, params)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// execute operation
	op := simulation.SimulateMsgCreateRatioPlan(app.AccountKeeper, app.BankKeeper, app.FarmingKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgCreateRatioPlan
	err = app.AppCodec().UnmarshalJSON(operationMsg.Msg, &msg)
	require.NoError(t, err)

	require.True(t, operationMsg.OK)
	require.Equal(t, types.TypeMsgCreateRatioPlan, msg.Type())
	require.Equal(t, "simulation", msg.Name)
	require.Equal(t, "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj3", msg.Creator)
	require.Equal(t, "1.000000000000000000stake", msg.StakingCoinWeights.String())
	require.Equal(t, "0.500000000000000000", msg.EpochRatio.String())
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgStake tests the normal scenario of a valid message of type TypeMsgStake.
// Abnormal scenarios, where the message are created by an errors are not tested here.
func TestSimulateMsgStake(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup a single account
	s := rand.NewSource(1)
	r := rand.New(s)

	accounts := getTestingAccounts(t, r, app, ctx, 1)

	// setup randomly generated staking creation fees
	feeCoins := simulation.GenStakingCreationFee(r)
	params := app.FarmingKeeper.GetParams(ctx)
	params.StakingCreationFee = feeCoins
	app.FarmingKeeper.SetParams(ctx, params)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// execute operation
	op := simulation.SimulateMsgStake(app.AccountKeeper, app.BankKeeper, app.FarmingKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgStake
	err = app.AppCodec().UnmarshalJSON(operationMsg.Msg, &msg)
	require.NoError(t, err)

	require.True(t, operationMsg.OK)
	require.Equal(t, types.TypeMsgStake, msg.Type())
	require.Equal(t, "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj3", msg.Farmer)
	require.Equal(t, "476941318stake", msg.StakingCoins.String())
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgUnstake tests the normal scenario of a valid message of type TypeMsgUnstake.
// Abnormal scenarios, where the message are created by an errors are not tested here.
func TestSimulateMsgUnstake(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup a single account
	s := rand.NewSource(1)
	r := rand.New(s)

	accounts := getTestingAccounts(t, r, app, ctx, 1)

	// setup randomly generated staking creation fees
	feeCoins := simulation.GenStakingCreationFee(r)
	params := app.FarmingKeeper.GetParams(ctx)
	params.StakingCreationFee = feeCoins
	app.FarmingKeeper.SetParams(ctx, params)

	// staking must exist in order to simulate unstake
	stakingCoins := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100_000_000))
	_, err := app.FarmingKeeper.Stake(ctx, accounts[0].Address, stakingCoins)
	require.NoError(t, err)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// execute operation
	op := simulation.SimulateMsgUnstake(app.AccountKeeper, app.BankKeeper, app.FarmingKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgUnstake
	err = app.AppCodec().UnmarshalJSON(operationMsg.Msg, &msg)
	require.NoError(t, err)

	require.True(t, operationMsg.OK)
	require.Equal(t, types.TypeMsgUnstake, msg.Type())
	require.Equal(t, "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj3", msg.Farmer)
	require.Equal(t, "89941318stake", msg.UnstakingCoins.String())
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgHarvest tests the normal scenario of a valid message of type TypeMsgHarvest.
// Abnormal scenarios, where the message are created by an errors are not tested here.
func TestSimulateMsgHarvest(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup a single account
	s := rand.NewSource(1)
	r := rand.New(s)

	accounts := getTestingAccounts(t, r, app, ctx, 1)

	// setup epoch days to 1
	params := app.FarmingKeeper.GetParams(ctx)
	params.EpochDays = 1
	app.FarmingKeeper.SetParams(ctx, params)

	// setup a fixed amount plan
	msgPlan := &types.MsgCreateFixedAmountPlan{
		Name:    "simulation",
		Creator: accounts[0].Address.String(),
		StakingCoinWeights: sdk.NewDecCoins(
			sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.NewDecWithPrec(10, 1)), // 100%
		),
		StartTime:   mustParseRFC3339("2021-08-01T00:00:00Z"),
		EndTime:     mustParseRFC3339("2021-08-31T00:00:00Z"),
		EpochAmount: sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 200_000_000)),
	}

	app.FarmingKeeper.CreateFixedAmountPlan(
		ctx,
		msgPlan,
		accounts[0].Address,
		accounts[0].Address,
		types.PlanTypePrivate,
	)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// set staking and the amount must be greater than the randomized value range for unharvest
	stakingCoins := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100_000_000))
	_, err := app.FarmingKeeper.Stake(ctx, accounts[0].Address, stakingCoins)
	require.NoError(t, err)

	app.FarmingKeeper.ProcessQueuedCoins(ctx)
	ctx = ctx.WithBlockTime(mustParseRFC3339("2021-08-20T00:00:00Z"))
	err = app.FarmingKeeper.DistributeRewards(ctx)
	require.NoError(t, err)

	// check that queue coins are moved to staked coins
	staking, found := app.FarmingKeeper.GetStaking(ctx, uint64(1))
	require.Equal(t, true, found)
	require.Equal(t, true, staking.QueuedCoins.IsZero())
	require.Equal(t, true, staking.StakedCoins.IsAllPositive())

	// execute operation
	op := simulation.SimulateMsgHarvest(app.AccountKeeper, app.BankKeeper, app.FarmingKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgHarvest
	err = app.AppCodec().UnmarshalJSON(operationMsg.Msg, &msg)
	require.NoError(t, err)

	require.True(t, operationMsg.OK)
	require.Equal(t, types.TypeMsgHarvest, msg.Type())
	require.Equal(t, "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj3", msg.Farmer)
	require.Equal(t, []string{"stake"}, msg.StakingCoinDenoms)
	require.Len(t, futureOperations, 0)
}

func createTestApp(isCheckTx bool) (*farmingapp.FarmingApp, sdk.Context) {
	app := farmingapp.Setup(isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.MintKeeper.SetParams(ctx, minttypes.DefaultParams())
	app.MintKeeper.SetMinter(ctx, minttypes.DefaultInitialMinter())

	return app, ctx
}

func getTestingAccounts(t *testing.T, r *rand.Rand, app *farmingapp.FarmingApp, ctx sdk.Context, n int) []simtypes.Account {
	accounts := simtypes.RandomAccounts(r, n)

	initAmt := app.StakingKeeper.TokensFromConsensusPower(ctx, 100_000_000_000)
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initAmt))

	// add coins to the accounts
	for _, account := range accounts {
		acc := app.AccountKeeper.NewAccountWithAddress(ctx, account.Address)
		app.AccountKeeper.SetAccount(ctx, acc)
		err := simapp.FundAccount(app.BankKeeper, ctx, account.Address, initCoins)
		require.NoError(t, err)
	}

	return accounts
}

func mustParseRFC3339(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestHarvest(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup a single account
	s := rand.NewSource(1)
	r := rand.New(s)

	accounts := getTestingAccounts(t, r, app, ctx, 1)

	// mint
	var mintCoins sdk.Coins
	for _, denom := range []string{
		"pool93E069B333B5ECEBFE24C6E1437E814003248E0DD7FF8B9F82119F4587449BA5",
	} {
		mintCoins = mintCoins.Add(sdk.NewInt64Coin(denom, int64(simtypes.RandIntBetween(r, 1e14, 1e15))))
	}

	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, mintCoins)
	require.NoError(t, err)

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accounts[0].Address, mintCoins)
	require.NoError(t, err)

	// setup epoch days to 1
	params := app.FarmingKeeper.GetParams(ctx)
	params.EpochDays = 1
	app.FarmingKeeper.SetParams(ctx, params)

	// setup a fixed amount plan
	msgPlan := &types.MsgCreateFixedAmountPlan{
		Name:    "simulation",
		Creator: accounts[0].Address.String(),
		StakingCoinWeights: sdk.NewDecCoins(
			sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.NewDecWithPrec(10, 1)), // 100%
		),
		StartTime:   ctx.BlockTime(),
		EndTime:     ctx.BlockTime().AddDate(0, 1, 0),
		EpochAmount: sdk.NewCoins(sdk.NewInt64Coin("pool93E069B333B5ECEBFE24C6E1437E814003248E0DD7FF8B9F82119F4587449BA5", 1)),
	}

	app.FarmingKeeper.CreateFixedAmountPlan(
		ctx,
		msgPlan,
		accounts[0].Address,
		accounts[0].Address,
		types.PlanTypePrivate,
	)

	// staking must exist in order to simulate unstake
	_, err = app.FarmingKeeper.Stake(ctx, accounts[0].Address, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100_000_000)))
	require.NoError(t, err)

	// increase epoch days and distribute rewards
	app.FarmingKeeper.ProcessQueuedCoins(ctx)
	ctx = ctx.WithBlockTime(ctx.BlockTime().AddDate(0, 0, 5))
	err = app.FarmingKeeper.DistributeRewards(ctx)
	require.NoError(t, err)
	app.FarmingKeeper.SetLastEpochTime(ctx, ctx.BlockTime())

	plans := app.FarmingKeeper.GetAllPlans(ctx)
	stakings := app.FarmingKeeper.GetAllStakings(ctx)
	rewards := app.FarmingKeeper.GetAllRewards(ctx)
	fmt.Println("<plan>")
	fmt.Println(plans[0])
	fmt.Println("<staking>")
	fmt.Println(stakings[0])
	fmt.Println("<rewards>")
	fmt.Println(rewards)
}
