package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	simapp "github.com/tendermint/farming/app"
	"github.com/tendermint/farming/x/farming/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestAllocationInfos() {
	normalPlans := []types.PlanI{
		types.NewFixedAmountPlan(
			types.NewBasePlan(
				1,
				"",
				types.PlanTypePrivate,
				suite.addrs[0].String(),
				suite.addrs[0].String(),
				sdk.NewDecCoins(sdk.NewDecCoinFromDec(denom1, sdk.NewDec(1))),
				types.ParseTime("2021-07-27T00:00:00Z"),
				types.ParseTime("2021-07-28T00:00:00Z"),
			),
			sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000))),
		types.NewFixedAmountPlan(
			types.NewBasePlan(
				2,
				"",
				types.PlanTypePrivate,
				suite.addrs[0].String(),
				suite.addrs[0].String(),
				sdk.NewDecCoins(sdk.NewDecCoinFromDec(denom1, sdk.NewDec(1))),
				types.ParseTime("2021-07-27T12:00:00Z"),
				types.ParseTime("2021-07-28T12:00:00Z"),
			),
			sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000))),
	}

	for _, tc := range []struct {
		name      string
		plans     []types.PlanI
		t         time.Time
		distrAmts map[uint64]sdk.Coins // planID => sdk.Coins
	}{
		{
			"insufficient farming pool balances",
			[]types.PlanI{
				types.NewFixedAmountPlan(
					types.NewBasePlan(
						1,
						"",
						types.PlanTypePrivate,
						suite.addrs[0].String(),
						suite.addrs[0].String(),
						sdk.NewDecCoins(sdk.NewDecCoinFromDec(denom1, sdk.NewDec(1))),
						types.ParseTime("2021-07-27T00:00:00Z"),
						types.ParseTime("2021-07-30T00:00:00Z"),
					),
					sdk.NewCoins(sdk.NewInt64Coin(denom3, 10_000_000_000))),
			},
			types.ParseTime("2021-07-28T00:00:00Z"),
			nil,
		},
		{
			"start time & end time edgecase #1",
			normalPlans,
			types.ParseTime("2021-07-26T23:59:59Z"),
			nil,
		},
		{
			"start time & end time edgecase #2",
			normalPlans,
			types.ParseTime("2021-07-27T00:00:00Z"),
			map[uint64]sdk.Coins{1: sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000))},
		},
		{
			"start time & end time edgecase #3",
			normalPlans,
			types.ParseTime("2021-07-27T11:59:59Z"),
			map[uint64]sdk.Coins{1: sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000))},
		},
		{
			"start time & end time edgecase #4",
			normalPlans,
			types.ParseTime("2021-07-27T12:00:00Z"),
			map[uint64]sdk.Coins{
				1: sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000)),
				2: sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000))},
		},
		{
			"start time & end time edgecase #5",
			normalPlans,
			types.ParseTime("2021-07-27T23:59:59Z"),
			map[uint64]sdk.Coins{
				1: sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000)),
				2: sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000))},
		},
		{
			"start time & end time edgecase #6",
			normalPlans,
			types.ParseTime("2021-07-28T00:00:00Z"),
			map[uint64]sdk.Coins{2: sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000))},
		},
		{
			"start time & end time edgecase #7",
			normalPlans,
			types.ParseTime("2021-07-28T11:59:59Z"),
			map[uint64]sdk.Coins{2: sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000))},
		},
		{
			"start time & end time edgecase #8",
			normalPlans,
			types.ParseTime("2021-07-28T12:00:00Z"),
			nil,
		},
	} {
		suite.Run(tc.name, func() {
			for _, plan := range tc.plans {
				suite.keeper.SetPlan(suite.ctx, plan)
			}

			suite.ctx = suite.ctx.WithBlockTime(tc.t)
			distrInfos := suite.keeper.AllocationInfos(suite.ctx)
			if suite.Len(distrInfos, len(tc.distrAmts)) {
				for _, distrInfo := range distrInfos {
					distrAmt, ok := tc.distrAmts[distrInfo.Plan.GetId()]
					if suite.True(ok) {
						suite.True(coinsEq(distrAmt, distrInfo.Amount))
					}
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestAllocateRewards() {
	for _, plan := range suite.sampleFixedAmtPlans {
		_ = plan.SetStartTime(types.ParseTime("0001-01-01T00:00:00Z"))
		_ = plan.SetEndTime(types.ParseTime("9999-12-31T00:00:00Z"))
		suite.keeper.SetPlan(suite.ctx, plan)
	}

	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000)))
	suite.Stake(suite.addrs[1], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000), sdk.NewInt64Coin(denom2, 1000000)))
	suite.keeper.ProcessQueuedCoins(suite.ctx)

	prevDistrCoins := map[uint64]sdk.Coins{}

	t := types.ParseTime("2021-09-01T00:00:00Z")
	for i := 0; i < 365; i++ {
		suite.ctx = suite.ctx.WithBlockTime(t)

		err := suite.keeper.AllocateRewards(suite.ctx)
		suite.Require().NoError(err)

		for _, plan := range suite.sampleFixedAmtPlans {
			plan, _ := suite.keeper.GetPlan(suite.ctx, plan.GetId())
			fixedAmtPlan := plan.(*types.FixedAmountPlan)

			dist := plan.GetDistributedCoins()
			suite.Require().True(coinsEq(prevDistrCoins[plan.GetId()].Add(fixedAmtPlan.EpochAmount...), dist))
			prevDistrCoins[plan.GetId()] = dist

			t2 := plan.GetLastDistributionTime()
			suite.Require().NotNil(t2)
			suite.Require().Equal(t, *t2)
		}

		t = t.AddDate(0, 0, 1)
	}
}

func (suite *KeeperTestSuite) TestAllocateRewards_FixedAmountPlanAllBalances() {
	farmingPoolAcc := simapp.AddTestAddrs(suite.app, suite.ctx, 1, sdk.ZeroInt())[0]
	err := simapp.FundAccount(suite.app.BankKeeper, suite.ctx, farmingPoolAcc, sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000000)))
	suite.Require().NoError(err)

	// The sum of epoch ratios is exactly 1.
	suite.SetFixedAmountPlan(1, farmingPoolAcc, map[string]string{denom1: "1.0"}, map[string]int64{denom3: 600000})
	suite.SetFixedAmountPlan(2, farmingPoolAcc, map[string]string{denom2: "1.0"}, map[string]int64{denom3: 400000})

	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000), sdk.NewInt64Coin(denom2, 1000000)))

	suite.AdvanceEpoch()
	suite.AdvanceEpoch()

	rewards := suite.Rewards(suite.addrs[0])
	suite.Require().True(coinsEq(sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000000)), rewards))
}

func (suite *KeeperTestSuite) TestAllocateRewards_RatioPlanAllBalances() {
	farmingPoolAcc := simapp.AddTestAddrs(suite.app, suite.ctx, 1, sdk.ZeroInt())[0]
	err := simapp.FundAccount(suite.app.BankKeeper, suite.ctx, farmingPoolAcc, sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000000)))
	suite.Require().NoError(err)

	// The sum of epoch ratios is exactly 1.
	suite.SetRatioPlan(1, farmingPoolAcc, map[string]string{denom1: "1.0"}, "0.5")
	suite.SetRatioPlan(2, farmingPoolAcc, map[string]string{denom2: "1.0"}, "0.5")

	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000), sdk.NewInt64Coin(denom2, 1000000)))

	suite.AdvanceEpoch()
	suite.AdvanceEpoch()

	rewards := suite.Rewards(suite.addrs[0])
	suite.Require().True(coinsEq(sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000000)), rewards))
}

func (suite *KeeperTestSuite) TestAllocateRewards_FixedAmountPlanOverBalances() {
	farmingPoolAcc := simapp.AddTestAddrs(suite.app, suite.ctx, 1, sdk.ZeroInt())[0]
	err := simapp.FundAccount(suite.app.BankKeeper, suite.ctx, farmingPoolAcc, sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000000)))
	suite.Require().NoError(err)

	// The sum of epoch amounts is over the balances the farming pool has,
	// so the reward allocation should never happen.
	suite.SetFixedAmountPlan(1, farmingPoolAcc, map[string]string{denom1: "1.0"}, map[string]int64{denom3: 700000})
	suite.SetFixedAmountPlan(2, farmingPoolAcc, map[string]string{denom2: "1.0"}, map[string]int64{denom3: 400000})

	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000), sdk.NewInt64Coin(denom2, 1000000)))

	suite.AdvanceEpoch()
	suite.AdvanceEpoch()

	rewards := suite.Rewards(suite.addrs[0])
	suite.Require().True(rewards.IsZero())
}

func (suite *KeeperTestSuite) TestAllocateRewards_RatioPlanOverBalances() {
	farmingPoolAcc := simapp.AddTestAddrs(suite.app, suite.ctx, 1, sdk.ZeroInt())[0]
	err := simapp.FundAccount(suite.app.BankKeeper, suite.ctx, farmingPoolAcc, sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000000)))
	suite.Require().NoError(err)

	// The sum of epoch ratios is over 1, so the reward allocation should never happen.
	suite.SetRatioPlan(1, farmingPoolAcc, map[string]string{denom1: "1.0"}, "0.8")
	suite.SetRatioPlan(2, farmingPoolAcc, map[string]string{denom2: "1.0"}, "0.5")

	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000), sdk.NewInt64Coin(denom2, 1000000)))

	suite.AdvanceEpoch()
	suite.AdvanceEpoch()

	rewards := suite.Rewards(suite.addrs[0])
	suite.Require().True(rewards.IsZero())
}

func (suite *KeeperTestSuite) TestOutstandingRewards() {
	// The block time here is not important, and has chosen randomly.
	suite.ctx = suite.ctx.WithBlockTime(types.ParseTime("2021-09-01T00:00:00Z"))

	suite.SetFixedAmountPlan(1, suite.addrs[4], map[string]string{
		denom1: "1",
	}, map[string]int64{
		denom3: 1000,
	})

	// Three farmers stake same amount of coins.
	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000)))
	suite.Stake(suite.addrs[1], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000)))
	suite.Stake(suite.addrs[2], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000)))

	// At first, the outstanding rewards shouldn't exist.
	_, found := suite.keeper.GetOutstandingRewards(suite.ctx, denom1)
	suite.Require().False(found)

	suite.AdvanceEpoch() // Queued staking coins have now staked.
	suite.AdvanceEpoch() // Allocate rewards for staked coins.

	// After the first allocation of rewards, the outstanding rewards should be 1000denom3.
	outstanding, found := suite.keeper.GetOutstandingRewards(suite.ctx, denom1)
	suite.Require().True(found)
	suite.Require().True(decCoinsEq(sdk.NewDecCoins(sdk.NewInt64DecCoin(denom3, 1000)), outstanding.Rewards))

	// All farmers harvest rewards, so the outstanding rewards should be (approximately)0.
	suite.Harvest(suite.addrs[0], []string{denom1})
	suite.Harvest(suite.addrs[1], []string{denom1})
	suite.Harvest(suite.addrs[2], []string{denom1})

	outstanding, _ = suite.keeper.GetOutstandingRewards(suite.ctx, denom1)
	truncatedOutstanding, _ := outstanding.Rewards.TruncateDecimal()
	suite.Require().True(truncatedOutstanding.IsZero())
}

func (suite *KeeperTestSuite) TestHarvest() {
	for _, plan := range suite.samplePlans {
		suite.keeper.SetPlan(suite.ctx, plan)
	}

	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1_000_000)))
	suite.keeper.ProcessQueuedCoins(suite.ctx)

	balancesBefore := suite.app.BankKeeper.GetAllBalances(suite.ctx, suite.addrs[0])
	suite.ctx = suite.ctx.WithBlockTime(types.ParseTime("2021-08-05T00:00:00Z"))
	err := suite.keeper.AllocateRewards(suite.ctx)
	suite.Require().NoError(err)

	rewards := suite.Rewards(suite.addrs[0])

	err = suite.keeper.Harvest(suite.ctx, suite.addrs[0], []string{denom1})
	suite.Require().NoError(err)

	balancesAfter := suite.app.BankKeeper.GetAllBalances(suite.ctx, suite.addrs[0])
	suite.Require().True(coinsEq(balancesBefore.Add(rewards...), balancesAfter))
	suite.Require().True(suite.app.BankKeeper.GetAllBalances(suite.ctx, suite.keeper.GetRewardsReservePoolAcc(suite.ctx)).IsZero())
	suite.Require().True(suite.Rewards(suite.addrs[0]).IsZero())
}

func (suite *KeeperTestSuite) TestMultipleHarvest() {
	// TODO: implement
}
