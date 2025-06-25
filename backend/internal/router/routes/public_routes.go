package routes

import (
	"go-backend-v2/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

type PublicRoutes struct {
	healthController *controllers.HealthController
}

func NewPublicRoutes() *PublicRoutes {
	return &PublicRoutes{
		healthController: controllers.NewHealthController(),
	}
}

func (r *PublicRoutes) GetPrefix() string {
	return ""
}

func (r *PublicRoutes) SetupRoutes(router fiber.Router) {
	router.Get("/", r.healthController.GetHealth)
	router.Get("/health", r.healthController.GetHealth)
	router.Get("/test", r.healthController.GetTest)
}
