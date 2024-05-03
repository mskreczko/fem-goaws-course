package api

import (
	"fmt"
	"lambda-v2/database"
	"lambda-v2/dto"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(event dto.RegisterUser) error {
	if event.Email == "" || event.Password == "" {
		return fmt.Errorf("request has empty parameters")
	}

	userExists, err := api.dbStore.DoesUserExist(event.Email)
	if err != nil {
		return fmt.Errorf("there was an error checking if user exists %w", err)
	}

	if userExists {
		return fmt.Errorf("a user with that email already exists")
	}

	err = api.dbStore.InsertUser(event)
	if err != nil {
		return fmt.Errorf("there was an error registering a user %w", err)
	}

	return nil
}
