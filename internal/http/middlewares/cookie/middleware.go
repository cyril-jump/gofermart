package cookie

import (
	"context"
	"github.com/cyril-jump/gofermart/internal/config"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (ck *Cookie) SessionWithCookies(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var login string
		var ok bool

		cookie, err := c.Cookie(config.TokenKey.String())
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		} else {
			login, ok, err = ck.CheckToken(cookie.Value)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError)
			} else {
				if !ok {
					return echo.NewHTTPError(http.StatusUnauthorized)
				}
			}
		}

		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), config.TokenKey, login)))

		return next(c)
	}
}
