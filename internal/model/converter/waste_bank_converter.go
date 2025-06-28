package converter

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func WasteBankToResponse(wasteBank *entity.WasteBankProfile) *model.WasteBankResponse {
	var openTime, closeTime string

	if !wasteBank.OpenTime.IsZero() {
		openTime = wasteBank.OpenTime.Format("15:04:05")
	}

	if !wasteBank.CloseTime.IsZero() {
		closeTime = wasteBank.CloseTime.Format("15:04:05")
	}

	var userResponse *model.UserResponse
	if wasteBank.User.ID != uuid.Nil {
		userResponse = UserToResponse(&wasteBank.User)
	}
	return &model.WasteBankResponse{
		ID:               wasteBank.ID.String(),
		UserID:           wasteBank.UserID.String(),
		TotalWasteWeight: wasteBank.TotalWasteWeight,
		TotalWorkers:     wasteBank.TotalWorkers,
		OpenTime:         openTime,
		CloseTime:        closeTime,
		User:             userResponse,
	}
}
