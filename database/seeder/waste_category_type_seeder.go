package seeder

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"gorm.io/gorm"
)

// SeedWasteCategories seeds the waste_categories table
func SeedWasteCategories(db *gorm.DB) error {
	categories := []entity.WasteCategory{
		{
			ID:          uuid.New(),
			Name:        "Plastic",
			Description: "All types of plastic waste including bottles, containers, and packaging",
		},
		{
			ID:          uuid.New(),
			Name:        "Paper",
			Description: "Paper waste including newspapers, cardboard, and office paper",
		},
		{
			ID:          uuid.New(),
			Name:        "Metal",
			Description: "Metal waste including aluminum cans, steel, and copper",
		},
		{
			ID:          uuid.New(),
			Name:        "Glass",
			Description: "Glass waste including bottles, jars, and window glass",
		},
		{
			ID:          uuid.New(),
			Name:        "Organic",
			Description: "Biodegradable waste including food scraps and garden waste",
		},
		{
			ID:          uuid.New(),
			Name:        "Electronic",
			Description: "Electronic waste including computers, phones, and appliances",
		},
		{
			ID:          uuid.New(),
			Name:        "Textile",
			Description: "Fabric and clothing waste",
		},
	}

	for _, category := range categories {
		var existing entity.WasteCategory
		if err := db.Where("name = ?", category.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&category).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedWasteTypes seeds the waste_types table
func SeedWasteTypes(db *gorm.DB) error {
	// First get all categories
	var categories []entity.WasteCategory
	if err := db.Find(&categories).Error; err != nil {
		return err
	}

	// Create a map for easy lookup
	categoryMap := make(map[string]uuid.UUID)
	for _, cat := range categories {
		categoryMap[cat.Name] = cat.ID
	}

	wasteTypes := []struct {
		CategoryName string
		Name         string
		Description  string
	}{
		// Plastic types
		{"Plastic", "PET Bottles", "Polyethylene terephthalate bottles"},
		{"Plastic", "HDPE Containers", "High-density polyethylene containers"},
		{"Plastic", "Plastic Bags", "Various plastic bags and films"},
		{"Plastic", "Styrofoam", "Expanded polystyrene foam"},

		// Paper types
		{"Paper", "Newspaper", "Daily newspapers and newsprint"},
		{"Paper", "Cardboard", "Corrugated cardboard boxes"},
		{"Paper", "Office Paper", "White and colored office paper"},
		{"Paper", "Magazines", "Glossy magazines and catalogs"},

		// Metal types
		{"Metal", "Aluminum Cans", "Beverage aluminum cans"},
		{"Metal", "Steel Cans", "Food steel cans"},
		{"Metal", "Copper Wire", "Electrical copper wiring"},
		{"Metal", "Iron Scrap", "Various iron and steel scrap"},

		// Glass types
		{"Glass", "Clear Glass Bottles", "Transparent glass bottles"},
		{"Glass", "Colored Glass Bottles", "Brown and green glass bottles"},
		{"Glass", "Glass Jars", "Food and beverage jars"},

		// Organic types
		{"Organic", "Food Waste", "Kitchen and food scraps"},
		{"Organic", "Garden Waste", "Leaves, branches, and yard trimmings"},

		// Electronic types
		{"Electronic", "Mobile Phones", "Old smartphones and feature phones"},
		{"Electronic", "Computers", "Desktop and laptop computers"},
		{"Electronic", "Televisions", "CRT and LCD televisions"},

		// Textile types
		{"Textile", "Cotton Clothing", "Used cotton garments"},
		{"Textile", "Synthetic Clothing", "Polyester and synthetic garments"},
	}

	for _, wt := range wasteTypes {
		categoryID, exists := categoryMap[wt.CategoryName]
		if !exists {
			continue
		}

		wasteType := entity.WasteType{
			ID:          uuid.New(),
			CategoryID:  categoryID,
			Name:        wt.Name,
			Description: wt.Description,
		}

		var existing entity.WasteType
		if err := db.Where("name = ? AND category_id = ?", wasteType.Name, wasteType.CategoryID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&wasteType).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
