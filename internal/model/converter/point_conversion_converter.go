package converter

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func PointConversionToSimpleResponse(pointConversion *entity.PointConversion) *model.PointConversionSimpleResponse {
	return &model.PointConversionSimpleResponse{
		ID:        pointConversion.ID.String(),
		UserID:    pointConversion.UserID.String(),
		Amount:    pointConversion.Amount,
		CreatedAt: pointConversion.CreatedAt.String(),
		Status:    pointConversion.Status,
		IsDeleted: pointConversion.IsDeleted,
	}
}

func PointConversionToResponse(pointConversion *entity.PointConversion) *model.PointConversionResponse {
	var user *model.UserResponse
	if pointConversion.UserID != uuid.Nil {
		user = UserToResponse(&pointConversion.User)
	}
	return &model.PointConversionResponse{
		ID:        pointConversion.ID.String(),
		UserID:    pointConversion.UserID.String(),
		Amount:    pointConversion.Amount,
		CreatedAt: pointConversion.CreatedAt.String(),
		Status:    pointConversion.Status,
		IsDeleted: pointConversion.IsDeleted,
		User:      user,
	}
}
