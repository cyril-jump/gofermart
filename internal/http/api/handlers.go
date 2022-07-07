package api

import (
	"encoding/json"
	"errors"
	"github.com/cyril-jump/gofermart/internal/config"
	"github.com/cyril-jump/gofermart/internal/dto"
	"github.com/cyril-jump/gofermart/internal/http/middlewares/cookie"
	"github.com/cyril-jump/gofermart/internal/service"
	"github.com/cyril-jump/gofermart/internal/storage"
	"github.com/cyril-jump/gofermart/internal/utils"
	"github.com/cyril-jump/gofermart/internal/utils/errs"
	"github.com/cyril-jump/gofermart/internal/workerpool/input"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

type Handler struct {
	usr      service.UsrService
	ord      service.OrdService
	acr      service.AcrService
	db       storage.DB
	ck       cookie.Cooker
	inWorker input.Worker
}

func New(db storage.DB, ck cookie.Cooker, inWorker input.Worker, usr service.UsrService, ord service.OrdService, acr service.AcrService) *Handler {
	return &Handler{
		db:       db,
		ck:       ck,
		inWorker: inWorker,
		usr:      usr,
		ord:      ord,
		acr:      acr,
	}
}

func (h Handler) PostUserRegister(c echo.Context) error {

	var user dto.NewUser
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

	if userID, err = h.usr.Register(user); err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			return c.NoContent(http.StatusConflict)
		}
		config.Logger.Warn("PostUserRegister", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	if err = h.ck.CreateCookie(c, userID); err != nil {
		config.Logger.Warn("PostUserRegister", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)

}

func (h Handler) PostUserLogin(c echo.Context) error {

	var user dto.NewUser
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

	if userID, err = h.usr.Login(user); err != nil {
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

	var userID string
	var orderNum string
	var err error

	if id := c.Request().Context().Value(config.TokenKey); id != nil {
		userID = id.(string)
	}

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
	orderNum = string(body)

	if err = h.ord.SetNewOrder(orderNum, userID); err != nil {
		if errors.Is(err, errs.ErrAlreadyUploadThisUser) {

			return c.NoContent(http.StatusOK)
		}
		if errors.Is(err, errs.ErrAlreadyUploadOtherUser) {
			return c.NoContent(http.StatusConflict)
		}
		config.Logger.Warn("PostUserOrders", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) GetUserOrders(c echo.Context) error {

	var orders []dto.Order1
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
	return c.JSON(http.StatusOK, orders)
}

func (h *Handler) GetUserBalance(c echo.Context) error {

	var useBalance dto.UserBalance1
	var err error
	var userID string

	if id := c.Request().Context().Value(config.TokenKey); id != nil {
		userID = id.(string)
	}

	config.Logger.Info(userID)

	if useBalance, err = h.ord.CheckBalance(userID); err != nil {
		config.Logger.Warn("GetUserOrders", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, useBalance)
}

func (h *Handler) PostUserBalanceWithdraw(c echo.Context) error {

	var withdrawals dto.Withdrawals1
	var userID string

	if id := c.Request().Context().Value(config.TokenKey); id != nil {
		userID = id.(string)
	}
	config.Logger.Info(userID)
	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	err = json.Unmarshal(body, &withdrawals)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err = h.ord.SetBalanceWithdraw(withdrawals, userID); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.NoContent(http.StatusUnprocessableEntity)
		}
		if errors.Is(err, errs.ErrInsufficientFunds) {
			return c.NoContent(http.StatusPaymentRequired)
		}
		config.Logger.Warn("PostUserBalanceWithdraw", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetUserBalanceWithdrawals(c echo.Context) error {
	var withdrawals []dto.Withdrawals1
	var err error
	var userID string

	if id := c.Request().Context().Value(config.TokenKey); id != nil {
		userID = id.(string)
	}
	config.Logger.Info(userID)

	if withdrawals, err = h.ord.CheckBalanceWithdraw(userID); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.NoContent(http.StatusNoContent)
		}
		config.Logger.Warn("GetUserOrders", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, withdrawals)
}
