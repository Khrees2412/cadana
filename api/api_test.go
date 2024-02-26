package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGetExchangeRate(t *testing.T) {

	// Mock API_KEY_ONE and API_KEY_TWO for testing
	os.Setenv("API_KEY_ONE", "mock_api_key_one")
	os.Setenv("API_KEY_TWO", "mock_api_key_two")

	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "This fetches exchange rate",
			route:        "/exchange-rate",
			expectedCode: 200,
		},
	}

	// Define Fiber app.
	app := fiber.New()

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route from the test case
		req := httptest.NewRequest("POST", test.route, strings.NewReader(`{"currency_pair":"USD-EUR"}`))

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
