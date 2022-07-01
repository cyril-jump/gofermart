package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/cyril-jump/gofermart/internal/config"
	"github.com/cyril-jump/gofermart/internal/dto"
	"github.com/cyril-jump/gofermart/internal/http/client"
	"github.com/cyril-jump/gofermart/internal/server"
	"github.com/cyril-jump/gofermart/internal/storage/postgres"
	"github.com/cyril-jump/gofermart/internal/workerpool/input"
	"github.com/cyril-jump/gofermart/internal/workerpool/output"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

func init() {

	//evn vars
	err := env.Parse(&config.EnvVar)
	if err != nil {
		config.Logger.Fatal("failed to parse env...", zap.Error(err))
	}

	//flags
	flag.StringVarP(&config.Flags.ServerAddress, "runAddress", "a", config.EnvVar.ServerAddress, "server address")
	flag.StringVarP(&config.Flags.DatabaseDSN, "psqlConn", "d", config.EnvVar.DatabaseDSN, "database URL conn")
	flag.StringVarP(&config.Flags.AccrualSystemAddress, "accrualAddress", "r", config.EnvVar.AccrualSystemAddress, "accrual system address")
	flag.Parse()

}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	// Init config
	cfg := config.NewConfig(config.Flags.ServerAddress, config.Flags.DatabaseDSN, config.Flags.AccrualSystemAddress)

	// Init DB
	psqlConn, err := cfg.Get("database_dsn_str")
	if err != nil {
		config.Logger.Fatal("failed to to get database_dsn_str: ", zap.Error(err))
	}

	db := postgres.New(ctx, psqlConn)

	// Init Workers
	g, _ := errgroup.WithContext(ctx)
	q := make(chan dto.Task, 100)
	rb := make(chan dto.Task, 100)
	mu := &sync.Mutex{}

	inWorker := input.NewWorker(ctx, mu, q, rb)

	// Init HTTP client
	accrualConn, err := cfg.Get("accrual_system_address")
	if err != nil {
		config.Logger.Fatal("failed to get accrual_system_address: ", zap.Error(err))
	}
	accrualClient := client.New(accrualConn, inWorker, db)

	// Init Workers
	for i := 1; i <= runtime.NumCPU(); i++ {
		outWorker := output.NewOutputWorker(ctx, mu, q, rb, accrualClient)
		g.Go(outWorker.Do)
	}

	// Init HTTPServer
	srv := server.InitSrv(ctx, db, inWorker)

	go func() {

		<-signalChan

		config.Logger.Info("Shutting down...")

		cancel()
		if err = srv.Shutdown(ctx); err != nil && err != ctx.Err() {
			srv.Logger.Fatal(err)
		}

		if err = db.Close(); err != nil {
			config.Logger.Fatal("Failed to close db...", zap.Error(err))
		}

		close(q)
		close(rb)
		err = g.Wait()
		if err != nil {
			config.Logger.Warn("err-group...", zap.Error(err))
		}

		config.Logger.Sync()
	}()

	serverAddress, err := cfg.Get("server_address_str")
	if err != nil {
		config.Logger.Fatal("failed to get server_address_str: ", zap.Error(err))
	}
	if err = srv.Start(serverAddress); err != nil && err != http.ErrServerClosed {
		srv.Logger.Fatal(err)
	}
}
