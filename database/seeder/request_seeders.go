package seeder

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/types"
	"gorm.io/gorm"
)

// Helper function to create TimeOnly from hour and minut

// SeedWasteDropRequests seeds waste drop requests
func SeedWasteDropRequests(db *gorm.DB) error {
	var customers []entity.User
	if err := db.Where("role = ?", "customer").Find(&customers).Error; err != nil {
		return err
	}

	var wasteBanks []entity.User
	if err := db.Where("role IN ?", []string{"waste_bank_unit", "waste_bank_central"}).Find(&wasteBanks).Error; err != nil {
		return err
	}

	var collectors []entity.User
	if err := db.Where("role IN ?", []string{"waste_collector_unit", "waste_collector_central"}).Find(&collectors).Error; err != nil {
		return err
	}

	if len(customers) == 0 || len(wasteBanks) == 0 {
		log.Println("Warning: No customers or waste banks found, skipping waste drop requests")
		return nil
	}

	tomorrow := time.Now().AddDate(0, 0, 1)
	nextWeek := time.Now().AddDate(0, 0, 7)

	var requests []entity.WasteDropRequest

	// Create requests based on available users
	if len(customers) >= 1 && len(wasteBanks) >= 1 {
		var assignedCollectorID *uuid.UUID
		if len(collectors) >= 1 {
			assignedCollectorID = &collectors[0].ID
		}

		requests = append(requests, entity.WasteDropRequest{
			ID:                   uuid.New(),
			DeliveryType:         "pickup",
			CustomerID:           customers[0].ID,
			UserPhoneNumber:      customers[0].PhoneNumber,
			WasteBankID:          &wasteBanks[0].ID,
			AssignedCollectorID:  assignedCollectorID,
			TotalPrice:           45000,
			Status:               "assigned",
			AppointmentLocation:  &types.Point{Lat: -7.2504, Lng: 112.7688},
			AppointmentDate:      tomorrow,
			AppointmentStartTime: createTimeOnly(10, 0),
			AppointmentEndTime:   createTimeOnly(11, 0),
			Notes:                "Please call before arriving",
		})
	}

	if len(customers) >= 2 && len(wasteBanks) >= 2 {
		requests = append(requests, entity.WasteDropRequest{
			ID:                   uuid.New(),
			DeliveryType:         "dropoff",
			CustomerID:           customers[1].ID,
			UserPhoneNumber:      customers[1].PhoneNumber,
			WasteBankID:          &wasteBanks[1].ID,
			TotalPrice:           62000,
			Status:               "pending",
			AppointmentLocation:  &types.Point{Lat: -7.2456, Lng: 112.7378},
			AppointmentDate:      time.Now().AddDate(0, 0, -2),
			AppointmentStartTime: createTimeOnly(14, 0),
			AppointmentEndTime:   createTimeOnly(15, 0),
			Notes:                "Regular customer dropoff",
		})
	}

	if len(customers) >= 3 && len(wasteBanks) >= 1 {
		var assignedCollectorID *uuid.UUID
		if len(collectors) >= 2 {
			assignedCollectorID = &collectors[1].ID
		} else if len(collectors) >= 1 {
			assignedCollectorID = &collectors[0].ID
		}

		customerIndex := 2
		if len(customers) < 3 {
			customerIndex = 0 // Fallback to first customer
		}

		requests = append(requests, entity.WasteDropRequest{
			ID:                   uuid.New(),
			DeliveryType:         "pickup",
			CustomerID:           customers[customerIndex].ID,
			UserPhoneNumber:      customers[customerIndex].PhoneNumber,
			WasteBankID:          &wasteBanks[0].ID,
			AssignedCollectorID:  assignedCollectorID,
			TotalPrice:           28000,
			Status:               "assigned",
			AppointmentLocation:  &types.Point{Lat: -7.2656, Lng: 112.7431},
			AppointmentDate:      nextWeek,
			AppointmentStartTime: createTimeOnly(9, 30),
			AppointmentEndTime:   createTimeOnly(10, 30),
			Notes:                "Large quantity pickup",
		})
	}

	// Central collector handling bulk request
	if len(customers) >= 1 && len(wasteBanks) >= 3 && len(collectors) >= 3 {
		requests = append(requests, entity.WasteDropRequest{
			ID:                   uuid.New(),
			DeliveryType:         "pickup",
			CustomerID:           customers[0].ID,
			UserPhoneNumber:      customers[0].PhoneNumber,
			WasteBankID:          &wasteBanks[2].ID, // Central waste bank
			AssignedCollectorID:  &collectors[2].ID, // Central collector
			TotalPrice:           125000,
			Status:               "pending",
			AppointmentLocation:  &types.Point{Lat: -7.2389, Lng: 112.7589},
			AppointmentDate:      time.Now().AddDate(0, 0, 3),
			AppointmentStartTime: createTimeOnly(8, 0),
			AppointmentEndTime:   createTimeOnly(10, 0),
			Notes:                "Large bulk collection by central team",
		})
	}

	for _, request := range requests {
		var existing entity.WasteDropRequest
		if err := db.Where("customer_id = ? AND appointment_date = ?", request.CustomerID, request.AppointmentDate).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&request).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedWasteDropRequestItems seeds items for waste drop requests
func SeedWasteDropRequestItems(db *gorm.DB) error {
	var requests []entity.WasteDropRequest
	if err := db.Find(&requests).Error; err != nil {
		return err
	}

	var wasteTypes []entity.WasteType
	if err := db.Find(&wasteTypes).Error; err != nil {
		return err
	}

	if len(requests) == 0 || len(wasteTypes) == 0 {
		return nil
	}

	// Create items for each request
	requestItems := []struct {
		RequestIndex     int
		WasteTypeName    string
		Quantity         int64
		VerifiedWeight   float64
		VerifiedSubtotal int64
	}{
		{0, "PET Bottles", 15, 4.5, 13500},
		{0, "Aluminum Cans", 8, 1.2, 18000},
		{0, "Cardboard", 3, 8.7, 13050},
		{1, "Office Paper", 10, 12.4, 24800},
		{1, "Steel Cans", 6, 4.6, 36800},
		{2, "HDPE Containers", 5, 6.2, 15500},
		{2, "Glass Jars", 4, 3.8, 12500},
	}

	for _, item := range requestItems {
		if item.RequestIndex >= len(requests) {
			continue
		}

		// Find the waste type
		var wasteType entity.WasteType
		if err := db.Where("name = ?", item.WasteTypeName).First(&wasteType).Error; err != nil {
			continue
		}

		requestItem := entity.WasteDropRequestItem{
			ID:               uuid.New(),
			RequestID:        requests[item.RequestIndex].ID,
			WasteTypeID:      wasteType.ID,
			Quantity:         item.Quantity,
			VerifiedWeight:   item.VerifiedWeight,
			VerifiedSubtotal: item.VerifiedSubtotal,
		}

		var existing entity.WasteDropRequestItem
		if err := db.Where("request_id = ? AND waste_type_id = ?", requestItem.RequestID, requestItem.WasteTypeID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&requestItem).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedWasteTransferRequests seeds waste transfer requests
func SeedWasteTransferRequests(db *gorm.DB) error {
	var wasteBanks []entity.User
	if err := db.Where("role IN ?", []string{"waste_bank_unit", "waste_bank_central"}).Find(&wasteBanks).Error; err != nil {
		return err
	}

	var industries []entity.User
	if err := db.Where("role = ?", "industry").Find(&industries).Error; err != nil {
		return err
	}

	if len(wasteBanks) < 2 && len(industries) == 0 {
		log.Println("Warning: Not enough waste banks or industries for transfer requests")
		return nil
	}

	nextWeek := time.Now().AddDate(0, 0, 7)
	nextMonth := time.Now().AddDate(0, 1, 0)

	var transfers []entity.WasteTransferRequest

	// Bank to bank
	if len(wasteBanks) >= 2 {
		transfers = append(transfers, entity.WasteTransferRequest{
			ID:                     uuid.New(),
			SourceUserID:           wasteBanks[0].ID,
			DestinationUserID:      wasteBanks[1].ID,
			FormType:               "waste_bank_request",
			TotalWeight:            250,
			TotalPrice:             750000,
			Status:                 "pending",
			SourcePhoneNumber:      wasteBanks[0].PhoneNumber,
			DestinationPhoneNumber: wasteBanks[1].PhoneNumber,
			AppointmentDate:        nextWeek,
			AppointmentStartTime:   createTimeOnly(10, 0),
			AppointmentEndTime:     createTimeOnly(12, 0),
			ImageURL:               "https://example.com/images/transfer1.jpg",
			Notes:                  "Weekly redistribution of collected recyclables",
			AppointmentLocation:    &types.Point{Lat: -6.200000, Lng: 106.816666}, // Jakarta
		})
	}

	// Bank to industry
	if len(wasteBanks) >= 1 && len(industries) >= 1 {
		bankIndex := 0
		if len(wasteBanks) >= 2 {
			bankIndex = 1
		}

		transfers = append(transfers, entity.WasteTransferRequest{
			ID:                     uuid.New(),
			SourceUserID:           wasteBanks[bankIndex].ID,
			DestinationUserID:      industries[0].ID,
			FormType:               "industry_request",
			TotalWeight:            500,
			TotalPrice:             1500000,
			Status:                 "cancelled",
			SourcePhoneNumber:      wasteBanks[bankIndex].PhoneNumber,
			DestinationPhoneNumber: industries[0].PhoneNumber,
			AppointmentDate:        nextMonth,
			AppointmentStartTime:   createTimeOnly(8, 0),
			AppointmentEndTime:     createTimeOnly(10, 0),
			ImageURL:               "https://example.com/images/transfer2.jpg",
			Notes:                  "Plastic packaging and bottles from city collection",
			AppointmentLocation:    &types.Point{Lat: -6.914744, Lng: 107.609810}, // Bandung
		})
	}

	// Central bank to industry
	if len(wasteBanks) >= 3 && len(industries) >= 1 {
		industryIndex := 0
		if len(industries) >= 2 {
			industryIndex = 1
		}

		transfers = append(transfers, entity.WasteTransferRequest{
			ID:                     uuid.New(),
			SourceUserID:           wasteBanks[2].ID,
			DestinationUserID:      industries[industryIndex].ID,
			FormType:               "waste_bank_request",
			TotalWeight:            1000,
			TotalPrice:             3000000,
			Status:                 "pending",
			SourcePhoneNumber:      wasteBanks[2].PhoneNumber,
			DestinationPhoneNumber: industries[industryIndex].PhoneNumber,
			AppointmentDate:        nextMonth.AddDate(0, 0, 7),
			AppointmentStartTime:   createTimeOnly(9, 0),
			AppointmentEndTime:     createTimeOnly(12, 0),
			ImageURL:               "https://example.com/images/transfer3.jpg",
			Notes:                  "Bulk transfer for recycling plant processing",
			AppointmentLocation:    &types.Point{Lat: -7.250445, Lng: 112.768845}, // Surabaya
		})
	}

	for _, transfer := range transfers {
		var existing entity.WasteTransferRequest
		if err := db.Where("source_user_id = ? AND destination_user_id = ? AND appointment_date = ?",
			transfer.SourceUserID, transfer.DestinationUserID, transfer.AppointmentDate).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&transfer).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedWasteTransferItemOfferings seeds items for waste transfer requests
func SeedWasteTransferItemOfferings(db *gorm.DB) error {
	var transfers []entity.WasteTransferRequest
	if err := db.Find(&transfers).Error; err != nil {
		return err
	}

	var wasteTypes []entity.WasteType
	if err := db.Find(&wasteTypes).Error; err != nil {
		return err
	}

	if len(transfers) == 0 || len(wasteTypes) == 0 {
		return nil
	}

	// Create items for each transfer
	transferItems := []struct {
		TransferIndex       int
		WasteTypeName       string
		OfferingWeight      float64
		OfferingPricePerKgs float64
		AcceptedWeight      float64
		AcceptedPricePerKgs float64
	}{
		{0, "PET Bottles", 120.5, 3000, 120.5, 3000},
		{0, "Aluminum Cans", 45.2, 15000, 45.2, 15000},
		{0, "Cardboard", 84.3, 1500, 84.3, 1500},
		{1, "HDPE Containers", 180.4, 2500, 180.4, 2500},
		{1, "Steel Cans", 95.6, 8000, 95.6, 8000},
		{1, "Office Paper", 224.0, 2000, 224.0, 2000},
	}

	for _, item := range transferItems {
		if item.TransferIndex >= len(transfers) {
			continue
		}

		// Find the waste type
		var wasteType entity.WasteType
		if err := db.Where("name = ?", item.WasteTypeName).First(&wasteType).Error; err != nil {
			continue
		}

		transferItem := entity.WasteTransferItemOffering{
			ID:                  uuid.New(),
			TransferFormID:      transfers[item.TransferIndex].ID,
			WasteTypeID:         wasteType.ID,
			OfferingWeight:      item.OfferingWeight,
			OfferingPricePerKgs: item.OfferingPricePerKgs,
			AcceptedWeight:      item.AcceptedWeight,
			AcceptedPricePerKgs: item.AcceptedPricePerKgs,
		}

		var existing entity.WasteTransferItemOffering
		if err := db.Where("transfer_request_id = ? AND waste_type_id = ?", transferItem.TransferFormID, transferItem.WasteTypeID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&transferItem).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
