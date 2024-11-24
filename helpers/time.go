package helpers

import (
	"database/sql"
	"fmt"
	"time"
)

func ParseDate(dateStr, format string) (time.Time, error) {
	parsedDate, err := time.Parse(format, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %w", err)
	}
	return parsedDate, nil
}

func NullTimeScan(t time.Time) sql.NullTime {
	// If the time is zero value, return an invalid sql.NullTime
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	// Otherwise, return a valid sql.NullTime
	return sql.NullTime{Time: t, Valid: true}
}
