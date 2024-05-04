package token

import (
	"lambda-v2/pkg/user"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(user user.User) string {
	now := time.Now()
	validUntil := now.Add(time.Hour * 1).Unix()

	claims := jwt.MapClaims{
		"user":    user.Email,
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)
	secret := "secret" // TODO move to vault

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Print("error occurred when signing token: ", err)
		return ""
	}

	return tokenString
}
