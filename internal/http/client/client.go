package client

import (
	"context"
	"encoding/json"
	"github.com/cyril-jump/gofermart/internal/config"
	"github.com/cyril-jump/gofermart/internal/dto"
	"github.com/cyril-jump/gofermart/internal/storage"
	"github.com/cyril-jump/gofermart/internal/workerpool/input"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
)

type Accrual struct {
	client         *resty.Client
	accrualAddress string
	inWorker       input.Worker
	db             storage.DB
}

func New(accrualAddress string, inWorker input.Worker, db storage.DB) *Accrual {
	accrualClient := resty.New()
	return &Accrual{
		accrualAddress: accrualAddress,
		client:         accrualClient,
		inWorker:       inWorker,
		db:             db,
	}
}

func (c *Accrual) GetAccrual(ctx context.Context, task dto.Task) (int, error) {
	var accrualResp dto.AccrualResponse
	task.IsNew = false
	resp, err := c.client.R().SetContext(ctx).SetPathParams(map[string]string{"orderNumber": task.NumOrder}).Get(c.accrualAddress + "/api/orders/{orderNumber}")
	if err != nil {
		c.inWorker.Do(task)
		return 0, err
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		err = json.Unmarshal(resp.Body(), &accrualResp)
		if err != nil {
			c.inWorker.Do(task)
			return 0, err
		}
		log.Print("GetAccrual   ", accrualResp)
		if accrualResp.OrderStatus == config.PROCESSED || accrualResp.OrderStatus == config.INVALID {
			if err = c.db.UpdateAccrualOrder(accrualResp, task.UserID); err != nil {
				c.inWorker.Do(task)
				return 0, err
			}
		}

		if accrualResp.OrderStatus == config.PROCESSING || accrualResp.OrderStatus == config.REGISTERED {
			if accrualResp.OrderStatus == config.PROCESSING {
				if err = c.db.UpdateAccrualOrder(accrualResp, task.UserID); err != nil {
					c.inWorker.Do(task)
					return 0, err
				}
			}
			c.inWorker.Do(task)
		}

	case http.StatusTooManyRequests:
		wait, _ := strconv.Atoi(resp.Header().Get("Retry-After"))
		return wait, nil
	}
	return 0, nil
}
