package migrate

import "api/internal/db"

func CreateMetricTable() error {
	query := `
		DROP TABLE IF EXISTS metrics CASCADE;

		CREATE TABLE IF NOT EXISTS metrics (
		id SERIAL PRIMARY KEY,
		code TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		ctime BIGINT,
		dtime BIGINT
	);`

	_, err := db.DB.Exec(query)
	return err
}
