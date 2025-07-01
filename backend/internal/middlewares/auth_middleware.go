package middlewares

import (
	"go-backend-v2/internal/common"
	"go-backend-v2/internal/services"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(authService services.AuthServiceInterface) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Cookies(common.JWTCookieName)
		if token == "" {
			return common.ErrTokenRequired
		}

		userID, err := authService.ValidateToken(token)
		if err != nil {
			return common.ErrTokenInvalid
		}

		ctx.Locals(common.ContextUserID, userID)

		return ctx.Next()
	}
}
