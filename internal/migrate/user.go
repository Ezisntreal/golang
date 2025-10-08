package migrate

import "api/internal/db"

func CreateUserTable() error {
	query := `
	DROP TABLE IF EXISTS users CASCADE;
	
	CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		fullname TEXT NOT NULL,
		phone TEXT,
		dtime BIGINT NULL
	);`
	
	_, err := db.DB.Exec(query)
	return err
}