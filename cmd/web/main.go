package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/wastetrack/wastetrack-backend/internal/config"
	"github.com/wastetrack/wastetrack-backend/pkg/timezone"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return origin == "http://localhost:3000" || origin == "https://wastetrack-staging.netlify.app"
		},
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,Access-Control-Request-Method,Access-Control-Request-Headers",
		AllowCredentials: true,
	}))

	timezone.InitTimeLocation()
	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})
	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
