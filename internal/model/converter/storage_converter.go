package converter

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func StorageToSimpleResponse(storage *entity.Storage) *model.StorageSimpleResponse {
	return &model.StorageSimpleResponse{
		ID:                    storage.ID.String(),
		UserID:                storage.UserID.String(),
		Length:                storage.Length,
		Width:                 storage.Width,
		Height:                storage.Height,
		IsForRecycledMaterial: storage.IsForRecycledMaterial,
	}
}

func StorageToResponse(storage *entity.Storage) *model.StorageResponse {
	var user *model.UserResponse
	if storage.UserID != uuid.Nil {
		user = UserToResponse(&storage.User)
	}
	return &model.StorageResponse{
		ID:                    storage.ID.String(),
		UserID:                storage.UserID.String(),
		Length:                storage.Length,
		Width:                 storage.Width,
		Height:                storage.Height,
		User:                  user,
		IsForRecycledMaterial: storage.IsForRecycledMaterial,
	}
}
