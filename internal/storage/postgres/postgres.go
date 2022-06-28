package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/cyril-jump/gofermart/internal/config"
	"github.com/cyril-jump/gofermart/internal/dto"
	"github.com/cyril-jump/gofermart/internal/utils"
	"github.com/cyril-jump/gofermart/internal/utils/errs"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"sync"
)

type DB struct {
	mu  sync.Mutex
	db  *sqlx.DB
	ctx context.Context
}

func New(ctx context.Context, psqlConn string) *DB {
	db, err := sqlx.Open("pgx", psqlConn)
	if err != nil {
		config.Logger.Fatal("Failed connect...", zap.Error(err))
	}

	// check db
	if err = db.Ping(); err != nil {
		config.Logger.Fatal("Failed ping...", zap.Error(err))
	}

	if _, err = db.Exec(schema); err != nil {
		config.Logger.Fatal("", zap.Error(err))
	}

	config.Logger.Info("Connected to DB!")

	return &DB{
		db:  db,
		ctx: ctx,
	}
}

func (db *DB) UserRegister(user dto.User) error {
	db.mu.Lock()
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	insertStmt, err := db.db.PrepareContext(db.ctx, "INSERT INTO users (login, password) VALUES ($1, $2)")
	if err != nil {
		return err
	}

	defer func() {
		insertStmt.Close()
		db.mu.Unlock()
	}()

	_, err = insertStmt.ExecContext(db.ctx, user.Login, hash)
	if err != nil {
		if pgerrcode.IsIntegrityConstraintViolation(err.(*pgconn.PgError).Code) {
			return errs.ErrAlreadyExists
		}
		return err
	}
	return nil

}

func (db *DB) UserLogin(user dto.User) error {
	db.mu.Lock()
	var usr dto.User
	selectStmt, err := db.db.PrepareContext(db.ctx, "SELECT login,password FROM users WHERE login=$1")
	if err != nil {
		return err
	}
	defer func() {
		selectStmt.Close()
		db.mu.Unlock()
	}()

	if err = selectStmt.QueryRowContext(db.ctx, user.Login).Scan(&usr.Login, &usr.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrBadLoginOrPass
		}
		return err
	}
	ok := utils.CheckPasswordHash(user.Password, usr.Password)
	if !ok {
		return errs.ErrBadLoginOrPass
	}
	return nil
}

func (db *DB) Ping() error {
	return db.db.Ping()
}

func (db *DB) Close() error {
	return db.db.Close()
}

var schema = `
	CREATE TABLE IF NOT EXISTS users (
		id serial primary key,
		login text not null unique,
		password text not null,
        "current" float not null default 0,
        withdrawn int not null  default 0
	);
`
