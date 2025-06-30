package middlewares

import (
	"errors"
	"net/http"

	"go-backend-v2/internal/common"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var apiErr *common.APIError
	if ok := errors.As(err, &apiErr); ok {
		return c.Status(apiErr.Status).JSON(fiber.Map{
			"error":   http.StatusText(apiErr.Status),
			"code":    apiErr.Code,
			"message": apiErr.Message,
			"path":    c.OriginalURL(),
		})
	}

	return c.Status(common.ErrInternalServer.Status).JSON(fiber.Map{
		"error":   http.StatusText(common.ErrInternalServer.Status),
		"code":    common.ErrInternalServer.Code,
		"message": common.ErrInternalServer.Message,
		"path":    c.OriginalURL(),
	})
}
