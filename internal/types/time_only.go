package types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type TimeOnly struct {
	sql.NullTime
}

// Scan implements the Scanner interface
func (t *TimeOnly) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		t.Valid = false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v
		t.Valid = true
		return nil
	case string:
		// Parse time string formats that PostgreSQL might return
		formats := []string{
			"15:04:05",
			"15:04:05.999999",
			"15:04:05-07",
			"15:04:05+07",
			"15:04:05-07:00",
			"15:04:05.999999-07:00",
			"15:04:05+07:00",
			"15:04:05.999999+07:00",
		}

		for _, format := range formats {
			if parsedTime, err := time.Parse(format, v); err == nil {
				t.Time = parsedTime
				t.Valid = true
				return nil
			}
		}
		return fmt.Errorf("cannot parse time string: %s", v)
	case []byte:
		return t.Scan(string(v))
	default:
		return fmt.Errorf("cannot scan %T into TimeOnly", value)
	}
}

// Value implements the driver Valuer interface
func (t TimeOnly) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time.Format("15:04:05"), nil
}

// IsZero returns true if the time is not valid or zero
func (t TimeOnly) IsZero() bool {
	return !t.Valid || t.Time.IsZero()
}

// Format formats the time if valid
func (t TimeOnly) Format(layout string) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(layout)
}

// NewTimeOnly creates a TimeOnly from time.Time
func NewTimeOnly(t time.Time) TimeOnly {
	if t.IsZero() {
		return TimeOnly{sql.NullTime{Valid: false}}
	}
	return TimeOnly{sql.NullTime{Time: t, Valid: true}}
}

// NewTimeOnlyFromString creates a TimeOnly from time string
func NewTimeOnlyFromString(timeStr string) (TimeOnly, error) {
	if timeStr == "" {
		return TimeOnly{sql.NullTime{Valid: false}}, nil
	}

	formats := []string{
		"15:04:05",
		"15:04",
		"3:04:05 PM",
		"3:04 PM",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return TimeOnly{sql.NullTime{Time: t, Valid: true}}, nil
		}
	}

	return TimeOnly{}, fmt.Errorf("cannot parse time string: %s", timeStr)
}
