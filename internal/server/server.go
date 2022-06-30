package server

import (
	"context"
	"github.com/cyril-jump/gofermart/internal/http/handlers"
	"github.com/cyril-jump/gofermart/internal/http/middlewares/cookie"
	"github.com/cyril-jump/gofermart/internal/storage"
	"github.com/cyril-jump/gofermart/internal/workerpool/input"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitSrv(ctx context.Context, db storage.DB, inWorker input.Worker) *echo.Echo {

	//new Echo instance
	e := echo.New()

	// Middleware
	ck := cookie.New(ctx)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Decompress())

	//Handler
	handler := handlers.New(db, ck, inWorker)

	//Restricted group
	r := e.Group("")
	r.Use(ck.SessionWithCookies)

	//Routes
	e.POST("/api/user/register", handler.PostUserRegister)
	e.POST("/api/user/login", handler.PostUserLogin)
	r.POST("/api/user/orders", handler.PostUserOrders)
	r.GET("/api/user/orders", handler.GetUserOrders)
	e.GET("/api/user/balance", handler.GetUserBalance)
	r.POST("/api/user/balance/withdraw", handler.PostUserBalanceWithdraw)
	r.GET("/api/user/balance/withdrawals", handler.GetUserBalanceWithdrawals)

	return e
}
