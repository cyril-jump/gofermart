package server

import (
	"github.com/cyril-jump/gofermart/internal/handlers"
	"github.com/cyril-jump/gofermart/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitSrv(db storage.DB) *echo.Echo {
	//server
	srv := handlers.New(db)

	//new Echo instance
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Decompress())
	e.POST("/api/user/register", srv.PostUserRegister)
	e.POST("/api/user/login", srv.PostUserLogin)
	e.POST("/api/user/orders", srv.PostUserOrders)
	e.GET("/api/user/orders", srv.GetUserOrders)
	e.GET("/api/user/balance", srv.GetUserBalance)
	e.POST("/api/user/balance/withdraw", srv.PostUserBalanceWithdraw)
	e.GET("/api/user/balance/withdrawals", srv.GetUserBalanceWithdrawals)

	return e
}
