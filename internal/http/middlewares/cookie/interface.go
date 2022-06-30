package cookie

import "github.com/labstack/echo/v4"

type Cooker interface {
	CreateCookie(c echo.Context, userID string) error
	CreateToken(userID string) (string, error)
	CheckToken(tokenString string) (string, bool, error)
}
