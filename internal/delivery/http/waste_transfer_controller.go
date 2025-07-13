package http

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type WasteTransferRequestController struct {
	Log                         *logrus.Logger
	WasteTransferRequestUsecase *usecase.WasteTransferRequestUsecase
}

func NewWasteTransferRequestController(usecase *usecase.WasteTransferRequestUsecase, logger *logrus.Logger) *WasteTransferRequestController {
	return &WasteTransferRequestController{
		Log:                         logger,
		WasteTransferRequestUsecase: usecase,
	}
}

func (c *WasteTransferRequestController) Create(ctx *fiber.Ctx) error {
	request := new(model.WasteTransferRequestRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.WasteTransferRequestUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create waste transfer request: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteTransferRequestSimpleResponse]{Data: response})
}

func (c *WasteTransferRequestController) Get(ctx *fiber.Ctx) error {
	request := &model.GetWasteTransferRequest{
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

	response, err := c.WasteTransferRequestUsecase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get waste transfer request: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteTransferRequestResponse]{Data: response})
}

func (c *WasteTransferRequestController) List(ctx *fiber.Ctx) error {
	request := &model.SearchWasteTransferRequest{
		SourceUserID:         ctx.Query("source_user_id"),
		DestinationUserID:    ctx.Query("destination_user_id"),
		AssignedCollectorID:  ctx.Query("assigned_collector_id"), // NEW: Support filtering by assigned collector
		FormType:             ctx.Query("form_type"),
		Status:               ctx.Query("status"),
		AppointmentDate:      ctx.Query("appointment_date"),
		AppointmentStartTime: ctx.Query("appointment_start_time"),
		AppointmentEndTime:   ctx.Query("appointment_end_time"),
		Page:                 ctx.QueryInt("page"),
		Size:                 ctx.QueryInt("size"),
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

	responses, total, err := c.WasteTransferRequestUsecase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search waste transfer requests")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.WasteTransferRequestSimpleResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *WasteTransferRequestController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateWasteTransferRequest)
	request.ID = ctx.Params("id")

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.WasteTransferRequestUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update waste transfer request: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteTransferRequestSimpleResponse]{Data: response})
}

func (c *WasteTransferRequestController) Delete(ctx *fiber.Ctx) error {
	request := &model.DeleteWasteTransferRequest{
		ID: ctx.Params("id"),
	}

	response, err := c.WasteTransferRequestUsecase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete waste transfer request: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteTransferRequestSimpleResponse]{Data: response})
}

func (c *WasteTransferRequestController) UpdateStatus(ctx *fiber.Ctx) error {
	request := &model.UpdateWasteTransferRequest{
		ID:     ctx.Params("id"),
		Status: ctx.Query("status"),
	}

	if request.Status == "" {
		c.Log.Warn("Status is required")
		return fiber.ErrBadRequest
	}
	if request.Status != "collecting" && request.Status != "cancelled" {
		c.Log.Warn("only collecting and cancelled are allowed")
		return fiber.NewError(fiber.StatusBadRequest, "only collecting and cancelled status are allowed")
	}

	updateRequest := &model.UpdateWasteTransferRequest{
		ID:     request.ID,
		Status: request.Status,
	}

	response, err := c.WasteTransferRequestUsecase.Update(ctx.UserContext(), updateRequest)
	if err != nil {
		c.Log.Warnf("Failed to update waste transfer request status: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteTransferRequestSimpleResponse]{Data: response})
}

func (c *WasteTransferRequestController) AssignCollectorByWasteType(ctx *fiber.Ctx) error {
	request := new(model.AssignCollectorByWasteTypeRequest)
	request.ID = ctx.Params("id")

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	// Ensure the ID from params is used
	request.ID = ctx.Params("id")

	response, err := c.WasteTransferRequestUsecase.AssignCollectorByWasteType(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to assign collector by waste type: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteTransferRequestSimpleResponse]{Data: response})
}
