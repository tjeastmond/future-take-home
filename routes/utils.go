package routes

import (
	"time"

	"github.com/tjeastmond/future-take-home/models"
)
func isValidTime(t time.Time) bool {
	_, offset := t.Zone()
	if offset != -8*3600 {
		return false
	}
	minutes := t.Minute()
	return minutes == 0 || minutes == 30
}

func getAvailableSlots(startTime, endTime time.Time, appointments []models.Appointment) []string {
	availableSlots := []string{}
	for t := startTime; t.Before(endTime); t = t.Add(30 * time.Minute) {
		// check if the time is either at :00 or :30
		if t.Minute() != 0 && t.Minute() != 30 {
			continue
		}

		// is the time slot is already taken
		slotTaken := false
		for _, app := range appointments {
			if t.Before(app.EndedAt) && t.Add(30*time.Minute).After(app.StartedAt) {
				slotTaken = true
				break
			}
		}

		if !slotTaken {
			availableSlots = append(availableSlots, t.Format(time.RFC3339))
		}
	}
	return availableSlots
}
