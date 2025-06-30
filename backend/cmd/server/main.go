package main

import (
	"fmt"
	"go-backend-v2/global"
	"go-backend-v2/internal/initialize"
	"go-backend-v2/internal/middlewares"
	"go-backend-v2/internal/router"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	initialize.Run()

	app := fiber.New(fiber.Config{
		AppName:      "IAM",
		ErrorHandler: middlewares.ErrorHandler,
	})

	app.Use(recover.New())

	router.SetupRoutes(app)

	port := fmt.Sprintf(":%d", global.Config.Server.Port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
