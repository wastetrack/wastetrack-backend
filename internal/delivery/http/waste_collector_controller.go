package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http/middleware"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type WasteCollectorController struct {
	Log                   *logrus.Logger
	WasteCollectorUsecase *usecase.WasteCollectorUseCase
}

func NewWasteCollectorController(wasteCollectorUsecase *usecase.WasteCollectorUseCase, logger *logrus.Logger) *WasteCollectorController {
	return &WasteCollectorController{
		Log:                   logger,
		WasteCollectorUsecase: wasteCollectorUsecase,
	}
}

func (c *WasteCollectorController) Get(ctx *fiber.Ctx) error {
	request := &model.GetWasteCollectorRequest{
		ID: ctx.Params("user_id"),
	}
	wasteCollectorResponse, err := c.WasteCollectorUsecase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get waste collector: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.WasteCollectorResponse]{Data: wasteCollectorResponse})
}

func (c *WasteCollectorController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateWasteCollectorRequest)
	request.ID = ctx.Params("id")
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	auth := middleware.GetUser(ctx)
	if auth == nil {
		return fiber.ErrUnauthorized
	}

	wasteCollectorResponse, err := c.WasteCollectorUsecase.Update(ctx.UserContext(), request, auth.ID, auth.Role)
	if err != nil {
		c.Log.Warnf("Failed to update waste collector: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteCollectorResponse]{Data: wasteCollectorResponse})
}

func (c *WasteCollectorController) Delete(ctx *fiber.Ctx) error {
	request := &model.DeleteWasteCollectorRequest{
		ID: ctx.Params("id"),
	}
	wasteCollectorResponse, err := c.WasteCollectorUsecase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete waste collector: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.WasteCollectorResponse]{Data: wasteCollectorResponse})
}
