package http

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http/middleware"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type StorageController struct {
	Log            *logrus.Logger
	StorageUsecase *usecase.StorageUsecase
}

func NewStorageController(usecase *usecase.StorageUsecase, logger *logrus.Logger) *StorageController {
	return &StorageController{
		Log:            logger,
		StorageUsecase: usecase,
	}
}

func (c *StorageController) Create(ctx *fiber.Ctx) error {
	request := new(model.StorageRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.StorageUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create storage: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.StorageSimpleResponse]{Data: response})
}

func (c *StorageController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	response, err := c.StorageUsecase.Get(ctx.UserContext(), id)
	if err != nil {
		c.Log.Warnf("Failed to get storage: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.StorageResponse]{Data: response})
}

func (c *StorageController) List(ctx *fiber.Ctx) error {
	var (
		page = ctx.QueryInt("page", 1)
		size = ctx.QueryInt("size", 10)
	)

	request := &model.SearchStorageRequest{
		UserID: ctx.Query("user_id"),
		Page:   page,
		Size:   size,
	}

	// Optional boolean and float filters from query
	if isForRecycled := ctx.Query("is_for_recycled_material"); isForRecycled != "" {
		val := isForRecycled == "true"
		request.IsForRecycledMaterial = &val
	}
	if minL := ctx.QueryFloat("min_length"); minL != 0 {
		request.MinLength = &minL
	}
	if maxL := ctx.QueryFloat("max_length"); maxL != 0 {
		request.MaxLength = &maxL
	}
	if minW := ctx.QueryFloat("min_width"); minW != 0 {
		request.MinWidth = &minW
	}
	if maxW := ctx.QueryFloat("max_width"); maxW != 0 {
		request.MaxWidth = &maxW
	}
	if minH := ctx.QueryFloat("min_height"); minH != 0 {
		request.MinHeight = &minH
	}
	if maxH := ctx.QueryFloat("max_height"); maxH != 0 {
		request.MaxHeight = &maxH
	}

	responses, total, err := c.StorageUsecase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search storages")
		return err
	}

	paging := &model.PageMetadata{
		Page:      page,
		Size:      size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(size))),
	}

	return ctx.JSON(model.WebResponse[[]model.StorageSimpleResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *StorageController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	authRequest := &model.GetUserRequest{
		ID: auth.ID,
	}
	request := new(model.UpdateStorageRequest)
	request.ID = ctx.Params("id")
	request.UserID = authRequest.ID
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.StorageUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update storage: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.StorageSimpleResponse]{Data: response})
}

func (c *StorageController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	response, err := c.StorageUsecase.Delete(ctx.UserContext(), id)
	if err != nil {
		c.Log.Warnf("Failed to delete storage: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.StorageSimpleResponse]{Data: response})
}
