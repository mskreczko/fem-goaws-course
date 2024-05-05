package user

import (
	"fmt"
	token "lambda-v2/internal/util"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	HandleRegisterUser(registerUser RegisterUser) (string, int)
	HandleLoginUser(loginUser LoginRequest) (string, int)
}

type UserHandler struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) UserService {
	return &UserHandler{repository: repository}
}

func (u *UserHandler) HandleRegisterUser(registerUser RegisterUser) (string, int) {
	if registerUser.Email == "" || registerUser.Password == "" {
		return "Invalid request", http.StatusBadRequest
	}

	userExists, err := u.repository.DoesUserExist(registerUser.Email)
	if err != nil {
		return "Internal server error", http.StatusInternalServerError
	}

	if userExists {
		return "User already exists", http.StatusConflict
	}

	user, err := newUser(registerUser)
	if err != nil {
		return "Internal server error", http.StatusInternalServerError
	}

	err = u.repository.InsertUser(user)
	if err != nil {
		return "Internal server error", http.StatusInternalServerError
	}

	return "Successfully registered user", http.StatusCreated
}

func (u *UserHandler) HandleLoginUser(loginUser LoginRequest) (string, int) {
	_user, err := u.repository.GetUser(loginUser.Email)
	if err != nil {
		return "Internal server error", http.StatusInternalServerError
	}

	if !validatePassword(_user.Password, loginUser.Password) {
		return "Invalid credentials", http.StatusUnauthorized
	}

	accessToken := token.CreateToken(map[string]string{"email": _user.Email})
	successMsg := fmt.Sprintf(`{access_token: "%s"}`, accessToken)
	return successMsg, http.StatusOK
}

func newUser(registerUser RegisterUser) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 10)
	if err != nil {
		return User{}, err
	}

	return User{
		Email:    registerUser.Email,
		Password: string(hashedPassword),
	}, nil
}

func validatePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
