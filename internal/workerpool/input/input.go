package input

import (
	"context"
	"github.com/cyril-jump/gofermart/internal/dto"
	"sync"
)

type Work struct {
	mu       *sync.Mutex
	ctx      context.Context
	queue    chan dto.Task
	ringBuff chan dto.Task
}

func NewWorker(ctx context.Context, mu *sync.Mutex, queue chan dto.Task, ringBuff chan dto.Task) *Work {
	return &Work{
		mu:       mu,
		ctx:      ctx,
		queue:    queue,
		ringBuff: ringBuff,
	}
}

func (q *Work) Do(task dto.Task) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if task.IsNew {
		q.queue <- task
	}
	q.ringBuff <- task
}
