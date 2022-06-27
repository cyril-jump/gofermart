package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/cyril-jump/gofermart/internal/config"
	"github.com/cyril-jump/gofermart/internal/logger/zaplog"
	"github.com/cyril-jump/gofermart/internal/server"
	"github.com/cyril-jump/gofermart/internal/storage/postgres"
	flag "github.com/spf13/pflag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	//evn vars
	err := env.Parse(&config.EnvVar)
	if err != nil {
		log.Fatal(err)
	}

	//flags
	flag.StringVarP(&config.Flags.ServerAddress, "address", "a", config.EnvVar.ServerAddress, "server address")
	flag.StringVarP(&config.Flags.DatabaseDSN, "psqlConn", "d", config.EnvVar.DatabaseDSN, "database URL conn")
	flag.Parse()

}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	// Init logger
	logger := zaplog.New()

	// Init config
	//config
	cfg := config.NewConfig(config.Flags.ServerAddress, config.Flags.DatabaseDSN)

	// Init DB
	psqlConn, err := cfg.Get("database_dsn_str")
	if err != nil {
		logger.Fatal("database_dsn_str: ", err)
	}

	db := postgres.New(ctx, logger, psqlConn)

	// Init HTTPServer
	srv := server.InitSrv(db)

	go func() {

		<-signalChan

		logger.Info("Shutting down...")

		cancel()
		if err := srv.Shutdown(ctx); err != nil && err != ctx.Err() {
			srv.Logger.Fatal(err)
		}

		if err = db.Close(); err != nil {
			logger.Fatal("Failed db...", err)
		}

		logger.Close()
	}()

	serverAddress, err := cfg.Get("server_address_str")
	if err != nil {
		logger.Fatal("server_address_str: ", err)
	}
	if err := srv.Start(serverAddress); err != nil && err != http.ErrServerClosed {
		srv.Logger.Fatal(err)
	}
}
