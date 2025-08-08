package helper

import (
	"fmt"
	"time"
)

func PrepareDateParams(monthStr string) (string, string) {
	if monthStr == "" {
		return "", ""
	}

	if parsed, err := time.Parse("2006-01", monthStr); err == nil {
		return monthStr, parsed.Format("2006-01-02 15:04:05")
	}

	return "", ""
}

// PrepareEndDateParams converts end month string to month string and end-of-month timestamp
// Used for end date filters to include the entire month
func PrepareEndDateParams(monthStr string) (string, string) {
	if monthStr == "" {
		return "", ""
	}

	if parsed, err := time.Parse("2006-01", monthStr); err == nil {
		// Set to last moment of the month: last day + 23:59:59
		endOfMonth := parsed.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		return monthStr, endOfMonth.Format("2006-01-02 15:04:05")
	}

	return "", ""
}

// ValidateMonthFormat validates that a string is in YYYY-MM format
func ValidateMonthFormat(monthStr string) error {
	if monthStr == "" {
		return nil // Empty is valid (optional parameter)
	}

	if _, err := time.Parse("2006-01", monthStr); err != nil {
		return fmt.Errorf("invalid month format '%s', expected YYYY-MM format", monthStr)
	}

	return nil
}

// ValidateDateRange validates that endMonth is after or equal to startMonth
func ValidateDateRange(startMonth, endMonth string) error {
	if startMonth == "" || endMonth == "" {
		return nil // Skip validation if either date is empty
	}

	startDate, err := time.Parse("2006-01", startMonth)
	if err != nil {
		return fmt.Errorf("invalid start_month format: %v", err)
	}

	endDate, err := time.Parse("2006-01", endMonth)
	if err != nil {
		return fmt.Errorf("invalid end_month format: %v", err)
	}

	if endDate.Before(startDate) {
		return fmt.Errorf("end_month (%s) must be after or equal to start_month (%s)", endMonth, startMonth)
	}

	return nil
}

// ParseMonthRange generates all months between start and end (inclusive)
// Returns slice of month strings in YYYY-MM format
func ParseMonthRange(startMonth, endMonth string) ([]string, error) {
	if startMonth == "" || endMonth == "" {
		return nil, fmt.Errorf("both start_month and end_month are required")
	}

	startDate, err := time.Parse("2006-01", startMonth)
	if err != nil {
		return nil, fmt.Errorf("invalid start_month format: %v", err)
	}

	endDate, err := time.Parse("2006-01", endMonth)
	if err != nil {
		return nil, fmt.Errorf("invalid end_month format: %v", err)
	}

	if endDate.Before(startDate) {
		return nil, fmt.Errorf("end_month must be after or equal to start_month")
	}

	var months []string
	currentMonth := startDate

	for currentMonth.Before(endDate) || currentMonth.Equal(endDate) {
		months = append(months, currentMonth.Format("2006-01"))
		currentMonth = currentMonth.AddDate(0, 1, 0)
	}

	return months, nil
}

// GetCurrentMonth returns current month in YYYY-MM format
func GetCurrentMonth() string {
	return time.Now().Format("2006-01")
}

// GetPreviousMonth returns previous month in YYYY-MM format
func GetPreviousMonth() string {
	return time.Now().AddDate(0, -1, 0).Format("2006-01")
}

// GetMonthsAgo returns month N months ago in YYYY-MM format
func GetMonthsAgo(months int) string {
	return time.Now().AddDate(0, -months, 0).Format("2006-01")
}
