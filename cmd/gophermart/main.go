package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/cyril-jump/gofermart/internal/config"
	"github.com/cyril-jump/gofermart/internal/server"
	"github.com/cyril-jump/gofermart/internal/storage/postgres"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
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
	flag.StringVarP(&config.Flags.AccrualSystemAddress, "accrualAddres", "r", config.EnvVar.AccrualSystemAddress, "accrual system address")
	flag.Parse()

}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	// Init config
	cfg := config.NewConfig(config.Flags.ServerAddress, config.Flags.DatabaseDSN)

	// Init DB
	psqlConn, err := cfg.Get("database_dsn_str")
	if err != nil {
		config.Logger.Fatal("failed to database_dsn_str: ", zap.Error(err))
	}

	db := postgres.New(ctx, psqlConn)

	// Init HTTPServer
	srv := server.InitSrv(ctx, db)

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
