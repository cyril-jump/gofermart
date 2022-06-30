package storage

import "github.com/cyril-jump/gofermart/internal/dto"

type DB interface {
	SetUserRegister(user dto.User) error
	GetUserLogin(user dto.User) (string, error)
	SetAccrualOrder(response dto.AccrualResponse) error
	UpdateAccrualOrder(response dto.AccrualResponse) error
	GetAccrualOrder(userID string) ([]dto.AccrualResponse, error)
	Ping() error
	Close() error
}
