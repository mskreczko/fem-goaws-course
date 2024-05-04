package token

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(requestedClaims map[string]string) string {
	now := time.Now()
	validUntil := now.Add(time.Hour * 1).Unix()

	claims := jwt.MapClaims{
		"expires": validUntil,
	}

	for k, v := range requestedClaims {
		claims[k] = v
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
