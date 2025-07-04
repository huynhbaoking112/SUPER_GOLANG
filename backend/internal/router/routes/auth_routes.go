package routes

import (
	"go-backend-v2/internal/controllers"
	"go-backend-v2/internal/middlewares"
	"go-backend-v2/internal/repo"
	"go-backend-v2/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthRoutes struct {
	controller  *controllers.AuthController
	authService services.AuthServiceInterface
}

func NewAuthRoutes() *AuthRoutes {
	userRepo := repo.NewUserRepository()
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	return &AuthRoutes{
		controller:  authController,
		authService: authService,
	}
}

func (r *AuthRoutes) GetPrefix() string {
	return "/auth"
}

func (r *AuthRoutes) SetupRoutes(router fiber.Router) {
	authGroup := router.Group(r.GetPrefix())

	authGroup.Post("/signup", r.controller.Signup)
	authGroup.Post("/login", r.controller.Login)
	authGroup.Post("/logout", middlewares.AuthMiddleware(r.authService), r.controller.Logout)
}
