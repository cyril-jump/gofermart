package user

import (
	"github.com/cyril-jump/gofermart/internal/dto"
	"github.com/cyril-jump/gofermart/internal/storage"
	"github.com/google/uuid"
)

type UsrService struct {
	db storage.UserDB
}

func New(db storage.UserDB) *UsrService {

	return &UsrService{
		db: db,
	}
}

func (u *UsrService) Register(user dto.NewUser) (string, error) {
	userID := uuid.New().String()

	if err := u.db.SetUserRegister(user, userID); err != nil {
		return "", err
	}
	return userID, nil
}

func (u *UsrService) Login(user dto.NewUser) (string, error) {
	var userID string
	var err error

	if userID, err = u.db.GetUserLogin(user); err != nil {
		return "", err
	}

	return userID, err

}
