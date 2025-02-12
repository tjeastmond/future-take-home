package models

import (
	"errors"
	"time"
)

type Appointment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id" binding:"required"`
	TrainerID int       `json:"trainer_id" binding:"required"`
	StartedAt time.Time `json:"started_at" binding:"required"`
	EndedAt   time.Time `json:"ended_at"`
}

func (a *Appointment) Validate() error {
	if a.StartedAt.IsZero() {
		return errors.New("started_at cannot be zero")
	}

	if a.EndedAt.IsZero() {
		return errors.New("ended_at cannot be zero")
	}

	if a.EndedAt.Before(a.StartedAt) {
		return errors.New("ended_at must be after started_at")
	}
	return nil
}

func NewAppointment(userID, trainerID int, startedAt, endedAt time.Time) (*Appointment, error) {
	appointment := &Appointment{
		UserID:    userID,
		TrainerID: trainerID,
		StartedAt: startedAt,
		EndedAt:   endedAt,
	}

	if err := appointment.Validate(); err != nil {
		return nil, err
	}

	return appointment, nil
}
