package accrual

import (
	"context"

	"github.com/cyril-jump/gofermart/internal/config"
	"github.com/cyril-jump/gofermart/internal/dto"
	"github.com/cyril-jump/gofermart/internal/http/client"
	"github.com/cyril-jump/gofermart/internal/storage"
	"github.com/cyril-jump/gofermart/internal/workerpool/input"
)

type AcrService struct {
	db        storage.AccrualDB
	inWorker  input.Worker
	arcClient client.Client
}

func New(db storage.AccrualDB, inWorker input.Worker, arcClient client.Client) *AcrService {

	return &AcrService{
		db:        db,
		inWorker:  inWorker,
		arcClient: arcClient,
	}
}

func (a *AcrService) UpdateOrder(ctx context.Context, event dto.Task) (int, error) {

	var accrualResp dto.AccrualResponse
	var wait int
	var err error

	event.IsNew = false

	wait, accrualResp, err = a.arcClient.GetAccrual(ctx, event)
	if err != nil {
		a.inWorker.Do(event)
		return 0, err
	}

	if accrualResp.OrderStatus == config.PROCESSED || accrualResp.OrderStatus == config.INVALID {
		if err = a.db.UpdateAccrualOrder(accrualResp, event.UserID); err != nil {
			a.inWorker.Do(event)
			return 0, err
		}
	}

	if accrualResp.OrderStatus == config.PROCESSING || accrualResp.OrderStatus == config.REGISTERED {
		if accrualResp.OrderStatus == config.PROCESSING {
			if err = a.db.UpdateAccrualOrder(accrualResp, event.UserID); err != nil {
				a.inWorker.Do(event)
				return 0, err
			}
		}
		a.inWorker.Do(event)
	}

	return wait, nil

}
