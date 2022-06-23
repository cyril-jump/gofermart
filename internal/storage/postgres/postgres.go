package postgres

import (
	"context"
	"github.com/cyril-jump/gofermart/internal/logger"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"sync"
)

type DB struct {
	mu     sync.Mutex
	db     *sqlx.DB
	ctx    context.Context
	logger *logger.Logger
}

func New(ctx context.Context, logger *logger.Logger, psqlConn string) *DB {
	db, err := sqlx.Open("pgx", psqlConn)
	if err != nil {
		logger.Zap.Fatal("Failed connect...", zap.Error(err))
	}

	// check db
	if err = db.Ping(); err != nil {
		logger.Zap.Fatal("Failed ping...", zap.Error(err))
	}

	logger.Zap.Info("Connected to DB!")

	return &DB{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (D *DB) Ping() error {
	return D.db.Ping()
}

func (D *DB) Close() error {
	return D.db.Close()
}
