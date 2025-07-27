package http

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http/middleware"
	"github.com/wastetrack/wastetrack-backend/internal/helper"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type PointConversionController struct {
	Log                    *logrus.Logger
	PointConversionUsecase *usecase.PointConversionUsecase
}

func NewPointConversionController(usecase *usecase.PointConversionUsecase, logger *logrus.Logger) *PointConversionController {
	return &PointConversionController{
		Log:                    logger,
		PointConversionUsecase: usecase,
	}
}

func (c *PointConversionController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.PointConversionRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	// Set UserID from authenticated user
	request.UserID = auth.ID
	response, err := c.PointConversionUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create point conversion: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.PointConversionSimpleResponse]{Data: response})
}

func (c *PointConversionController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	response, err := c.PointConversionUsecase.Get(ctx.UserContext(), id)
	if err != nil {
		c.Log.Warnf("Failed to get point conversion: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.PointConversionResponse]{Data: response})
}

func (c *PointConversionController) List(ctx *fiber.Ctx) error {
	request := &model.SearchPointConversionRequest{
		UserID:    ctx.Query("user_id"),
		Amount:    int64(ctx.QueryInt("amount", 0)),
		Status:    ctx.Query("status"),
		IsDeleted: helper.ParseBoolQuery(ctx, "is_deleted"),
		OrderBy:   ctx.Query("order_by"),
		OrderDir:  ctx.Query("order_dir"),
		Page:      ctx.QueryInt("page", 1),
		Size:      ctx.QueryInt("size", 10),
	}

	responses, total, err := c.PointConversionUsecase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search point conversions")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.PointConversionSimpleResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *PointConversionController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdatePointConversionRequest)
	request.ID = ctx.Params("id")

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.PointConversionUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update point conversion: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.PointConversionSimpleResponse]{Data: response})
}

func (c *PointConversionController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	response, err := c.PointConversionUsecase.Delete(ctx.UserContext(), id)
	if err != nil {
		c.Log.Warnf("Failed to delete point conversion: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.PointConversionSimpleResponse]{Data: response})
}
