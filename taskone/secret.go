package taskone

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

type Secret struct {
	ApiKeyOne string `json:"cadana-service-one"`
	ApiKeyTwo string `json:"cadana-service-two"`
}

func GetSecret() Secret {
	secretName := os.Getenv("SECRET_NAME")
	region := "eu-north-1"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {

		log.Fatal(err.Error())
	}

	var secret Secret
	// Decrypts secret using the associated KMS key.
	var res = *result.SecretString
	byteValue := []byte(res)

	err = json.Unmarshal(byteValue, &secret)
	if err != nil {
		fmt.Println(err)
	}
	return secret
}

func WithSecret(fn func(ctx *fiber.Ctx, secret Secret) error, secret Secret) func(c *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return fn(ctx, secret)
	}
}
