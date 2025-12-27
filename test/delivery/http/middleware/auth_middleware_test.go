package middleware_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"challenge-backend-1/internal/delivery/http/middleware"
	"challenge-backend-1/internal/model"
	"challenge-backend-1/test/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenProvider := mocks.NewMockTokenProvider(ctrl)
	logger := zap.NewNop().Sugar()

	app := fiber.New()
	app.Use(middleware.NewAuthMiddleware(mockTokenProvider, logger))
	app.Get("/v1/test", func(c *fiber.Ctx) error {
		auth := middleware.GetUser(c)
		return c.JSON(auth)
	})

	t.Run("Success", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": float64(123),
		}

		mockTokenProvider.EXPECT().ValidateToken("valid_token").Return(&claims, nil)

		req := httptest.NewRequest("GET", "/v1/test", nil)
		req.Header.Set("Authorization", "Bearer valid_token")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("MissingHeader", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/test", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("InvalidFormat", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/test", nil)
		req.Header.Set("Authorization", "InvalidFormat")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("InvalidToken", func(t *testing.T) {
		mockTokenProvider.EXPECT().ValidateToken("invalid_token").Return(nil, errors.New("invalid token"))

		req := httptest.NewRequest("GET", "/v1/test", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("InvalidSubClaimType", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "string_id",
		}

		mockTokenProvider.EXPECT().ValidateToken("valid_token").Return(&claims, nil)

		req := httptest.NewRequest("GET", "/v1/test", nil)
		req.Header.Set("Authorization", "Bearer valid_token")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}

func TestGetUser(t *testing.T) {
	app := fiber.New()
	app.Get("/user", func(c *fiber.Ctx) error {
		auth := &model.Auth{ID: 1}
		c.Locals("auth", auth)
		user := middleware.GetUser(c)
		return c.JSON(user)
	})

	req := httptest.NewRequest("GET", "/user", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
