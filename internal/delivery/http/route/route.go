package route

import (
	"challenge-backend-1/internal/delivery/http"

	_ "challenge-backend-1/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"
)

type RouteConfig struct {
	App            *fiber.App
	UserController *http.UserController
	Log            *zap.SugaredLogger
	AuthMiddleware fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/session", c.UserController.Login)

	c.App.Get("/swagger/*", swagger.HandlerDefault)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)
}
