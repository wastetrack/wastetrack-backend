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

// SeedAllUserProfiles creates profiles for all users based on their roles
func SeedAllUserProfiles(db *gorm.DB) error {
	log.Println("Starting to seed user profiles...")

	// Seed all profile types in dependency order
	if err := SeedCustomerProfiles(db); err != nil {
		log.Printf("Error seeding customer profiles: %v", err)
		return err
	}

	if err := SeedWasteBankProfiles(db); err != nil {
		log.Printf("Error seeding waste bank profiles: %v", err)
		return err
	}

	if err := SeedWasteCollectorProfiles(db); err != nil {
		log.Printf("Error seeding waste collector profiles: %v", err)
		return err
	}

	if err := SeedIndustryProfiles(db); err != nil {
		log.Printf("Error seeding industry profiles: %v", err)
		return err
	}

	if err := SeedGovernmentProfiles(db); err != nil {
		log.Printf("Error seeding government profiles: %v", err)
		return err
	}

	// Seed management relationships after all profiles are created
	if err := SeedCollectorManagement(db); err != nil {
		log.Printf("Error seeding collector management: %v", err)
		return err
	}

	log.Println("Successfully seeded all user profiles")
	return nil
}

// SeedCustomerProfiles creates profiles for all customer users
func SeedCustomerProfiles(db *gorm.DB) error {
	log.Println("Starting to seed customer profiles...")

	var customers []entity.User
	if err := db.Where("role = ?", "customer").Find(&customers).Error; err != nil {
		log.Printf("Error fetching customers: %v", err)
		return err
	}

	if len(customers) == 0 {
		log.Println("Warning: No customer users found, skipping customer profiles")
		return nil
	}

	// Sample data for customer profiles - can be expanded as needed
	customerData := []struct {
		CarbonDeficit int64
		WaterSaved    int64
		BagsStored    int64
		Trees         int64
	}{
		{45, 120, 8, 3},
		{67, 180, 12, 5},
		{23, 75, 4, 2},
		{89, 240, 15, 7},
		{34, 95, 6, 3},
		{56, 150, 10, 4},
		{78, 200, 14, 6},
		{41, 110, 7, 3},
	}

	for i, customer := range customers {
		// Check if profile already exists
		var existing entity.CustomerProfile
		if err := db.Where("user_id = ?", customer.ID).First(&existing).Error; err != gorm.ErrRecordNotFound {
			log.Printf("Customer profile already exists for user: %s", customer.Username)
			continue
		}

		// Use modulo to cycle through data if we have more customers than data entries
		dataIndex := i % len(customerData)
		data := customerData[dataIndex]

		profile := entity.CustomerProfile{
			ID:            uuid.New(),
			UserID:        customer.ID,
			CarbonDeficit: data.CarbonDeficit,
			WaterSaved:    data.WaterSaved,
			BagsStored:    data.BagsStored,
			Trees:         data.Trees,
		}

		if err := db.Create(&profile).Error; err != nil {
			log.Printf("Error creating customer profile for user %s: %v", customer.Username, err)
			return err
		}
		log.Printf("Created customer profile for user: %s", customer.Username)
	}

	log.Printf("Successfully processed %d customer profiles", len(customers))
	return nil
}

// SeedWasteBankProfiles creates profiles for waste bank unit and central users
func SeedWasteBankProfiles(db *gorm.DB) error {
	log.Println("Starting to seed waste bank profiles...")

	var wasteBankUsers []entity.User
	if err := db.Where("role IN ?", []string{"waste_bank_unit", "waste_bank_central"}).Find(&wasteBankUsers).Error; err != nil {
		log.Printf("Error fetching waste bank users: %v", err)
		return err
	}

	if len(wasteBankUsers) == 0 {
		log.Println("Warning: No waste bank users found, skipping waste bank profiles")
		return nil
	}

	// Sample data for waste bank profiles aligned with your specific users
	profileData := []struct {
		TotalWasteWeight float64
		TotalWorkers     int64
		OpenHour         int
		OpenMinute       int
		CloseHour        int
		CloseMinute      int
	}{
		{2840.6, 12, 8, 0, 17, 0},   // BSU Hijau Berkah (adi_suberkah)
		{8500.2, 45, 6, 0, 20, 0},   // BSI Raya Surabaya (mina_astiya) - Central
		{12500.8, 65, 6, 0, 22, 0},  // BSI Ageng Buana Surabaya (andi_diata) - Central
		{4200.3, 25, 8, 30, 16, 30}, // Additional unit banks
		{6800.9, 35, 7, 0, 19, 0},   // Additional central banks
		{3200.5, 18, 8, 0, 17, 30},  // Additional unit banks
		{9800.4, 50, 7, 0, 21, 0},   // Additional central banks
	}

	for i, user := range wasteBankUsers {
		// Check if profile already exists
		var existing entity.WasteBankProfile
		if err := db.Where("user_id = ?", user.ID).First(&existing).Error; err != gorm.ErrRecordNotFound {
			log.Printf("Waste bank profile already exists for user: %s", user.Username)
			continue
		}

		// Use modulo to cycle through data if we have more users than data entries
		dataIndex := i % len(profileData)
		data := profileData[dataIndex]

		// Adjust data based on role
		adjustedData := data
		if user.Role == "waste_bank_central" {
			// Central banks typically have more capacity
			adjustedData.TotalWasteWeight *= 1.5
			adjustedData.TotalWorkers = int64(float64(adjustedData.TotalWorkers) * 1.3)
		}

		profile := entity.WasteBankProfile{
			ID:               uuid.New(),
			UserID:           user.ID,
			TotalWasteWeight: adjustedData.TotalWasteWeight,
			TotalWorkers:     adjustedData.TotalWorkers,
			OpenTime:         createTimeOnly(adjustedData.OpenHour, adjustedData.OpenMinute),
			CloseTime:        createTimeOnly(adjustedData.CloseHour, adjustedData.CloseMinute),
		}

		if err := db.Create(&profile).Error; err != nil {
			log.Printf("Error creating waste bank profile for user %s: %v", user.Username, err)
			return err
		}
		log.Printf("Created waste bank profile for user: %s (%s) - Workers: %d, Waste: %.1f kg",
			user.Username, user.Role, profile.TotalWorkers, profile.TotalWasteWeight)
	}

	log.Printf("Successfully processed %d waste bank profiles", len(wasteBankUsers))
	return nil
}

// SeedWasteCollectorProfiles creates profiles for waste collector unit and central users
func SeedWasteCollectorProfiles(db *gorm.DB) error {
	log.Println("Starting to seed waste collector profiles...")

	var collectorUsers []entity.User
	if err := db.Where("role IN ?", []string{"waste_collector_unit", "waste_collector_central"}).Find(&collectorUsers).Error; err != nil {
		log.Printf("Error fetching waste collector users: %v", err)
		return err
	}

	if len(collectorUsers) == 0 {
		log.Println("Warning: No waste collector users found, skipping collector profiles")
		return nil
	}

	// Sample data for waste collector profiles aligned with your users
	collectorData := []struct {
		TotalWasteWeight float64
		Role             string // To match specific roles
	}{
		{1250.4, "waste_collector_unit"},    // ega_basuka
		{3500.8, "waste_collector_central"}, // awai_sina
		{4200.5, "waste_collector_central"}, // asi_bunaya
		{1800.3, "waste_collector_unit"},    // Additional unit collectors
		{2750.9, "waste_collector_central"}, // Additional central collectors
		{1450.7, "waste_collector_unit"},    // Additional unit collectors
	}

	for i, user := range collectorUsers {
		// Check if profile already exists
		var existing entity.WasteCollectorProfile
		if err := db.Where("user_id = ?", user.ID).First(&existing).Error; err != gorm.ErrRecordNotFound {
			log.Printf("Waste collector profile already exists for user: %s", user.Username)
			continue
		}

		// Use modulo to cycle through data, preferring role matches
		var weight float64
		dataIndex := i % len(collectorData)

		// Try to find matching role data first
		for j, data := range collectorData {
			if data.Role == user.Role {
				dataIndex = j
				break
			}
		}

		weight = collectorData[dataIndex].TotalWasteWeight

		// Add variation based on user index to avoid identical weights
		variation := float64(i) * 100.5
		weight += variation

		profile := entity.WasteCollectorProfile{
			ID:               uuid.New(),
			UserID:           user.ID,
			TotalWasteWeight: weight,
		}

		if err := db.Create(&profile).Error; err != nil {
			log.Printf("Error creating waste collector profile for user %s: %v", user.Username, err)
			return err
		}
		log.Printf("Created waste collector profile for user: %s (%s) - Waste: %.1f kg",
			user.Username, user.Role, weight)
	}

	log.Printf("Successfully processed %d waste collector profiles", len(collectorUsers))
	return nil
}

// SeedIndustryProfiles creates profiles for industry users
func SeedIndustryProfiles(db *gorm.DB) error {
	log.Println("Starting to seed industry profiles...")

	var industryUsers []entity.User
	if err := db.Where("role = ?", "industry").Find(&industryUsers).Error; err != nil {
		log.Printf("Error fetching industry users: %v", err)
		return err
	}

	if len(industryUsers) == 0 {
		log.Println("Warning: No industry users found, skipping industry profiles")
		return nil
	}

	// Sample data for industry profiles
	profileData := []struct {
		TotalWasteWeight    float64
		TotalRecycledWeight float64
	}{
		{5420.5, 4876.3}, // ofi_takena - Offtaker Eko Subur Langgeng
		{7800.9, 7020.8}, // Additional industries
		{4650.2, 4185.1},
		{6200.4, 5580.3},
		{8900.7, 8010.6},
		{3280.7, 2952.4},
	}

	for i, user := range industryUsers {
		// Check if profile already exists
		var existing entity.IndustryProfile
		if err := db.Where("user_id = ?", user.ID).First(&existing).Error; err != gorm.ErrRecordNotFound {
			log.Printf("Industry profile already exists for user: %s", user.Username)
			continue
		}

		// Use modulo to cycle through data if we have more users than data entries
		dataIndex := i % len(profileData)
		data := profileData[dataIndex]

		profile := entity.IndustryProfile{
			ID:                  uuid.New(),
			UserID:              user.ID,
			TotalWasteWeight:    data.TotalWasteWeight,
			TotalRecycledWeight: data.TotalRecycledWeight,
		}

		if err := db.Create(&profile).Error; err != nil {
			log.Printf("Error creating industry profile for user %s: %v", user.Username, err)
			return err
		}

		recyclingRate := (data.TotalRecycledWeight / data.TotalWasteWeight) * 100
		log.Printf("Created industry profile for user: %s - Total: %.1f kg, Recycled: %.1f kg (%.1f%%)",
			user.Username, data.TotalWasteWeight, data.TotalRecycledWeight, recyclingRate)
	}

	log.Printf("Successfully processed %d industry profiles", len(industryUsers))
	return nil
}

// SeedGovernmentProfiles creates profiles for government users
func SeedGovernmentProfiles(db *gorm.DB) error {
	log.Println("Starting to seed government profiles...")

	var govUsers []entity.User
	if err := db.Where("role = ?", "government").Find(&govUsers).Error; err != nil {
		log.Printf("Error fetching government users: %v", err)
		return err
	}

	if len(govUsers) == 0 {
		log.Println("Warning: No government users found, skipping government profiles")
		return nil
	}

	for _, user := range govUsers {
		// Check if profile already exists
		var existing entity.GovernmentProfile
		if err := db.Where("user_id = ?", user.ID).First(&existing).Error; err != gorm.ErrRecordNotFound {
			log.Printf("Government profile already exists for user: %s", user.Username)
			continue
		}

		profile := entity.GovernmentProfile{
			ID:     uuid.New(),
			UserID: user.ID,
		}

		if err := db.Create(&profile).Error; err != nil {
			log.Printf("Error creating government profile for user %s: %v", user.Username, err)
			return err
		}
		log.Printf("Created government profile for user: %s", user.Username)
	}

	log.Printf("Successfully processed %d government profiles", len(govUsers))
	return nil
}

// SeedCollectorManagement creates management relationships between waste banks and collectors
func SeedCollectorManagement(db *gorm.DB) error {
	log.Println("Starting to seed collector management relationships...")

	// First try specific relationships based on your user data
	if err := SeedSpecificCollectorManagement(db); err != nil {
		log.Printf("Error seeding specific collector management: %v", err)
		return err
	}

	// Then fill any gaps with institution-based matching
	if err := SeedInstitutionBasedCollectorManagement(db); err != nil {
		log.Printf("Error seeding institution-based collector management: %v", err)
		return err
	}

	log.Println("Successfully seeded collector management relationships")
	return nil
}

// SeedSpecificCollectorManagement creates management based on your specific user data
func SeedSpecificCollectorManagement(db *gorm.DB) error {
	log.Println("Seeding specific collector management relationships...")

	// Define the specific relationships based on your user seeder data
	managementSpecs := []struct {
		WasteBankUsername string
		CollectorUsername string
		Status            string
		Description       string
	}{
		// Primary relationships (same institution)
		{"adi_suberkah", "ega_basuka", "active", "BSU Hijau Berkah primary relationship"},
		{"mina_astiya", "awai_sina", "active", "BSI Raya Surabaya primary relationship"},
		{"andi_diata", "asi_bunaya", "active", "BSI Ageng Buana primary relationship"},
	}

	for _, spec := range managementSpecs {
		// Find waste bank
		var wasteBank entity.User
		if err := db.Where("username = ?", spec.WasteBankUsername).First(&wasteBank).Error; err != nil {
			log.Printf("Warning: Waste bank user '%s' not found, skipping", spec.WasteBankUsername)
			continue
		}

		// Find collector
		var collector entity.User
		if err := db.Where("username = ?", spec.CollectorUsername).First(&collector).Error; err != nil {
			log.Printf("Warning: Collector user '%s' not found, skipping", spec.CollectorUsername)
			continue
		}

		// Check if relationship already exists
		var existing entity.CollectorManagement
		if err := db.Where("waste_bank_id = ? AND collector_id = ?",
			wasteBank.ID, collector.ID).First(&existing).Error; err != gorm.ErrRecordNotFound {
			log.Printf("Relationship already exists: %s -> %s", wasteBank.Username, collector.Username)
			continue
		}

		// Create management relationship
		management := entity.CollectorManagement{
			ID:          uuid.New(),
			WasteBankID: wasteBank.ID,
			CollectorID: collector.ID,
			Status:      spec.Status,
		}

		if err := db.Create(&management).Error; err != nil {
			log.Printf("Error creating specific management relationship: %v", err)
			return err
		}

		log.Printf("Created specific relationship: %s manages %s (%s) - %s",
			wasteBank.Username, collector.Username, spec.Status, spec.Description)
	}

	return nil
}

// SeedInstitutionBasedCollectorManagement fills gaps with institution-based matching
func SeedInstitutionBasedCollectorManagement(db *gorm.DB) error {
	log.Println("Seeding institution-based collector management relationships...")

	// Get all waste banks and collectors
	var wasteBanks []entity.User
	if err := db.Where("role IN ?", []string{"waste_bank_unit", "waste_bank_central"}).Find(&wasteBanks).Error; err != nil {
		return err
	}

	var collectors []entity.User
	if err := db.Where("role IN ?", []string{"waste_collector_unit", "waste_collector_central"}).Find(&collectors).Error; err != nil {
		return err
	}

	// Find unassigned collectors
	var assignedCollectorIDs []uuid.UUID
	db.Model(&entity.CollectorManagement{}).Pluck("collector_id", &assignedCollectorIDs)

	assignedMap := make(map[uuid.UUID]bool)
	for _, id := range assignedCollectorIDs {
		assignedMap[id] = true
	}

	// Assign unassigned collectors to waste banks with same institution
	for _, collector := range collectors {
		if assignedMap[collector.ID] {
			continue // Already assigned
		}

		// Find waste bank with same institution
		for _, wasteBank := range wasteBanks {
			if collector.Institution == wasteBank.Institution && collector.Institution != "" {
				management := entity.CollectorManagement{
					ID:          uuid.New(),
					WasteBankID: wasteBank.ID,
					CollectorID: collector.ID,
					Status:      "active",
				}

				if err := db.Create(&management).Error; err != nil {
					log.Printf("Error creating institution-based relationship: %v", err)
					continue
				}

				log.Printf("Created institution-based relationship: %s manages %s (Institution: %s)",
					wasteBank.Username, collector.Username, collector.Institution)
				assignedMap[collector.ID] = true
				break
			}
		}
	}

	// Assign remaining unassigned collectors using round-robin
	unassignedCount := 0
	for _, collector := range collectors {
		if !assignedMap[collector.ID] {
			if len(wasteBanks) > 0 {
				wasteBank := wasteBanks[unassignedCount%len(wasteBanks)]

				management := entity.CollectorManagement{
					ID:          uuid.New(),
					WasteBankID: wasteBank.ID,
					CollectorID: collector.ID,
					Status:      "inactive", // Inactive since no institution match
				}

				if err := db.Create(&management).Error; err != nil {
					log.Printf("Error creating round-robin relationship: %v", err)
					continue
				}

				log.Printf("Created round-robin relationship: %s manages %s (Status: inactive - no institution match)",
					wasteBank.Username, collector.Username)
				unassignedCount++
			}
		}
	}

	return nil
}

// Utility function to get management summary
func GetCollectorManagementSummary(db *gorm.DB) error {
	log.Println("=== Collector Management Summary ===")

	var wasteBanks []entity.User
	if err := db.Where("role IN ?", []string{"waste_bank_unit", "waste_bank_central"}).Find(&wasteBanks).Error; err != nil {
		return err
	}

	for _, wasteBank := range wasteBanks {
		var managements []entity.CollectorManagement
		db.Where("waste_bank_id = ?", wasteBank.ID).Find(&managements)

		activeCount := 0
		inactiveCount := 0
		for _, mgmt := range managements {
			if mgmt.Status == "active" {
				activeCount++
			} else {
				inactiveCount++
			}
		}

		log.Printf("Waste Bank: %s (%s, %s) manages %d collectors (%d active, %d inactive)",
			wasteBank.Username, wasteBank.Role, wasteBank.Institution,
			len(managements), activeCount, inactiveCount)
	}

	var totalCollectors int64
	var totalManagements int64
	db.Model(&entity.User{}).Where("role IN ?", []string{"waste_collector_unit", "waste_collector_central"}).Count(&totalCollectors)
	db.Model(&entity.CollectorManagement{}).Count(&totalManagements)

	log.Printf("Total collectors: %d, Total management relationships: %d", totalCollectors, totalManagements)
	log.Println("=== End Summary ===")

	return nil
}
