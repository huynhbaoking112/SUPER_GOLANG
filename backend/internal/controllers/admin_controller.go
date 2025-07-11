package controllers

import (
	"go-backend-v2/internal/common"
	"go-backend-v2/internal/dto"
	"go-backend-v2/internal/services"
	"go-backend-v2/pkg/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AdminController struct {
	workspaceService services.WorkspaceServiceInterface
	validator        *validator.Validate
}

func NewAdminController(workspaceService services.WorkspaceServiceInterface) *AdminController {
	v := validator.New()
	utils.SetupCustomValidators(v)

	return &AdminController{
		workspaceService: workspaceService,
		validator:        v,
	}
}

func (c *AdminController) CreateWorkspace(ctx *fiber.Ctx) error {
	userID := ctx.Locals(common.ContextUserID)
	if userID == nil {
		return common.ErrUnauthorized
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return common.ErrUnauthorized
	}

	var req dto.CreateWorkspaceRequest
	if err := ctx.BodyParser(&req); err != nil {
		return common.ErrInvalidRequestBody
	}

	if err := c.validator.Struct(&req); err != nil {
		return common.ErrValidationFailed
	}

	workspace, err := c.workspaceService.CreateWorkspace(userIDStr, &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Workspace created successfully",
		"data":    workspace,
		"meta": fiber.Map{
			"timestamp": time.Now(),
			"path":      ctx.Path(),
		},
	})
}
