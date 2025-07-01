package http

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type WasteTypeController struct {
	Log              *logrus.Logger
	WasteTypeUsecase *usecase.WasteTypeUsecase
}

func NewWasteTypeController(usecase *usecase.WasteTypeUsecase, logger *logrus.Logger) *WasteTypeController {
	return &WasteTypeController{
		Log:              logger,
		WasteTypeUsecase: usecase,
	}
}

func (c *WasteTypeController) Create(ctx *fiber.Ctx) error {
	request := new(model.WasteTypeRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.WasteTypeUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create waste type: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteTypeResponse]{Data: response})
}

func (c *WasteTypeController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	response, err := c.WasteTypeUsecase.Get(ctx.UserContext(), id)
	if err != nil {
		c.Log.Warnf("Failed to get waste type: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteTypeResponse]{Data: response})
}

func (c *WasteTypeController) List(ctx *fiber.Ctx) error {
	request := &model.SearchWasteTypeRequest{
		Name:          ctx.Query("name"),
		CategoryID:    ctx.Query("category_id"),
		SubCategoryID: ctx.Query("subcategory_id"),
		Page:          ctx.QueryInt("page", 1),
		Size:          ctx.QueryInt("size", 10),
	}

	responses, total, err := c.WasteTypeUsecase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search waste types")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.WasteTypeResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *WasteTypeController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateWasteTypeRequest)
	request.ID = ctx.Params("id")

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.WasteTypeUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update waste type: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteTypeResponse]{Data: response})
}

func (c *WasteTypeController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	response, err := c.WasteTypeUsecase.Delete(ctx.UserContext(), id)
	if err != nil {
		c.Log.Warnf("Failed to delete waste type: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteTypeResponse]{Data: response})
}
