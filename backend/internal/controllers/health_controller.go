package controllers

import (
	"context"
	"go-backend-v2/global"

	"github.com/gofiber/fiber/v2"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (h *HealthController) GetHealth(c *fiber.Ctx) error {
	redisClient := global.RedisClient
	if redisClient == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"message": "Redis connection failed",
		})
	}
	redisClient.Set(context.Background(), "test", "test", 0)
	return c.JSON(fiber.Map{
		"status":  "OK",
		"message": "Server is running! 🚀",
		"service": "go-backend-v2",
		"path":    c.Path(),
	})
}

func (h *HealthController) GetTest(c *fiber.Ctx) error {
	response := fiber.Map{
		"message": "Test endpoint working! 👋",
		"path":    c.Path(),
		"method":  c.Method(),
	}

	// Kiểm tra có auth info không
	if userID := c.Locals("user_id"); userID != nil {
		response["authenticated"] = true
		response["user_id"] = userID
		response["username"] = c.Locals("username")
	} else {
		response["authenticated"] = false
	}

	return c.JSON(response)
}
