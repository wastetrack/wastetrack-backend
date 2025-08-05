package timezone

import (
	"time"
)

var (
	WIB  *time.Location // Western Indonesia Time (UTC+7) - Jakarta, Sumatra, Java
	WITA *time.Location // Central Indonesia Time (UTC+8) - Bali, Kalimantan, Sulawesi
	WIT  *time.Location // Eastern Indonesia Time (UTC+9) - Papua, Maluku
)

func InitTimeLocation() {
	var err error

	// Load WIB (Western Indonesia Time - UTC+7)
	WIB, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic("failed to load Asia/Jakarta timezone: " + err.Error())
	}

	// Load WITA (Central Indonesia Time - UTC+8)
	WITA, err = time.LoadLocation("Asia/Makassar")
	if err != nil {
		panic("failed to load Asia/Makassar timezone: " + err.Error())
	}

	// Load WIT (Eastern Indonesia Time - UTC+9)
	WIT, err = time.LoadLocation("Asia/Jayapura")
	if err != nil {
		panic("failed to load Asia/Jayapura timezone: " + err.Error())
	}
}

// GetLocationByOffset returns the appropriate Indonesia timezone based on UTC offset
func GetLocationByOffset(offset int) *time.Location {
	switch offset {
	case 7:
		return WIB
	case 8:
		return WITA
	case 9:
		return WIT
	default:
		return WIB // Default to WIB
	}
}

// GetNowInTimezone returns current time in the specified Indonesia timezone
func GetNowInTimezone(timezone *time.Location) time.Time {
	return time.Now().In(timezone)
}

// IsDateTimeInPast checks if the given date and time is in the past
func IsDateTimeInPast(appointmentDate time.Time, appointmentTime time.Time, timezone *time.Location) bool {
	// Combine date and time
	combined := time.Date(
		appointmentDate.Year(),
		appointmentDate.Month(),
		appointmentDate.Day(),
		appointmentTime.Hour(),
		appointmentTime.Minute(),
		appointmentTime.Second(),
		0,
		timezone,
	)

	now := GetNowInTimezone(timezone)
	return combined.Before(now)
}

// IsDateTimeInPastFromParsed checks if the given date and parsed time string is in the past
func IsDateTimeInPastFromParsed(appointmentDate time.Time, timeStr string) (bool, error) {
	parsedTime, location, err := ParseTimeWithTimezone(timeStr)
	if err != nil {
		return false, err
	}

	return IsDateTimeInPast(appointmentDate, parsedTime, location), nil
}

// ParseTimeWithTimezone parses time string with timezone offset and returns the appropriate timezone
func ParseTimeWithTimezone(timeStr string) (time.Time, *time.Location, error) {
	parsedTime, err := time.Parse("15:04:05Z07:00", timeStr)
	if err != nil {
		return time.Time{}, nil, err
	}

	// Get the offset from the parsed time
	_, offset := parsedTime.Zone()
	offsetHours := offset / 3600

	// Get the appropriate Indonesia timezone
	location := GetLocationByOffset(offsetHours)

	return parsedTime, location, nil
}
