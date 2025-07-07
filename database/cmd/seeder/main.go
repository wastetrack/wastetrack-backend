package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/wastetrack/wastetrack-backend/database/seeder"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	var (
		dbURL     = flag.String("db", "", "Database URL (required)")
		clear     = flag.Bool("clear", false, "Clear all data before seeding")
		onlyClear = flag.Bool("only-clear", false, "Only clear data and exit")
		onlySeed  = flag.Bool("only-seed", false, "Only run seeders (skip clearing)")
		help      = flag.Bool("help", false, "Show help message")
		verbose   = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *dbURL == "" {
		// Try to get from environment variable
		*dbURL = os.Getenv("DATABASE_URL")
		if *dbURL == "" {
			log.Fatal("Database URL is required. Use -db flag or set DATABASE_URL environment variable")
		}
	}

	// Configure GORM logger
	logLevel := logger.Silent
	if *verbose {
		logLevel = logger.Info
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(*dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate all entities (optional, useful for development)
	if err := autoMigrate(db); err != nil {
		log.Printf("Warning: Auto-migration failed: %v", err)
	}

	// -- Handle only-clear flag
	if *onlyClear {
		log.Println("Clearing existing data...")
		if err := seeder.ClearAllData(db); err != nil {
			log.Fatalf("Failed to clear data: %v", err)
		}
		log.Println("Data cleared successfully!")
		return
	}

	// -- Handle clear flag unless -only-seed is present
	if *clear && !*onlySeed {
		log.Println("Clearing existing data...")
		if err := seeder.ClearAllData(db); err != nil {
			log.Fatalf("Failed to clear data: %v", err)
		}
		log.Println("Data cleared successfully!")
	}

	// -- Run seeders (default or when -only-seed is passed)
	if *onlySeed || !*onlyClear {
		if err := seeder.RunAllSeeders(db); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
		log.Println("All seeders completed successfully!")
	}
}

func showHelp() {
	fmt.Println("WasteTrack Database Seeder")
	fmt.Println("Usage: go run cmd/seeder/main.go [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -db string        Database URL (can also use DATABASE_URL env var)")
	fmt.Println("  -clear            Clear all existing data before seeding")
	fmt.Println("  -only-clear       Only clear data and exit (skip seeding)")
	fmt.Println("  -only-seed        Only run seeders (skip clearing)")
	fmt.Println("  -verbose          Enable verbose database logging")
	fmt.Println("  -help             Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/seeder/main.go -db 'postgres://user:pass@localhost/wastetrack'")
	fmt.Println("  go run cmd/seeder/main.go -clear")
	fmt.Println("  go run cmd/seeder/main.go -only-clear")
	fmt.Println("  go run cmd/seeder/main.go -only-seed")
	fmt.Println("  DATABASE_URL='postgres://user:pass@localhost/wastetrack' go run cmd/seeder/main.go -clear")
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.WasteCategory{},
		&entity.WasteType{},
		&entity.User{},
		&entity.CustomerProfile{},
		&entity.GovernmentProfile{},
		&entity.IndustryProfile{},
		&entity.WasteBankProfile{},
		&entity.WasteCollectorProfile{},
		&entity.CollectorManagement{},
		&entity.RefreshToken{},
		&entity.Storage{},
		&entity.StorageItem{},
		&entity.WasteBankPricedType{},
		&entity.WasteDropRequest{},
		&entity.WasteDropRequestItem{},
		&entity.WasteTransferRequest{},
		&entity.WasteTransferItemOffering{},
		&entity.SalaryTransaction{},
	)
}
