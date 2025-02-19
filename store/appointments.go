package store

import (
	"errors"
	"time"

	"github.com/tjeastmond/future-take-home/config"
	"github.com/tjeastmond/future-take-home/utils"
)

type Appointment struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id" binding:"required"`
	TrainerID int       `db:"trainer_id" json:"trainer_id" binding:"required"`
	StartsAt  time.Time `db:"starts_at" json:"starts_at" binding:"required"`
	EndsAt    time.Time `db:"ends_at" json:"ends_at" binding:"required"`
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

func (a *Appointment) ToPst() error {
	location, err := time.LoadLocation(config.LocationName)
	if err != nil {
		return err
	}

	a.StartsAt = a.StartsAt.In(location)
	a.EndsAt = a.EndsAt.In(location)

	return nil
}

type Appointments struct{}

func (app *Appointments) GetAppointments(trainerID int) ([]*Appointment, error) {
	var appointments []*Appointment
	var err error

	query := `SELECT * FROM appointments WHERE trainer_id = $1`
	err = db.Select(&appointments, query, trainerID)
	if err != nil {
		return nil, err
	}

	for _, appointment := range appointments {
		if err = appointment.ToPst(); err != nil {
			return nil, err
		}
	}

	return appointments, nil
}

func (app *Appointments) CreateAppointment(appointment *Appointment) error {
	if err := appointment.Validate(); err != nil {
		return err
	}

	appointment.ToPst()

	_ = db.MustExec(
		`INSERT INTO appointments (user_id, trainer_id, starts_at, ends_at) VALUES ($1, $2, $3, $4)`,
		appointment.UserID, appointment.TrainerID, appointment.StartsAt, appointment.EndsAt,
	)

	return nil
}

func (app *Appointments) GetAvailableSlots(trainer_id int, start, end time.Time) ([]time.Time, error) {
	var slots []time.Time

	query := `
	WITH slots AS (
		SELECT generate_series(
			$1::timestamptz,
			$2::timestamptz,
			'30 minutes'::interval
		) AS slot_start
	)

	SELECT s.slot_start
	FROM slots s
	WHERE NOT EXISTS (
		SELECT 1
		FROM appointments a
		WHERE a.trainer_id = $3
			AND s.slot_start < a.ends_at
			AND (s.slot_start + interval '30 minutes') > a.starts_at
	)
	ORDER BY s.slot_start;
	`

	err := db.Select(&slots, query, start, end, trainer_id)
	if err != nil {
		return nil, err
	}

	var validSlots []time.Time
	for _, slot := range slots {
		parsedSlot, _ := utils.ParseTimeToPST(slot)
		if valid := utils.IsValidTime(parsedSlot); valid {
			validSlots = append(validSlots, parsedSlot)
		}
	}

	return validSlots, nil
}

func (app *Appointments) IsAvailable(appointment *Appointment) bool {
	var count int
	_ = db.Get(
		&count,
		`SELECT COUNT(*) FROM appointments WHERE trainer_id = $1 AND starts_at < $2 AND ends_at > $3`,
		appointment.TrainerID, appointment.EndsAt, appointment.StartsAt,
	)

	return count == 0
}
