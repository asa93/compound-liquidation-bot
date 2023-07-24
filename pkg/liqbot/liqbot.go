package liqbot

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gitlab.com/q-dev/exchange-rate-oracle/pkg/config"
	"gitlab.com/q-dev/exchange-rate-oracle/pkg/subgraph"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/q-dev/exchange-rate-oracle/pkg/contracts"
)

// Oracle represents price feed oracle
type Liqbot interface {
	Start(ctx context.Context)
}

// New creates new price feed oracle
func New(logger log.Logger, cfg config.Config) Liqbot {
	return &liqbot{
		logger: logger,
		cfg:    cfg,
	}
}

type liqbot struct {
	logger log.Logger
	cfg    config.Config
}

func (o *liqbot) Start(ctx context.Context) {
	comptroller, ctoken, err := o.getInitialSources()
	if err != nil {
		level.Error(o.logger).Log("msg", err.Error())
		return
	}

	if err != nil {
		level.Error(o.logger).Log("msg", err.Error())
		return
	}

	//to remove
	callerOpts := &bind.CallOpts{
		Pending: false,
		Context: ctx,
	}

	foo, err := comptroller.Admin(callerOpts)
	if err != nil {
		level.Error(o.logger).Log("error comptroller", err.Error())
		return
	}
	level.Info(o.logger).Log("✅ SUCCESS COMPTROLLER CALL", foo)

	for {
		select {
		case <-time.After(o.cfg.UpdateInterval()):
			o.logger.Log("msg", "==== LIQBOT")
			subgraph := subgraph.NewSubgraph()

			// test subgraph to remove
			accounts, err := subgraph.GetAccounts(ctx)

			if err != nil {
				level.Error(o.logger).Log("ERROR SUBGRAPH", err)
			} else {
				level.Info(o.logger).Log("msg", "✅ SUCCESS FETCHING SUBGRAPH")
			}

			//search
			level.Info(o.logger).Log("msg", "🔎 Searching unhealthy positions")

			for i, a := range accounts {

				fmt.Println(" account ", i, " -", a.Id)

				if a.IsLiquidable() {

					fmt.Println(" 🗡️ liquidating account ")
					tx, err := liquidateBorrow(a.Id, ctoken, ctx)
					if err != nil {
						level.Error(o.logger).Log("msg", "❌ Error calling liquidateBorrow method")
						level.Error(o.logger).Log("msg", err)
					} else {
						fmt.Println("✅ Account liquidated :", tx.Hash().Hex())
					}

				}
			}

			break
		case <-ctx.Done():
			return
		}
	}
}

func liquidateBorrow(borrowerAddrStr string, ctoken *contracts.CToken, ctx context.Context) (*types.Transaction, error) {

	borrowerAddress := common.HexToAddress(borrowerAddrStr)
	cTokenCollateral := common.HexToAddress("0x0")
	amount := new(big.Int)

	gasPrice := new(big.Int).SetInt64(240736218990)

	txOps := &bind.TransactOpts{
		Context:  ctx,
		GasPrice: gasPrice,
	}

	tx, err := ctoken.LiquidateBorrow(txOps, borrowerAddress, amount, cTokenCollateral)
	if err != nil {
		return nil, err
	}

	return tx, nil

}

func (o *liqbot) getInitialSources() (*contracts.Comptroller, *contracts.CToken, error) {
	cl, err := ethclient.Dial(o.cfg.RPCURL().String())
	if err != nil {
		return nil, nil, errors.New("Setting ethclient: " + err.Error())
	}
	ctoken, err := contracts.NewCToken(o.cfg.ContractCusdcAddress(), cl)
	if err != nil {
		return nil, nil, errors.New("Setting ctoken: " + err.Error())
	}
	comptroller, err := contracts.NewComptroller(o.cfg.ContractComptrollerAddress(), cl)
	if err != nil {
		return nil, nil, errors.New("Setting comptroller: " + err.Error())
	}

	return comptroller, ctoken, nil
}

const (
	updatePriceTimeout = time.Second * 10
)
