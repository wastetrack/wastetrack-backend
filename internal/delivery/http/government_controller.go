package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type GovernmentController struct {
	Log               *logrus.Logger
	GovernmentUsecase *usecase.GovernmentUseCase
}

func NewGovernmentController(governmentUsecase *usecase.GovernmentUseCase, logger *logrus.Logger) *GovernmentController {
	return &GovernmentController{
		Log:               logger,
		GovernmentUsecase: governmentUsecase,
	}
}

// GetDashboard handles GET /api/government/dashboard
func (c *GovernmentController) GetDashboard(ctx *fiber.Ctx) error {
	request := &model.GovernmentDashboardRequest{
		StartMonth: ctx.Query("start_month"),
		EndMonth:   ctx.Query("end_month"),
		Province:   ctx.Query("province"),
		City:       ctx.Query("city"),
	}

	response, err := c.GovernmentUsecase.GetDashboard(request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to get government dashboard")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.GovernmentDashboardResponse]{
		Data: response,
	})
}
