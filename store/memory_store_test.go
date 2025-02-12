package store

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/tjeastmond/future-take-home/models"
)

func TestMemoryStore_AddAppointment(t *testing.T) {
	testFile := "appointments.json"
	defer os.Remove(testFile)

	store := NewMemoryStore(testFile)

	appointment := models.Appointment{
		TrainerID: 1,
		StartedAt: time.Now(),
		EndedAt:   time.Now().Add(1 * time.Hour),
	}

	id, err := store.AddAppointment(appointment)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if id == 0 {
		t.Fatal("expected non-zero ID")
	}

	// Check if appointment is stored correctly
	storedAppointment, exists := store.GetAppointmentByID(id)
	if !exists {
		t.Fatal("expected appointment to exist")
	}

	if !reflect.DeepEqual(appointment.TrainerID, storedAppointment.TrainerID) ||
		!reflect.DeepEqual(appointment.StartedAt, storedAppointment.StartedAt) ||
		!reflect.DeepEqual(appointment.EndedAt, storedAppointment.EndedAt) {
		t.Errorf("stored appointment does not match original, got %+v", storedAppointment)
	}
}

func TestMemoryStore_UpdateAppointment(t *testing.T) {
	testFile := "test_data.json"
	defer os.Remove(testFile)

	store := NewMemoryStore(testFile)

	appointment := models.Appointment{
		TrainerID: 1,
		StartedAt: time.Now(),
		EndedAt:   time.Now().Add(1 * time.Hour),
	}

	id, _ := store.AddAppointment(appointment)

	newAppointment := models.Appointment{
		TrainerID: 2,
		StartedAt: time.Now(),
		EndedAt:   time.Now().Add(2 * time.Hour),
	}

	updated := store.UpdateAppointment(id, newAppointment)

	if !updated {
		t.Fatal("expected update to succeed")
	}

	storedAppointment, _ := store.GetAppointmentByID(id)
	if !reflect.DeepEqual(newAppointment, storedAppointment) {
		t.Error("expected updated appointment to match new appointment")
	}
}

func TestMemoryStore_DeleteAppointment(t *testing.T) {
	testFile := "test_data.json"
	defer os.Remove(testFile)

	store := NewMemoryStore(testFile)

	appointment := models.Appointment{
		TrainerID: 1,
		StartedAt: time.Now(),
		EndedAt:   time.Now().Add(1 * time.Hour),
	}

	id, _ := store.AddAppointment(appointment)

	// Ensure it exists before deletion
	if _, exists := store.GetAppointmentByID(id); !exists {
		t.Fatal("expected appointment to exist before deletion")
	}

	deleted := store.DeleteAppointment(id)
	if !deleted {
		t.Fatal("expected deletion to succeed")
	}

	// Check if appointment is deleted
	if _, exists := store.GetAppointmentByID(id); exists {
		t.Fatal("expected appointment to not exist after deletion")
	}
}

func TestMemoryStore_GetAvailability(t *testing.T) {
	testFile := "test_data.json"
	defer os.Remove(testFile)

	store := NewMemoryStore(testFile)

	now := time.Now()
	appointment1 := models.Appointment{
		TrainerID: 1,
		StartedAt: now,
		EndedAt:   now.Add(30 * time.Minute),
	}

	appointment2 := models.Appointment{
		TrainerID: 1,
		StartedAt: now.Add(30 * time.Minute),
		EndedAt:   now.Add(60 * time.Minute),
	}

	store.AddAppointment(appointment1)
	store.AddAppointment(appointment2)

	startTime := now
	endTime := now.Add(2 * time.Hour)

	availableSlots := store.GetTrainerAvailability(1, startTime, endTime)

	expectedSlots := []string{
		now.Add(90 * time.Minute).Format(time.RFC3339),
	}

	if !reflect.DeepEqual(expectedSlots, availableSlots) {
		t.Errorf("expected available slots to match, got %+v", availableSlots)
	}
}

func TestMemoryStore_SortedAppointments(t *testing.T) {
	testFile := "test_data.json"
	defer os.Remove(testFile)

	store := NewMemoryStore(testFile)

	appt1 := models.Appointment{
		TrainerID: 1,
		StartedAt: time.Now(),
		EndedAt:   time.Now().Add(time.Hour),
	}

	appt2 := models.Appointment{
		TrainerID: 1,
		StartedAt: time.Now().Add(2 * time.Hour),
		EndedAt:   time.Now().Add(3 * time.Hour),
	}

	store.AddAppointment(appt1)
	store.AddAppointment(appt2)

	sortedAppointments := store.GetAllAppointments()

	if len(sortedAppointments) != 2 {
		t.Fatalf("expected 2 appointments, got %d", len(sortedAppointments))
	}

	if !reflect.DeepEqual(sortedAppointments[0], appt1) && !reflect.DeepEqual(sortedAppointments[1], appt2) {
		t.Error("appointments are not sorted correctly")
	}
}
