package storage

import "github.com/cyril-jump/gofermart/internal/dto"

type DB interface {
	SetUserRegister(user dto.NewUser, id string) error
	GetUserLogin(user dto.NewUser) (string, error)
	SetAccrualOrder(response dto.AccrualResponse, userID string) error
	UpdateAccrualOrder(response dto.AccrualResponse, userID string) error
	GetAccrualOrder(userID string) ([]dto.Order1, error)
	GetUserBalance(userID string) (dto.UserBalance1, error)
	SetBalanceWithdraw(userID string, withdraw dto.Withdrawals1) error
	GetBalanceWithdrawals(userID string) ([]dto.Withdrawals1, error)
	Ping() error
	Close() error
}

type UserDB interface {
	SetUserRegister(user dto.NewUser, userID string) error
	GetUserLogin(user dto.NewUser) (string, error)
}

type AccrualDB interface {
	UpdateAccrualOrder(response dto.AccrualResponse, userID string) error
}

type OrderDB interface {
	SetAccrualOrder(response dto.AccrualResponse, userID string) error
	GetAccrualOrder(userID string) ([]dto.Order1, error)
	GetUserBalance(userID string) (dto.UserBalance1, error)
	SetBalanceWithdraw(userID string, withdraw dto.Withdrawals1) error
	GetBalanceWithdrawals(userID string) ([]dto.Withdrawals1, error)
}
