package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"time"
)

func CreateToken(customerID, name string, viperCfg *viper.Viper) (token string, err error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"customer_id": customerID,
		"name":        name,
		"exp":         now.Add(time.Hour * 168).Unix(), // 7 days
		"iat":         now.Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = jwtToken.SignedString([]byte(viperCfg.GetString("JWT_SECRET_KEY")))
	if err != nil {
		return
	}

	return
}

func GetTokenClaims(ctx echo.Context) (claims jwt.MapClaims, err error) {
	token, ok := ctx.Get("user").(*jwt.Token)
	if !ok {
		err = errors.New("missing or invalid jwt token")
		return
	}

	claims, ok = token.Claims.(jwt.MapClaims)
	if !ok || token == nil {
		err = errors.New("invalid token")
		return
	}

	return
}
