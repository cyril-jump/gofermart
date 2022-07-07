package service

import "github.com/cyril-jump/gofermart/internal/dto"

type UsrService interface {
	Register(user dto.NewUser) (string, error)
	Login(user dto.NewUser) (string, error)
}

type AcrService interface {
}

type OrdService interface {
}
