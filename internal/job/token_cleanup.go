package job

import (
	"fmt"
	"time"

	"github.com/wastetrack/wastetrack-backend/internal/helper"
	"gorm.io/gorm"
)

func StartTokenCleanupJob(db *gorm.DB, jwtHelper *helper.JWTHelper) {
	ticker := time.NewTicker(24 * time.Hour) // Run daily
	go func() {
		for range ticker.C {
			if err := jwtHelper.CleanupExpiredTokens(db); err != nil {
				fmt.Println("Error cleaning up expired tokens:", err)
			}
		}
	}()
}
