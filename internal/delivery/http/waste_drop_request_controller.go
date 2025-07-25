package http

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/helper"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type WasteDropRequestController struct {
	Log                     *logrus.Logger
	WasteDropRequestUsecase *usecase.WasteDropRequestUsecase
}

func NewWasteDropRequestController(usecase *usecase.WasteDropRequestUsecase, logger *logrus.Logger) *WasteDropRequestController {
	return &WasteDropRequestController{
		Log:                     logger,
		WasteDropRequestUsecase: usecase,
	}
}

func (c *WasteDropRequestController) Create(ctx *fiber.Ctx) error {
	request := new(model.WasteDropRequestRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.WasteDropRequestUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create waste drop request: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteDropRequestSimpleResponse]{Data: response})
}

func (c *WasteDropRequestController) Get(ctx *fiber.Ctx) error {
	request := &model.GetWasteDropRequest{
		ID: ctx.Params("id"),
	}

	// Parse optional latitude and longitude query parameters
	if latStr := ctx.Query("latitude"); latStr != "" {
		if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
			request.Latitude = &lat
		} else {
			c.Log.Warnf("Invalid latitude parameter: %v", err)
			return fiber.ErrBadRequest
		}
	}

	if lngStr := ctx.Query("longitude"); lngStr != "" {
		if lng, err := strconv.ParseFloat(lngStr, 64); err == nil {
			request.Longitude = &lng
		} else {
			c.Log.Warnf("Invalid longitude parameter: %v", err)
			return fiber.ErrBadRequest
		}
	}

	response, err := c.WasteDropRequestUsecase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get waste drop request: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteDropRequestResponse]{Data: response})
}

func (c *WasteDropRequestController) List(ctx *fiber.Ctx) error {
	request := &model.SearchWasteDropRequest{
		DeliveryType:         ctx.Query("delivery_type"),
		CustomerID:           ctx.Query("customer_id"),
		WasteBankID:          ctx.Query("waste_bank_id"),
		AssignedCollectorID:  ctx.Query("assigned_collector_id"),
		Status:               ctx.Query("status"),
		AppointmentDate:      ctx.Query("appointment_date"),
		AppointmentStartTime: ctx.Query("appointment_start_time"),
		AppointmentEndTime:   ctx.Query("appointment_end_time"),
		IsDeleted:            helper.ParseBoolQuery(ctx, "is_deleted"),
		// Parse order direction for created_at
		OrderDir: ctx.Query("order_dir"), // "asc" or "desc" (default: "desc")
		Page:     ctx.QueryInt("page"),
		Size:     ctx.QueryInt("size"),
	}

	// Parse optional latitude and longitude query parameters for distance calculation
	if latStr := ctx.Query("latitude"); latStr != "" {
		if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
			request.Latitude = &lat
		} else {
			c.Log.Warnf("Invalid latitude parameter: %v", err)
			return fiber.ErrBadRequest
		}
	}

	if lngStr := ctx.Query("longitude"); lngStr != "" {
		if lng, err := strconv.ParseFloat(lngStr, 64); err == nil {
			request.Longitude = &lng
		} else {
			c.Log.Warnf("Invalid longitude parameter: %v", err)
			return fiber.ErrBadRequest
		}
	}

	// Set default values for pagination
	if request.Page == 0 {
		request.Page = 1
	}
	if request.Size == 0 {
		request.Size = 10
	}

	// Validate order direction if provided
	if request.OrderDir != "" && request.OrderDir != "asc" && request.OrderDir != "desc" {
		c.Log.Warnf("Invalid order direction: %s", request.OrderDir)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid order direction. Use 'asc' or 'desc'")
	}

	responses, total, err := c.WasteDropRequestUsecase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search waste drop requests")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.WasteDropRequestSimpleResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *WasteDropRequestController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateWasteDropRequest)
	request.ID = ctx.Params("id")

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.WasteDropRequestUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update waste drop request: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteDropRequestSimpleResponse]{Data: response})
}

func (c *WasteDropRequestController) Delete(ctx *fiber.Ctx) error {
	request := &model.DeleteWasteDropRequest{
		ID: ctx.Params("id"),
	}

	response, err := c.WasteDropRequestUsecase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete waste drop request: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteDropRequestSimpleResponse]{Data: response})
}

// Additional controller methods for specific operations
func (c *WasteDropRequestController) UpdateStatus(ctx *fiber.Ctx) error {
	request := &model.UpdateWasteDropRequest{
		ID:     ctx.Params("id"),
		Status: ctx.Query("status"),
	}

	if request.Status == "" || request.Status == "completed" {
		c.Log.Warn("Status is not valid")
		return fiber.ErrBadRequest
	}

	updateRequest := &model.UpdateWasteDropRequest{
		ID:     request.ID,
		Status: request.Status,
	}

	response, err := c.WasteDropRequestUsecase.Update(ctx.UserContext(), updateRequest)
	if err != nil {
		c.Log.Warnf("Failed to update waste drop request status: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteDropRequestSimpleResponse]{Data: response})
}

// UPDATED COMPLETE METHOD - Now handles item verification
func (c *WasteDropRequestController) Complete(ctx *fiber.Ctx) error {
	request := new(model.CompleteWasteDropRequest)
	request.ID = ctx.Params("id")

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.WasteDropRequestUsecase.Complete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to complete waste drop request: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteDropRequestSimpleResponse]{Data: response})
}

func (c *WasteDropRequestController) AssignCollector(ctx *fiber.Ctx) error {
	request := &model.UpdateWasteDropRequest{
		ID:                  ctx.Params("id"),
		AssignedCollectorID: ctx.Query("collector_id"),
	}

	if request.AssignedCollectorID == "" {
		c.Log.Warn("Collector ID is required")
		return fiber.ErrBadRequest
	}

	updateRequest := &model.UpdateWasteDropRequest{
		ID:                  request.ID,
		AssignedCollectorID: request.AssignedCollectorID,
		Status:              "assigned",
	}

	response, err := c.WasteDropRequestUsecase.Update(ctx.UserContext(), updateRequest)
	if err != nil {
		c.Log.Warnf("Failed to assign collector to waste drop request: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteDropRequestSimpleResponse]{Data: response})
}
