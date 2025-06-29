package usecase

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/model/converter"
	"github.com/wastetrack/wastetrack-backend/internal/repository"
	"gorm.io/gorm"
)

type CustomerUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	CustomerRepository *repository.CustomerRepository
}

func NewCustomerUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, customerRepository *repository.CustomerRepository) *CustomerUseCase {
	return &CustomerUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		CustomerRepository: customerRepository,
	}
}

// TODO: Implement Search

func (c *CustomerUseCase) Get(ctx context.Context, request *model.GetCustomerRequest) (*model.CustomerResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	customer := new(entity.CustomerProfile)
	if err := c.CustomerRepository.FindByUserID(tx, customer, request.ID); err != nil {
		c.Log.Warnf("Failed find profile by user id : %+v", err)
		return nil, fiber.ErrNotFound
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CustomerToResponse(customer), nil

}

func (c *CustomerUseCase) Update(ctx context.Context, request *model.UpdateCustomerRequest, userID string, userRole string) (*model.CustomerResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	customer := new(entity.CustomerProfile)
	if err := c.CustomerRepository.FindById(tx, customer, request.ID); err != nil {
		c.Log.Warnf("Failed find subject by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Authorization check: Skip for admin, otherwise check ownership
	if userRole != "admin" && customer.UserID != uuid.MustParse(userID) {
		c.Log.Warnf("User %s is not authorized to update customer %s", userID, request.ID)
		return nil, fiber.ErrForbidden
	}

	if request.BagsStored != nil {
		customer.BagsStored = *request.BagsStored
	}

	if request.CarbonDeficit != nil {
		customer.CarbonDeficit = *request.CarbonDeficit
	}

	if request.Trees != nil {
		customer.Trees = *request.Trees
	}

	if request.WaterSaved != nil {
		customer.WaterSaved = *request.WaterSaved
	}

	if err := c.CustomerRepository.Update(tx, customer); err != nil {
		c.Log.Warnf("Failed to update customer: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CustomerToResponse(customer), nil
}

func (c *CustomerUseCase) Delete(ctx context.Context, request *model.DeleteWasteBankRequest) (*model.CustomerResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find customer by id
	customer := new(entity.CustomerProfile)
	if err := c.CustomerRepository.FindById(tx, customer, request.ID); err != nil {
		c.Log.Warnf("Failed find customer by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Delete customer
	if err := c.CustomerRepository.Delete(tx, customer); err != nil {
		c.Log.Warnf("Failed delete customer : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CustomerToResponse(customer), nil
}
