package jwt

import (
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

func CreateToken(customerID, name string) (token string, err error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"customer_id": customerID,
		"name":        name,
		"exp":         now.Add(time.Minute * 43200).Unix(),
		"iat":         now.Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return
	}

	return
}
