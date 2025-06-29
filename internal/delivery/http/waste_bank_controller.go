package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http/middleware"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type WasteBankController struct {
	Log              *logrus.Logger
	WasteBankUsecase *usecase.WasteBankUseCase
}

func NewWasteBankController(wasteBankUsecase *usecase.WasteBankUseCase, logger *logrus.Logger) *WasteBankController {
	return &WasteBankController{
		Log:              logger,
		WasteBankUsecase: wasteBankUsecase,
	}
}

func (c *WasteBankController) Get(ctx *fiber.Ctx) error {
	request := &model.GetWasteBankRequest{
		ID: ctx.Params("user_id"),
	}
	wasteBankResponse, err := c.WasteBankUsecase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get waste bank: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.WasteBankResponse]{Data: wasteBankResponse})
}

func (c *WasteBankController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateWasteBankRequest)
	request.ID = ctx.Params("id")
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}
	// Extract user info from auth middleware

	auth := middleware.GetUser(ctx)
	if auth == nil {
		return fiber.ErrUnauthorized
	}

	wasteBankResponse, err := c.WasteBankUsecase.Update(ctx.UserContext(), request, auth.ID, auth.Role)
	if err != nil {
		c.Log.Warnf("Failed to update waste bank: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteBankResponse]{Data: wasteBankResponse})
}

func (c *WasteBankController) Delete(ctx *fiber.Ctx) error {
	request := &model.DeleteWasteBankRequest{
		ID: ctx.Params("id"),
	}
	wasteBankResponse, err := c.WasteBankUsecase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete waste bank: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.WasteBankResponse]{Data: wasteBankResponse})
}
