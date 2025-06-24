package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (h *HealthController) GetHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Server is running!",
		"service": "IAM",
	})
}

func (h *HealthController) GetTest(c *fiber.Ctx) error {
	return c.SendString("Hello, World ! This is a test endpoint")
}
