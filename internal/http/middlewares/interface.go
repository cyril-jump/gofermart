package middlewares

import "github.com/cyril-jump/gofermart/internal/http/middlewares/cookie"

type Middleware interface {
	cookie.Cooker
}
