package router

import (
	"go-backend-v2/internal/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App) {
	//Middleware
	app.Use(cors.New())
	// app.Use(logger.New())

	//Initialize controller
	healthController := controllers.NewHealthController()

	// API routes
	api := app.Group("/api")

	// Health check routes
	app.Get("/", healthController.GetHealth)
	app.Get("/test", healthController.GetTest)
	app.Get("/healths", healthController.GetHealth)

	// API routes group
	api.Get("/health", healthController.GetHealth)

}
