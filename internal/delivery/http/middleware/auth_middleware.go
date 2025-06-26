package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/wastetrack/wastetrack-backend/internal/helper"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func NewJWTAuth(jwtHelper *helper.JWTHelper) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return fiber.ErrUnauthorized
		}

		// Extract Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return fiber.ErrUnauthorized
		}

		token := tokenParts[1]
		claims, err := jwtHelper.ValidateAccessToken(token)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		auth := &model.Auth{
			ID:              claims.UserID,
			Role:            claims.Role,
			IsEmailVerified: claims.IsEmailVerified,
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}

func RequireRoles(roles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := GetUser(ctx)
		if auth == nil {
			return fiber.ErrForbidden
		}
		for _, r := range roles {
			if auth.Role == r {
				return ctx.Next()
			}
		}
		return fiber.ErrForbidden
	}
}

func RequireEmailVerification() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := GetUser(ctx)
		if auth == nil || !auth.IsEmailVerified {
			return fiber.NewError(fiber.StatusForbidden, "Email verification required")
		}
		return ctx.Next()
	}
}
