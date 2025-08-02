package seeder

import (
	"log"

	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"gorm.io/gorm"
)

// SeedWasteBankPricedTypes seeds custom pricing for waste banks
func SeedWasteBankPricedTypes(db *gorm.DB) error {
	log.Println("Starting to seed waste bank priced types...")

	var wasteBanks []entity.User
	if err := db.Where("role IN ?", []string{"waste_bank_unit", "waste_bank_central"}).Find(&wasteBanks).Error; err != nil {
		return err
	}

	var wasteTypes []entity.WasteType
	if err := db.Find(&wasteTypes).Error; err != nil {
		return err
	}

	if len(wasteBanks) == 0 {
		log.Println("Warning: No waste banks found, skipping pricing")
		return nil
	}

	if len(wasteTypes) == 0 {
		log.Println("Warning: No waste types found, skipping pricing")
		return nil
	}

	// Create a map for easier waste type lookup
	wasteTypeMap := make(map[string]entity.WasteType)
	for _, wt := range wasteTypes {
		wasteTypeMap[wt.Name] = wt
	}

	// Define pricing data with realistic prices (in rupiah per kg)
	// Aligned with your exact waste types
	pricingData := []struct {
		WasteTypeName string
		BasePrice     int64 // Base price that will be varied per waste bank
	}{
		// Plastic types (highest value recyclables)
		{"PET Bottles", 3000},
		{"HDPE Containers", 2500},
		{"Plastic Bags", 800},
		{"Styrofoam", 500},

		// Paper types
		{"Cardboard", 1500},
		{"Office Paper", 2000},
		{"Newspaper", 1200},
		{"Magazines", 1000},

		// Metal types (highest prices)
		{"Aluminum Cans", 15000},
		{"Steel Cans", 8000},
		{"Copper Wire", 25000},
		{"Iron Scrap", 3500},

		// Glass types
		{"Clear Glass Bottles", 1000},
		{"Colored Glass Bottles", 800},
		{"Glass Jars", 900},

		// Electronic types (vary widely)
		{"Mobile Phones", 50000},
		{"Computers", 80000},
		{"Televisions", 30000},
	}

	// Create pricing for each waste bank
	for bankIndex, wasteBank := range wasteBanks {
		for _, priceData := range pricingData {
			wasteType, exists := wasteTypeMap[priceData.WasteTypeName]
			if !exists {
				log.Printf("Warning: Waste type '%s' not found, skipping", priceData.WasteTypeName)
				continue
			}

			// Vary prices based on waste bank index and type
			priceVariation := int64(bankIndex * 200) // Each bank has slightly different prices
			if wasteBank.Role == "waste_bank_central" {
				priceVariation += 500 // Central banks offer higher prices
			}

			finalPrice := priceData.BasePrice + priceVariation

			pricedType := entity.WasteBankPricedType{
				ID:                uuid.New(),
				WasteBankID:       wasteBank.ID,
				WasteTypeID:       wasteType.ID,
				CustomPricePerKgs: finalPrice,
			}

			var existing entity.WasteBankPricedType
			if err := db.Where("waste_bank_id = ? AND waste_type_id = ?",
				pricedType.WasteBankID, pricedType.WasteTypeID).First(&existing).Error; err == gorm.ErrRecordNotFound {
				if err := db.Create(&pricedType).Error; err != nil {
					log.Printf("Error creating priced type for %s at %s: %v",
						wasteType.Name, wasteBank.Username, err)
					return err
				}
				log.Printf("Created pricing for %s at %s: Rp %d/kg",
					wasteType.Name, wasteBank.Username, finalPrice)
			}
		}
	}

	log.Printf("Successfully seeded pricing for %d waste banks", len(wasteBanks))
	return nil
}

// SeedStorages seeds storage facilities for waste banks and industries
func SeedStorages(db *gorm.DB) error {
	log.Println("Starting to seed storage facilities...")

	var wasteBanks []entity.User
	if err := db.Where("role IN ?", []string{"waste_bank_unit", "waste_bank_central"}).Find(&wasteBanks).Error; err != nil {
		return err
	}

	var industries []entity.User
	if err := db.Where("role = ?", "industry").Find(&industries).Error; err != nil {
		return err
	}

	allUsers := append(wasteBanks, industries...)

	if len(allUsers) == 0 {
		log.Println("Warning: No waste banks or industries found, skipping storage creation")
		return nil
	}

	var storages []entity.Storage

	// Create storages for each user based on their role
	for _, user := range allUsers {
		var userStorages []entity.Storage

		if user.Role == "waste_bank_central" {
			// Central waste banks get larger storage facilities
			userStorages = []entity.Storage{
				{
					ID:     uuid.New(),
					UserID: user.ID,
					Length: 3000, // 30m
					Width:  2000, // 20m
					Height: 800,  // 8m
				},
				{
					ID:                    uuid.New(),
					UserID:                user.ID,
					Length:                1500, // 15m
					Width:                 1000, // 10m
					Height:                500,  // 5m
					IsForRecycledMaterial: true,
				},
			}
		} else if user.Role == "waste_bank_unit" {
			// Unit waste banks get medium-sized storage
			userStorages = []entity.Storage{
				{
					ID:     uuid.New(),
					UserID: user.ID,
					Length: 1500, // 15m
					Width:  1000, // 10m
					Height: 500,  // 5m
				},
				{
					ID:                    uuid.New(),
					UserID:                user.ID,
					Length:                800, // 8m
					Width:                 600, // 6m
					Height:                400, // 4m
					IsForRecycledMaterial: true,
				},
			}
		} else if user.Role == "industry" {
			// Industries get very large storage facilities
			userStorages = []entity.Storage{
				{
					ID:     uuid.New(),
					UserID: user.ID,
					Length: 4000, // 40m
					Width:  3000, // 30m
					Height: 1000, // 10m
				},
				{
					ID:                    uuid.New(),
					UserID:                user.ID,
					Length:                2000, // 20m
					Width:                 1500, // 15m
					Height:                600,  // 6m
					IsForRecycledMaterial: true,
				},
			}
		}

		storages = append(storages, userStorages...)
		log.Printf("Prepared %d storage facilities for %s (%s)",
			len(userStorages), user.Username, user.Role)
	}

	// Create the storage records
	for _, storage := range storages {
		var existing entity.Storage
		if err := db.Where("user_id = ? AND length = ? AND width = ? AND height = ? AND is_for_recycled_material = ?",
			storage.UserID, storage.Length, storage.Width, storage.Height, storage.IsForRecycledMaterial).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&storage).Error; err != nil {
				log.Printf("Error creating storage: %v", err)
				return err
			}
			log.Printf("Created storage facility: %dx%dx%d cm", storage.Length, storage.Width, storage.Height)
		}
	}

	log.Printf("Successfully seeded %d storage facilities", len(storages))
	return nil
}

// SeedStorageItems seeds items stored in storage facilities
func SeedStorageItems(db *gorm.DB) error {
	log.Println("Starting to seed storage items...")

	var storages []entity.Storage
	if err := db.Find(&storages).Error; err != nil {
		return err
	}

	var wasteTypes []entity.WasteType
	if err := db.Find(&wasteTypes).Error; err != nil {
		return err
	}

	if len(storages) == 0 {
		log.Println("Warning: No storages found, skipping storage items")
		return nil
	}

	if len(wasteTypes) == 0 {
		log.Println("Warning: No waste types found, skipping storage items")
		return nil
	}

	// Create a map for easier waste type lookup
	wasteTypeMap := make(map[string]entity.WasteType)
	for _, wt := range wasteTypes {
		wasteTypeMap[wt.Name] = wt
	}

	// Define storage item data aligned with your waste types
	storageItemData := []struct {
		WasteTypeName string
		BaseWeight    float64 // Base weight that will be varied per storage
	}{
		// Plastic types (commonly stored)
		{"PET Bottles", 150.5},
		{"HDPE Containers", 198.5},
		{"Plastic Bags", 89.4},
		{"Styrofoam", 34.6},

		// Paper types (bulk items)
		{"Cardboard", 289.3},
		{"Office Paper", 120.7},
		{"Newspaper", 156.2},
		{"Magazines", 76.8},

		// Metal types (high value, smaller quantities)
		{"Aluminum Cans", 45.2},
		{"Steel Cans", 67.8},
		{"Copper Wire", 12.5},
		{"Iron Scrap", 234.7},

		// Glass types (heavy items)
		{"Clear Glass Bottles", 187.3},
		{"Colored Glass Bottles", 156.9},
		{"Glass Jars", 98.4},

		// Electronic types (processed in specialized facilities)
		{"Mobile Phones", 5.2},
		{"Computers", 15.8},
		{"Televisions", 8.7},

		// Textile types (bulky but light)
		{"Cotton Clothing", 67.3},
		{"Synthetic Clothing", 54.8},
	}

	// Create storage items for each storage facility
	for storageIndex, storage := range storages {
		// Get the user who owns this storage to determine storage capacity
		var user entity.User
		if err := db.Where("id = ?", storage.UserID).First(&user).Error; err != nil {
			log.Printf("Warning: Could not find user for storage, skipping")
			continue
		}

		// Determine how many waste types to store based on storage size and user role
		var itemsToCreate int
		var allowedCategories []string

		if user.Role == "industry" {
			// Industries focus on specific categories they can process
			allowedCategories = []string{"Plastic", "Metal", "Paper", "Glass"}
			itemsToCreate = 8 + (storageIndex % 4) // Industries store 8-11 types
		} else if user.Role == "waste_bank_central" {
			// Central banks handle most categories except organic
			allowedCategories = []string{"Plastic", "Paper", "Metal", "Glass", "Electronic"}
			itemsToCreate = 6 + (storageIndex % 3) // Central banks store 6-8 types
		} else {
			// Unit banks focus on common recyclables
			allowedCategories = []string{"Plastic", "Paper", "Metal", "Glass"}
			itemsToCreate = 4 + (storageIndex % 3) // Unit banks store 4-6 types
		}

		// Filter storage items by allowed categories
		var filteredItems []struct {
			WasteTypeName string
			BaseWeight    float64
		}

		for _, itemData := range storageItemData {
			wasteType, exists := wasteTypeMap[itemData.WasteTypeName]
			if !exists {
				continue
			}

			// Get category for this waste type
			var category entity.WasteCategory
			if err := db.Where("id = ?", wasteType.CategoryID).First(&category).Error; err != nil {
				continue
			}

			// Check if this category is allowed for this user role
			for _, allowedCat := range allowedCategories {
				if category.Name == allowedCat {
					filteredItems = append(filteredItems, itemData)
					break
				}
			}
		}

		// Ensure we don't exceed available filtered waste types
		if itemsToCreate > len(filteredItems) {
			itemsToCreate = len(filteredItems)
		}

		// Create items for this storage from filtered list
		for i := 0; i < itemsToCreate; i++ {
			itemData := filteredItems[i]

			wasteType, exists := wasteTypeMap[itemData.WasteTypeName]
			if !exists {
				log.Printf("Warning: Waste type '%s' not found, skipping", itemData.WasteTypeName)
				continue
			}

			// Calculate weight based on storage capacity and user role
			weightMultiplier := 1.0
			if user.Role == "industry" {
				weightMultiplier = 2.5 // Industries store more
			} else if user.Role == "waste_bank_central" {
				weightMultiplier = 1.8 // Central banks store more
			}

			// Add some variation based on storage index
			weightVariation := 1.0 + (float64(storageIndex%5) * 0.2)
			finalWeight := itemData.BaseWeight * weightMultiplier * weightVariation

			storageItem := entity.StorageItem{
				ID:          uuid.New(),
				StorageID:   storage.ID,
				WasteTypeID: wasteType.ID,
				WeightKgs:   finalWeight,
			}

			var existing entity.StorageItem
			if err := db.Where("storage_id = ? AND waste_type_id = ?",
				storageItem.StorageID, storageItem.WasteTypeID).First(&existing).Error; err == gorm.ErrRecordNotFound {
				if err := db.Create(&storageItem).Error; err != nil {
					log.Printf("Error creating storage item: %v", err)
					return err
				}
				log.Printf("Created storage item: %.1f kg of %s for %s",
					finalWeight, wasteType.Name, user.Username)
			}
		}
	}

	log.Println("Successfully seeded storage items")
	return nil
}

// SeedAllStorageData seeds all storage-related data
func SeedAllStorageData(db *gorm.DB) error {
	log.Println("Starting to seed all storage-related data...")

	// First, seed waste bank pricing
	if err := SeedWasteBankPricedTypes(db); err != nil {
		log.Printf("Error seeding waste bank priced types: %v", err)
		return err
	}

	// Then, seed storage facilities
	if err := SeedStorages(db); err != nil {
		log.Printf("Error seeding storages: %v", err)
		return err
	}

	// Finally, seed storage items
	if err := SeedStorageItems(db); err != nil {
		log.Printf("Error seeding storage items: %v", err)
		return err
	}

	log.Println("Successfully seeded all storage-related data")
	return nil
}

// Helper function to get storage capacity utilization
func GetStorageUtilization(db *gorm.DB, storageID uuid.UUID) (float64, error) {
	var items []entity.StorageItem
	if err := db.Where("storage_id = ?", storageID).Find(&items).Error; err != nil {
		return 0, err
	}

	var totalWeight float64
	for _, item := range items {
		totalWeight += item.WeightKgs
	}

	return totalWeight, nil
}
