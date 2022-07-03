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
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"sync"
	"time"
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

func (db *DB) SetUserRegister(user dto.User) error {
	db.mu.Lock()
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	insertStmt, err := db.db.PrepareContext(db.ctx, "INSERT INTO users (id, login, password) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}

	defer func() {
		insertStmt.Close()
		db.mu.Unlock()
	}()

	_, err = insertStmt.ExecContext(db.ctx, user.UserID, user.Login, hash)
	if err != nil {
		if pgerrcode.IsIntegrityConstraintViolation(err.(*pgconn.PgError).Code) {
			return errs.ErrAlreadyExists
		}
		return err
	}
	return nil

}

func (db *DB) GetUserLogin(user dto.User) (string, error) {
	db.mu.Lock()
	var usr dto.User
	selectStmt, err := db.db.PrepareContext(db.ctx, "SELECT id, login,password FROM users WHERE login=$1")
	if err != nil {
		return "", err
	}
	defer func() {
		selectStmt.Close()
		db.mu.Unlock()
	}()

	if err = selectStmt.QueryRowContext(db.ctx, user.Login).Scan(&usr.UserID, &usr.Login, &usr.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrBadLoginOrPass
		}
		return "", err
	}
	ok := utils.CheckPasswordHash(user.Password, usr.Password)
	if !ok {
		return "", errs.ErrBadLoginOrPass
	}
	return usr.UserID, nil
}

func (db *DB) SetAccrualOrder(resp dto.AccrualResponse, usrID string) error {
	db.mu.Lock()
	log.Print("SetAccrualOrder   ", resp)
	var userID string
	insertStmt, err := db.db.PrepareContext(db.ctx, "INSERT INTO orders (user_id, number, status, accrual, uploaded_at) VALUES ($1, $2, $3, $4, $5) RETURNING (user_id)")
	if err != nil {
		return err
	}

	selectStmt, err := db.db.PrepareContext(db.ctx, "SELECT user_id FROM orders WHERE number=$1")
	if err != nil {
		return err
	}
	defer func() {
		insertStmt.Close()
		selectStmt.Close()
		db.mu.Unlock()
	}()
	uploadedAt := time.Now().Format(time.RFC3339)
	_, err = insertStmt.ExecContext(db.ctx, usrID, resp.NumOrder, resp.OrderStatus, resp.Accrual, uploadedAt)
	if err != nil {
		if pgerrcode.IsIntegrityConstraintViolation(err.(*pgconn.PgError).Code) {
			if err = selectStmt.QueryRowContext(db.ctx, resp.NumOrder).Scan(&userID); err != nil {
				return err
			}
			if userID == usrID {
				config.Logger.Warn(userID, zap.Error(err))
				return errs.ErrAlreadyUploadThisUser
			}
			config.Logger.Warn("", zap.Error(err))
			return errs.ErrAlreadyUploadOtherUser
		}
		return err
	}
	return nil
}

func (db *DB) UpdateAccrualOrder(resp dto.AccrualResponse, userID string) error {
	db.mu.Lock()
	log.Print("UpdateAccrualOrder   ", resp)
	updateStmt1, err := db.db.PrepareContext(db.ctx, "UPDATE orders SET status = $1, accrual = $2 WHERE number = $3")
	if err != nil {
		return err
	}

	updateStmt2, err := db.db.PrepareContext(db.ctx, "UPDATE users SET current = current + $1 WHERE id = $2")
	if err != nil {
		return err
	}

	tx, err := db.db.BeginTx(db.ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		updateStmt1.Close()
		updateStmt2.Close()
		tx.Rollback()
		db.mu.Unlock()
	}()

	_, err = tx.StmtContext(db.ctx, updateStmt1).ExecContext(db.ctx, resp.OrderStatus, resp.Accrual, resp.NumOrder)
	if err != nil {
		return err
	}

	_, err = tx.StmtContext(db.ctx, updateStmt2).ExecContext(db.ctx, userID, resp.Accrual)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetAccrualOrder(userID string) ([]dto.Order, error) {
	db.mu.Lock()
	log.Print("GetAccrual   ", userID)
	orders := make([]dto.Order, 0, 100)
	var order dto.Order
	selectStmt, err := db.db.PrepareContext(db.ctx, "SELECT number, status, accrual, uploaded_at  FROM orders WHERE user_id=$1 ORDER BY uploaded_at DESC")
	if err != nil {
		return nil, err
	}
	rows, err := selectStmt.QueryContext(db.ctx, userID)
	if err != nil {
		return nil, err
	}

	defer func() {
		selectStmt.Close()
		rows.Close()
		db.mu.Unlock()
	}()

	if err = rows.Err(); err != nil {
		return nil, err
	}
	for rows.Next() {
		if err = rows.Scan(&order.NumOrder, &order.OrderStatus, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return nil, errs.ErrNotFound
	}
	log.Print("GetAccrual   ", orders)
	return orders, nil
}

func (db *DB) GetUserBalance(userID string) (*dto.UserBalance, error) {
	db.mu.Lock()
	log.Print("GetUserBalance   ", userID)
	var usrBalance dto.UserBalance

	selectStmt, err := db.db.PrepareContext(db.ctx, "SELECT current ,withdrawn FROM users WHERE id=$1")
	if err != nil {
		return nil, err
	}

	defer func() {
		selectStmt.Close()
		db.mu.Unlock()
	}()

	err = selectStmt.QueryRowContext(db.ctx, userID).Scan(&usrBalance.Current, &usrBalance.Withdrawn)
	if err != nil {
		return nil, err
	}

	return &usrBalance, nil

}

func (db *DB) Ping() error {
	return db.db.Ping()
}

func (db *DB) Close() error {
	return db.db.Close()
}

var schema = `
	CREATE TABLE IF NOT EXISTS users (
		id text primary key,
		login text not null unique,
		password text not null,
        "current" float not null default 0,
        withdrawn int not null  default 0
	);
	CREATE TABLE IF NOT EXISTS orders (
	  	"number" text primary key unique,
	  	user_id text not null references users(id),
	    status text not null,
	    accrual float not null,
	    uploaded_at timestamp
	);
`
