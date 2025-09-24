package migrate

import "day1/internal/db"

func createUserTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		fullname TEXT NOT NULL,
		phone TEXT
	);`

	_, err := db.DB.Exec(query)
	return err
}