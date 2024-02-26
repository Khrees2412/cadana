package taskone

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/khrees2412/cadana/util"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func Start() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Get("/v1/exchange-rate", handler)

	port := "5001"
	if err := app.Listen(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatal(fmt.Sprintf("listen: %s\n", err))
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
		return nil
	}
	jsonFile, err := os.Open("api/servicea.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer jsonFile.Close()

	byteValue, e := ioutil.ReadAll(jsonFile)
	if e != nil {
		fmt.Println(e)
		return nil
	}

	var quotes []Quotes
	err = json.Unmarshal(byteValue, &quotes)
	if err != nil {
		fmt.Println(err)
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
		return nil
	}
	jsonFile, err := os.Open("api/serviceb.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer jsonFile.Close()

	byteValue, e := ioutil.ReadAll(jsonFile)
	if e != nil {
		fmt.Println(e)
		return nil
	}

	var quotes []Quotes
	err = json.Unmarshal(byteValue, &quotes)
	if err != nil {
		fmt.Println(err)
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

	apiKeys := util.GetSecrets()

	// Channels to receive results from the functions
	result1 := make(chan *float64)
	result2 := make(chan *float64)

	// Start Goroutines to execute both functions concurrently
	go func() {
		result1 <- exchangeRateAPIOne(currencyPair, apiKeys.ApiKeyOne)
	}()
	go func() {
		result2 <- exchangeRateAPITwo(currencyPair, apiKeys.ApiKeyTwo)
	}()

	var res *float64

	// Select the first result that is available
	select {
	case res = <-result1:
		fmt.Println("Result from API one", *res)
	case res = <-result2:
		fmt.Println("Result from API two:", *res)
	}

	if res == nil {
		return ctx.JSON(fiber.Map{
			"message":    "Currency pair not found",
			currencyPair: nil,
		})
	}

	return ctx.JSON(
		fiber.Map{
			currencyPair: *res,
		},
	)
}
