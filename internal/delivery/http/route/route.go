package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http/middleware"
)

type RouteConfig struct {
	App                           *fiber.App
	UserController                *http.UserController
	CustomerController            *http.CustomerController
	WasteBankController           *http.WasteBankController
	WasteCollectorController      *http.WasteCollectorController
	IndustryController            *http.IndustryController
	WasteCategoryController       *http.WasteCategoryController
	WasteSubCategoryController    *http.WasteSubCategoryController
	WasteTypeController           *http.WasteTypeController
	WasteBankPricedTypeController *http.WasteBankPricedTypeController
	AuthMiddleware                fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
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
	auth.Get("/users/current", c.UserController.Current)
	auth.Post("/auth/logout", c.UserController.Logout)
	auth.Post("/auth/logout-all-devices", c.UserController.LogoutAllDevices)

	// Customer endpoints
	customerOnly := c.App.Group("/api/customer", c.AuthMiddleware, middleware.RequireRoles("admin", "customer"))
	// Profiles
	customerOnly.Get("/profiles/:user_id", c.CustomerController.Get)
	customerOnly.Put("/profiles/:id", c.CustomerController.Update)

	// WasteBank endpoints
	wasteBankOnly := c.App.Group("/api/waste-bank", c.AuthMiddleware, middleware.RequireRoles("admin", "waste_bank_unit", "waste_bank_central"))
	// Profiles
	wasteBankOnly.Get("/profiles/:user_id", c.WasteBankController.Get)
	wasteBankOnly.Put("/profiles/:id", c.WasteBankController.Update)
	// Waste Type Prices
	wasteBankOnly.Get("/waste-type-prices", c.WasteBankPricedTypeController.List)
	wasteBankOnly.Get("/waste-type-prices/:id", c.WasteBankPricedTypeController.Get)
	wasteBankOnly.Post("/batch-waste-type-prices", c.WasteBankPricedTypeController.CreateBatch)
	wasteBankOnly.Post("/waste-type-prices", c.WasteBankPricedTypeController.Create)
	wasteBankOnly.Put("/waste-type-prices/:id", c.WasteBankPricedTypeController.Update)
	wasteBankOnly.Delete("/waste-type-prices/:id", c.WasteBankPricedTypeController.Delete)

	// WasteCollector endpoints
	wasteCollectorOnly := c.App.Group("/api/waste-collector", c.AuthMiddleware, middleware.RequireRoles("admin", "waste_collector_unit", "waste_collector_central"))
	// Profiles
	wasteCollectorOnly.Get("/profiles/:user_id", c.WasteCollectorController.Get)
	wasteCollectorOnly.Put("/profiles/:id", c.WasteCollectorController.Update)

	// Industry endpoints
	industryOnly := c.App.Group("/api/industry", c.AuthMiddleware, middleware.RequireRoles("admin", "industry"))
	// Profiles
	industryOnly.Get("/profiles/:user_id", c.IndustryController.Get)
	industryOnly.Put("/profiles/:id", c.IndustryController.Update)

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
	adminOnly.Get("/waste-categories", c.WasteCategoryController.List)
	adminOnly.Get("/waste-categories/:id", c.WasteCategoryController.Get)
	adminOnly.Post("/waste-categories", c.WasteCategoryController.Create)
	adminOnly.Put("/waste-categories/:id", c.WasteCategoryController.Update)
	adminOnly.Delete("/waste-categories/:id", c.WasteCategoryController.Delete)

	// Waste SubCategories
	adminOnly.Get("/waste-subcategories", c.WasteSubCategoryController.List)
	adminOnly.Get("/waste-subcategories/:id", c.WasteSubCategoryController.Get)
	adminOnly.Post("/waste-subcategories", c.WasteSubCategoryController.Create)
	adminOnly.Put("/waste-subcategories/:id", c.WasteSubCategoryController.Update)
	adminOnly.Delete("/waste-subcategories/:id", c.WasteSubCategoryController.Delete)

	// Waste Types
	adminOnly.Get("/waste-types", c.WasteTypeController.List)
	adminOnly.Get("/waste-types/:id", c.WasteTypeController.Get)
	adminOnly.Post("/waste-types", c.WasteTypeController.Create)
	adminOnly.Put("/waste-types/:id", c.WasteTypeController.Update)
	adminOnly.Delete("/waste-types/:id", c.WasteTypeController.Delete)

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
