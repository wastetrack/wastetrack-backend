package config

import (
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http/middleware"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http/route"
	"github.com/wastetrack/wastetrack-backend/internal/helper"
	"github.com/wastetrack/wastetrack-backend/internal/repository"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {

	// Setup repositories
	userRepository := repository.NewUserRepository(config.Log)
	refreshTokenRepository := repository.NewRefreshTokenRepository(config.Log)
	customerRepository := repository.NewCustomerRepository(config.Log)
	wasteBankRepository := repository.NewWasteBankRepository(config.Log)
	wasteCollectorRepository := repository.NewWasteCollectorRepository(config.Log)
	industryRepository := repository.NewIndustryRepository(config.Log)
	wasteCategoryRepository := repository.NewWasteCategoryRepository(config.Log)
	wasteTypeRepository := repository.NewWasteTypeRepository(config.Log)
	wasteBankPricedTypeRepository := repository.NewWasteBankPricedTypeRepository(config.Log)
	wasteDropRequestRepository := repository.NewWasteDropRequestRepository(config.Log)
	wasteDropRequesItemRepository := repository.NewWasteDropRequestItemRepository(config.Log)
	collectorManagementRepository := repository.NewCollectorManagementRepository(config.Log)

	// Setup JWT Helper
	jwtHelper := helper.NewJWTHelper(
		config.Config.GetString("jwt.secret_key"),                     // JWT secret for access tokens
		config.Config.GetString("jwt.refresh_secret_key"),             // JWT secret for refresh tokens
		config.Config.GetDuration("jwt.access_token_ttl")*time.Minute, // Access token TTL (e.g., 15 minutes)
		config.Config.GetDuration("jwt.refresh_token_ttl")*time.Hour,  // Refresh token TTL (e.g., 24 hours)
		refreshTokenRepository,
	)

	// Setup Email Helper
	emailHelper := helper.NewEmailHelper(
		config.Config.GetString("email.smtp_host"),     // SMTP server host (e.g., smtp.gmail.com)
		config.Config.GetString("email.smtp_port"),     // SMTP server port (e.g., 587)
		config.Config.GetString("email.smtp_username"), // SMTP username (email address)
		config.Config.GetString("email.smtp_password"), // SMTP password (app password for Gmail)
		config.Config.GetString("email.from_email"),    // From email address
	)

	// Setup use cases
	userUseCase := usecase.NewUserUseCase(
		config.DB,
		config.Log,
		config.Validate,
		userRepository,
		customerRepository,
		wasteBankRepository,
		wasteCollectorRepository,
		industryRepository,
		collectorManagementRepository,
		jwtHelper,
		emailHelper,
		config.Config.GetString("app.base_url"), // Base URL for email links
	)
	customerUseCase := usecase.NewCustomerUseCase(config.DB, config.Log, config.Validate, customerRepository)
	wasteBankUseCase := usecase.NewWasteBankUseCase(config.DB, config.Log, config.Validate, wasteBankRepository)
	wasteCollectorUseCase := usecase.NewWasteCollectorUseCase(config.DB, config.Log, config.Validate, wasteCollectorRepository)
	industryUseCase := usecase.NewIndustryUseCase(config.DB, config.Log, config.Validate, industryRepository)
	wasteCategoryUseCase := usecase.NewWasteCategoryUsecase(config.DB, config.Log, config.Validate, wasteCategoryRepository)
	wasteTypeUseCase := usecase.NewWasteTypeUsecase(config.DB, config.Log, config.Validate, wasteCategoryRepository, wasteTypeRepository)
	wasteBankPricedTypeUseCase := usecase.NewWasteBankPricedTypeUsecase(config.DB, config.Log, config.Validate, wasteBankPricedTypeRepository, wasteTypeRepository)
	wasteDropRequestUseCase := usecase.NewWasteDropRequestUsecase(config.DB, config.Log, config.Validate, wasteDropRequestRepository, userRepository, wasteTypeRepository, wasteDropRequesItemRepository, wasteBankPricedTypeRepository, customerRepository, wasteBankRepository, wasteCollectorRepository)
	wasteDropRequestItemUseCase := usecase.NewWasteDropRequestItemUsecase(config.DB, config.Log, config.Validate, wasteDropRequesItemRepository, wasteDropRequestRepository, wasteTypeRepository)

	// Setup controllers
	userController := http.NewUserController(
		userUseCase,
		config.Log,
	)
	customerController := http.NewCustomerController(customerUseCase, config.Log)
	wasteBankController := http.NewWasteBankController(wasteBankUseCase, config.Log)
	wasteCollectorController := http.NewWasteCollectorController(wasteCollectorUseCase, config.Log)
	industryController := http.NewIndustryController(industryUseCase, config.Log)
	wasteCategoryController := http.NewWasteCategoryController(wasteCategoryUseCase, config.Log)
	wasteTypeController := http.NewWasteTypeController(wasteTypeUseCase, config.Log)
	wasteBankPricedTypeController := http.NewWasteBankPricedTypeController(wasteBankPricedTypeUseCase, config.Log)
	wasteDropRequestController := http.NewWasteDropRequestController(wasteDropRequestUseCase, config.Log)
	wasteDropRequestItemController := http.NewWasteDropRequestItemController(wasteDropRequestItemUseCase, config.Log)

	// Setup middlewares
	authMiddleware := middleware.NewJWTAuth(
		jwtHelper,
	)

	routeConfig := route.RouteConfig{
		App:                            config.App,
		UserController:                 userController,
		CustomerController:             customerController,
		WasteBankController:            wasteBankController,
		WasteCollectorController:       wasteCollectorController,
		IndustryController:             industryController,
		WasteCategoryController:        wasteCategoryController,
		WasteTypeController:            wasteTypeController,
		WasteBankPricedTypeController:  wasteBankPricedTypeController,
		WasteDropRequestController:     wasteDropRequestController,
		WasteDropRequestItemController: wasteDropRequestItemController,
		AuthMiddleware:                 authMiddleware,
	}

	routeConfig.Setup()
}
