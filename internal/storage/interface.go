package storage

import "github.com/cyril-jump/gofermart/internal/dto"

type DB interface {
	SetUserRegister(user dto.User) error
	GetUserLogin(user dto.User) (string, error)
	SetAccrualOrder(response dto.AccrualResponse, userID string) error
	UpdateAccrualOrder(response dto.AccrualResponse, userID string) error
	GetAccrualOrder(userID string) ([]dto.Order, error)
	GetUserBalance(userID string) (*dto.UserBalance, error)
	Ping() error
	Close() error
}
