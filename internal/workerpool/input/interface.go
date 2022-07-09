package input

import "github.com/cyril-jump/gofermart/internal/dto"

type Worker interface {
	Do(task dto.Task)
}
