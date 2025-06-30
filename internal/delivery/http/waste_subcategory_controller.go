package http

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type WasteSubCategoryController struct {
	Log                     *logrus.Logger
	WasteSubCategoryUsecase *usecase.WasteSubCategoryUsecase
}

func NewWasteSubCategoryController(usecase *usecase.WasteSubCategoryUsecase, logger *logrus.Logger) *WasteSubCategoryController {
	return &WasteSubCategoryController{
		Log:                     logger,
		WasteSubCategoryUsecase: usecase,
	}
}

func (c *WasteSubCategoryController) Create(ctx *fiber.Ctx) error {
	request := new(model.WasteSubCategoryRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.WasteSubCategoryUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create waste subcategory: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteSubCategoryResponse]{Data: response})
}

func (c *WasteSubCategoryController) Get(ctx *fiber.Ctx) error {
	request := &model.GetWasteCategoryRequest{
		ID: ctx.Params("id"),
	}

	response, err := c.WasteSubCategoryUsecase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get waste subcategory: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteSubCategoryResponse]{Data: response})
}

func (c *WasteSubCategoryController) List(ctx *fiber.Ctx) error {
	request := &model.SearchWasteSubCategoryRequest{
		Name:       ctx.Query("name"),
		CategoryID: ctx.Query("category_id"),
		Page:       ctx.QueryInt("page"),
		Size:       ctx.QueryInt("size"),
	}

	responses, total, err := c.WasteSubCategoryUsecase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search waste subcategories")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.WasteSubCategoryResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *WasteSubCategoryController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateWasteSubCategoryRequest)
	request.ID = ctx.Params("id")

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.WasteSubCategoryUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update waste subcategory: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteSubCategoryResponse]{Data: response})
}

func (c *WasteSubCategoryController) Delete(ctx *fiber.Ctx) error {
	request := &model.DeleteWasteSubCategoryRequest{
		ID: ctx.Params("id"),
	}

	response, err := c.WasteSubCategoryUsecase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete waste subcategory: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteSubCategoryResponse]{Data: response})
}
