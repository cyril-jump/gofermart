package postgres

import (
	"context"
	"github.com/cyril-jump/gofermart/internal/dto"
	"github.com/cyril-jump/gofermart/internal/logger"
	"github.com/cyril-jump/gofermart/internal/utils/errs"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"sync"
)

type DB struct {
	mu     sync.Mutex
	db     *sqlx.DB
	ctx    context.Context
	logger logger.Logger
}

func New(ctx context.Context, logger logger.Logger, psqlConn string) *DB {
	db, err := sqlx.Open("pgx", psqlConn)
	if err != nil {
		logger.Fatal("Failed connect...", err)
	}

	// check db
	if err = db.Ping(); err != nil {
		logger.Fatal("Failed ping...", err)
	}

	if _, err = db.Exec(schema); err != nil {
		logger.Fatal("", err)
	}

	logger.Info("Connected to DB!")

	return &DB{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (D *DB) UserRegister(user dto.User) error {
	D.mu.Lock()

	insertStmt, err := D.db.PrepareContext(D.ctx, "INSERT INTO users (login, password) VALUES ($1, $2)")
	if err != nil {
		return err
	}

	defer func() {
		insertStmt.Close()
		D.mu.Unlock()
	}()

	_, err = insertStmt.ExecContext(D.ctx, user.Login, user.Password)
	if err != nil {
		if pgerrcode.IsIntegrityConstraintViolation(err.(*pgconn.PgError).Code) {
			return errs.ErrAlreadyExists
		}
		return err
	}
	return nil

}

func (D *DB) Ping() error {
	return D.db.Ping()
}

func (D *DB) Close() error {
	return D.db.Close()
}

var schema = `
	CREATE TABLE IF NOT EXISTS users (
		id serial primary key,
		login text not null unique,
		password text not null
	);
	CREATE TABLE IF NOT EXISTS user_balance(
	    id  int not null references users(id),
	    user_id text not null,
	    status text not null,
	    current float not null default 0,
	    withdrawn float not null default 0
	);
`
