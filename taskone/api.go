package taskone

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
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

func Start() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	secret := GetSecret()
	app.Get("/v1/exchange-rate", WithSecret(handler, secret))

	port := "5001"
	if err := app.Listen(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatal(fmt.Sprintf("listen: %s\n", err))
	}
}

type Quote struct {
	Pair string  `json:"pair"`
	Rate float64 `json:"rate"`
}

func exchangeRateAPIOne(pair string, apiKey string) *float64 {
	if apiKey != os.Getenv("API_KEY_ONE") {
		log.Fatalf("invalid api key supplied")
		return nil
	}
	jsonFile, err := os.Open("taskone/servicea.json")
	if err != nil {
		log.Fatalf("unable to open json file: %v", err)
		return nil
	}
	defer jsonFile.Close()

	byteValue, e := ioutil.ReadAll(jsonFile)
	if e != nil {
		log.Fatalf("unable to read json file: %v", e)
		return nil
	}

	var quotes []Quote
	err = json.Unmarshal(byteValue, &quotes)
	if err != nil {
		log.Fatalf("unable to unmarshal json file: %v", err)
		return nil
	}

	for _, v := range quotes {
		if v.Pair == pair {
			return &v.Rate
		}
	}
	return nil
}

func exchangeRateAPITwo(pair string, apiKey string) *float64 {
	if apiKey != os.Getenv("API_KEY_TWO") {
		log.Fatalf("invalid api key supplied")
		return nil
	}
	jsonFile, err := os.Open("taskone/serviceb.json")
	if err != nil {
		log.Fatalf("unable to open json file: %v", err)
		return nil
	}
	defer jsonFile.Close()

	byteValue, e := ioutil.ReadAll(jsonFile)
	if e != nil {
		log.Fatalf("unable to read json file: %v", e)
		return nil
	}

	var quotes []Quote
	err = json.Unmarshal(byteValue, &quotes)
	if err != nil {
		log.Fatalf("unable to unmarshal json file: %v", err)
		return nil
	}

	for _, v := range quotes {
		if v.Pair == pair {
			return &v.Rate
		}
	}
	return nil
}

func handler(ctx *fiber.Ctx, secret Secret) error {
	currencyPair := ctx.Query("currency_pair")
	if currencyPair == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": "Currency pair is required",
		})
	}

	apiKeyOne := secret.ApiKeyOne
	apiKeyTwo := secret.ApiKeyTwo

	appCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Channels to receive results from the functions
	result1 := make(chan *float64)
	result2 := make(chan *float64)

	// Start Goroutines to execute both functions concurrently
	go func() {
		result1 <- exchangeRateAPIOne(currencyPair, apiKeyOne)
	}()
	go func() {
		result2 <- exchangeRateAPITwo(currencyPair, apiKeyTwo)
	}()

	var res *float64

	// Select the first result that is available
	select {
	case res = <-result1:
		if res == nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
			})
		}
		fmt.Println("Result from API one", *res)
	case res = <-result2:
		if res == nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
			})
		}
		fmt.Println("Result from API two:", *res)
	case <-appCtx.Done():
		fmt.Println("Context cancelled")
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			currencyPair: *res,
		},
	)
}
