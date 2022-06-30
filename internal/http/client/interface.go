package client

import (
	"context"
	"github.com/cyril-jump/gofermart/internal/dto"
)

type Client interface {
	GetAccrual(ctx context.Context, task dto.Task) (int, error)
}
