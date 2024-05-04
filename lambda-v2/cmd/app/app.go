package app

import (
	api "lambda-v2/internal/api/http/user"
	"lambda-v2/pkg/user"
)

type App struct {
	ApiHandler api.ApiHandler
}

func NewApp() App {
	db, err := user.NewDynamoDBClient()
	if err != nil {
		return App{}
	}
	apiHandler := api.NewApiHandler(db)

	return App{
		ApiHandler: apiHandler,
	}
}
