package taskone

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExchangeAPI(t *testing.T) {
	pair := "USD-EUR"
	_ = Secret{
		ApiKeyOne: os.Getenv("API_KEY_ONE"),
		ApiKeyTwo: os.Getenv("API_KEY_TWO"),
	}

	t.Run("returns nil for invalid api key", func(t *testing.T) {
		rate := exchangeRateAPI(pair, "wrong-key", "one")
		assert.Nil(t, rate)
	})
}
