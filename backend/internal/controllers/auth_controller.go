package controllers

import (
	"fmt"
	"go-backend-v2/global"
	"go-backend-v2/internal/common"
	"go-backend-v2/internal/dto"
	"go-backend-v2/internal/services"
	"go-backend-v2/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService services.AuthServiceInterface
	validator   *validator.Validate
}

func NewAuthController(authService services.AuthServiceInterface) *AuthController {
	v := validator.New()
	utils.SetupCustomValidators(v)

	return &AuthController{
		authService: authService,
		validator:   v,
	}
}

func (c *AuthController) Signup(ctx *fiber.Ctx) error {
	var req dto.SignupRequest

	if err := ctx.BodyParser(&req); err != nil {
		return common.ErrInvalidRequestBody
	}

	if err := c.validator.Struct(&req); err != nil {
		return common.ErrValidationFailed
	}

	if err := c.authService.Signup(&req); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.MessageResponse{
		Message: "Signup successful",
	})
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := ctx.BodyParser(&req); err != nil {
		return common.ErrInvalidRequestBody
	}

	if err := c.validator.Struct(&req); err != nil {
		return common.ErrValidationFailed
	}

	loginResponse, err := c.authService.Login(&req)
	if err != nil {
		return err
	}

	c.setJWTCookie(ctx, loginResponse.AccessToken)
	c.setEncryptedTokenCookie(ctx, loginResponse.EncryptedToken)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"user":    loginResponse.User,
	})
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	userID := ctx.Locals(common.ContextUserID)
	encryptedToken := ctx.Cookies(common.EncryptedTokenCookieName)

	if userID != nil && encryptedToken != "" {
		userIDStr, ok := userID.(string)
		if ok {
			err := c.authService.Logout(userIDStr, encryptedToken)
			if err != nil {
				fmt.Println("Failed to logout", err)
				// Log error but don't fail logout
				// TODO: Use proper logger instead of fmt.Printf
			}
		}
	}

	// Clear both cookies
	c.clearJWTCookie(ctx)
	c.clearEncryptedTokenCookie(ctx)

	return ctx.Status(fiber.StatusOK).JSON(dto.MessageResponse{
		Message: "Logout successful",
	})
}

func (c *AuthController) setJWTCookie(ctx *fiber.Ctx, token string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     common.JWTCookieName,
		Value:    token,
		MaxAge:   int(global.Config.JWT.ExpirationTime.Seconds()),
		HTTPOnly: global.Config.Cookie.HttpOnly,
		Secure:   global.Config.Cookie.Secure,
		SameSite: c.getSameSiteValue(global.Config.Cookie.SameSite),
		Domain:   global.Config.Cookie.Domain,
	})
}

func (c *AuthController) clearJWTCookie(ctx *fiber.Ctx) {
	ctx.Cookie(&fiber.Cookie{
		Name:     common.JWTCookieName,
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: global.Config.Cookie.HttpOnly,
		Secure:   global.Config.Cookie.Secure,
		SameSite: c.getSameSiteValue(global.Config.Cookie.SameSite),
		Domain:   global.Config.Cookie.Domain,
	})
}

func (c *AuthController) setEncryptedTokenCookie(ctx *fiber.Ctx, token string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     common.EncryptedTokenCookieName,
		Value:    token,
		MaxAge:   int(global.Config.JWT.ExpirationTime.Seconds()),
		HTTPOnly: global.Config.Cookie.HttpOnly,
		Secure:   global.Config.Cookie.Secure,
		SameSite: c.getSameSiteValue(global.Config.Cookie.SameSite),
		Domain:   global.Config.Cookie.Domain,
	})
}

func (c *AuthController) clearEncryptedTokenCookie(ctx *fiber.Ctx) {
	ctx.Cookie(&fiber.Cookie{
		Name:     common.EncryptedTokenCookieName,
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: global.Config.Cookie.HttpOnly,
		Secure:   global.Config.Cookie.Secure,
		SameSite: c.getSameSiteValue(global.Config.Cookie.SameSite),
		Domain:   global.Config.Cookie.Domain,
	})
}

func (c *AuthController) getSameSiteValue(sameSite string) string {
	switch sameSite {
	case common.CookieSameSiteStrict:
		return "Strict"
	case common.CookieSameSiteLax:
		return "Lax"
	case common.CookieSameSiteNone:
		return "None"
	default:
		return "Strict"
	}
}
