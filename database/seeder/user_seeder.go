package seeder

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedUsers seeds the users table with various roles
func SeedUsers(db *gorm.DB) error {
	// Hash password for all users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []entity.User{
		// Customers
		{
			ID:              uuid.New(),
			Username:        "john_customer",
			Email:           "john.customer@example.com",
			Password:        string(hashedPassword),
			Role:            "customer",
			PhoneNumber:     "+628123456789",
			Institution:     "",
			Address:         "Jl. Merdeka No. 1",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          150,
			Balance:         50000,
			Location:        &types.Point{Lat: -7.2504, Lng: 112.7688},
			IsEmailVerified: true,
		},
		{
			ID:              uuid.New(),
			Username:        "jane_customer",
			Email:           "jane.customer@example.com",
			Password:        string(hashedPassword),
			Role:            "customer",
			PhoneNumber:     "+628123456790",
			Institution:     "",
			Address:         "Jl. Pemuda No. 15",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          200,
			Balance:         75000,
			Location:        &types.Point{Lat: -7.2575, Lng: 112.7521},
			IsEmailVerified: true,
		},
		{
			ID:              uuid.New(),
			Username:        "bob_customer",
			Email:           "bob.customer@example.com",
			Password:        string(hashedPassword),
			Role:            "customer",
			PhoneNumber:     "+628123456791",
			Institution:     "",
			Address:         "Jl. Diponegoro No. 25",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          80,
			Balance:         25000,
			Location:        &types.Point{Lat: -7.2656, Lng: 112.7431},
			IsEmailVerified: true,
		},

		// Waste Bank Units
		{
			ID:              uuid.New(),
			Username:        "green_waste_bank",
			Email:           "info@greenwaste.com",
			Password:        string(hashedPassword),
			Role:            "waste_bank_unit",
			PhoneNumber:     "+628123456792",
			Institution:     "Green Waste Bank",
			Address:         "Jl. Kertajaya No. 10",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          0,
			Balance:         500000,
			Location:        &types.Point{Lat: -7.2819, Lng: 112.7958},
			IsEmailVerified: true,
		},
		{
			ID:              uuid.New(),
			Username:        "eco_waste_center",
			Email:           "contact@ecowaste.com",
			Password:        string(hashedPassword),
			Role:            "waste_bank_unit",
			PhoneNumber:     "+628123456793",
			Institution:     "Eco Waste Center",
			Address:         "Jl. Ahmad Yani No. 50",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          0,
			Balance:         750000,
			Location:        &types.Point{Lat: -7.2456, Lng: 112.7378},
			IsEmailVerified: true,
		},

		// Waste Bank Central
		{
			ID:              uuid.New(),
			Username:        "central_waste_hub",
			Email:           "central@wastehub.com",
			Password:        string(hashedPassword),
			Role:            "waste_bank_central",
			PhoneNumber:     "+628123456799",
			Institution:     "Central Waste Management Hub",
			Address:         "Jl. Industri Central No. 1",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          0,
			Balance:         2000000,
			Location:        &types.Point{Lat: -7.2389, Lng: 112.7589},
			IsEmailVerified: true,
		},

		// Waste Collector Units
		{
			ID:              uuid.New(),
			Username:        "collector_ahmad",
			Email:           "ahmad.collector@example.com",
			Password:        string(hashedPassword),
			Role:            "waste_collector_unit",
			PhoneNumber:     "+628123456794",
			Institution:     "",
			Address:         "Jl. Gubeng No. 5",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          0,
			Balance:         100000,
			Location:        &types.Point{Lat: -7.2653, Lng: 112.7536},
			IsEmailVerified: true,
		},
		{
			ID:              uuid.New(),
			Username:        "collector_siti",
			Email:           "siti.collector@example.com",
			Password:        string(hashedPassword),
			Role:            "waste_collector_unit",
			PhoneNumber:     "+628123456795",
			Institution:     "",
			Address:         "Jl. Raya Darmo No. 12",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          0,
			Balance:         85000,
			Location:        &types.Point{Lat: -7.2733, Lng: 112.7319},
			IsEmailVerified: true,
		},

		// Waste Collector Central
		{
			ID:              uuid.New(),
			Username:        "central_collector_team",
			Email:           "team@centralcollect.com",
			Password:        string(hashedPassword),
			Role:            "waste_collector_central",
			PhoneNumber:     "+628123456800",
			Institution:     "Central Collection Team",
			Address:         "Jl. Central Operations No. 8",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          0,
			Balance:         250000,
			Location:        &types.Point{Lat: -7.2289, Lng: 112.7489},
			IsEmailVerified: true,
		},

		// Government
		{
			ID:              uuid.New(),
			Username:        "surabaya_gov",
			Email:           "waste.management@surabaya.go.id",
			Password:        string(hashedPassword),
			Role:            "government",
			PhoneNumber:     "+628123456796",
			Institution:     "Surabaya City Government",
			Address:         "Jl. Taman Surya No. 1",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          0,
			Balance:         0,
			Location:        &types.Point{Lat: -7.2456, Lng: 112.7378},
			IsEmailVerified: true,
		},

		// Industry
		{
			ID:              uuid.New(),
			Username:        "plastic_industry",
			Email:           "sustainability@plasticorp.com",
			Password:        string(hashedPassword),
			Role:            "industry",
			PhoneNumber:     "+628123456797",
			Institution:     "Plasticorp Industries",
			Address:         "Jl. Industri Raya No. 100",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          0,
			Balance:         1000000,
			Location:        &types.Point{Lat: -7.3049, Lng: 112.7378},
			IsEmailVerified: true,
		},
		{
			ID:              uuid.New(),
			Username:        "metal_recycling",
			Email:           "ops@metalrecycle.com",
			Password:        string(hashedPassword),
			Role:            "industry",
			PhoneNumber:     "+628123456798",
			Institution:     "Metal Recycling Co.",
			Address:         "Jl. Logam No. 25",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          0,
			Balance:         800000,
			Location:        &types.Point{Lat: -7.2189, Lng: 112.6319},
			IsEmailVerified: true,
		},
	}

	for _, user := range users {
		var existing entity.User
		if err := db.Where("email = ?", user.Email).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&user).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
