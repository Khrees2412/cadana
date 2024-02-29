package taskone

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

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
		result1 <- exchangeRateAPI(currencyPair, apiKeyOne, "one")
	}()
	go func() {
		result2 <- exchangeRateAPI(currencyPair, apiKeyTwo, "two")
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
