package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gitlab.com/q-dev/exchange-rate-oracle/pkg/config"
	"gitlab.com/q-dev/exchange-rate-oracle/pkg/liqbot"
)

func main() {

	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	level.Info(logger).Log("msg", "...initializing liqbot")

	level.Debug(logger).Log(
		"msg", "config values",
		"rpc url", cfg.RPCURL(),
		"account address", cfg.AccountAddress().Hex(),
		"contract address", cfg.ContractComptrollerAddress(),
		"update interval", cfg.UpdateInterval().Seconds(),
	)

	liqbot_ := liqbot.New(logger, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		level.Info(logger).Log("msg", "🛑 stopping liquidation bot")
		cancel()
	}()

	level.Info(logger).Log("msg", " 🚀 starting liquidation bot")

	//start bot
	liqbot_.Start(ctx)

}
