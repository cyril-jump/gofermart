package handlers

import (
	"encoding/json"
	"errors"
	"github.com/cyril-jump/gofermart/internal/dto"
	"github.com/cyril-jump/gofermart/internal/http/middlewares/cookie"
	"github.com/cyril-jump/gofermart/internal/storage"
	"github.com/cyril-jump/gofermart/internal/utils/errs"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

type Handler struct {
	db storage.DB
	ck cookie.Cooker
}

func New(db storage.DB, ck cookie.Cooker) *Handler {
	return &Handler{
		db: db,
		ck: ck,
	}
}

func (h *Handler) PostUserRegister(c echo.Context) error {

	var user dto.User

	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if user.Login == "" || user.Password == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	if err = h.db.UserRegister(user); err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			return c.NoContent(http.StatusConflict)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) PostUserLogin(c echo.Context) error {

	var user dto.User

	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if user.Login == "" || user.Password == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	if err = h.db.UserLogin(user); err != nil {
		if errors.Is(err, errs.ErrBadLoginOrPass) {
			return c.NoContent(http.StatusUnauthorized)
		}
		return c.NoContent(http.StatusInternalServerError)
	}
	if err = h.ck.CreateCookie(c, user.Login); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) PostUserOrders(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetUserOrders(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetUserBalance(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (h *Handler) PostUserBalanceWithdraw(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetUserBalanceWithdrawals(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
