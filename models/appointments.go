package models

import (
	"errors"
	"time"
)

type Appointment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id" binding:"required"`
	TrainerID int       `json:"trainer_id" binding:"required"`
	StartsAt  time.Time `json:"starts_at" binding:"required"`
	EndsAt    time.Time `json:"ends_at"`
}

func (a *Appointment) Validate() error {
	if a.StartsAt.IsZero() {
		return errors.New("started_at cannot be zero")
	}

	if a.EndsAt.IsZero() {
		return errors.New("ended_at cannot be zero")
	}

	if a.EndsAt.Before(a.StartsAt) {
		return errors.New("ended_at must be after started_at")
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
