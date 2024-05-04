package api

import (
	"encoding/json"
	"net/http"

	"lambda-v2/pkg/user"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	userHandler user.UserService
}

func NewApiHandler(dbStore user.UserRepository) ApiHandler {
	return ApiHandler{
		userHandler: user.NewUserService(dbStore),
	}
}

func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerUser user.RegisterUser
	err := json.Unmarshal([]byte(request.Body), &registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	message, status := api.userHandler.HandleRegisterUser(registerUser)
	return events.APIGatewayProxyResponse{
		Body:       message,
		StatusCode: status,
	}, nil
}

func (api ApiHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var loginRequest user.LoginRequest
	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	message, status := api.userHandler.HandleLoginUser(loginRequest)
	return events.APIGatewayProxyResponse{
		Body:       message,
		StatusCode: status,
	}, nil
}
