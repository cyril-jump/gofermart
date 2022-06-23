package handlers

import (
	"github.com/cyril-jump/gofermart/internal/storage"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Server struct {
	db storage.DB
}

func New(db storage.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s Server) PostUserRegister(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (s Server) PostUserLogin(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (s Server) PostUserOrders(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (s Server) GetUserOrders(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (s Server) GetUserBalance(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (s Server) PostUserBalanceWithdraw(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (s Server) GetUserBalanceWithdrawals(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
