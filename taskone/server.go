package taskone

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
)

func Start() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	secret := GetSecret()
	app.Get("/v1/exchange-rate", WithSecret(handler, secret))

	port := "5001"
	if err := app.Listen(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatal(fmt.Sprintf("listen: %s\n", err))
	}
}
