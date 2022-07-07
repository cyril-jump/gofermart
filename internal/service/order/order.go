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

	task.UserID = userID
	task.NumOrder = orderNum
	task.IsNew = true

	if err := o.db.SetAccrualOrder(order, userID); err != nil {
		return err
	}

	o.inWorker.Do(task)

	return nil
}

func (o *OrdService) GetAllUserOrders(userID string) ([]dto.Order1, error) {
	var orders []dto.Order1
	var err error

	if orders, err = o.db.GetAccrualOrder(userID); err != nil {
		return nil, err
	}
	return orders, nil

}

func (o *OrdService) CheckBalance(userID string) (dto.UserBalance1, error) {

	var useBalance dto.UserBalance1
	var err error

	if useBalance, err = o.db.GetUserBalance(userID); err != nil {
		return dto.UserBalance1{}, err
	}

	return useBalance, nil
}

func (o *OrdService) SetBalanceWithdraw(withdrawals dto.Withdrawals1, userID string) error {

	if err := o.db.SetBalanceWithdraw(userID, withdrawals); err != nil {
		return err
	}
	return nil
}

func (o *OrdService) CheckBalanceWithdraw(userID string) ([]dto.Withdrawals1, error) {
	var err error
	var withdrawals []dto.Withdrawals1

	if withdrawals, err = o.db.GetBalanceWithdrawals(userID); err != nil {
		return nil, err
	}
	return withdrawals, nil
}
