package store

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/tjeastmond/future-take-home/models"
)

var (
	appointments = make(map[int]models.Appointment)
	mutex        = &sync.Mutex{}
	filePath     = "docs/appointments.json"
)

func LoadData() error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.Open(filePath)
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
		appointments[appointment.ID] = appointment
	}

	return nil
}

func SaveData() error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var appointmentsList []models.Appointment
	for _, app := range appointments {
		appointmentsList = append(appointmentsList, app)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(appointmentsList); err != nil {
		return err
	}

	return nil
}

func persistData() {
	if err := SaveData(); err != nil {
		log.Printf("Error saving data: %v", err)
	}
}

func AddAppointment(a models.Appointment) int {
	mutex.Lock()
	defer mutex.Unlock()

	// this is fine for a small in-memory store
	appointments[a.ID] = a
	persistData()

	return a.ID
}

// GetAllAppointments returns all appointments
func GetAllAppointments() []models.Appointment {
	mutex.Lock()
	defer mutex.Unlock()

	var allAppointments []models.Appointment
	for _, app := range appointments {
		allAppointments = append(allAppointments, app)
	}
	return allAppointments
}

func GetAppointmentByID(id int) (models.Appointment, bool) {
	mutex.Lock()
	defer mutex.Unlock()
	appointment, exists := appointments[id]
	return appointment, exists
}

func UpdateAppointment(id int, a models.Appointment) bool {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := appointments[id]; exists {
		appointments[id] = a
		persistData()
		return true
	}

	return false
}

// DeleteAppointment deletes an appointment by ID
func DeleteAppointment(id int) bool {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := appointments[id]; exists {
		delete(appointments, id)
		persistData()
		return true
	}

	return false
}

func GetAppointmentsForTrainer(trainerID string, startTime, endTime time.Time) []models.Appointment {
	matchedAppointments := []models.Appointment{}
	trainerIDInt, err := strconv.Atoi(trainerID)
	if err != nil {
		return nil
	}

	mutex.Lock()
	defer mutex.Unlock()

	for _, app := range appointments {
		if app.TrainerID == trainerIDInt &&
			app.StartedAt.Before(endTime) &&
			app.EndedAt.After(startTime) {
			matchedAppointments = append(matchedAppointments, app)
		}
	}

	return matchedAppointments
}
