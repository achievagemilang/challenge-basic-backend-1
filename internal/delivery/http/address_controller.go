package http

import (
	"challenge-backend-1/internal/delivery/http/middleware"
	"challenge-backend-1/internal/model"
	"challenge-backend-1/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AddressController struct {
	UseCase *usecase.AddressUseCase
	Log     *zap.SugaredLogger
}

func NewAddressController(useCase *usecase.AddressUseCase, log *zap.SugaredLogger) *AddressController {
	return &AddressController{
		Log:     log,
		UseCase: useCase,
	}
}

// Create godoc
// @Summary Create new address
// @Description Create new address
// @Tags Address API
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param contactId path string true "Contact ID"
// @Param request body model.CreateAddressRequest true "Create Address Request"
// @Success 200 {object} model.WebResponse[model.AddressResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/contacts/{contactId}/addresses [post]
func (c *AddressController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.CreateAddressRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Errorw("failed to parse request body", "error", err)
		return fiber.ErrBadRequest
	}

	request.UserId = auth.ID
	request.ContactId = ctx.Params("contactId")

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Errorw("failed to create address", "error", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AddressResponse]{Data: response})
}

// List godoc
// @Summary List addresses
// @Description List addresses
// @Tags Address API
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param contactId path string true "Contact ID"
// @Success 200 {object} model.WebResponse[[]model.AddressResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/contacts/{contactId}/addresses [get]
func (c *AddressController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Params("contactId")

	request := &model.ListAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
	}

	responses, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.Errorw("failed to list addresses", "error", err)
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.AddressResponse]{Data: responses})
}

// Get godoc
// @Summary Get address
// @Description Get address
// @Tags Address API
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param contactId path string true "Contact ID"
// @Param addressId path string true "Address ID"
// @Success 200 {object} model.WebResponse[model.AddressResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/contacts/{contactId}/addresses/{addressId} [get]
func (c *AddressController) Get(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Params("contactId")
	addressId := ctx.Params("addressId")

	request := &model.GetAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
		ID:        addressId,
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Errorw("failed to get address", "error", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AddressResponse]{Data: response})
}

// Update godoc
// @Summary Update address
// @Description Update address
// @Tags Address API
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param contactId path string true "Contact ID"
// @Param addressId path string true "Address ID"
// @Param request body model.UpdateAddressRequest true "Update Address Request"
// @Success 200 {object} model.WebResponse[model.AddressResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/contacts/{contactId}/addresses/{addressId} [put]
func (c *AddressController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateAddressRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Errorw("failed to parse request body", "error", err)
		return fiber.ErrBadRequest
	}

	request.UserId = auth.ID
	request.ContactId = ctx.Params("contactId")
	request.ID = ctx.Params("addressId")

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Errorw("failed to update address", "error", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AddressResponse]{Data: response})
}

// Delete godoc
// @Summary Delete address
// @Description Delete address
// @Tags Address API
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param contactId path string true "Contact ID"
// @Param addressId path string true "Address ID"
// @Success 200 {object} model.WebResponse[bool]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/contacts/{contactId}/addresses/{addressId} [delete]
func (c *AddressController) Delete(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Params("contactId")
	addressId := ctx.Params("addressId")

	request := &model.DeleteAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
		ID:        addressId,
	}

	if err := c.UseCase.Delete(ctx.UserContext(), request); err != nil {
		c.Log.Errorw("failed to delete address", "error", err)
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}
