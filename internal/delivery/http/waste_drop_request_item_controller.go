package http

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type WasteDropRequestItemController struct {
	Log                         *logrus.Logger
	WasteDropRequestItemUsecase *usecase.WasteDropRequestItemUsecase
}

func NewWasteDropRequestItemController(usecase *usecase.WasteDropRequestItemUsecase, logger *logrus.Logger) *WasteDropRequestItemController {
	return &WasteDropRequestItemController{
		Log:                         logger,
		WasteDropRequestItemUsecase: usecase,
	}
}

func (c *WasteDropRequestItemController) Create(ctx *fiber.Ctx) error {
	request := new(model.WasteDropRequestItemRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.WasteDropRequestItemUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create waste drop request item: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteDropRequestItemSimpleResponse]{Data: response})
}

func (c *WasteDropRequestItemController) Get(ctx *fiber.Ctx) error {
	request := &model.GetWasteDropRequestItemRequest{
		ID: ctx.Params("id"),
	}

	response, err := c.WasteDropRequestItemUsecase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get waste drop request item: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteDropRequestItemResponse]{Data: response})
}

func (c *WasteDropRequestItemController) List(ctx *fiber.Ctx) error {
	request := &model.SearchWasteDropRequestItemRequest{
		RequestID:   ctx.Query("request_id"),
		WasteTypeID: ctx.Query("waste_type_id"),
		Page:        ctx.QueryInt("page"),
		Size:        ctx.QueryInt("size"),
	}

	responses, total, err := c.WasteDropRequestItemUsecase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search waste drop request items")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.WasteDropRequestItemSimpleResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *WasteDropRequestItemController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateWasteDropRequestItemRequest)
	request.ID = ctx.Params("id")

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.WasteDropRequestItemUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update waste drop request item: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteDropRequestItemSimpleResponse]{Data: response})
}

func (c *WasteDropRequestItemController) Delete(ctx *fiber.Ctx) error {
	request := &model.DeleteWasteDropRequestItemRequest{
		ID: ctx.Params("id"),
	}

	response, err := c.WasteDropRequestItemUsecase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete waste drop request item: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteDropRequestItemSimpleResponse]{Data: response})
}
