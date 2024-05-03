package app

import (
	"lambda-v2/api"
	"lambda-v2/database"
)

type App struct {
	ApiHandler api.ApiHandler
}

func NewApp() App {
	db := database.NewDynamoDBClient()
	apiHandler := api.NewApiHandler(db)

	return App{
		ApiHandler: apiHandler,
	}
}