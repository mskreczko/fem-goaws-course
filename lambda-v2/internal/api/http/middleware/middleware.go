package middleware

import (
	token "lambda-v2/internal/util"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

func ValidateJWTMiddleware(next func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		tokenString := token.ExtractTokenFromHeaders(request.Headers)
		if tokenString == "" {
			return events.APIGatewayProxyResponse{
				Body:       "Missing Auth Token",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		claims, err := token.ParseToken(tokenString)
		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       "User unathorized",
				StatusCode: http.StatusUnauthorized,
			}, err
		}

		expires := int64(claims["expires"].(float64))
		if time.Now().Unix() > expires {
			return events.APIGatewayProxyResponse{
				Body:       "Token expired",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		return next(request)
	}
}
