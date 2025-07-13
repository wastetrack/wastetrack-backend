package seeder

import (
	"log"

	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"gorm.io/gorm"
)

// SeedWasteBankPricedTypes seeds custom pricing for waste banks
func SeedWasteBankPricedTypes(db *gorm.DB) error {
	var wasteBanks []entity.User
	if err := db.Where("role IN ?", []string{"waste_bank_unit", "waste_bank_central"}).Find(&wasteBanks).Error; err != nil {
		return err
	}

	var wasteTypes []entity.WasteType
	if err := db.Find(&wasteTypes).Error; err != nil {
		return err
	}

	if len(wasteBanks) == 0 || len(wasteTypes) == 0 {
		return nil
	}

	// Create pricing for common waste types
	pricings := []struct {
		WasteBankIndex int
		WasteTypeName  string
		PricePerKg     int64
	}{
		{0, "PET Bottles", 3000},
		{0, "Aluminum Cans", 15000},
		{0, "Cardboard", 1500},
		{0, "Office Paper", 2000},
		{1, "PET Bottles", 3200},
		{1, "Aluminum Cans", 14500},
		{1, "Cardboard", 1600},
		{1, "Steel Cans", 8000},
		{2, "PET Bottles", 3500}, // Central bank (higher prices)
		{2, "Aluminum Cans", 16000},
		{2, "Cardboard", 1800},
		{2, "Steel Cans", 8500},
	}

	for _, pricing := range pricings {
		if pricing.WasteBankIndex >= len(wasteBanks) {
			continue
		}

		// Find the waste type
		var wasteType entity.WasteType
		if err := db.Where("name = ?", pricing.WasteTypeName).First(&wasteType).Error; err != nil {
			continue
		}

		pricedType := entity.WasteBankPricedType{
			ID:                uuid.New(),
			WasteBankID:       wasteBanks[pricing.WasteBankIndex].ID,
			WasteTypeID:       wasteType.ID,
			CustomPricePerKgs: pricing.PricePerKg,
		}

		var existing entity.WasteBankPricedType
		if err := db.Where("waste_bank_id = ? AND waste_type_id = ?", pricedType.WasteBankID, pricedType.WasteTypeID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&pricedType).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedStorages seeds storage facilities
func SeedStorages(db *gorm.DB) error {
	var wasteBanks []entity.User
	if err := db.Where("role IN ?", []string{"waste_bank_unit", "waste_bank_central"}).Find(&wasteBanks).Error; err != nil {
		return err
	}

	var industries []entity.User
	if err := db.Where("role = ?", "industry").Find(&industries).Error; err != nil {
		return err
	}

	var allUsers []entity.User
	allUsers = append(allUsers, wasteBanks...)
	allUsers = append(allUsers, industries...)

	if len(allUsers) == 0 {
		log.Println("Warning: No waste banks or industries found, skipping storage creation")
		return nil
	}

	storages := []entity.Storage{
		{
			UserID: allUsers[0].ID,
			Length: 1000, // 10m
			Width:  800,  // 8m
			Height: 400,  // 4m
		},
		{
			UserID: allUsers[0].ID,
			Length: 600, // 6m
			Width:  600, // 6m
			Height: 300, // 3m
		},
	}

	// Add more storages if we have more users
	if len(allUsers) > 1 {
		storages = append(storages, entity.Storage{
			UserID: allUsers[1].ID,
			Length: 1500, // 15m
			Width:  1000, // 10m
			Height: 500,  // 5m
		})
	}

	if len(allUsers) > 2 {
		storages = append(storages, entity.Storage{
			UserID: allUsers[2].ID,
			Length: 2000, // 20m
			Width:  1500, // 15m
			Height: 600,  // 6m
		})
	}

	// Central waste bank gets larger storage
	if len(wasteBanks) >= 3 {
		storages = append(storages, entity.Storage{
			UserID: wasteBanks[2].ID, // Central waste bank
			Length: 3000,             // 30m
			Width:  2000,             // 20m
			Height: 800,              // 8m
		})
	}

	for _, storage := range storages {
		var existing entity.Storage
		if err := db.Where("user_id = ? AND length = ? AND width = ? AND height = ?",
			storage.UserID, storage.Length, storage.Width, storage.Height).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&storage).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedStorageItems seeds items stored in storage facilities
func SeedStorageItems(db *gorm.DB) error {
	var storages []entity.Storage
	if err := db.Find(&storages).Error; err != nil {
		return err
	}

	var wasteTypes []entity.WasteType
	if err := db.Find(&wasteTypes).Error; err != nil {
		return err
	}

	if len(storages) == 0 || len(wasteTypes) == 0 {
		return nil
	}

	// Create storage items for each storage
	storageItems := []struct {
		StorageIndex  int
		WasteTypeName string
		QuantityKgs   float64
	}{
		{0, "PET Bottles", 150.5},
		{0, "Cardboard", 89.3},
		{0, "Aluminum Cans", 45.2},
		{1, "Office Paper", 120.7},
		{1, "Steel Cans", 67.8},
		{2, "PET Bottles", 340.2},
		{2, "HDPE Containers", 198.5},
		{2, "Glass Jars", 156.3},
		{3, "Plastic Bags", 89.4},
		{3, "Styrofoam", 34.6},
	}

	for _, item := range storageItems {
		if item.StorageIndex >= len(storages) {
			continue
		}

		// Find the waste type
		var wasteType entity.WasteType
		if err := db.Where("name = ?", item.WasteTypeName).First(&wasteType).Error; err != nil {
			continue
		}

		storageItem := entity.StorageItem{
			StorageID:   uuid.New(),
			WasteTypeID: wasteType.ID,
			QuantityKgs: item.QuantityKgs,
		}

		var existing entity.StorageItem
		if err := db.Where("storage_id = ? AND waste_type_id = ?", storageItem.StorageID, storageItem.WasteTypeID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&storageItem).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
