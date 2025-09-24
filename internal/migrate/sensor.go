package migrate

import "day1/internal/db"

func createSensorTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS sensors (
		id SERIAL PRIMARY KEY,
		code TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL
	);`

	_, err := db.DB.Exec(query)
	return err
}