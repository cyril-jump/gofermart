package service

import "github.com/cyril-jump/gofermart/internal/dto"

type UsrService interface {
	Register(user dto.NewUser) (string, error)
	Login(user dto.NewUser) (string, error)
}

type AcrService interface {
}

type OrdService interface {
	SetNewOrder(orderNum, userID string) error
	GetAllUserOrders(userID string) ([]dto.Order1, error)
	CheckBalance(userID string) (dto.UserBalance1, error)
	SetBalanceWithdraw(withdrawals dto.Withdrawals1, userID string) error
	CheckBalanceWithdraw(userID string) ([]dto.Withdrawals1, error)
}
