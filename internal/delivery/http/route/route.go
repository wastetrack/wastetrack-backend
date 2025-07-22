package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http/middleware"
)

type RouteConfig struct {
	App                                 *fiber.App
	UserController                      *http.UserController
	CustomerController                  *http.CustomerController
	WasteBankController                 *http.WasteBankController
	WasteCollectorController            *http.WasteCollectorController
	IndustryController                  *http.IndustryController
	WasteCategoryController             *http.WasteCategoryController
	WasteTypeController                 *http.WasteTypeController
	WasteBankPricedTypeController       *http.WasteBankPricedTypeController
	WasteDropRequestController          *http.WasteDropRequestController
	WasteDropRequestItemController      *http.WasteDropRequestItemController
	WasteTransferController             *http.WasteTransferRequestController
	WasteTransferItemOfferingController *http.WasteTransferItemOfferingController
	CollectorManagementController       *http.CollectorManagementController
	SalaryTransactionController         *http.SalaryTransactionController
	StorageController                   *http.StorageController
	StorageItemController               *http.StorageItemController
	AuthMiddleware                      fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	// Users
	c.App.Get("/api/users", c.UserController.List)
	// Authentication routes (no auth required)
	c.App.Post("/api/auth/register", c.UserController.Register)
	c.App.Post("/api/auth/login", c.UserController.Login)
	c.App.Post("/api/auth/verify-email", c.UserController.VerifyEmail)
	c.App.Post("/api/auth/resend-verification", c.UserController.ResendVerification)
	c.App.Post("/api/auth/forgot-password", c.UserController.ForgotPassword)
	c.App.Post("/api/auth/reset-password", c.UserController.ResetPassword)
	c.App.Post("/api/auth/refresh-token", c.UserController.RefreshToken)
}

func (c *RouteConfig) SetupAuthRoute() {
	// Apply JWT authentication middleware to protected routes
	auth := c.App.Group("/api", c.AuthMiddleware)

	// Authenticated user endpoints
	// Auth
	auth.Get("/users/current", c.UserController.Current)
	auth.Post("/auth/logout", c.UserController.Logout)
	auth.Post("/auth/logout-all-devices", c.UserController.LogoutAllDevices)
	// Profiles
	// Waste Bank
	auth.Get("/waste-bank/profiles/:user_id", c.WasteBankController.Get)
	// Waste Categories
	auth.Get("/waste-categories", c.WasteCategoryController.List)
	auth.Get("/waste-categories/:id", c.WasteCategoryController.Get)
	// Waste Type
	auth.Get("/waste-types", c.WasteTypeController.List)
	auth.Get("/waste-types/:id", c.WasteTypeController.Get)
	// Waste Type Prices
	auth.Get("/waste-type-prices", c.WasteBankPricedTypeController.List)
	auth.Get("/waste-type-prices/:id", c.WasteBankPricedTypeController.Get)
	// Waste Drop Requests
	auth.Get("/waste-drop-requests", c.WasteDropRequestController.List)
	auth.Get("/waste-drop-requests/:id", c.WasteDropRequestController.Get)
	// Waste Drop Request Items
	auth.Get("/waste-drop-request-items", c.WasteDropRequestItemController.List)
	auth.Get("/waste-drop-request-items/:id", c.WasteDropRequestItemController.Get)
	// Waste Transfer Requests
	auth.Get("/waste-transfer-requests", c.WasteTransferController.List)
	auth.Get("/waste-transfer-requests/:id", c.WasteTransferController.Get)
	// Waste Transfer Item Offerings
	auth.Get("/waste-transfer-item-offerings", c.WasteTransferItemOfferingController.List)
	auth.Get("/waste-transfer-item-offerings/:id", c.WasteTransferItemOfferingController.Get)
	// Salary Transactions
	auth.Get("/salary-transactions", c.SalaryTransactionController.List)
	auth.Get("/salary-transactions/:id", c.SalaryTransactionController.Get)
	// Storage
	auth.Get("/storages", c.StorageController.List)
	auth.Get("/storages/:id", c.StorageController.Get)
	// Storage Items
	auth.Get("/storage-items", c.StorageItemController.List)
	auth.Get("/storage-items/:id", c.StorageItemController.Get)

	// Customer endpoints
	customerOnly := c.App.Group("/api/customer", c.AuthMiddleware, middleware.RequireRoles("admin", "customer"))
	// Profiles
	customerOnly.Get("/profiles/:user_id", c.CustomerController.Get)
	customerOnly.Put("/profiles/:id", c.CustomerController.Update)
	// Waste Drop Requests
	customerOnly.Post("/waste-drop-requests", c.WasteDropRequestController.Create)

	// WasteBank endpoints
	wasteBankOnly := c.App.Group("/api/waste-bank", c.AuthMiddleware, middleware.RequireRoles("admin", "waste_bank_unit", "waste_bank_central"))
	// Profiles
	wasteBankOnly.Put("/profiles/:id", c.WasteBankController.Update)
	// Waste Type Prices
	wasteBankOnly.Post("/batch-waste-type-prices", c.WasteBankPricedTypeController.CreateBatch)
	wasteBankOnly.Post("/waste-type-prices", c.WasteBankPricedTypeController.Create)
	wasteBankOnly.Put("/waste-type-prices/:id", c.WasteBankPricedTypeController.Update)
	wasteBankOnly.Delete("/waste-type-prices/:id", c.WasteBankPricedTypeController.Delete)
	// Waste Drop Requests
	wasteBankOnly.Put("/waste-drop-requests/:id", c.WasteDropRequestController.UpdateStatus)
	wasteBankOnly.Put("/waste-drop-requests/:id/assign-collector", c.WasteDropRequestController.AssignCollector)
	// Waste Transfer
	wasteBankOnly.Post("/waste-transfer-requests", c.WasteTransferController.Create)
	wasteBankOnly.Put("/waste-transfer-requests/:id/assign-collector", c.WasteTransferController.AssignCollectorByWasteType)
	wasteBankOnly.Put("/waste-transfer-requests/:id", c.WasteTransferController.UpdateStatus)
	wasteBankOnly.Put("/waste-transfer-requests/:id/complete", c.WasteTransferController.CompleteRequest)
	// Collector Management
	wasteBankOnly.Get("/collector-management", c.CollectorManagementController.List)
	wasteBankOnly.Get("/collector-management/:id", c.CollectorManagementController.Get)
	wasteBankOnly.Post("/collector-management", c.CollectorManagementController.Create)
	wasteBankOnly.Put("/collector-management/:id", c.CollectorManagementController.Update)
	wasteBankOnly.Delete("/collector-management/:id", c.CollectorManagementController.Delete)
	// Salary Transactions
	wasteBankOnly.Post("/salary-transactions", c.SalaryTransactionController.Create)
	wasteBankOnly.Put("/salary-transactions/:id", c.SalaryTransactionController.Update)
	// Storage
	wasteBankOnly.Put("/storages/:id", c.StorageController.Update)
	// Storage Items
	wasteBankOnly.Post("/storage-items", c.StorageItemController.Create)
	wasteBankOnly.Put("/storage-items/:id", c.StorageItemController.Update)
	wasteBankOnly.Put("/storage-items/:id/deduct-weight", c.StorageItemController.DeductStorageItem)
	wasteBankOnly.Delete("/storage-items/:id", c.StorageItemController.Delete)

	// WasteCollector endpoints
	wasteCollectorOnly := c.App.Group("/api/waste-collector", c.AuthMiddleware, middleware.RequireRoles("admin", "waste_collector_unit", "waste_collector_central", "waste_bank_unit", "waste_bank_central"))
	// Profiles
	wasteCollectorOnly.Get("/profiles/:user_id", c.WasteCollectorController.Get)
	wasteCollectorOnly.Put("/profiles/:id", c.WasteCollectorController.Update)
	// Waste Drop Requests
	wasteCollectorOnly.Put("/waste-drop-requests/:id/complete", c.WasteDropRequestController.Complete)

	// Industry endpoints
	industryOnly := c.App.Group("/api/industry", c.AuthMiddleware, middleware.RequireRoles("admin", "industry"))
	// Profiles
	industryOnly.Get("/profiles/:user_id", c.IndustryController.Get)
	industryOnly.Put("/profiles/:id", c.IndustryController.Update)
	// Recycle Waste Transfer
	industryOnly.Put("/waste-transfer-requests/:id", c.WasteTransferController.UpdateStatus)
	industryOnly.Put("/waste-transfer-requests/:id/complete", c.WasteTransferController.CompleteRequest)
	industryOnly.Put("waste-transfer-requests/:id/assign-collector", c.WasteTransferController.AssignCollectorByWasteType)
	// Storage
	industryOnly.Put("/storages/:id", c.StorageController.Update)
	// Storage Items
	industryOnly.Post("/storage-items", c.StorageItemController.Create)
	industryOnly.Put("/storage-items/:id", c.StorageItemController.Update)
	industryOnly.Put("/storage-items/:id/deduct-weight", c.StorageItemController.DeductStorageItem)
	industryOnly.Delete("/storage-items/:id", c.StorageItemController.Delete)
	// Admin endpoints
	adminOnly := c.App.Group("/api/admin", c.AuthMiddleware, middleware.RequireRoles("admin"))
	// Customer profiles
	adminOnly.Delete("/customer/profiles/:id", c.CustomerController.Delete)
	// Wastebank profiles
	adminOnly.Delete("/waste-bank/profiles/:id", c.WasteBankController.Delete)
	// Waste collector profiles
	adminOnly.Delete("/waste-collector/profiles/:id", c.WasteCollectorController.Delete)
	// Industry profiles
	adminOnly.Delete("/industry/profiles/:id", c.IndustryController.Delete)
	// Waste Categories
	adminOnly.Post("/waste-categories", c.WasteCategoryController.Create)
	adminOnly.Put("/waste-categories/:id", c.WasteCategoryController.Update)
	adminOnly.Delete("/waste-categories/:id", c.WasteCategoryController.Delete)
	// Waste Types
	adminOnly.Post("/waste-types", c.WasteTypeController.Create)
	adminOnly.Put("/waste-types/:id", c.WasteTypeController.Update)
	adminOnly.Delete("/waste-types/:id", c.WasteTypeController.Delete)
	// Waste Drop Requests
	adminOnly.Put("/waste-drop-requests/:id", c.WasteDropRequestController.Update)
	adminOnly.Delete("/waste-drop-requests/:id", c.WasteDropRequestController.Delete)
	// Waste Transfer
	adminOnly.Put("/waste-transfer-requests/:id", c.WasteTransferController.Update)
	adminOnly.Delete("/waste-transfer-requests/:id", c.WasteTransferController.Delete)
	// Salary Transactions
	adminOnly.Delete("/salary-transactions/:id", c.SalaryTransactionController.Delete)
	// Storage
	adminOnly.Delete("/storages/:id", c.StorageController.Delete)

}

// auth.Put("/users", c.UserController.Update)
// auth.Delete("/users", c.UserController.Delete)

//Email verification required
// emailVerified := c.App.Group("", middleware.RequireEmailVerification())
// Admin-only routes (requires admin role + email verification)
// adminOnly := c.App.Group("/api/admin",
// 	middleware.RequireEmailVerification(),
// 	middleware.RequireRole("admin"),
// )
