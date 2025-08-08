package usecase

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/helper"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/repository"
	"gorm.io/gorm"
)

type GovernmentUseCase struct {
	DB                             *gorm.DB
	Log                            *logrus.Logger
	Validate                       *validator.Validate
	UserRepository                 *repository.UserRepository
	WasteDropRequestItemRepository *repository.WasteDropRequestItemRepository
	WasteTransferItemRepository    *repository.WasteTransferItemOfferingRepository
	WasteTransferRequestRepository *repository.WasteTransferRequestRepository
	StorageRepository              *repository.StorageRepository
}

func NewGovernmentUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	userRepository *repository.UserRepository,
	wasteDropRequestItemRepository *repository.WasteDropRequestItemRepository,
	wasteTransferItemRepository *repository.WasteTransferItemOfferingRepository,
	wasteTransferRequestRepository *repository.WasteTransferRequestRepository,
	storageRepository *repository.StorageRepository,
) *GovernmentUseCase {
	return &GovernmentUseCase{
		DB:                             db,
		Log:                            log,
		Validate:                       validate,
		UserRepository:                 userRepository,
		WasteDropRequestItemRepository: wasteDropRequestItemRepository,
		WasteTransferItemRepository:    wasteTransferItemRepository,
		WasteTransferRequestRepository: wasteTransferRequestRepository,
		StorageRepository:              storageRepository,
	}
}

func (uc *GovernmentUseCase) GetDashboard(request *model.GovernmentDashboardRequest) (*model.GovernmentDashboardResponse, error) {
	// Validate request using helper functions
	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.Warnf("Failed to validate government dashboard request: %v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("validation failed: %v", err))
	}

	// Validate date formats using helper
	if err := helper.ValidateMonthFormat(request.StartMonth); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := helper.ValidateMonthFormat(request.EndMonth); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Validate date range using helper
	if err := helper.ValidateDateRange(request.StartMonth, request.EndMonth); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Log the dashboard request
	uc.Log.Infof("Getting government dashboard data for period %s to %s", request.StartMonth, request.EndMonth)

	// Get basic statistics (total bank sampah and total offtakers)
	totalBankSampah, err := uc.UserRepository.CountWasteBanks(uc.DB, request)
	if err != nil {
		uc.Log.Errorf("Failed to count waste banks: %v", err)

		// Check if it's a validation error from repository
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "format") {
			return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		// Database or other internal errors
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve waste bank data")
	}

	totalOfftakers, err := uc.UserRepository.CountOfftakers(uc.DB, request)
	if err != nil {
		uc.Log.Errorf("Failed to count offtakers: %v", err)

		// Check if it's a validation error from repository
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "format") {
			return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		// Database or other internal errors
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve offtaker data")
	}

	// Get total collected waste
	totalCollected, err := uc.WasteDropRequestItemRepository.GetTotalCollectedWaste(uc.DB, request)
	if err != nil {
		uc.Log.Errorf("Failed to get total collected waste: %v", err)

		// Check if it's a validation error from repository
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "format") {
			return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		// Database or other internal errors
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve collected waste data")
	}

	// Get collection trends by role (combining waste drop and waste transfer)
	var collectionTrends []model.CollectionTrendByRole
	if request.StartMonth != "" && request.EndMonth != "" {
		// Get waste drop trends
		wasteDropTrends, err := uc.WasteDropRequestItemRepository.GetCollectionTrendsByRole(uc.DB, request)
		if err != nil {
			uc.Log.Errorf("Failed to get waste drop collection trends: %v", err)
			if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "format") {
				return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve waste drop trends data")
		}

		// Get waste transfer trends
		wasteTransferTrends, err := uc.WasteTransferItemRepository.GetWasteTransferTrendsByRole(uc.DB, request)
		if err != nil {
			uc.Log.Errorf("Failed to get waste transfer collection trends: %v", err)
			if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "format") {
				return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve waste transfer trends data")
		}

		// Get waste drop verified weights by month
		wasteDropWeights, err := uc.WasteDropRequestItemRepository.GetWasteDropVerifiedWeightByMonth(uc.DB, request)
		if err != nil {
			uc.Log.Errorf("Failed to get waste drop verified weights: %v", err)
			if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "format") {
				return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve waste drop weight data")
		}

		// Get waste transfer verified weights by month
		wasteTransferWeights, err := uc.WasteTransferItemRepository.GetTransferVerifiedWeightByMonth(uc.DB, request)
		if err != nil {
			uc.Log.Errorf("Failed to get waste transfer verified weights: %v", err)
			if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "format") {
				return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve waste transfer weight data")
		}

		// Combine trends from both sources
		combinedTrends := uc.combineTrends(wasteDropTrends, wasteTransferTrends)
		collectionTrends = uc.formatTrendsResponse(combinedTrends, wasteDropWeights, wasteTransferWeights, request.StartMonth, request.EndMonth)
	}

	// Get top offtakers
	topOfftakers, err := uc.WasteTransferRequestRepository.GetTopOfftakers(uc.DB, request)
	if err != nil {
		uc.Log.Errorf("Failed to get top offtakers: %v", err)
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "format") {
			return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve top offtakers data")
	}
	largestBanks, err := uc.GetLargestBanks(request)
	if err != nil {
		uc.Log.Errorf("Failed to get largest waste banks: %v", err)
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "format") {
			return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve largest banks data")
	}

	// Build complete response
	response := &model.GovernmentDashboardResponse{
		TotalBankSampah:  totalBankSampah,
		TotalOfftaker:    totalOfftakers,
		TotalCollected:   totalCollected,
		CollectionTrends: collectionTrends,
		TopOfftakers:     topOfftakers,
		LargestBanks:     largestBanks, // NOW PROPERLY ASSIGNED
	}

	uc.Log.Infof("Successfully retrieved complete dashboard data: %d bank sampah, %d offtakers, %.2f kg collected, %d trend months, %d top offtakers, %d largest banks",
		totalBankSampah, totalOfftakers, totalCollected, len(collectionTrends), len(topOfftakers), len(largestBanks))

	return response, nil
}

// BankScore holds temporary calculation data
type BankScore struct {
	User        *entity.User
	Volume      float64
	TotalWeight float64
	Score       float64
}

func (uc *GovernmentUseCase) GetLargestBanks(request *model.GovernmentDashboardRequest) ([]model.LargestBank, error) {
	// 1. Get all waste bank users (single source of truth for user data)
	users, err := uc.UserRepository.GetWasteBankUsers(uc.DB, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get waste bank users: %v", err)
	}

	// 2. Get storage volumes (storage repository responsibility)
	storageVolumes, err := uc.StorageRepository.GetWasteBankStorageVolumes(uc.DB, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage volumes: %v", err)
	}

	// 3. Get waste drop totals (waste drop repository responsibility)
	wasteDropTotals, err := uc.WasteDropRequestItemRepository.GetWasteBankWasteDropTotals(uc.DB, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get waste drop totals: %v", err)
	}

	// 4. Get waste transfer totals (waste transfer repository responsibility)
	wasteTransferTotals, err := uc.WasteTransferRequestRepository.GetWasteBankTransferTotals(uc.DB, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get waste transfer totals: %v", err)
	}

	// 5. BUSINESS LOGIC: Combine data and calculate scores (UseCase responsibility)
	bankScores := make([]BankScore, 0)

	for userID, user := range users {
		volume := storageVolumes[userID]
		dropWeight := wasteDropTotals[userID]
		transferWeight := wasteTransferTotals[userID]
		totalWeight := dropWeight + transferWeight

		// Only include banks with some activity or storage
		if volume > 0 || totalWeight > 0 {
			// BUSINESS RULE: 70% processing capacity + 30% storage volume
			score := totalWeight*0.3 + volume*0.7

			bankScores = append(bankScores, BankScore{
				User:        user,
				Volume:      volume,
				TotalWeight: totalWeight,
				Score:       score,
			})
		}
	}

	// 6. BUSINESS LOGIC: Sort by composite score (UseCase responsibility)
	sort.Slice(bankScores, func(i, j int) bool {
		if bankScores[i].Score != bankScores[j].Score {
			return bankScores[i].Score > bankScores[j].Score // Higher score first
		}
		if bankScores[i].Volume != bankScores[j].Volume {
			return bankScores[i].Volume > bankScores[j].Volume // Higher volume first
		}
		return bankScores[i].TotalWeight > bankScores[j].TotalWeight // Higher weight first
	})

	// 7. BUSINESS RULE: Return top 3 (UseCase responsibility)
	limit := 3
	if len(bankScores) < limit {
		limit = len(bankScores)
	}

	// 8. Convert to response model (UseCase responsibility)
	largestBanks := make([]model.LargestBank, limit)
	for i := 0; i < limit; i++ {
		bank := bankScores[i]
		largestBanks[i] = model.LargestBank{
			ID:          bank.User.ID.String(),
			Name:        uc.getUserName(bank.User), // Use helper method
			Institution: bank.User.Institution,
			City:        bank.User.City,
			Province:    bank.User.Province,
			TotalWeight: bank.TotalWeight, // Total waste processed
			Volume:      bank.Volume,      // Storage volume
		}
	}

	return largestBanks, nil
}

// Helper method to safely get user name
func (uc *GovernmentUseCase) getUserName(user *entity.User) string {
	if user.Username != "" {
		return user.Username
	}
	if user.Username != "" {
		return user.Username
	}
	return "Unknown"
}
func (uc *GovernmentUseCase) combineTrends(wasteDropTrends, wasteTransferTrends map[string]map[string]float64) map[string]map[string]float64 {
	combined := make(map[string]map[string]float64)

	// Add waste drop trends
	for month, roles := range wasteDropTrends {
		if combined[month] == nil {
			combined[month] = make(map[string]float64)
		}
		for role, weight := range roles {
			combined[month][role] += weight
		}
	}

	// Add waste transfer trends
	for month, roles := range wasteTransferTrends {
		if combined[month] == nil {
			combined[month] = make(map[string]float64)
		}
		for role, weight := range roles {
			combined[month][role] += weight
		}
	}

	return combined
}

// Helper method to safely get user name
func (uc *GovernmentUseCase) formatTrendsResponse(
	combinedTrends map[string]map[string]float64,
	wasteDropWeights map[string]float64,
	wasteTransferWeights map[string]float64,
	startMonth, endMonth string,
) []model.CollectionTrendByRole {
	var trends []model.CollectionTrendByRole

	// Use helper to generate month range
	months, err := helper.ParseMonthRange(startMonth, endMonth)
	if err != nil {
		uc.Log.Errorf("Failed to parse month range: %v", err)
		return trends
	}

	// Generate trends for each month in range
	for _, monthStr := range months {
		trend := model.CollectionTrendByRole{
			Month:              monthStr,
			WasteBankUnit:      0,
			WasteBankCentral:   0,
			Industry:           0,
			CollectionRequests: 0,
			TransferRequests:   0,
			TotalAmount:        0, // Initialize TotalAmount
		}

		// Fill in actual weight data by role if available
		if monthData, exists := combinedTrends[monthStr]; exists {
			trend.WasteBankUnit = monthData["waste_bank_unit"]
			trend.WasteBankCentral = monthData["waste_bank_central"]
			trend.Industry = monthData["industry"]
		}

		// Fill in collection and transfer weights
		if dropWeight, exists := wasteDropWeights[monthStr]; exists {
			trend.CollectionRequests = dropWeight
		}

		if transferWeight, exists := wasteTransferWeights[monthStr]; exists {
			trend.TransferRequests = transferWeight
		}

		// Calculate total amount
		trend.TotalAmount = trend.CollectionRequests + trend.TransferRequests

		trends = append(trends, trend)
	}

	return trends
}
