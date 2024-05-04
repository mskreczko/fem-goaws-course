package user

type RegisterUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
