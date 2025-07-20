package http

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http/middleware"
	"github.com/wastetrack/wastetrack-backend/internal/helper"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type SalaryTransactionController struct {
	Log                      *logrus.Logger
	SalaryTransactionUsecase *usecase.SalaryTransactionUsecase
}

func NewSalaryTransactionController(usecase *usecase.SalaryTransactionUsecase, logger *logrus.Logger) *SalaryTransactionController {
	return &SalaryTransactionController{
		Log:                      logger,
		SalaryTransactionUsecase: usecase,
	}
}

func (c *SalaryTransactionController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	authRequest := &model.GetUserRequest{
		ID: auth.ID,
	}
	request := new(model.SalaryTransactionRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	request.SenderID = authRequest.ID
	response, err := c.SalaryTransactionUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create salary transaction: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.SalaryTransactionSimpleResponse]{Data: response})
}

func (c *SalaryTransactionController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	response, err := c.SalaryTransactionUsecase.Get(ctx.UserContext(), id)
	if err != nil {
		c.Log.Warnf("Failed to get salary transaction: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.SalaryTransactionResponse]{Data: response})
}

func (c *SalaryTransactionController) List(ctx *fiber.Ctx) error {
	request := &model.SearchSalaryTransactionRequest{
		SenderID:        ctx.Query("sender_id"),
		ReceiverID:      ctx.Query("receiver_id"),
		TransactionType: ctx.Query("transaction_type"),
		Status:          ctx.Query("status"),
		IsDeleted:       helper.ParseBoolQuery(ctx, "is_deleted"),
		Page:            ctx.QueryInt("page", 1),
		Size:            ctx.QueryInt("size", 10),
	}

	responses, total, err := c.SalaryTransactionUsecase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search salary transactions")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.SalaryTransactionSimpleResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *SalaryTransactionController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateSalaryTransactionRequest)
	request.ID = ctx.Params("id")

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.SalaryTransactionUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update salary transaction: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.SalaryTransactionSimpleResponse]{Data: response})
}

func (c *SalaryTransactionController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	response, err := c.SalaryTransactionUsecase.Delete(ctx.UserContext(), id)
	if err != nil {
		c.Log.Warnf("Failed to delete salary transaction: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.SalaryTransactionSimpleResponse]{Data: response})
}
