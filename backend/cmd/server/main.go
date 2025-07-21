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

	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins:     "http://localhost:5173",
	// 	AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
	// 	AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	// 	AllowCredentials: true, // Chỉ bật khi cần thiết
	// 	ExposeHeaders:    "Set-Cookie, Authorization",
	// 	MaxAge:           86400,
	// }))
	// app.Options("/*", func(c *fiber.Ctx) error {
	// 	return c.Status(204).Send(nil)
	// })
	app.Use(recover.New())

	router.SetupRoutes(app)

	port := fmt.Sprintf(":%d", global.Config.Server.Port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
