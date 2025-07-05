package converter

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func CollectorManagementToSimpleResponse(collectorManagement *entity.CollectorManagement) *model.CollectorManagementSimpleResponse {
	return &model.CollectorManagementSimpleResponse{
		ID:          collectorManagement.ID.String(),
		WasteBankID: collectorManagement.WasteBankID.String(),
		CollectorID: collectorManagement.CollectorID.String(),
		Status:      collectorManagement.Status,
	}
}

func CollectorManagementToResponse(collectorManagement *entity.CollectorManagement) *model.CollectorManagementResponse {
	var wasteBank *model.UserResponse
	if collectorManagement.WasteBankID != uuid.Nil {
		wasteBank = UserToResponse(&collectorManagement.WasteBank)
	}
	var collector *model.UserResponse
	if collectorManagement.CollectorID != uuid.Nil {
		collector = UserToResponse(&collectorManagement.Collector)
	}
	return &model.CollectorManagementResponse{
		ID:          collectorManagement.ID.String(),
		WasteBankID: collectorManagement.WasteBankID.String(),
		CollectorID: collectorManagement.CollectorID.String(),
		Status:      collectorManagement.Status,
		WasteBank:   wasteBank,
		Collector:   collector,
	}
}
