package config

import (
	"challenge-backend-1/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(config *viper.Viper) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      config.GetString("app.name"),
		ErrorHandler: NewErrorHandler(),
		Prefork:      config.GetBool("web.prefork"),
	})

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		if code == fiber.StatusBadRequest {
			return ctx.Status(code).JSON(model.ErrorResponse{
				Ok:  false,
				Err: "ERR_BAD_REQUEST",
				Msg: err.Error(),
			})
		}

		if code == fiber.StatusUnauthorized {
			return ctx.Status(code).JSON(model.ErrorResponse{
				Ok:  false,
				Err: "ERR_INVALID_ACCESS_TOKEN",
				Msg: "invalid access token",
			})
		}

		if code == fiber.StatusForbidden {
			return ctx.Status(code).JSON(model.ErrorResponse{
				Ok:  false,
				Err: "ERR_FORBIDDEN_ACCESS",
				Msg: "user doesn't have enough authorization",
			})
		}

		if code == fiber.StatusNotFound {
			return ctx.Status(code).JSON(model.ErrorResponse{
				Ok:  false,
				Err: "ERR_NOT_FOUND",
				Msg: "resource is not found",
			})
		}

		return ctx.Status(code).JSON(model.ErrorResponse{
			Ok:  false,
			Err: "ERR_INTERNAL_ERROR",
			Msg: err.Error(),
		})
	}
}
