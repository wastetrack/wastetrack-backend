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
		{
			ID:              uuid.New(),
			Username:        "wasty_si_nasabah",
			Email:           "nasabah@wt.id",
			Password:        string(hashedPassword),
			Role:            "customer",
			PhoneNumber:     "+628123456789",
			Institution:     "",
			Address:         "Asrama ITS, Jalan Teknik Elektro, RW 04, Keputih, Sukolilo, Surabaya, East Java, Java, 60111, Indonesia",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          150,
			Balance:         50000,
			Location:        &types.Point{Lat: -7.2504, Lng: 112.7688},
			IsEmailVerified: true,
		},
		{
			ID:              uuid.New(),
			Username:        "adi_suberkah",
			Email:           "bsu@wt.id",
			Password:        string(hashedPassword),
			Role:            "waste_bank_unit",
			PhoneNumber:     "+628123456789",
			Institution:     "BSU Hijau Berkah",
			Address:         "Institut Teknologi Sepuluh Nopember, Jalan Keputih Perintis I A, RW 03, Keputih, Sukolilo, Surabaya, Jawa Timur, 60111, Indonesia",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          0,
			Balance:         500000,
			Location:        &types.Point{Lat: -7.2819, Lng: 112.7958},
			IsEmailVerified: true,
		},
		{
			ID:              uuid.New(),
			Username:        "ega_basuka",
			Email:           "pegawaibsu@wt.id",
			Password:        string(hashedPassword),
			Role:            "waste_collector_unit",
			PhoneNumber:     "+628123456789",
			Institution:     "BSU Hijau Berkah",
			Address:         "Institut Teknologi Sepuluh Nopember, Jalan Kejawen Putih Tambak II, RW 03, Kejawen Putih Tambak, Mulyorejo, Surabaya, East Java, 60112, Indonesia",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          0,
			Balance:         0,
			Location:        &types.Point{Lat: -7.2456, Lng: 112.7378},
			IsEmailVerified: true,
		},
		{
			ID:              uuid.New(),
			Username:        "mina_astiya",
			Email:           "bsi@wt.id",
			Password:        string(hashedPassword),
			Role:            "waste_bank_central",
			PhoneNumber:     "+628123456789",
			Institution:     "BSI Raya Surabaya",
			Address:         "Dinas Kebersihan dan Ruang Terbuka Hijau Kota Surabaya, Jalan Sukodami III, RW 07, Manyar Sabrangan, Mulyorejo, Surabaya, East Java, 60282, Indonesia",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          0,
			Balance:         750000,
			Location:        &types.Point{Lat: -7.2456, Lng: 112.7378},
			IsEmailVerified: true,
		},
		{
			ID:              uuid.New(),
			Username:        "andi_diata",
			Email:           "bsi2@wt.id",
			Password:        string(hashedPassword),
			Role:            "waste_bank_central",
			PhoneNumber:     "+628123456789",
			Institution:     "BSI Ageng Buana Surabaya",
			Address:         "Jalan Puri Jambangan I, RW 11, Karah, Jambangan, Surabaya, East Java, 60223, Indonesia",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          0,
			Balance:         2000000,
			Location:        &types.Point{Lat: -7.2389, Lng: 112.7589},
			IsEmailVerified: true,
		},
		{
			ID:              uuid.New(),
			Username:        "awai_sina",
			Email:           "pegawaibsi@wt.id",
			Password:        string(hashedPassword),
			Role:            "waste_collector_central",
			PhoneNumber:     "+628123456789",
			Institution:     "BSI Raya Surabaya",
			Address:         "",
			City:            "",
			Province:        "East Java",
			Points:          0,
			Balance:         100000,
			Location:        &types.Point{Lat: -7.2653, Lng: 112.7536},
			IsEmailVerified: true,
		},
		{
			ID:              uuid.New(),
			Username:        "asi_bunaya",
			Email:           "pegawaibsi2@wt.id",
			Password:        string(hashedPassword),
			Role:            "waste_collector_central",
			PhoneNumber:     "+628123456789",
			Institution:     "BSI Ageng Buana Surabaya",
			Address:         "Jalan Royal Ketintang Regency Blok G‑H, Royal Ketintang Regency, Karah, Jambangan, Surabaya, East Java, 60223, Indonesia",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          0,
			Balance:         250000,
			Location:        &types.Point{Lat: -7.2733, Lng: 112.7319},
			IsEmailVerified: true,
		},
		{
			ID:              uuid.New(),
			Username:        "ofi_takena",
			Email:           "offtaker@wt.id",
			Password:        string(hashedPassword),
			Role:            "industry",
			PhoneNumber:     "+628123456789",
			Institution:     "Offtaker Eko Subur Langgeng",
			Address:         "Jalan Pandugo Baru V, Wisma Penjaringan Sari, Penjaringan Sari, Rungkut, Surabaya, East Java, 60297, Indonesia",
			City:            "Surabaya",
			Province:        "East Java",
			Points:          0,
			Balance:         250000,
			Location:        &types.Point{Lat: -7.2289, Lng: 112.7489},
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
