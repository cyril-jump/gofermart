package accrual

import "github.com/cyril-jump/gofermart/internal/storage"

type AcrService struct {
	db storage.AccrualDB
}

func New(db storage.AccrualDB) *AcrService {

	return &AcrService{
		db: db,
	}
}
