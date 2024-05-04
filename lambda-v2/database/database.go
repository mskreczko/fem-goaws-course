package database

import (
	"context"
	"fmt"
	"lambda-v2/dto"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	TABLE_NAME = "users"
)

type UserStore interface {
	DoesUserExist(username string) (bool, error)
	InsertUser(user dto.User) error
	GetUser(username string) (dto.User, error)
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

func (db DynamoDBClient) InsertUser(user dto.User) error {
	item := map[string]types.AttributeValue{"email": &types.AttributeValueMemberS{Value: user.Email}, "password": &types.AttributeValueMemberS{Value: user.Password}}

	_, err := db.client.PutItem(context.TODO(), &dynamodb.PutItemInput{TableName: aws.String(TABLE_NAME), Item: item})
	if err != nil {
		return err
	}
	return nil
}

func (db DynamoDBClient) GetUser(email string) (dto.User, error) {
	var user dto.User

	result, err := db.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key:       map[string]types.AttributeValue{"email": &types.AttributeValueMemberS{Value: email}},
	})

	if err != nil {
		return user, err
	}

	if result.Item == nil {
		return user, fmt.Errorf("user not found")
	}

	err = attributevalue.UnmarshalMap(result.Item, &user)
	log.Print(user)
	if err != nil {
		return user, err
	}

	return user, nil
}
