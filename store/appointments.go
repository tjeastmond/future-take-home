package appointments

import (
	"errors"
	"fmt"
	"time"
)

type Appointment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id" binding:"required"`
	TrainerID int       `json:"trainer_id" binding:"required"`
	StartsAt  time.Time `json:"starts_at" binding:"required"`
	EndsAt    time.Time `json:"ends_at" binding:"required"`
}

func (a *Appointment) Validate() error {
	if a.StartsAt.IsZero() {
		return errors.New("starts_at cannot be zero")
	}

	if a.EndsAt.IsZero() {
		return errors.New("ends_at cannot be zero")
	}

	if a.EndsAt.Before(a.StartsAt) {
		return errors.New("ends_at must be after starts_at")
	}

	return nil
}

func NewAppointment(userID, trainerID int, startsAt, endsAt time.Time) (*Appointment, error) {
	appointment := &Appointment{
		UserID:    userID,
		TrainerID: trainerID,
		StartsAt:  startsAt,
		EndsAt:    endsAt,
	}

	if err := appointment.Validate(); err != nil {
		return nil, err
	}

	return appointment, nil
}

type Appointments struct{}

func (as *Appointments) GetAppointments(trainer_id int) ([]Appointment, error) {
	query := `SELECT * FROM appointments WHERE trainer_id = $1`
	results, err := db.Query(query, trainer_id)
	if err != nil {
		return nil, err
	}

	fmt.Println(results)

	return nil, nil
}
