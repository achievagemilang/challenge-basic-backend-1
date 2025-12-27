package http

import (
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
// @Success 200 {object} model.WebResponse[model.LoginResponse]
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /session [post]
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

	return ctx.JSON(model.WebResponse[*model.LoginResponse]{
		Ok:   true,
		Data: response,
	})
}
