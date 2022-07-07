package order

import "github.com/cyril-jump/gofermart/internal/storage"

type OrdService struct {
	db storage.OrderDB
}

func New(db storage.OrderDB) *OrdService {
	return &OrdService{
		db: db,
	}
}
