package taskone

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {

	// Test data
	pair := "USD-EUR"
	//rate := 0.92
	secret := Secret{
		ApiKeyOne: os.Getenv("API_KEY_ONE"),
		ApiKeyTwo: os.Getenv("API_KEY_TWO"),
	}

	// Mock API call
	callCount := 0
	exchangeRateAPI(pair, secret.ApiKeyOne, "one")

	// Init fiber app
	app := fiber.New()

	// Send test request
	req := httptest.NewRequest("GET", "/?currency_pair="+pair, nil)
	c := http.Client{}
	c.Do(req)
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().URI().SetQueryString("currency_pair=" + pair)

	// Call handler
	err := handler(ctx, secret)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, 2, callCount) // API called twice

	resp := ctx.Response()
	assert.Equal(t, fiber.StatusOK, resp.StatusCode())
	assert.JSONEq(t,
		`{"USD-EUR": 0.92}`, string(resp.Body()))

}
