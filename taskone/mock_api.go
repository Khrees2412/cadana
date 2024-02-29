package taskone

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Quote struct {
	Pair string  `json:"pair"`
	Rate float64 `json:"rate"`
}

func exchangeRateAPI(pair string, apiKey string, service string) *float64 {
	var filepath string
	var envKey string

	switch service {
	case "one":
		filepath = "taskone/servicea.json"
		envKey = os.Getenv("API_KEY_ONE")
	case "two":
		filepath = "taskone/serviceb.json"
		envKey = os.Getenv("API_KEY_TWO")
	default:
		log.Fatalf("invalid service passed %v", service)
	}

	if apiKey != envKey {
		log.Fatalf("invalid api key supplied")
		return nil
	}
	jsonFile, err := os.Open(filepath)
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
