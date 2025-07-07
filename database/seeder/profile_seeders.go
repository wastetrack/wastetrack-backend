package seeder

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/types"
	"gorm.io/gorm"
)

// Helper function to create TimeOnly from hour and minute
func createTimeOnly(hour, minute int) types.TimeOnly {
	return types.NewTimeOnly(time.Date(0, 1, 1, hour, minute, 0, 0, time.UTC))
}

// SeedCustomerProfiles seeds customer profiles
func SeedCustomerProfiles(db *gorm.DB) error {
	var customers []entity.User
	if err := db.Where("role = ?", "customer").Find(&customers).Error; err != nil {
		return err
	}

	if len(customers) == 0 {
		log.Println("Warning: No customer users found, skipping customer profiles")
		return nil
	}

	customerData := []struct {
		CarbonDeficit int64
		WaterSaved    int64
		BagsStored    int64
		Trees         int64
	}{
		{45, 120, 8, 3},
		{67, 180, 12, 5},
		{23, 75, 4, 2},
	}

	for i, customer := range customers {
		if i >= len(customerData) {
			break
		}
		data := customerData[i]

		profile := entity.CustomerProfile{
			ID:            uuid.New(),
			UserID:        customer.ID,
			CarbonDeficit: data.CarbonDeficit,
			WaterSaved:    data.WaterSaved,
			BagsStored:    data.BagsStored,
			Trees:         data.Trees,
		}

		var existing entity.CustomerProfile
		if err := db.Where("user_id = ?", profile.UserID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&profile).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedGovernmentProfiles seeds government profiles
func SeedGovernmentProfiles(db *gorm.DB) error {
	var govUsers []entity.User
	if err := db.Where("role = ?", "government").Find(&govUsers).Error; err != nil {
		return err
	}

	for _, user := range govUsers {
		profile := entity.GovernmentProfile{
			ID:     uuid.New(),
			UserID: user.ID,
		}

		var existing entity.GovernmentProfile
		if err := db.Where("user_id = ?", profile.UserID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&profile).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedIndustryProfiles seeds industry profiles
func SeedIndustryProfiles(db *gorm.DB) error {
	var industryUsers []entity.User
	if err := db.Where("role = ?", "industry").Find(&industryUsers).Error; err != nil {
		return err
	}

	if len(industryUsers) == 0 {
		log.Println("Warning: No industry users found, skipping industry profiles")
		return nil
	}

	profileData := []struct {
		TotalWasteWeight    float64
		TotalRecycledWeight float64
	}{
		{5420.5, 4876.3},
		{3280.7, 2952.4},
	}

	for i, user := range industryUsers {
		if i >= len(profileData) {
			break
		}
		data := profileData[i]

		profile := entity.IndustryProfile{
			ID:                  uuid.New(),
			UserID:              user.ID,
			TotalWasteWeight:    data.TotalWasteWeight,
			TotalRecycledWeight: data.TotalRecycledWeight,
		}

		var existing entity.IndustryProfile
		if err := db.Where("user_id = ?", profile.UserID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&profile).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedWasteBankProfiles seeds waste bank profiles
func SeedWasteBankProfiles(db *gorm.DB) error {
	var wasteBankUsers []entity.User
	if err := db.Where("role IN ?", []string{"waste_bank_unit", "waste_bank_central"}).Find(&wasteBankUsers).Error; err != nil {
		return err
	}

	if len(wasteBankUsers) == 0 {
		log.Println("Warning: No waste bank users found, skipping waste bank profiles")
		return nil
	}

	profileData := []struct {
		TotalWasteWeight float64
		TotalWorkers     int64
		OpenHour         int
		OpenMinute       int
		CloseHour        int
		CloseMinute      int
	}{
		{2840.6, 12, 8, 0, 17, 0},
		{3567.8, 18, 7, 30, 18, 0},
		{8500.2, 45, 6, 0, 20, 0},
	}

	for i, user := range wasteBankUsers {
		if i >= len(profileData) {
			break
		}
		data := profileData[i]

		profile := entity.WasteBankProfile{
			ID:               uuid.New(),
			UserID:           user.ID,
			TotalWasteWeight: data.TotalWasteWeight,
			TotalWorkers:     data.TotalWorkers,
			OpenTime:         createTimeOnly(data.OpenHour, data.OpenMinute),
			CloseTime:        createTimeOnly(data.CloseHour, data.CloseMinute),
		}

		var existing entity.WasteBankProfile
		if err := db.Where("user_id = ?", profile.UserID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&profile).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedWasteCollectorProfiles seeds waste collector profiles
func SeedWasteCollectorProfiles(db *gorm.DB) error {
	var collectorUsers []entity.User
	if err := db.Where("role IN ?", []string{"waste_collector_unit", "waste_collector_central"}).Find(&collectorUsers).Error; err != nil {
		return err
	}

	if len(collectorUsers) == 0 {
		log.Println("Warning: No waste collector users found, skipping collector profiles")
		return nil
	}

	collectorData := []float64{1250.4, 980.7, 3500.8}

	for i, user := range collectorUsers {
		weight := collectorData[0]
		if i < len(collectorData) {
			weight = collectorData[i]
		}

		profile := entity.WasteCollectorProfile{
			ID:               uuid.New(),
			UserID:           user.ID,
			TotalWasteWeight: weight,
		}

		var existing entity.WasteCollectorProfile
		if err := db.Where("user_id = ?", profile.UserID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&profile).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedCollectorManagement seeds collector management relationships
func SeedCollectorManagement(db *gorm.DB) error {
	var wasteBanks []entity.User
	if err := db.Where("role IN ?", []string{"waste_bank_unit", "waste_bank_central"}).Find(&wasteBanks).Error; err != nil {
		return err
	}

	var collectors []entity.User
	if err := db.Where("role IN ?", []string{"waste_collector_unit", "waste_collector_central"}).Find(&collectors).Error; err != nil {
		return err
	}

	if len(wasteBanks) == 0 {
		log.Println("Warning: No waste banks found, skipping collector management")
		return nil
	}
	if len(collectors) == 0 {
		log.Println("Warning: No collectors found, skipping collector management")
		return nil
	}

	var managements []entity.CollectorManagement

	if len(wasteBanks) >= 1 && len(collectors) >= 1 {
		managements = append(managements, entity.CollectorManagement{
			ID:          uuid.New(),
			WasteBankID: wasteBanks[0].ID,
			CollectorID: collectors[0].ID,
			Status:      "active",
		})
	}
	if len(wasteBanks) >= 1 && len(collectors) >= 2 {
		managements = append(managements, entity.CollectorManagement{
			ID:          uuid.New(),
			WasteBankID: wasteBanks[0].ID,
			CollectorID: collectors[1].ID,
			Status:      "active",
		})
	}
	if len(wasteBanks) >= 2 && len(collectors) >= 2 {
		managements = append(managements, entity.CollectorManagement{
			ID:          uuid.New(),
			WasteBankID: wasteBanks[1].ID,
			CollectorID: collectors[1].ID,
			Status:      "inactive",
		})
	}
	if len(wasteBanks) >= 3 && len(collectors) >= 3 {
		managements = append(managements, entity.CollectorManagement{
			ID:          uuid.New(),
			WasteBankID: wasteBanks[2].ID,
			CollectorID: collectors[2].ID,
			Status:      "active",
		})
		for i := 0; i < len(collectors) && i < 2; i++ {
			managements = append(managements, entity.CollectorManagement{
				ID:          uuid.New(),
				WasteBankID: wasteBanks[2].ID,
				CollectorID: collectors[i].ID,
				Status:      "active",
			})
		}
	}

	for _, management := range managements {
		var existing entity.CollectorManagement
		if err := db.Where("waste_bank_id = ? AND collector_id = ?", management.WasteBankID, management.CollectorID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&management).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
