package converter

import (
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:              user.ID,
		Username:        user.Username,
		Email:           user.Email,
		Role:            user.Role,
		PhoneNumber:     user.PhoneNumber,
		Institution:     user.Institution,
		Address:         user.Address,
		City:            user.City,
		Province:        user.Province,
		Points:          user.Points,
		Balance:         user.Balance,
		IsEmailVerified: user.IsEmailVerified,
		Location:        parseLocation(*user.Location),
	}
}
