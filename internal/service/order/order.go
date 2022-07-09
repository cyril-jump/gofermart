package order

import (
	"github.com/cyril-jump/gofermart/internal/config"
	"github.com/cyril-jump/gofermart/internal/dto"
	"github.com/cyril-jump/gofermart/internal/storage"
	"github.com/cyril-jump/gofermart/internal/workerpool/input"
)

type OrdService struct {
	db       storage.OrderDB
	inWorker input.Worker
}

func New(db storage.OrderDB, in input.Worker) *OrdService {
	return &OrdService{
		db:       db,
		inWorker: in,
	}
}

func (o *OrdService) SetNewOrder(orderNum, userID string) error {

	var order dto.AccrualResponse
	var task dto.Task

	order.OrderStatus = config.NEW
	order.Accrual = 0.0
	order.NumOrder = orderNum

	task.UserID = userID
	task.NumOrder = orderNum
	task.IsNew = true

	if err := o.db.SetAccrualOrder(order, userID); err != nil {
		return err
	}

	o.inWorker.Do(task)

	return nil
}

func (o *OrdService) GetAllUserOrders(userID string) ([]dto.Order, error) {
	var orders []dto.Order
	var err error

	if orders, err = o.db.GetAccrualOrder(userID); err != nil {
		return nil, err
	}
	return orders, nil

}

func (o *OrdService) CheckBalance(userID string) (dto.UserBalance, error) {

	var useBalance dto.UserBalance
	var err error

	if useBalance, err = o.db.GetUserBalance(userID); err != nil {
		return dto.UserBalance{}, err
	}

	return useBalance, nil
}

func (o *OrdService) SetBalanceWithdraw(withdrawals dto.Withdrawals, userID string) error {

	if err := o.db.SetBalanceWithdraw(userID, withdrawals); err != nil {
		return err
	}
	return nil
}

func (o *OrdService) CheckBalanceWithdraw(userID string) ([]dto.Withdrawals, error) {
	var err error
	var withdrawals []dto.Withdrawals

	if withdrawals, err = o.db.GetBalanceWithdrawals(userID); err != nil {
		return nil, err
	}
	return withdrawals, nil
}
