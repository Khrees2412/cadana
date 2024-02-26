package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

func ValidateStruct(str interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(str)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

type Secret struct {
	ApiKeyOne string `json:"cadana-service-one"`
	ApiKeyTwo string `json:"cadana-service-two"`
}

func GetSecrets() Secret {
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
