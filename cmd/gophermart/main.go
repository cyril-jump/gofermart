package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/cyril-jump/gofermart/internal/config"
	logger "github.com/cyril-jump/gofermart/internal/logger"
	"github.com/cyril-jump/gofermart/internal/server"
	"github.com/cyril-jump/gofermart/internal/storage/postgres"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
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
	l := logger.New()

	// Init config
	//config
	cfg := config.NewConfig(config.Flags.ServerAddress, config.Flags.DatabaseDSN)

	// Init DB
	psqlConn, err := cfg.Get("database_dsn_str")
	if err != nil {
		l.Zap.Error("", zap.Error(err))
	}

	db := postgres.New(ctx, l, psqlConn)

	// Init HTTPServer
	srv := server.InitSrv(db)

	go func() {

		<-signalChan

		l.Zap.Info("Shutting down...")

		cancel()
		if err := srv.Shutdown(ctx); err != nil && err != ctx.Err() {
			srv.Logger.Fatal(err)
		}

		if err = db.Close(); err != nil {
			l.Zap.Fatal("Failed db...", zap.Error(err))
		}

		l.Close()
	}()

	if err := srv.Start(":8080"); err != nil && err != http.ErrServerClosed {
		srv.Logger.Fatal(err)
	}
}
