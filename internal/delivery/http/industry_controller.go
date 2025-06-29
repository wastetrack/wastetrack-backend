package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http/middleware"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type IndustryController struct {
	Log             *logrus.Logger
	IndustryUsecase *usecase.IndustryUseCase
}

func NewIndustryController(industryUsecase *usecase.IndustryUseCase, logger *logrus.Logger) *IndustryController {
	return &IndustryController{
		Log:             logger,
		IndustryUsecase: industryUsecase,
	}
}

func (c *IndustryController) Get(ctx *fiber.Ctx) error {
	request := &model.GetIndustryRequest{
		ID: ctx.Params("user_id"),
	}
	industryResponse, err := c.IndustryUsecase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get industry profile: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.IndustryResponse]{Data: industryResponse})
}

func (c *IndustryController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateIndustryRequest)
	request.ID = ctx.Params("id")
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	auth := middleware.GetUser(ctx)
	if auth == nil {
		return fiber.ErrUnauthorized
	}

	industryResponse, err := c.IndustryUsecase.Update(ctx.UserContext(), request, auth.ID, auth.Role)
	if err != nil {
		c.Log.Warnf("Failed to update industry profile: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.IndustryResponse]{Data: industryResponse})
}

func (c *IndustryController) Delete(ctx *fiber.Ctx) error {
	request := &model.DeleteIndustryRequest{
		ID: ctx.Params("id"),
	}
	industryResponse, err := c.IndustryUsecase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete industry profile: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.IndustryResponse]{Data: industryResponse})
}
