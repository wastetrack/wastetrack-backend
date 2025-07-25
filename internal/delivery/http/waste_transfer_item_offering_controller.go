package http

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type WasteTransferItemOfferingController struct {
	Log                              *logrus.Logger
	WasteTransferItemOfferingUsecase *usecase.WasteTransferItemOfferingUsecase
}

func NewWasteTransferItemOfferingController(usecase *usecase.WasteTransferItemOfferingUsecase, logger *logrus.Logger) *WasteTransferItemOfferingController {
	return &WasteTransferItemOfferingController{
		Log:                              logger,
		WasteTransferItemOfferingUsecase: usecase,
	}
}

func (c *WasteTransferItemOfferingController) Get(ctx *fiber.Ctx) error {
	request := &model.GetWasteTransferItemOfferingRequest{
		ID: ctx.Params("id"),
	}

	response, err := c.WasteTransferItemOfferingUsecase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get waste transfer item offering: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WasteTransferItemOfferingResponse]{Data: response})
}

func (c *WasteTransferItemOfferingController) List(ctx *fiber.Ctx) error {
	request := &model.SearchWasteTransferItemOfferingRequest{
		TransferFormID:  ctx.Query("transfer_form_id"),
		WasteCategoryID: ctx.Query("waste_category_id"),
		WasteTypeID:     ctx.Query("waste_type_id"),
		Page:            ctx.QueryInt("page"),
		Size:            ctx.QueryInt("size"),
	}

	// Set default values for pagination
	if request.Page == 0 {
		request.Page = 1
	}
	if request.Size == 0 {
		request.Size = 10
	}

	responses, total, err := c.WasteTransferItemOfferingUsecase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search waste transfer item offerings")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.WasteTransferItemOfferingSimpleResponse]{
		Data:   responses,
		Paging: paging,
	})
}
