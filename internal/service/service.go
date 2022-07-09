package service

import (
	"context"
	"github.com/cyril-jump/gofermart/internal/dto"
)

type UsrService interface {
	Register(user dto.NewUser) (string, error)
	Login(user dto.NewUser) (string, error)
}

type AcrService interface {
	UpdateOrder(ctx context.Context, event dto.Task) (int, error)
}

type OrdService interface {
	SetNewOrder(orderNum, userID string) error
	GetAllUserOrders(userID string) ([]dto.Order, error)
	CheckBalance(userID string) (dto.UserBalance, error)
	SetBalanceWithdraw(withdrawals dto.Withdrawals, userID string) error
	CheckBalanceWithdraw(userID string) ([]dto.Withdrawals, error)
}
