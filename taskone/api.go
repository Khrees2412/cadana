package taskone

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/khrees2412/cadana/util"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func Start() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Get("/v1/exchange-rate", secretMiddleware(), handler)

	port := "5001"
	if err := app.Listen(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatal(fmt.Sprintf("listen: %s\n", err))
	}
}

func secretMiddleware() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("secret_key_one", util.GetSecret().ApiKeyOne)
		c.Set("secret_key_two", util.GetSecret().ApiKeyTwo)
		return c.Next()
	}
}

type Quotes struct {
	Pair string  `json:"pair"`
	Rate float64 `json:"rate"`
}
type CurrencyPairRequest struct {
	CurrencyPair string `json:"currency_pair" validate:"required"`
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

	var quotes []Quotes
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

	var quotes []Quotes
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

func handler(ctx *fiber.Ctx) error {
	currencyPair := ctx.Query("currency_pair")
	if currencyPair == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": "Currency pair is required",
		})
	}

	apiKeyOne := ctx.GetRespHeader("secret_key_one")
	apiKeyTwo := ctx.GetRespHeader("secret_key_two")

	appCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	return ctx.JSON(
		fiber.Map{
			currencyPair: *res,
		},
	)
}
