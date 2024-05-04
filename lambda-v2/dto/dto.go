package dto

import (
	"log"

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
