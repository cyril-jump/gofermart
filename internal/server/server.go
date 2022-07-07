package server

import (
	"context"
	"fmt"
	"github.com/cyril-jump/gofermart/internal/http/api"
	"github.com/cyril-jump/gofermart/internal/http/middlewares/cookie"
	"github.com/cyril-jump/gofermart/internal/service"
	"github.com/cyril-jump/gofermart/internal/storage"
	"github.com/cyril-jump/gofermart/internal/workerpool/input"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"os"
)

func InitSrv(ctx context.Context, db storage.DB, inWorker input.Worker, usr service.UsrService, ord service.OrdService, acr service.AcrService) *echo.Echo {

	//new Echo instance
	e := echo.New()

	g := e.Group("")

	// Middleware
	ck := cookie.New(ctx)
	g.Use(echomiddleware.Logger())
	g.Use(echomiddleware.Recover())
	g.Use(echomiddleware.Gzip())
	g.Use(echomiddleware.Decompress())
	//g.Use(ck.SessionWithCookies)

	//Handler
	//handler := handlers.New(db, ck, inWorker)

	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}
	swagger.Servers = nil

	server := api.New(db, ck, inWorker, usr, ord, acr)

	validator := middleware.OapiRequestValidatorWithOptions(swagger,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: ck.Authenticator,
			},
			ErrorHandler: ck.ErrorHandler,
			Skipper:      ck.Skipper,
		})
	g.Use(validator)

	api.RegisterHandlers(g, server)

	/*	g.Use(middleware.OapiRequestValidator(swagger))
		api.RegisterHandlers(g, server)*/
	//Restricted group
	/*	r := e.Group("")
		r.Use(ck.SessionWithCookies)

		//Routes
		e.POST("/api/user/register", handler.PostUserRegister)
		e.POST("/api/user/login", handler.PostUserLogin)
		r.POST("/api/user/orders", handler.PostUserOrders)
		r.GET("/api/user/orders", handler.GetUserOrders)
		r.GET("/api/user/balance", handler.GetUserBalance)
		r.POST("/api/user/balance/withdraw", handler.PostUserBalanceWithdraw)
		r.GET("/api/user/balance/withdrawals", handler.GetUserBalanceWithdrawals)*/

	return e
}
