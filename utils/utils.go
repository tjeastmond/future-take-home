package utils

import (
	"time"
)

func IsValidTime(t time.Time) bool {
	_, offset := t.Zone()
	if offset != -8*3600 {
		return false
	}

	hours := t.Hour()
	minutes := t.Minute()
	weekday := t.Weekday()

	return (minutes == 0 || minutes == 30) &&
		hours >= 8 && hours < 17 &&
		weekday != time.Saturday &&
		weekday != time.Sunday
}

func IsValidTimes(startTime, endTime *time.Time) bool {
	return IsValidTime(*startTime) && IsValidTime(*endTime)
}
