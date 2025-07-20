package repository

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
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

func (r *UserRepository) Search(db *gorm.DB, request *model.SearchUserRequest) ([]entity.User, int64, error) {
	var users []entity.User

	// Build the query with distance calculation if coordinates provided
	query := db.Scopes(r.FilterUser(request))

	// If latitude and longitude are provided, calculate distance and order by it
	if request.Latitude != nil && request.Longitude != nil {
		distanceSelect := fmt.Sprintf(`*, 
			CASE 
				WHEN location IS NOT NULL THEN 
					ST_Distance(
						location, 
						ST_SetSRID(ST_MakePoint(%f, %f), 4326)
					)
				ELSE NULL 
			END as distance`,
			*request.Longitude, *request.Latitude)

		// Get radius in meters (default 10km = 10000m if not specified)
		radiusMeters := 10000
		if request.RadiusMeters != nil && *request.RadiusMeters > 0 {
			radiusMeters = *request.RadiusMeters
		}

		// Add distance filter for specified radius
		distanceFilter := fmt.Sprintf(`(location IS NULL OR ST_Distance(
			location, 
			ST_SetSRID(ST_MakePoint(%f, %f), 4326)
		) <= %d)`, *request.Longitude, *request.Latitude, radiusMeters)

		query = query.Select(distanceSelect).
			Where(distanceFilter).
			Order("distance ASC NULLS LAST")
	}

	// Apply pagination and execute query
	if err := query.Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// Count total records with same filters (without distance calculation for performance)
	var total int64 = 0
	countQuery := db.Model(&entity.User{}).Scopes(r.FilterUser(request))

	// Apply same distance filter for count when coordinates provided
	if request.Latitude != nil && request.Longitude != nil {
		radiusMeters := 10000
		if request.RadiusMeters != nil && *request.RadiusMeters > 0 {
			radiusMeters = *request.RadiusMeters
		}

		distanceFilter := fmt.Sprintf(`(location IS NULL OR ST_Distance(
			location, 
			ST_SetSRID(ST_MakePoint(%f, %f), 4326)
		) <= %d)`, *request.Longitude, *request.Latitude, radiusMeters)
		countQuery = countQuery.Where(distanceFilter)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) FilterUser(request *model.SearchUserRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if username := request.Username; username != "" {
			tx = tx.Where("username ILIKE ?", "%"+username+"%")
		}
		if email := request.Email; email != "" {
			tx = tx.Where("email ILIKE ?", "%"+email+"%")
		}
		if role := request.Role; role != "" {
			tx = tx.Where("role = ?", role)
		}
		if institution := request.Institution; institution != "" {
			tx = tx.Where("institution ILIKE ?", "%"+institution+"%")
		}
		if address := request.Address; address != "" {
			tx = tx.Where("address ILIKE ?", "%"+address+"%")
		}
		if city := request.City; city != "" {
			tx = tx.Where("city ILIKE ?", "%"+city+"%")
		}
		if province := request.Province; province != "" {
			tx = tx.Where("province ILIKE ?", "%"+province+"%")
		}
		// Fixed: Now properly handle the boolean filter
		if request.IsAcceptingCustomer != nil {
			tx = tx.Where("is_accepting_customer = ?", *request.IsAcceptingCustomer)
		}
		return tx
	}
}

// FindByIDWithDistance finds a user by ID and calculates distance if coordinates are provided
func (r *UserRepository) FindByIDWithDistance(db *gorm.DB, user *entity.User, id string, lat, lng *float64) error {
	query := db.Where("id = ?", id)

	if lat != nil && lng != nil {
		distanceSelect := fmt.Sprintf(`*, 
			CASE 
				WHEN location IS NOT NULL THEN 
					ST_Distance(
						location, 
						ST_SetSRID(ST_MakePoint(%f, %f), 4326)
					)
				ELSE NULL 
			END as distance`,
			*lng, *lat)

		// Also apply 10km radius filter for single user lookup
		distanceFilter := fmt.Sprintf(`(location IS NULL OR ST_Distance(
			location, 
			ST_SetSRID(ST_MakePoint(%f, %f), 4326)
		) <= 10000)`, *lng, *lat)

		query = query.Select(distanceSelect).Where(distanceFilter)
	}

	return query.First(user).Error
}
