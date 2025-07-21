package http

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http/middleware"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type StorageItemController struct {
	Log                *logrus.Logger
	StorageItemUsecase *usecase.StorageItemUsecase
}

func NewStorageItemController(usecase *usecase.StorageItemUsecase, logger *logrus.Logger) *StorageItemController {
	return &StorageItemController{
		Log:                logger,
		StorageItemUsecase: usecase,
	}
}

func (c *StorageItemController) Create(ctx *fiber.Ctx) error {
	request := new(model.StorageItemRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.StorageItemUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create storage item: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.StorageItemSimpleResponse]{Data: response})
}

func (c *StorageItemController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	response, err := c.StorageItemUsecase.Get(ctx.UserContext(), id)
	if err != nil {
		c.Log.Warnf("Failed to get storage item: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.StorageItemResponse]{Data: response})
}

func (c *StorageItemController) List(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	size := ctx.QueryInt("size", 10)

	request := &model.SearchStorageItemRequest{
		StorageID:        ctx.Query("storage_id"),
		WasteTypeID:      ctx.Query("waste_type_id"),
		OrderByWeightKgs: ctx.Query("order_by_weight_kgs"),
		Page:             page,
		Size:             size,
	}

	responses, total, err := c.StorageItemUsecase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warn("Failed to search storage items")
		return err
	}

	paging := &model.PageMetadata{
		Page:      page,
		Size:      size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(size))),
	}

	return ctx.JSON(model.WebResponse[[]model.StorageItemSimpleResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *StorageItemController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	authRequest := &model.GetUserRequest{
		ID: auth.ID,
	}
	request := new(model.UpdateStorageItemRequest)
	request.ID = ctx.Params("id")
	request.UserID = authRequest.ID

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.StorageItemUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update storage item: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.StorageItemSimpleResponse]{Data: response})
}

func (c *StorageItemController) DeductStorageItem(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	authRequest := &model.GetUserRequest{
		ID: auth.ID,
	}
	request := new(model.DeductStorageItemRequest)
	request.ID = ctx.Params("id")
	request.UserID = authRequest.ID

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.StorageItemUsecase.DeductFromStorage(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update storage item: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.StorageItemSimpleResponse]{Data: response})
}

func (c *StorageItemController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	response, err := c.StorageItemUsecase.Delete(ctx.UserContext(), id)
	if err != nil {
		c.Log.Warnf("Failed to delete storage item: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.StorageItemSimpleResponse]{Data: response})
}
