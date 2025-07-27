package controllers

import (
	"go-backend-v2/internal/common"
	"go-backend-v2/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService services.UserServiceInterface
}

func NewUserController(userService services.UserServiceInterface) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) GetCurrentUser(ctx *fiber.Ctx) error {
	userID := ctx.Locals(common.ContextUserID)
	if userID == nil {
		return common.ErrUnauthorized
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return common.ErrUnauthorized
	}

	user, err := c.userService.GetUserWithWorkspaces(userIDStr)
	if err != nil {
		return common.ErrUserNotFound
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User retrieved successfully",
		"user":    user,
	})
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	userID := ctx.Locals(common.ContextUserID)
	if userID == nil {
		return common.ErrUnauthorized
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return common.ErrUnauthorized
	}

	err := c.userService.DeleteUser(userIDStr)
	if err != nil {
		return common.ErrUserNotFound
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
