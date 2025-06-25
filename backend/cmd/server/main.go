package main

import (
	"fmt"
	"go-backend-v2/global"
	"go-backend-v2/internal/initialize"
	"go-backend-v2/internal/router"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	initialize.Run()

	app := fiber.New(fiber.Config{
		AppName: "IAM",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error":  err.Error(),
				"code":   code,
				"path":   c.Path(),
				"method": c.Method(),
			})
		},
	})

	router.SetupRoutes(app)

	port := fmt.Sprintf(":%d", global.Config.Server.Port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
