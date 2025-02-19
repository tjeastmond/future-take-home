package utils

import (
	"fmt"
	"time"

	"github.com/tjeastmond/future-take-home/config"
)

func IsValidTime(t time.Time) bool {
	location, err := time.LoadLocation(config.LocationName)
	if err != nil {
		return false
	}

	if t.Location().String() != location.String() {
		fmt.Println(t.Location().String())
		fmt.Println(location.String())
		fmt.Println("Time is not in the correct location")
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

func IsValidTimes(start, end *time.Time) bool {
	return IsValidTime(*start) && IsValidTime(*end)
}

func IsValidSlot(start, end time.Time) bool {
	fmt.Println(start)
	fmt.Println(end)

	if IsValidTimes(&start, &end) {
		fmt.Println("Times are valid")
		return end.Sub(start) == 30*time.Minute
	}

	return false
}

func ParseTimeToPST(t time.Time) (time.Time, error) {
	location, err := time.LoadLocation(config.LocationName)
	if err != nil {
		return time.Time{}, err
	}

	return t.In(location), nil
}

func ValidateAndParse(startStr, endStr string) (time.Time, time.Time, bool) {
	start, _ := StrToTimestamp(startStr)
	end, _ := StrToTimestamp(endStr)

	if !IsValidTimes(&start, &end) {
		return time.Time{}, time.Time{}, false
	}

	startPST, err := ParseTimeToPST(start)
	if err != nil {
		return time.Time{}, time.Time{}, false
	}

	endPST, err := ParseTimeToPST(end)
	if err != nil {
		return time.Time{}, time.Time{}, false
	}

	return startPST, endPST, true
}

func StrToTimestamp(timeStr string) (time.Time, error) {
	location, err := time.LoadLocation(config.LocationName)
	if err != nil {
		return time.Time{}, err
	}

	t, err := time.ParseInLocation(config.TimeLayout, timeStr, location)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
