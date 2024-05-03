package database

import (
	"context"
	"lambda-v2/dto"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	TABLE_NAME = "users"
)

type UserStore interface {
	DoesUserExist(username string) (bool, error)
	InsertUser(user dto.RegisterUser) error
}

type DynamoDBClient struct {
	cfg    aws.Config
	client *dynamodb.Client
}

func NewDynamoDBClient(options ...func(*config.LoadOptions) error) (*DynamoDBClient, error) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		options...,
	)

	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)

	return &DynamoDBClient{
		client: client,
	}, nil
}

func (db DynamoDBClient) DoesUserExist(email string) (bool, error) {
	result, err := db.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{"email": &types.AttributeValueMemberS{Value: email}}, TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		return true, err
	}

	if result.Item == nil {
		return false, nil
	}

	return true, nil
}

func (db DynamoDBClient) InsertUser(user dto.RegisterUser) error {
	newUser, err := dto.NewUser(user)
	if err != nil {
		log.Print(err)
		return err
	}
	item := map[string]types.AttributeValue{"email": &types.AttributeValueMemberS{Value: newUser.Email}, "password": &types.AttributeValueMemberS{Value: newUser.PasswordHash}}

	_, err = db.client.PutItem(context.TODO(), &dynamodb.PutItemInput{TableName: aws.String(TABLE_NAME), Item: item})
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
