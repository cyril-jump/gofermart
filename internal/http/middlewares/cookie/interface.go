package cookie

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
)

type Cooker interface {
	CreateCookie(c echo.Context, userID string) error
	CreateToken(userID string) (string, error)
	CheckToken(tokenString string) (string, bool, error)
	Authenticator(ctx context.Context, input *openapi3filter.AuthenticationInput) error
	Skipper(c echo.Context) bool
}
