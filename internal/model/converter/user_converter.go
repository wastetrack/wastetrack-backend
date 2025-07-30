// internal/model/converter/user_converter.go - Updated to include distance

package converter

import (
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	var location *model.LocationResponse
	if user.Location != nil {
		location = &model.LocationResponse{
			Latitude:  user.Location.Lat,
			Longitude: user.Location.Lng,
		}
	}

	return &model.UserResponse{
		ID:                  user.ID,
		Username:            user.Username,
		Email:               user.Email,
		Role:                user.Role,
		AvatarURL:           user.AvatarURL,
		PhoneNumber:         user.PhoneNumber,
		Institution:         user.Institution,
		Address:             user.Address,
		City:                user.City,
		Province:            user.Province,
		Points:              user.Points,
		Balance:             user.Balance,
		IsEmailVerified:     user.IsEmailVerified,
		IsAcceptingCustomer: user.IsAcceptingCustomer,
		Location:            location,
		CreatedAt:           &user.CreatedAt,
		UpdatedAt:           &user.UpdatedAt,
		Distance:            user.Distance, // Now directly accessible from entity
		IsAgreeedToTerms:    user.IsAgreeedToTerms,
	}
}

func UserToListResponse(user *entity.User) *model.UserListResponse {
	var location *model.LocationResponse
	if user.Location != nil {
		location = &model.LocationResponse{
			Latitude:  user.Location.Lat,
			Longitude: user.Location.Lng,
		}
	}

	return &model.UserListResponse{
		ID:                  user.ID,
		Username:            user.Username,
		Email:               user.Email,
		Role:                user.Role,
		AvatarURL:           user.AvatarURL,
		PhoneNumber:         user.PhoneNumber,
		Institution:         user.Institution,
		Address:             user.Address,
		City:                user.City,
		Province:            user.Province,
		Points:              user.Points,
		Balance:             user.Balance,
		IsEmailVerified:     user.IsEmailVerified,
		IsAcceptingCustomer: user.IsAcceptingCustomer,
		Location:            location,
		CreatedAt:           &user.CreatedAt,
		UpdatedAt:           &user.UpdatedAt,
		Distance:            user.Distance,
		IsAgreeedToTerms:    user.IsAgreeedToTerms,
	}
}
