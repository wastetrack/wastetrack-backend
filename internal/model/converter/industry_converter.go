package converter

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func IndustryToResponse(industry *entity.IndustryProfile) *model.IndustryResponse {
	var userResponse *model.UserResponse
	if industry.User.ID != uuid.Nil {
		userResponse = UserToResponse(&industry.User)
	}
	return &model.IndustryResponse{
		ID:                  industry.ID.String(),
		UserID:              industry.UserID.String(),
		TotalWasteWeight:    industry.TotalWasteWeight,
		TotalRecycledWeight: industry.TotalRecycledWeight,
		User:                userResponse,
	}
}
