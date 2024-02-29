package taskone

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExchangeAPI(t *testing.T) {
	pair := "USD-EUR"
	testRate := 0.92
	secret := Secret{
		ApiKeyOne: os.Getenv("API_KEY_ONE"),
		ApiKeyTwo: os.Getenv("API_KEY_TWO"),
	}

	rate := exchangeRateAPI(pair, secret.ApiKeyOne, "one")

	assert.Equal(t, rate, testRate)
}
