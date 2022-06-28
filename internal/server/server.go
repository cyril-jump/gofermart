package server

import (
	"context"
	"github.com/cyril-jump/gofermart/internal/http/handlers"
	"github.com/cyril-jump/gofermart/internal/http/middlewares/cookie"
	"github.com/cyril-jump/gofermart/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitSrv(ctx context.Context, db storage.DB) *echo.Echo {

	//new Echo instance
	e := echo.New()

	// Middleware
	ck := cookie.New(ctx)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Decompress())

	//Handler
	handler := handlers.New(db, ck)

	//Routes
	e.POST("/api/user/register", handler.PostUserRegister)
	e.POST("/api/user/login", handler.PostUserLogin)
	e.POST("/api/user/orders", handler.PostUserOrders)
	e.GET("/api/user/orders", handler.GetUserOrders)
	e.GET("/api/user/balance", handler.GetUserBalance)
	e.POST("/api/user/balance/withdraw", handler.PostUserBalanceWithdraw)
	e.GET("/api/user/balance/withdrawals", handler.GetUserBalanceWithdrawals)

	return e
}
