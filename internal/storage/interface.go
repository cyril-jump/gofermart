package storage

import "github.com/cyril-jump/gofermart/internal/dto"

type Logger interface {
	Close()
}

type DB interface {
	UserRegister(user dto.User) error
	Ping() error
	Close() error
}
