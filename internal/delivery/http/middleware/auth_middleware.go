package middleware

import (
	"strings"

	"challenge-backend-1/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewAuthMiddleware(config *viper.Viper, log *zap.SugaredLogger) fiber.Handler {
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

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(config.GetString("jwt.secret")), nil
		})

		if err != nil || !token.Valid {
			log.Warnf("Invalid token: %v", err)
			return fiber.ErrUnauthorized
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.ErrUnauthorized
		}

		var userID int64
		switch v := claims["sub"].(type) {
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
