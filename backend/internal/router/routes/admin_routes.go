package routes

import (
	"go-backend-v2/internal/controllers"
	"go-backend-v2/internal/middlewares"
	"go-backend-v2/internal/repo"
	"go-backend-v2/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AdminRoutes struct {
	adminController *controllers.AdminController
	authService     services.AuthServiceInterface
}

func NewAdminRoutes() *AdminRoutes {
	workspaceRepo := repo.NewWorkspaceRepository()
	userRepo := repo.NewUserRepository()

	workspaceService := services.NewWorkspaceService(workspaceRepo, userRepo)
	authService := services.NewAuthService(userRepo)

	adminController := controllers.NewAdminController(workspaceService)

	return &AdminRoutes{
		adminController: adminController,
		authService:     authService,
	}
}

func (r *AdminRoutes) GetPrefix() string {
	return "/admin"
}

func (r *AdminRoutes) SetupRoutes(router fiber.Router) {
	adminGroup := router.Group(r.GetPrefix())
	adminGroup.Use(middlewares.RequireSuperAdmin(r.authService))

	workspacesGroup := adminGroup.Group("/workspaces")
	workspacesGroup.Post("/", r.adminController.CreateWorkspace)
}
