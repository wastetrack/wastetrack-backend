package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
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

func (r *UserRepository) Search(db *gorm.DB, request *model.SearchUserRequest) ([]entity.User, map[string]*entity.CustomerProfile, map[string]*entity.WasteBankProfile, map[string]*entity.IndustryProfile, map[string]*entity.GovernmentProfile, int64, error) {
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
		return nil, nil, nil, nil, nil, 0, err
	}

	// Count total records with same filters
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
		return nil, nil, nil, nil, nil, 0, err
	}

	// Collect user IDs for profile loading
	var userIDs []string
	userRoleMap := make(map[string]string)

	for _, user := range users {
		userID := user.ID.String()
		userIDs = append(userIDs, userID)
		userRoleMap[userID] = user.Role
	}

	// Load profiles based on roles
	customerProfiles := make(map[string]*entity.CustomerProfile)
	wasteBankProfiles := make(map[string]*entity.WasteBankProfile)
	industryProfiles := make(map[string]*entity.IndustryProfile)
	governmentProfiles := make(map[string]*entity.GovernmentProfile)

	if len(userIDs) > 0 {
		// Load customer profiles
		var customers []entity.CustomerProfile
		if err := db.Where("user_id IN ?", userIDs).Find(&customers).Error; err != nil {
			r.Log.Warnf("Failed to load customer profiles: %v", err)
		} else {
			for _, customer := range customers {
				customerProfiles[customer.UserID.String()] = &customer
			}
		}

		// Load waste bank profiles
		var wasteBanks []entity.WasteBankProfile
		if err := db.Where("user_id IN ?", userIDs).Find(&wasteBanks).Error; err != nil {
			r.Log.Warnf("Failed to load waste bank profiles: %v", err)
		} else {
			for _, wasteBank := range wasteBanks {
				wasteBankProfiles[wasteBank.UserID.String()] = &wasteBank
			}
		}

		// Load industry profiles
		var industries []entity.IndustryProfile
		if err := db.Where("user_id IN ?", userIDs).Find(&industries).Error; err != nil {
			r.Log.Warnf("Failed to load industry profiles: %v", err)
		} else {
			for _, industry := range industries {
				industryProfiles[industry.UserID.String()] = &industry
			}
		}

		// Load government profiles
		var governments []entity.GovernmentProfile
		if err := db.Where("user_id IN ?", userIDs).Find(&governments).Error; err != nil {
			r.Log.Warnf("Failed to load government profiles: %v", err)
		} else {
			for _, government := range governments {
				governmentProfiles[government.UserID.String()] = &government
			}
		}
	}

	return users, customerProfiles, wasteBankProfiles, industryProfiles, governmentProfiles, total, nil
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

func (r *UserRepository) FindByEmailChangeToken(db *gorm.DB, user *entity.User, token string) error {
	return db.Where("email_change_token = ? AND email_change_expiry > NOW()", token).First(user).Error
}

func (r *UserRepository) CountByEmailExcludingUser(db *gorm.DB, email string, userID uuid.UUID) (int64, error) {
	var total int64
	err := db.Model(new(entity.User)).Where("email = ? AND id != ?", email, userID).Count(&total).Error
	return total, err
}

// Government
func (r *UserRepository) CountWasteBanks(db *gorm.DB, request *model.GovernmentDashboardRequest) (int64, error) {
	var count int64

	query := db.Model(&entity.User{}).
		Where("role IN ?", []string{"waste_bank_unit", "waste_bank_central"})

	if request.EndMonth != "" {
		endDate, err := time.Parse("2006-01", request.EndMonth)
		if err != nil {
			return 0, fmt.Errorf("invalid end_month format: %v", err)
		}
		endDate = endDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		query = query.Where("created_at <= ?", endDate)
	}

	// Location filters
	if request.Province != "" {
		query = query.Where("province ILIKE ?", "%"+request.Province+"%")
	}

	if request.City != "" {
		query = query.Where("city ILIKE ?", "%"+request.City+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count waste banks: %v", err)
	}

	return count, nil
}

func (r *UserRepository) CountOfftakers(db *gorm.DB, request *model.GovernmentDashboardRequest) (int64, error) {
	var count int64

	query := db.Model(&entity.User{}).
		Where("role = ?", "industry")

	// Apply date filters only if provided
	if request.EndMonth != "" {
		endDate, err := time.Parse("2006-01", request.EndMonth)
		if err != nil {
			return 0, fmt.Errorf("invalid end_month format: %v", err)
		}
		endDate = endDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		query = query.Where("created_at <= ?", endDate)
	}

	// Location filters
	if request.Province != "" {
		query = query.Where("province ILIKE ?", "%"+request.Province+"%")
	}

	if request.City != "" {
		query = query.Where("city ILIKE ?", "%"+request.City+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count offtakers: %v", err)
	}

	return count, nil
}

func (r *UserRepository) GetWasteBankUsers(db *gorm.DB, request *model.GovernmentDashboardRequest) (map[string]*entity.User, error) {
	var users []entity.User

	query := db.Model(&entity.User{}).
		Where("role IN ? ", []string{"waste_bank_unit", "waste_bank_central"})

	// Apply location filters
	if request.Province != "" {
		query = query.Where("province ILIKE ?", "%"+request.Province+"%")
	}
	if request.City != "" {
		query = query.Where("city ILIKE ?", "%"+request.City+"%")
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get waste bank users: %v", err)
	}

	// Convert to map for easy lookup
	userMap := make(map[string]*entity.User)
	for i := range users {
		userMap[users[i].ID.String()] = &users[i]
	}

	return userMap, nil
}
