package dto

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(registerUser RegisterUser) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 10)
	if err != nil {
		log.Print("error occured ", err)
		return User{}, err
	}

	return User{
		Email:    registerUser.Email,
		Password: string(hashedPassword),
	}, nil
}

func ValidatePassword(hashedPassword, plainPassword string) bool {
	log.Print(hashedPassword)
	log.Print(plainPassword)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func CreateToken(user User) string {
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
