package token

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetTokenSecret() (string, error) {
	SECRET_ID := "jwt-secret"
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
	)
	if err != nil {
		return "", err
	}

	client := secretsmanager.NewFromConfig(cfg)
	result, err := client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{SecretId: &SECRET_ID})
	if err != nil {
		return "", err
	}

	return *result.SecretString, nil
}
