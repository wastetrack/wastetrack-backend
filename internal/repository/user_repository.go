package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) FindByEmail(db *gorm.DB, user *entity.User, email string) error {
	return db.Where("email = ?", email).First(user).Error
}

func (r *UserRepository) FindByEmailVerificationToken(db *gorm.DB, user *entity.User, token string) error {
	return db.Where("email_verification_token = ?", token).First(user).Error
}

func (r *UserRepository) FindByResetPasswordToken(db *gorm.DB, user *entity.User, token string) error {
	return db.Where("reset_password_token = ? AND reset_password_expiry > NOW()", token).First(user).Error
}

func (r *UserRepository) CountByEmail(db *gorm.DB, email string) (int64, error) {
	var total int64
	err := db.Model(new(entity.User)).Where("email = ?", email).Count(&total).Error
	return total, err
}

func (r *UserRepository) CountByUsername(db *gorm.DB, username string) (int64, error) {
	var total int64
	err := db.Model(new(entity.User)).Where("username = ?", username).Count(&total).Error
	return total, err
}

// func (r *UserRepository) Search(db *gorm.DB, request *model.SearchUserRequest) ([]entity.User, int64, error) {
// 	var users []entity.User
// 	if err := db.Scopes(r.FilterUser(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&users).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	var total int64 = 0
// 	if err := db.Model(&entity.User{}).Scopes(r.FilterUser(request)).Count(&total).Error; err != nil {
// 		return nil, 0, err
// 	}
// 	return users, total, nil
// }

// func (r *UserRepository) FilterUser(request *model.SearchUserRequest) func(tx *gorm.DB) *gorm.DB {
// 	return func(tx *gorm.DB) *gorm.DB {
// 		if username := request.Username; username != "" {
// 			tx = tx.Where("username LIKE ?", "%"+username+"%")
// 		}
// 		if email := request.Email; email != "" {
// 			tx = tx.Where("email LIKE ?", "%"+email+"%")
// 		}
// 		if phoneNumber := request.PhoneNumber; phoneNumber != "" {
// 			tx = tx.Where("phone_number LIKE ?", "%"+phoneNumber+"%")
// 		}
// 		if gradeLevel := request.GradeLevel; gradeLevel != 0 {
// 			tx = tx.Where("grade_level = ?", gradeLevel)
// 		}
// 		if birthDate := request.BirthDate; birthDate != nil {
// 			tx = tx.Where("birth_date = ?", birthDate)
// 		}
// 		return tx
// 	}
// }
