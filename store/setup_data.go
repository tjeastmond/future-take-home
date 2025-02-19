package store

// func SetupData() {
// 	var err error

// 	createTableSQL := `
// 	CREATE TABLE IF NOT EXISTS appointments (
// 		id SERIAL PRIMARY KEY,
// 		trainer_id INT NOT NULL,
// 		user_id INT NOT NULL,
// 		starts_at TIMESTAMPTZ NOT NULL,
// 		ends_at TIMESTAMPTZ NOT NULL,
// 		CONSTRAINT appointments_trainer_start UNIQUE (trainer_id, starts_at)
// 	);
// 	`

// 	_, err = db.Exec(createTableSQL)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	createIndexesSQL := `
// 	CREATE INDEX IF NOT EXISTS idx_appointments_id ON appointments(id);
// 	CREATE INDEX IF NOT EXISTS idx_appointments_trainer_id ON appointments(trainer_id);
// 	CREATE INDEX IF NOT EXISTS idx_appointments_user_id ON appointments(user_id);
// 	CREATE INDEX IF NOT EXISTS idx_appointments_starts_at ON appointments(starts_at);
// 	`

// 	_, err = db.Exec(createIndexesSQL)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	pst, err := time.LoadLocation("America/Los_Angeles")
// 	if err != nil {
// 		log.Fatalf("Error loading PST location: %v", err)
// 	}

// 	// appointments := []Appointment{
// 	// 	{{REWRITTEN_CODE}}
// 	// 			{ID: 1, UserID: 1, TrainerID: 1, StartsAt: "2019-01-24 09:00:00-08", EndsAt: "2019-01-24 09:30:00-08"},
// 	// 			{ID: 2, UserID: 2, TrainerID: 1, StartsAt: "2019-01-24 10:00:00-08", EndsAt: "2019-01-24 10:30:00-08"},
// 	// 			{ID: 3, UserID: 3, TrainerID: 1, StartsAt: "2019-01-25 10:00:00-08", EndsAt: "2019-01-25 10:30:00-08"},
// 	// 			{ID: 4, UserID: 4, TrainerID: 1, StartsAt: "2019-01-25 10:30:00-08", EndsAt: "2019-01-25 11:00:00-08"},
// 	// 			{ID: 5, UserID: 5, TrainerID: 1, StartsAt: "2019-01-26 10:00:00-08", EndsAt: "2019-01-26 10:30:00-08"},
// 	// 			{ID: 6, UserID: 6, TrainerID: 2, StartsAt: "2019-01-24 09:00:00-08", EndsAt: "2019-01-24 09:30:00-08"},
// 	// 			{ID: 7, UserID: 7, TrainerID: 2, StartsAt: "2019-01-26 10:00:00-08", EndsAt: "2019-01-26 10:30:00-08"},
// 	// 			{ID: 8, UserID: 8, TrainerID: 3, StartsAt: "2019-01-26 12:00:00-08", EndsAt: "2019-01-26 12:30:00-08"},
// 	// 			{ID: 9, UserID: 9, TrainerID: 3, StartsAt: "2019-01-26 13:00:00-08", EndsAt: "2019-01-26 13:30:00-08"},
// 	// 			{ID: 10, UserID: 10, TrainerID: 3, StartsAt: "2019-01-26 14:00:00-08", EndsAt: "2019-01-26 14:30:00-08"},
// 	// }

// 	// for _, appointment := range appointments {
// 	// 	insertSQL := `
// 	// 		INSERT INTO appointments (id, trainer_id, user_id, starts_at, ends_at)
// 	// 		VALUES ($1, $2, $3, $4, $5)
// 	// 	`

// 	// 	fmt.Println(insertSQL)

// 		_, err = db.Exec(insertSQL, appointment.ID, appointment.TrainerID, appointment.UserID, appointment.StartsAt, appointment.EndsAt)
// 		if err != nil {
// 			log.Fatalf("Error inserting appointment with ID %d: %v", appointment.ID, err)
// 		}
// 	}

// 	// fmt.Println("Table, indexes, and data inserted successfully.")
// }

// func Down() {
// 	dropTableSQL := `DROP TABLE IF EXISTS appointments;`
// 	_, err := db.Exec(dropTableSQL)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func parseTime(t string, loc *time.Location) time.Time {
// 	parsedTime, err := time.ParseInLocation(time.RFC3339, t, loc)
// 	if err != nil {
// 		log.Fatalf("Error parsing time: %v", err)
// 	}

// 	return parsedTime
// }
