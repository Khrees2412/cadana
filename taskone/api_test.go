package taskone

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestGetExchangeRate(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "This fetches an exchange rate",
			route:        "/exchange-rate",
			expectedCode: 200,
		},
	}

	app := fiber.New()

	for _, test := range tests {
		// Create a new http request with the route from the test case
		req := httptest.NewRequest("GET", "/v1"+test.route, strings.NewReader(`{"currency_pair":"USD-EUR"}`))

		// Perform the request plain with the app,
		// the second argument is a request latency
		// (set to -1 for no latency)
		resp, _ := app.Test(req, 1)

		// Check response code
		if resp.StatusCode != test.expectedCode {
			t.Errorf("Expected status code %d, got %d", test.expectedCode, resp.StatusCode)
		}

		// Verify, if the status code is as expected
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}

}
func TestExchangeRateAPIOne(t *testing.T) {
	key := os.Getenv("API_KEY_ONE")
	t.Run("valid api key and pair", func(t *testing.T) {
		rate := exchangeRateAPIOne("USD-EUR", key)
		log.Println(rate)
		assert.NotNil(t, rate)
		value := 0.85
		valueType := reflect.TypeOf(value)
		assert.Equal(t, valueType, *rate)
	})

	t.Run("invalid api key", func(t *testing.T) {
		rate := exchangeRateAPIOne("USD-EUR", "wrong_key")
		assert.Nil(t, rate)
	})

	t.Run("pair not found", func(t *testing.T) {
		rate := exchangeRateAPIOne("USD-XYZ", key)
		assert.Nil(t, rate)
	})
}

func TestExchangeRateAPITwo(t *testing.T) {
	key := os.Getenv("API_KEY_Two")
	t.Run("valid api key and pair", func(t *testing.T) {
		rate := exchangeRateAPITwo("USD-EUR", key)
		log.Println(rate)
		assert.NotNil(t, rate)
		value := 0.85
		valueType := reflect.TypeOf(value)
		assert.Equal(t, valueType, *rate)
	})

	t.Run("invalid api key", func(t *testing.T) {
		rate := exchangeRateAPIOne("USD-EUR", "wrong_key")
		assert.Nil(t, rate)
	})

	t.Run("pair not found", func(t *testing.T) {
		rate := exchangeRateAPIOne("USD-XYZ", key)
		assert.Nil(t, rate)
	})
}
