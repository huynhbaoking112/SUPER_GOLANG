package routes

import (
	"go-backend-v2/internal/controllers"
	"go-backend-v2/internal/middlewares"
	"go-backend-v2/internal/repo"
	"go-backend-v2/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserRoutes struct {
	controller  *controllers.UserController
	authService services.AuthServiceInterface
}

func NewUserRoutes() *UserRoutes {
	userRepo := repo.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	authService := services.NewAuthService(userRepo)

	return &UserRoutes{
		controller:  userController,
		authService: authService,
	}
}

func (r *UserRoutes) GetPrefix() string {
	return "/users"
}

func (r *UserRoutes) SetupRoutes(router fiber.Router) {
	userGroup := router.Group(r.GetPrefix())
	userGroup.Use(middlewares.AuthMiddleware(r.authService))

	userGroup.Get("/me", r.controller.GetCurrentUser)
	userGroup.Delete("/me", r.controller.DeleteUser)
}
