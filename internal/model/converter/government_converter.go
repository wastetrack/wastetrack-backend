package converter

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func GovernmentToResponse(government *entity.GovernmentProfile) *model.GovernmentResponse {
	var userResponse *model.UserResponse
	if government.User.ID != uuid.Nil {
		userResponse = UserToResponse(&government.User)
	}
	return &model.GovernmentResponse{
		ID:     government.ID.String(),
		UserID: government.UserID.String(),
		User:   userResponse,
	}
}
func GovernmentDashboardToResponse(dashboard *model.GovernmentDashboardResponse) *model.GovernmentDashboardResponse {
	if dashboard == nil {
		return &model.GovernmentDashboardResponse{}
	}

	return &model.GovernmentDashboardResponse{
		TotalBankSampah:  dashboard.TotalBankSampah,
		TotalOfftaker:    dashboard.TotalOfftaker,
		TotalCollected:   dashboard.TotalCollected,
		CollectionTrends: dashboard.CollectionTrends,
		TopOfftakers:     dashboard.TopOfftakers,
		LargestBanks:     dashboard.LargestBanks,
	}
}
