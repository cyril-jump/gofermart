package storage

import "github.com/cyril-jump/gofermart/internal/dto"

type DB interface {
	UserRegister(user dto.User) error
	UserLogin(user dto.User) error
	Ping() error
	Close() error
}
