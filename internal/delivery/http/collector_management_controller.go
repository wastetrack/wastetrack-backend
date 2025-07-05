package http

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http/middleware"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type CollectorManagementController struct {
	Log                        *logrus.Logger
	CollectorManagementUsecase *usecase.CollectorManagementUsecase
}

func NewCollectorManagementController(usecase *usecase.CollectorManagementUsecase, logger *logrus.Logger) *CollectorManagementController {
	return &CollectorManagementController{
		Log:                        logger,
		CollectorManagementUsecase: usecase,
	}
}

func (c *CollectorManagementController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	authRequest := &model.GetUserRequest{
		ID: auth.ID,
	}
	request := new(model.CollectorManagementRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	request.WasteBankID = authRequest.ID

	response, err := c.CollectorManagementUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create collector management: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.CollectorManagementSimpleResponse]{Data: response})
}

func (c *CollectorManagementController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	response, err := c.CollectorManagementUsecase.Get(ctx.UserContext(), id)
	if err != nil {
		c.Log.Warnf("Failed to get collector management: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.CollectorManagementResponse]{Data: response})
}

func (c *CollectorManagementController) List(ctx *fiber.Ctx) error {
	request := &model.SearchCollectorManagementRequest{
		WasteBankID: ctx.Query("waste_bank_id"),
		CollectorID: ctx.Query("collector_id"),
		Status:      ctx.Query("status"),
		Page:        ctx.QueryInt("page", 1),
		Size:        ctx.QueryInt("size", 10),
	}

	responses, total, err := c.CollectorManagementUsecase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search collector managements")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.CollectorManagementSimpleResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *CollectorManagementController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateCollectorManagementRequest)
	request.ID = ctx.Params("id")

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.CollectorManagementUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update collector management: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.CollectorManagementSimpleResponse]{Data: response})
}

func (c *CollectorManagementController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	response, err := c.CollectorManagementUsecase.Delete(ctx.UserContext(), id)
	if err != nil {
		c.Log.Warnf("Failed to delete collector management: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.CollectorManagementSimpleResponse]{Data: response})
}
