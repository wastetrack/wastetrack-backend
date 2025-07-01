package http

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http/middleware"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type WasteBankPricedTypeController struct {
	Log                        *logrus.Logger
	WasteBankPricedTypeUsecase *usecase.WasteBankPricedTypeUsecase
}

func NewWasteBankPricedTypeController(wptUsecase *usecase.WasteBankPricedTypeUsecase, logger *logrus.Logger) *WasteBankPricedTypeController {
	return &WasteBankPricedTypeController{
		Log:                        logger,
		WasteBankPricedTypeUsecase: wptUsecase,
	}
}
func (c *WasteBankPricedTypeController) CreateBatch(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	var requests []model.WasteBankPricedTypeRequest
	if err := ctx.BodyParser(&requests); err != nil {
		c.Log.Warnf("Failed to parse batch create request: %v", err)
		return fiber.ErrBadRequest
	}

	for i := range requests {
		requests[i].WasteBankID = auth.ID
	}

	responses, err := c.WasteBankPricedTypeUsecase.CreateBatch(ctx.UserContext(), requests)
	if err != nil {
		c.Log.Warnf("Failed to create batch priced types: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[[]*model.WasteBankPricedTypeResponse]{Data: responses})
}

func (c *WasteBankPricedTypeController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	authRequest := &model.GetUserRequest{
		ID: auth.ID,
	}
	c.Log.Infof("Creating priced type for WasteBankID: %s", authRequest.ID)
	request := new(model.WasteBankPricedTypeRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse create request: %v", err)
		return fiber.ErrBadRequest
	}

	request.WasteBankID = authRequest.ID

	result, err := c.WasteBankPricedTypeUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create priced type: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteBankPricedTypeResponse]{Data: result})
}

func (c *WasteBankPricedTypeController) List(ctx *fiber.Ctx) error {
	request := &model.SearchWasteBankPricedTypeRequest{
		WasteBankID: ctx.Query("waste_bank_id"),
		WasteTypeID: ctx.Query("waste_type_id"),
		Page:        ctx.QueryInt("page"),
		Size:        ctx.QueryInt("size"),
	}

	responses, total, err := c.WasteBankPricedTypeUsecase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search waste bank priced types")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.WasteBankPricedTypeSimpleResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *WasteBankPricedTypeController) Get(ctx *fiber.Ctx) error {
	request := &model.GetWasteBankPricedTypeRequest{
		ID: ctx.Params("id"),
	}

	result, err := c.WasteBankPricedTypeUsecase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get priced type: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteBankPricedTypeResponse]{Data: result})
}

func (c *WasteBankPricedTypeController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateWasteBankPricedTypeRequest)
	request.ID = ctx.Params("id")
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse update request: %v", err)
		return fiber.ErrBadRequest
	}

	result, err := c.WasteBankPricedTypeUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update priced type: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteBankPricedTypeResponse]{Data: result})
}

func (c *WasteBankPricedTypeController) Delete(ctx *fiber.Ctx) error {
	request := &model.DeleteWasteBankPricedTypeRequest{
		ID: ctx.Params("id"),
	}

	result, err := c.WasteBankPricedTypeUsecase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete priced type: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteBankPricedTypeResponse]{Data: result})
}
