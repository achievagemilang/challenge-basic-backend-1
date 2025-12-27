package http

import (
	"math"

	"challenge-backend-1/internal/delivery/http/middleware"
	"challenge-backend-1/internal/model"
	"challenge-backend-1/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ContactController struct {
	UseCase *usecase.ContactUseCase
	Log     *zap.SugaredLogger
}

func NewContactController(useCase *usecase.ContactUseCase, log *zap.SugaredLogger) *ContactController {
	return &ContactController{
		UseCase: useCase,
		Log:     log,
	}
}

// Create godoc
// @Summary Create new contact
// @Description Create new contact
// @Tags Contact API
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body model.CreateContactRequest true "Create Contact Request"
// @Success 200 {object} model.WebResponse[model.ContactResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/contacts [post]
func (c *ContactController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.CreateContactRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Errorw("error parsing request body", "error", err)
		return fiber.ErrBadRequest
	}
	request.UserId = auth.ID

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Errorw("error creating contact", "error", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ContactResponse]{Data: response})
}

// List godoc
// @Summary List contacts
// @Description List contacts
// @Tags Contact API
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string false "Name"
// @Param email query string false "Email"
// @Param phone query string false "Phone"
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Success 200 {object} model.WebResponse[[]model.ContactResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/contacts [get]
func (c *ContactController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.SearchContactRequest{
		UserId: auth.ID,
		Name:   ctx.Query("name", ""),
		Email:  ctx.Query("email", ""),
		Phone:  ctx.Query("phone", ""),
		Page:   ctx.QueryInt("page", 1),
		Size:   ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.Errorw("error searching contact", "error", err)
		return err
	}

	paging := model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.PageResponse[model.ContactResponse]{
		Data:   responses,
		Paging: paging,
	})
}

// Get godoc
// @Summary Get contact
// @Description Get contact
// @Tags Contact API
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param contactId path string true "Contact ID"
// @Success 200 {object} model.WebResponse[model.ContactResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/contacts/{contactId} [get]
func (c *ContactController) Get(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetContactRequest{
		UserId: auth.ID,
		ID:     ctx.Params("contactId"),
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Errorw("error getting contact", "error", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ContactResponse]{Data: response})
}

// Update godoc
// @Summary Update contact
// @Description Update contact
// @Tags Contact API
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param contactId path string true "Contact ID"
// @Param request body model.UpdateContactRequest true "Update Contact Request"
// @Success 200 {object} model.WebResponse[model.ContactResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/contacts/{contactId} [put]
func (c *ContactController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateContactRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Errorw("error parsing request body", "error", err)
		return fiber.ErrBadRequest
	}

	request.UserId = auth.ID
	request.ID = ctx.Params("contactId")

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Errorw("error updating contact", "error", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ContactResponse]{Data: response})
}

// Delete godoc
// @Summary Delete contact
// @Description Delete contact
// @Tags Contact API
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param contactId path string true "Contact ID"
// @Success 200 {object} model.WebResponse[bool]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/contacts/{contactId} [delete]
func (c *ContactController) Delete(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Params("contactId")

	request := &model.DeleteContactRequest{
		UserId: auth.ID,
		ID:     contactId,
	}

	if err := c.UseCase.Delete(ctx.UserContext(), request); err != nil {
		c.Log.Errorw("error deleting contact", "error", err)
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}
