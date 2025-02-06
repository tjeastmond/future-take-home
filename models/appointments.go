package models

import "time"

type Appointment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	TrainerID int       `json:"trainer_id"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}
