package seeder

import (
	"log"

	"gorm.io/gorm"
)

// RunAllSeeders runs all seeders in the correct order
func RunAllSeeders(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Seed in dependency order
	seeders := []func(*gorm.DB) error{
		SeedWasteCategories,
		SeedWasteTypes,
		SeedUsers,
		SeedCustomerProfiles,
		SeedGovernmentProfiles,
		SeedIndustryProfiles,
		SeedWasteBankProfiles,
		SeedWasteCollectorProfiles,
		SeedCollectorManagement,
		SeedWasteBankPricedTypes,
		SeedStorages,
		SeedStorageItems,
		SeedWasteDropRequests,
		SeedWasteDropRequestItems,
		SeedWasteTransferRequests,
		SeedWasteTransferItemOfferings,
		SeedSalaryTransactions,
	}

	for i, seeder := range seeders {
		log.Printf("Running seeder %d/%d...", i+1, len(seeders))
		if err := seeder(db); err != nil {
			return err
		}
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

// ClearAllData clears all data from tables (useful for testing)
func ClearAllData(db *gorm.DB) error {
	tables := []string{
		"salary_transactions",
		"waste_transfer_items",
		"waste_transfer_requests",
		"waste_drop_request_items",
		"waste_drop_requests",
		"storage_items",
		"storage",
		"waste_bank_priced_types",
		"collector_managements",
		"waste_collector_profiles",
		"waste_bank_profiles",
		"industry_profiles",
		"government_profiles",
		"customer_profiles",
		"refresh_tokens",
		"users",
		"waste_types",
		"waste_categories",
	}

	for _, table := range tables {
		if err := db.Exec("DELETE FROM " + table).Error; err != nil {
			return err
		}
	}

	return nil
}
