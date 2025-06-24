package main

import (
	"go-backend-v2/internal/router"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:               "IAM",
		DisableStartupMessage: false,
	})

	// Setup routes
	router.SetupRoutes(app)

	// Start server
	log.Println("Server is running on port 8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
