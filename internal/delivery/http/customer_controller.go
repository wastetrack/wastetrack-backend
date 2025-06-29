package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http/middleware"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type CustomerController struct {
	Log             *logrus.Logger
	CustomerUsecase *usecase.CustomerUseCase
}

func NewCustomerController(customerUsecase *usecase.CustomerUseCase, logger *logrus.Logger) *CustomerController {
	return &CustomerController{
		Log:             logger,
		CustomerUsecase: customerUsecase,
	}
}

func (c *CustomerController) Get(ctx *fiber.Ctx) error {
	request := &model.GetCustomerRequest{
		ID: ctx.Params("user_id"),
	}
	customerResponse, err := c.CustomerUsecase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get customer: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.CustomerResponse]{Data: customerResponse})
}

func (c *CustomerController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateCustomerRequest)
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

	customerResponse, err := c.CustomerUsecase.Update(ctx.UserContext(), request, auth.ID, auth.Role)
	if err != nil {
		c.Log.Warnf("Failed to update customer: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.CustomerResponse]{Data: customerResponse})
}

func (c *CustomerController) Delete(ctx *fiber.Ctx) error {
	request := &model.DeleteWasteBankRequest{
		ID: ctx.Params("id"),
	}
	customerResponse, err := c.CustomerUsecase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete customer: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.CustomerResponse]{Data: customerResponse})
}
