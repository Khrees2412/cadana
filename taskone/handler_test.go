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
	// Init fiber app
	app := fiber.New()

	pair := "USD-EUR"

	secret := Secret{
		ApiKeyOne: os.Getenv("API_KEY_ONE"),
		ApiKeyTwo: os.Getenv("API_KEY_TWO"),
	}

	callCount := 0

	// Send test request
	req := httptest.NewRequest("GET", "/v1/exchange-rate?currency_pair="+pair, nil)
	c := http.Client{}
	c.Do(req)
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().URI().SetQueryString("currency_pair=" + pair)

	w := httptest.NewRecorder()

	// Call handler
	err := handler(ctx, secret)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, 2, callCount) // API called twice

	resp := ctx.Response()
	assert.Equal(t, fiber.StatusOK, resp.StatusCode())
	assert.JSONEq(t,
		`{"USD-EUR": 0.92}`, string(resp.Body()))

	if w.Code == http.StatusOK {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}
}
