package store

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/tjeastmond/future-take-home/models"
)

type MemoryStore struct {
	appointments map[int]models.Appointment
	filePath     string
	mutex        *sync.Mutex
	nextID       int
}

func NewMemoryStore(filePath string) *MemoryStore {
	return &MemoryStore{
		appointments: make(map[int]models.Appointment),
		filePath:     filePath,
		mutex:        &sync.Mutex{},
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

func (ms *MemoryStore) SaveData() error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	file, err := os.Create(ms.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var appointmentsList []models.Appointment
	for _, app := range ms.appointments {
		appointmentsList = append(appointmentsList, app)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(appointmentsList)
}

func (ms *MemoryStore) SortedAppointments() []models.Appointment {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	keys := make([]int, 0, len(ms.appointments))
	for k := range ms.appointments {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return ms.appointments[keys[i]].StartedAt.Before(ms.appointments[keys[j]].StartedAt)
	})

	var sortedAppointments []models.Appointment
	for _, k := range keys {
		sortedAppointments = append(sortedAppointments, ms.appointments[k])
	}

	return sortedAppointments
}

func (ms *MemoryStore) PersistDataAsync() {
	if err := ms.SaveData(); err != nil {
		log.Printf("Error saving data: %v", err)
	}
}

func (ms *MemoryStore) AddAppointment(a models.Appointment) (int, error) {
	if err := a.Validate(); err != nil {
		return 0, err
	}

	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	a.ID = ms.nextID
	ms.appointments[a.ID] = a
	ms.nextID++

	go ms.PersistDataAsync()
	return a.ID, nil
}

func (ms *MemoryStore) GetAllAppointments() []models.Appointment {
	return ms.SortedAppointments()
}

func (ms *MemoryStore) GetAppointmentByID(id int) (models.Appointment, bool) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	appointment, exists := ms.appointments[id]
	return appointment, exists
}

func (ms *MemoryStore) DeleteAppointment(id int) bool {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	if _, exists := ms.appointments[id]; exists {
		delete(ms.appointments, id)
		ms.PersistDataAsync()
		return true
	}

	return false
}

func (ms *MemoryStore) GetAppointmentsForTrainer(trainerID int, startTime, endTime time.Time) []models.Appointment {
	matchedAppointments := []models.Appointment{}
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	for _, app := range ms.appointments {
		if app.TrainerID == trainerID &&
			app.StartedAt.Before(endTime) &&
			app.EndedAt.After(startTime) {
			matchedAppointments = append(matchedAppointments, app)
		}
	}

	return matchedAppointments
}

func (ms *MemoryStore) GetTrainerAvailability(trainerID int, startTime, endTime time.Time) []string {
	appointments := ms.GetAppointmentsForTrainer(trainerID, startTime, endTime)

	availableSlots := []string{}
	for t := startTime; t.Before(endTime); t = t.Add(30 * time.Minute) {
		if t.Minute() != 0 && t.Minute() != 30 {
			continue
		}

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

func (ms *MemoryStore) IsAvailable(appointment models.Appointment) bool {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	for _, app := range ms.appointments {
		if app.TrainerID == appointment.TrainerID && appointment.StartedAt == app.StartedAt {
			return false
		}
	}

	return true
}
