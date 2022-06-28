package cookie

import "github.com/labstack/echo/v4"

type Cooker interface {
	CreateCookie(c echo.Context, login string) error
	CreateToken(login string) (string, error)
	CheckToken(tokenString string) (string, bool, error)
}
