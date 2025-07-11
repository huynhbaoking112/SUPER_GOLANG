package routes

import "github.com/gofiber/fiber/v2"

type RouteModule interface {
	SetupRoutes(router fiber.Router)
	GetPrefix() string
}

type RouteConfig struct {
	APIVersion string
	BaseURL    string
}

func NewRouteConfig() *RouteConfig {
	return &RouteConfig{
		APIVersion: "v1",
		BaseURL:    "/api",
	}
}
