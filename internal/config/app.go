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
		jwtHelper,
		emailHelper,
		config.Config.GetString("app.base_url"), // Base URL for email links
	)
	wasteBankUseCase := usecase.NewWasteBankUseCase(config.DB, config.Log, config.Validate, wasteBankRepository)
	// Setup controllers
	userController := http.NewUserController(
		userUseCase,
		config.Log,
	)
	wasteBankController := http.NewWasteBankController(wasteBankUseCase, config.Log)

	// Setup middlewares
	authMiddleware := middleware.NewJWTAuth(
		jwtHelper,
	)

	routeConfig := route.RouteConfig{
		App:                 config.App,
		UserController:      userController,
		WasteBankController: wasteBankController,
		AuthMiddleware:      authMiddleware,
	}

	routeConfig.Setup()
}
