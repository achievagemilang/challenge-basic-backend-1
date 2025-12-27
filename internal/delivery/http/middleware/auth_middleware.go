package middleware

import (
	"strings"

	"challenge-backend-1/internal/model"
	"challenge-backend-1/internal/security"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewAuthMiddleware(tokenProvider security.TokenProvider, log *zap.SugaredLogger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return fiber.ErrUnauthorized
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.ErrUnauthorized
		}
		tokenString := parts[1]

		claims, err := tokenProvider.ValidateToken(tokenString)
		if err != nil {
			log.Warnf("Invalid token: %v", err)
			return fiber.ErrUnauthorized
		}

		var userID int64
		switch v := (*claims)["sub"].(type) {
		case float64:
			userID = int64(v)
		case string:
			return fiber.ErrUnauthorized
		default:
			return fiber.ErrUnauthorized
		}

		auth := &model.Auth{
			ID: userID,
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
