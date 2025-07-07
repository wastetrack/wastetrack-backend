package seeder

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"gorm.io/gorm"
)

// SeedSalaryTransactions seeds salary transactions
func SeedSalaryTransactions(db *gorm.DB) error {
	var wasteBanks []entity.User
	if err := db.Where("role IN ?", []string{"waste_bank_unit", "waste_bank_central"}).Find(&wasteBanks).Error; err != nil {
		return err
	}

	var collectors []entity.User
	if err := db.Where("role IN ?", []string{"waste_collector_unit", "waste_collector_central"}).Find(&collectors).Error; err != nil {
		return err
	}

	var industries []entity.User
	if err := db.Where("role = ?", "industry").Find(&industries).Error; err != nil {
		return err
	}

	if len(wasteBanks) == 0 || len(collectors) == 0 {
		log.Println("Warning: No waste banks or collectors found, skipping salary transactions")
		return nil
	}

	lastMonth := time.Now().AddDate(0, -1, 0)
	lastWeek := time.Now().AddDate(0, 0, -7)
	yesterday := time.Now().AddDate(0, 0, -1)

	var transactions []entity.SalaryTransaction

	// Create transactions based on available users
	// First waste bank pays first collector unit
	if len(wasteBanks) >= 1 && len(collectors) >= 1 {
		transactions = append(transactions, entity.SalaryTransaction{
			ID:              uuid.New(),
			SenderID:        wasteBanks[0].ID,
			ReceiverID:      collectors[0].ID,
			Amount:          2500000,
			TransactionType: "salary",
			CreatedAt:       lastMonth,
			Status:          "completed",
			Notes:           "Monthly salary payment for waste collection services",
		})
	}

	// First waste bank pays second collector unit (if available)
	if len(wasteBanks) >= 1 && len(collectors) >= 2 {
		transactions = append(transactions, entity.SalaryTransaction{
			ID:              uuid.New(),
			SenderID:        wasteBanks[0].ID,
			ReceiverID:      collectors[1].ID,
			Amount:          2300000,
			TransactionType: "salary",
			CreatedAt:       lastMonth,
			Status:          "completed",
			Notes:           "Monthly salary payment for waste collection services",
		})
	}

	// Second waste bank pays bonus (if available)
	if len(wasteBanks) >= 2 && len(collectors) >= 1 {
		collectorIndex := 0
		if len(collectors) >= 2 {
			collectorIndex = 1
		}

		transactions = append(transactions, entity.SalaryTransaction{
			ID:              uuid.New(),
			SenderID:        wasteBanks[1].ID,
			ReceiverID:      collectors[collectorIndex].ID,
			Amount:          500000,
			TransactionType: "salary",
			CreatedAt:       lastWeek,
			Status:          "completed",
			Notes:           "Performance bonus for exceeding collection targets",
		})
	}

	// Central waste bank pays central collector (if available)
	if len(wasteBanks) >= 3 && len(collectors) >= 3 {
		transactions = append(transactions, entity.SalaryTransaction{
			ID:              uuid.New(),
			SenderID:        wasteBanks[2].ID, // Central waste bank
			ReceiverID:      collectors[2].ID, // Central collector
			Amount:          4500000,          // Higher salary for central operations
			TransactionType: "salary",
			CreatedAt:       lastMonth,
			Status:          "completed",
			Notes:           "Monthly salary for central collection team leader",
		})
	}

	// Current month salary
	if len(wasteBanks) >= 1 && len(collectors) >= 1 {
		transactions = append(transactions, entity.SalaryTransaction{
			ID:              uuid.New(),
			SenderID:        wasteBanks[0].ID,
			ReceiverID:      collectors[0].ID,
			Amount:          2500000,
			TransactionType: "salary",
			CreatedAt:       yesterday,
			Status:          "failed",
			Notes:           "Current month salary payment",
		})
	}

	// Central bank supervision bonus to unit collectors
	if len(wasteBanks) >= 3 && len(collectors) >= 2 {
		for i := 0; i < 2 && i < len(collectors); i++ { // Unit collectors only
			transactions = append(transactions, entity.SalaryTransaction{
				ID:              uuid.New(),
				SenderID:        wasteBanks[2].ID, // Central bank
				ReceiverID:      collectors[i].ID,
				Amount:          300000,
				TransactionType: "points_conversion",
				CreatedAt:       lastWeek,
				Status:          "completed",
				Notes:           "Bonus for cooperation with central coordination",
			})
		}
	}

	// Add industry to waste bank transactions if available
	if len(industries) > 0 && len(wasteBanks) > 0 {
		industryTransactions := []entity.SalaryTransaction{
			{
				ID:              uuid.New(),
				SenderID:        industries[0].ID,
				ReceiverID:      wasteBanks[0].ID,
				Amount:          15000000,
				TransactionType: "waste_payment",
				CreatedAt:       lastWeek,
				Status:          "completed",
				Notes:           "Payment for bulk waste material purchase",
			},
		}

		// Second transaction to second waste bank if available
		if len(wasteBanks) >= 2 {
			industryTransactions = append(industryTransactions, entity.SalaryTransaction{
				ID:              uuid.New(),
				SenderID:        industries[0].ID,
				ReceiverID:      wasteBanks[1].ID,
				Amount:          22000000,
				TransactionType: "waste_payment",
				CreatedAt:       time.Now().AddDate(0, 0, -3),
				Status:          "completed",
				Notes:           "Payment for recycled plastic materials",
			})
		}

		// Large central bank transaction if available
		if len(wasteBanks) >= 3 {
			industryTransactions = append(industryTransactions, entity.SalaryTransaction{
				ID:              uuid.New(),
				SenderID:        industries[0].ID,
				ReceiverID:      wasteBanks[2].ID, // Central bank
				Amount:          50000000,
				TransactionType: "waste_payment",
				CreatedAt:       time.Now().AddDate(0, 0, -5),
				Status:          "completed",
				Notes:           "Large bulk purchase from central waste bank",
			})
		}

		// Second industry to central bank if available
		if len(industries) >= 2 && len(wasteBanks) >= 3 {
			industryTransactions = append(industryTransactions, entity.SalaryTransaction{
				ID:              uuid.New(),
				SenderID:        industries[1].ID,
				ReceiverID:      wasteBanks[2].ID, // Central bank
				Amount:          35000000,
				TransactionType: "waste_payment",
				CreatedAt:       time.Now().AddDate(0, 0, -1),
				Status:          "completed",
				Notes:           "Metal recycling bulk purchase from central facility",
			})
		}

		transactions = append(transactions, industryTransactions...)
	}

	for _, transaction := range transactions {
		var existing entity.SalaryTransaction
		if err := db.Where("sender_id = ? AND receiver_id = ? AND created_at = ? AND amount = ?",
			transaction.SenderID, transaction.ReceiverID, transaction.CreatedAt, transaction.Amount).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&transaction).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
