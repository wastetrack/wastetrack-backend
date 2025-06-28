package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http"
)

type RouteConfig struct {
	App            *fiber.App
	UserController *http.UserController
	AuthMiddleware fiber.Handler
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
	// auth.Put("/users", c.UserController.Update)
	// auth.Delete("/users", c.UserController.Delete)

	//Email verification required
	// emailVerified := c.App.Group("", middleware.RequireEmailVerification())
	// Admin-only routes (requires admin role + email verification)
	// adminOnly := c.App.Group("/api/admin",
	// 	middleware.RequireEmailVerification(),
	// 	middleware.RequireRole("admin"),
	// )
}
