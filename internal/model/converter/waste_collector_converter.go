package converter

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func WasteCollectorToResponse(wasteCollector *entity.WasteCollectorProfile) *model.WasteCollectorResponse {
	var userResponse *model.UserResponse
	if wasteCollector.User.ID != uuid.Nil {
		userResponse = UserToResponse(&wasteCollector.User)
	}
	return &model.WasteCollectorResponse{
		ID:               wasteCollector.ID.String(),
		UserID:           wasteCollector.UserID.String(),
		TotalWasteWeight: wasteCollector.TotalWasteWeight,
		User:             userResponse,
	}
}
