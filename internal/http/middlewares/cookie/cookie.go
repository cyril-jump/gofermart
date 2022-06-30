package cookie

import (
	"context"
	"fmt"
	"github.com/cyril-jump/gofermart/internal/config"
	"github.com/cyril-jump/gofermart/internal/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Cookie struct {
	randNum []byte
	ctx     context.Context
}

func New(ctx context.Context) *Cookie {

	key, err := utils.GenerateRandom(16)
	if err != nil {
		config.Logger.Fatal("generate random...", zap.Error(err))
	}

	return &Cookie{
		randNum: key,
		ctx:     ctx,
	}
}

func (ck *Cookie) CreateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user": userID})
	tokenString, _ := token.SignedString(ck.randNum)
	return tokenString, nil
}

func (ck *Cookie) CheckToken(tokenString string) (string, bool, error) {

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexected signing method: %v", token.Header["alg"])
		}
		return ck.randNum, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return fmt.Sprintf("%s", claims["user"]), ok, nil
	}
	return "", false, nil
}

func (ck *Cookie) CreateCookie(c echo.Context, userID string) error {
	var err error
	cookie := new(http.Cookie)
	cookie.Path = "/"
	cookie.Value, err = ck.CreateToken(userID)
	if err != nil {
		return err
	}
	cookie.Name = config.TokenKey.String()
	c.SetCookie(cookie)
	c.Request().AddCookie(cookie)
	return nil
}
