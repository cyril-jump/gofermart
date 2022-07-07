package client

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/gommon/log"

	"github.com/cyril-jump/gofermart/internal/dto"
	"github.com/cyril-jump/gofermart/internal/storage"
	"github.com/cyril-jump/gofermart/internal/workerpool/input"
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

func (c *Accrual) GetAccrual(ctx context.Context, task dto.Task) (int, dto.AccrualResponse, error) {
	var accrualResp dto.AccrualResponse

	resp, err := c.client.R().SetContext(ctx).SetPathParams(map[string]string{"orderNumber": task.NumOrder}).Get(c.accrualAddress + "/api/orders/{orderNumber}")
	if err != nil {
		return 0, dto.AccrualResponse{}, err
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		err = json.Unmarshal(resp.Body(), &accrualResp)
		if err != nil {
			return 0, dto.AccrualResponse{}, err
		}
		log.Info("GetAccrual   ", accrualResp)

		return 0, accrualResp, nil

	case http.StatusTooManyRequests:
		wait, _ := strconv.Atoi(resp.Header().Get("Retry-After"))
		return wait, dto.AccrualResponse{}, nil
	}
	return 0, dto.AccrualResponse{}, nil
}
