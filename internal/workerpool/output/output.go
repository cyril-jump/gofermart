package output

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/cyril-jump/gofermart/internal/config"
	"github.com/cyril-jump/gofermart/internal/dto"
	"github.com/cyril-jump/gofermart/internal/service"
)

type Work struct {
	mu       *sync.Mutex
	ctx      context.Context
	queue    chan dto.Task
	ringBuff chan dto.Task
	ticker   *time.Ticker
	acr      service.AcrService
}

func NewOutputWorker(ctx context.Context, mu *sync.Mutex, queue chan dto.Task, ringBuff chan dto.Task, acr service.AcrService) *Work {
	ticker := time.NewTicker(10 * time.Second)
	return &Work{
		mu:       mu,
		ctx:      ctx,
		queue:    queue,
		ringBuff: ringBuff,
		ticker:   ticker,
		acr:      acr,
	}
}

func (w *Work) Do() error {
	for {
		select {
		case <-w.ctx.Done():
			w.ticker.Stop()
			return nil
		case eventNew := <-w.queue:
			wait, err := w.acr.UpdateOrder(w.ctx, eventNew)
			if err != nil {
				config.Logger.Warn("", zap.Error(err))
			}
			if wait != 0 {
				time.Sleep(time.Duration(wait) * time.Second)
			}
		case <-w.ticker.C:
			if len(w.ringBuff) == 0 {
				break
			}
			for oldEvent := range w.ringBuff {
				wait, err := w.acr.UpdateOrder(w.ctx, oldEvent)
				if err != nil {
					config.Logger.Warn("", zap.Error(err))
				}
				if wait != 0 {
					time.Sleep(time.Duration(wait) * time.Second)
				}
				if len(w.ringBuff) == 0 {
					break
				}
			}
		}
	}
}
