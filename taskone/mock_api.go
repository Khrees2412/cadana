package taskone

import (
	"log"
	"os"
)

type Quote struct {
	Pair string
	Rate float64
}

var serviceOneData = []Quote{
	{
		Pair: "USD-EUR",
		Rate: 0.92,
	},
	{
		Pair: "USD-GBP",
		Rate: 0.75,
	},
	{
		Pair: "EUR-GBP",
		Rate: 0.81,
	},
	{
		Pair: "USD-JPY",
		Rate: 110.21,
	},
}

var serviceTwoData = []Quote{
	{
		Pair: "USD-EUR",
		Rate: 0.95,
	},
	{
		Pair: "USD-GBP",
		Rate: 0.77,
	},
	{
		Pair: "EUR-GBP",
		Rate: 0.90,
	},
	{
		Pair: "USD-JPY",
		Rate: 111.31,
	},
}

func exchangeRateAPI(pair string, apiKey string, service string) *float64 {
	var data []Quote
	var envKey string

	switch service {
	case "one":
		data = serviceOneData
		envKey = os.Getenv("API_KEY_ONE")
	case "two":
		data = serviceTwoData
		envKey = os.Getenv("API_KEY_TWO")
	default:
		log.Fatalf("invalid service passed %v", service)
	}

	if apiKey != envKey {
		log.Println("invalid api key supplied")
		return nil
	}

	for _, v := range data {
		if v.Pair == pair {
			return &v.Rate
		}
	}
	return nil
}
