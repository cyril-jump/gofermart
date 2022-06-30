package handlers

import (
	"encoding/json"
	"errors"
	"github.com/cyril-jump/gofermart/internal/config"
	"github.com/cyril-jump/gofermart/internal/dto"
	"github.com/cyril-jump/gofermart/internal/http/middlewares/cookie"
	"github.com/cyril-jump/gofermart/internal/storage"
	"github.com/cyril-jump/gofermart/internal/utils/errs"
	"github.com/cyril-jump/gofermart/internal/workerpool/input"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Handler struct {
	db       storage.DB
	ck       cookie.Cooker
	inWorker input.Worker
}

func New(db storage.DB, ck cookie.Cooker, inWorker input.Worker) *Handler {
	return &Handler{
		db:       db,
		ck:       ck,
		inWorker: inWorker,
	}
}

func (h *Handler) PostUserRegister(c echo.Context) error {

	var user dto.User
	user.UserID = uuid.New().String()
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

	if err = h.db.SetUserRegister(user); err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			return c.NoContent(http.StatusConflict)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	if err = h.ck.CreateCookie(c, user.UserID); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) PostUserLogin(c echo.Context) error {

	var user dto.User
	var userID string

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

	if userID, err = h.db.GetUserLogin(user); err != nil {
		if errors.Is(err, errs.ErrBadLoginOrPass) {
			return c.NoContent(http.StatusUnauthorized)
		}
		return c.NoContent(http.StatusInternalServerError)
	}
	if err = h.ck.CreateCookie(c, userID); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) PostUserOrders(c echo.Context) error {
	var order dto.AccrualResponse
	var task dto.Task

	if id := c.Request().Context().Value(config.TokenKey); id != nil {
		order.UserID = id.(string)
	}

	if c.Request().Header.Get("Content-Type") != "text/plain" {
		config.Logger.Info(c.Response().Header().Get("Content-Type"))
		return c.NoContent(http.StatusBadRequest)
	}
	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	order.NumOrder = string(body)
	order.OrderStatus = config.REGISTERED

	task.NumOrder = string(body)
	task.IsNew = true

	if err = h.db.SetAccrualOrder(order); err != nil {
		if errors.Is(err, errs.ErrAlreadyUploadThisUser) {
			return c.NoContent(http.StatusOK)
		} else if errors.Is(err, errs.ErrAlreadyUploadOtherUser) {
			return c.NoContent(http.StatusConflict)
		} else {
			config.Logger.Warn("", zap.Error(err))
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	h.inWorker.Do(task)
	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) GetUserOrders(c echo.Context) error {

	orders := make([]dto.AccrualResponse, 0, 100)
	var err error
	var userID string

	if id := c.Request().Context().Value(config.TokenKey); id != nil {
		userID = id.(string)
	}

	if orders, err = h.db.GetAccrualOrder(userID); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.NoContent(http.StatusNoContent)
		}
		config.Logger.Warn("", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, orders)
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
