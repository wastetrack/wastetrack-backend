package converter

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func CustomerToResponse(customer *entity.CustomerProfile) *model.CustomerResponse {
	var userResponse *model.UserResponse
	if customer.User.ID != uuid.Nil {
		userResponse = UserToResponse(&customer.User)
	}
	return &model.CustomerResponse{
		ID:            customer.ID.String(),
		UserID:        customer.UserID.String(),
		CarbonDeficit: customer.CarbonDeficit,
		WaterSaved:    customer.WaterSaved,
		BagsStored:    customer.BagsStored,
		Trees:         customer.Trees,
		User:          userResponse,
	}
}
