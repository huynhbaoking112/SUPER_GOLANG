package router

import (
	"go-backend-v2/internal/router/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type RouteManager struct {
	config       *routes.RouteConfig
	routeModules []routes.RouteModule
}

func NewRouteManager() *RouteManager {
	config := routes.NewRouteConfig()

	routeModules := []routes.RouteModule{
		routes.NewPublicRoutes(),
	}

	return &RouteManager{
		config:       config,
		routeModules: routeModules,
	}
}

func (rm *RouteManager) SetupRoutes(app *fiber.App) {

	app.Use(helmet.New())
	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path} - ${ip}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// ======================
	// SETUP ROUTE MODULES
	// ======================
	publicRoutes := routes.NewPublicRoutes()
	publicRoutes.SetupRoutes(app)

	api := app.Group(rm.config.BaseURL)
	version := api.Group("/" + rm.config.APIVersion)

	for _, module := range rm.routeModules {
		if module.GetPrefix() != "" {
			module.SetupRoutes(version)
		}
	}

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":  "Route not found",
			"path":   c.Path(),
			"method": c.Method(),
		})
	})
}

func SetupRoutes(app *fiber.App) {
	manager := NewRouteManager()
	manager.SetupRoutes(app)
}
