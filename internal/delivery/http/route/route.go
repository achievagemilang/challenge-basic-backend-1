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
	v1 := c.App.Group("/api/v1")
	c.SetupGuestRoute(v1)
	c.SetupAuthRoute(v1)
}

func (c *RouteConfig) SetupGuestRoute(r fiber.Router) {
	r.Post("/session", c.UserController.Login)
	r.Put("/session", c.UserController.Refresh)

	c.App.Get("/swagger/*", swagger.HandlerDefault)
}

func (c *RouteConfig) SetupAuthRoute(r fiber.Router) {
	r.Use(c.AuthMiddleware)
}
