package http

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type WasteCategoryController struct {
	Log                  *logrus.Logger
	WasteCategoryUsecase *usecase.WasteCategoryUsecase
}

func NewWasteCategoryController(wasteCategoryUsecase *usecase.WasteCategoryUsecase, logger *logrus.Logger) *WasteCategoryController {
	return &WasteCategoryController{
		Log:                  logger,
		WasteCategoryUsecase: wasteCategoryUsecase,
	}
}

func (c *WasteCategoryController) Create(ctx *fiber.Ctx) error {
	request := new(model.WasteCategoryRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}
	response, err := c.WasteCategoryUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create waste category: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.WasteCategoryResponse]{Data: response})
}

func (c *WasteCategoryController) Get(ctx *fiber.Ctx) error {
	request := &model.GetWasteCategoryRequest{
		ID: ctx.Params("id"),
	}
	response, err := c.WasteCategoryUsecase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get waste category: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.WasteCategoryResponse]{Data: response})
}

func (c *WasteCategoryController) List(ctx *fiber.Ctx) error {
	request := &model.SearchWasteCategoryRequest{
		Name: ctx.Query("name"),
		Page: ctx.QueryInt("page"),
		Size: ctx.QueryInt("size"),
	}

	responses, total, err := c.WasteCategoryUsecase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search waste categories")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.WasteCategoryResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *WasteCategoryController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateWasteCategoryRequest)
	request.ID = ctx.Params("id")
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}
	response, err := c.WasteCategoryUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update waste category: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.WasteCategoryResponse]{Data: response})
}

func (c *WasteCategoryController) Delete(ctx *fiber.Ctx) error {
	request := &model.DeleteWasteCategoryRequest{
		ID: ctx.Params("id"),
	}
	response, err := c.WasteCategoryUsecase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete waste category: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.WasteCategoryResponse]{Data: response})
}
