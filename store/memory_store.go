package store

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/tjeastmond/future-take-home/models"
	"github.com/tjeastmond/future-take-home/utils"
)

type MemoryStore struct {
	appointments map[int]models.Appointment
	filePath     string
	mutex        sync.RWMutex
	nextID       int
}

func NewMemoryStore(filePath string) *MemoryStore {
	return &MemoryStore{
		appointments: make(map[int]models.Appointment),
		filePath:     filePath,
		mutex:        sync.RWMutex{},
	}
}

func (ms *MemoryStore) LoadData() error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	file, err := os.Open(ms.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var loadedAppointments []models.Appointment
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&loadedAppointments); err != nil {
		return err
	}

	for _, appointment := range loadedAppointments {
		ms.appointments[appointment.ID] = appointment
		if appointment.ID >= ms.nextID {
			ms.nextID = appointment.ID + 1
		}
	}

	return nil
}

func (ms *MemoryStore) SortAppointmentsByID(appointments []models.Appointment) {
	sort.SliceStable(appointments, func(i, j int) bool {
		return appointments[i].ID < appointments[j].ID
	})
}

func (ms *MemoryStore) SortedAppointments() []models.Appointment {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	appointments := make([]models.Appointment, 0, len(ms.appointments))
	for _, appointment := range ms.appointments {
		appointments = append(appointments, appointment)
	}

	ms.SortAppointmentsByID(appointments)

	return appointments
}

func (ms *MemoryStore) AddAppointment(a models.Appointment) (int, error) {
	if err := a.Validate(); err != nil {
		return 0, err
	}

	ms.mutex.Lock()
	a.ID = ms.nextID
	ms.appointments[a.ID] = a
	ms.nextID++
	ms.mutex.Unlock()

	return a.ID, nil
}

func (ms *MemoryStore) GetAllAppointments() []models.Appointment {
	return ms.SortedAppointments()
}

func (ms *MemoryStore) GetAppointmentByID(id int) (models.Appointment, bool) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()
	appointment, exists := ms.appointments[id]
	return appointment, exists
}

func (ms *MemoryStore) GetAllAppointmentsForTrainer(trainerID int) []models.Appointment {
	trainerAppointments := []models.Appointment{}
	for _, app := range ms.appointments {
		if app.TrainerID == trainerID {
			trainerAppointments = append(trainerAppointments, app)
		}
	}

	ms.SortAppointmentsByID(trainerAppointments)

	return trainerAppointments
}

func (ms *MemoryStore) GetAppointmentsForTrainer(trainerID int, startTime, endTime *time.Time) []models.Appointment {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	matchedAppointments := []models.Appointment{}
	for _, app := range ms.appointments {
		if app.TrainerID == trainerID {
			if startTime != nil && endTime != nil {
				if app.StartsAt.Before(*endTime) && app.EndsAt.After(*startTime) {
					matchedAppointments = append(matchedAppointments, app)
				}
			} else {
				matchedAppointments = append(matchedAppointments, app)
			}
		}
	}

	ms.SortAppointmentsByID(matchedAppointments)

	return matchedAppointments
}

func (ms *MemoryStore) GetTrainerAvailability(trainerID int, startTime, endTime *time.Time) []string {
	appointments := ms.GetAppointmentsForTrainer(trainerID, startTime, endTime)
	availableSlots := []string{}
	for t := *startTime; t.Before(*endTime); t = t.Add(30 * time.Minute) {
		if !utils.IsValidTime(t) {
			continue
		}

		if t.Minute() != 0 && t.Minute() != 30 {
			continue
		}

		slotTaken := false
		for _, app := range appointments {
			if t.Before(app.EndsAt) && t.Add(30*time.Minute).After(app.StartsAt) {
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

func (ms *MemoryStore) IsAvailable(appointment models.Appointment) bool {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	for _, app := range ms.appointments {
		if app.TrainerID == appointment.TrainerID {
			// check if there is an exact match for start and end times
			if appointment.StartsAt.Equal(app.StartsAt) && appointment.EndsAt.Equal(app.EndsAt) {
				return false
			}

			// check for overlapping times
			if appointment.StartsAt.Before(app.EndsAt) && appointment.EndsAt.After(app.StartsAt) {
				return false
			}
		}
	}

	return true
}
