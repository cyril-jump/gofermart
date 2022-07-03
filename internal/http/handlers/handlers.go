package handlers

import (
	"encoding/json"
	"errors"
	"github.com/cyril-jump/gofermart/internal/config"
	"github.com/cyril-jump/gofermart/internal/dto"
	"github.com/cyril-jump/gofermart/internal/http/middlewares/cookie"
	"github.com/cyril-jump/gofermart/internal/storage"
	"github.com/cyril-jump/gofermart/internal/utils"
	"github.com/cyril-jump/gofermart/internal/utils/errs"
	"github.com/cyril-jump/gofermart/internal/workerpool/input"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"strconv"
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
		config.Logger.Warn("PostUserRegister", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	if err = h.ck.CreateCookie(c, user.UserID); err != nil {
		config.Logger.Warn("PostUserRegister", zap.Error(err))
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
		config.Logger.Warn("PostUserLogin", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	config.Logger.Info(userID)

	if err = h.ck.CreateCookie(c, userID); err != nil {
		config.Logger.Warn("PostUserLogin", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) PostUserOrders(c echo.Context) error {
	var order dto.AccrualResponse
	var task dto.Task
	var userID string

	if id := c.Request().Context().Value(config.TokenKey); id != nil {
		userID = id.(string)
	}
	log.Println(order)
	if c.Request().Header.Get("Content-Type") != "text/plain" {
		config.Logger.Info(c.Response().Header().Get("Content-Type"))
		return c.NoContent(http.StatusBadRequest)
	}
	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	num, err := strconv.Atoi(string(body))
	if err != nil {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	if ok := utils.ValidOrder(num); !ok {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	order.NumOrder = string(body)
	order.OrderStatus = config.NEW
	order.Accrual = 0.0

	task.UserID = userID
	task.NumOrder = string(body)
	task.IsNew = true

	if err = h.db.SetAccrualOrder(order, userID); err != nil {
		if errors.Is(err, errs.ErrAlreadyUploadThisUser) {

			return c.NoContent(http.StatusOK)
		}
		if errors.Is(err, errs.ErrAlreadyUploadOtherUser) {
			return c.NoContent(http.StatusConflict)
		}
		config.Logger.Warn("PostUserOrders", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	h.inWorker.Do(task)
	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) GetUserOrders(c echo.Context) error {

	var orders []dto.Order
	var err error
	var userID string

	if id := c.Request().Context().Value(config.TokenKey); id != nil {
		userID = id.(string)
	}

	if orders, err = h.db.GetAccrualOrder(userID); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.NoContent(http.StatusNoContent)
		}
		config.Logger.Warn("GetUserOrders", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}
	log.Println(orders)
	return c.JSON(http.StatusOK, orders)
}

func (h *Handler) GetUserBalance(c echo.Context) error {

	var useBalance *dto.UserBalance
	var err error
	var userID string

	if id := c.Request().Context().Value(config.TokenKey); id != nil {
		userID = id.(string)
	}

	config.Logger.Info(userID)

	if useBalance, err = h.db.GetUserBalance(userID); err != nil {
		config.Logger.Warn("GetUserOrders", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, &useBalance)
}

func (h *Handler) PostUserBalanceWithdraw(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetUserBalanceWithdrawals(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
