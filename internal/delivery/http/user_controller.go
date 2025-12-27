package http

import (
	"challenge-backend-1/internal/delivery/http/middleware"
	"challenge-backend-1/internal/model"
	"challenge-backend-1/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserController struct {
	Log     *zap.SugaredLogger
	UseCase *usecase.UserUseCase
}

func NewUserController(useCase *usecase.UserUseCase, logger *zap.SugaredLogger) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
	}
}

// Login godoc
// @Summary Login user
// @Description Login user
// @Tags User API
// @Accept json
// @Produce json
// @Param request body model.LoginUserRequest true "Login User Request"
// @Success 200 {object} model.WebResponse[model.UserResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/users/_login [post]
func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Current godoc
// @Summary Get current user
// @Description Get current user
// @Tags User API
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.WebResponse[model.UserResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/users/_current [get]
func (c *UserController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetUserRequest{
		ID: auth.ID,
	}

	response, err := c.UseCase.Current(ctx.UserContext(), request)
	if err != nil {
		c.Log.Errorw("Failed to get current user", "error", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Logout godoc
// @Summary Logout user
// @Description Logout user
// @Tags User API
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.WebResponse[bool]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/users [delete]
func (c *UserController) Logout(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.LogoutUserRequest{
		ID: auth.ID,
	}

	response, err := c.UseCase.Logout(ctx.UserContext(), request)
	if err != nil {
		c.Log.Errorw("Failed to logout user", "error", err)
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: response})
}

// Update godoc
// @Summary Update user
// @Description Update user
// @Tags User API
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body model.UpdateUserRequest true "Update User Request"
// @Success 200 {object} model.WebResponse[model.UserResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/users/_current [patch]
func (c *UserController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	request.ID = auth.ID
	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Errorw("Failed to update user", "error", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}
