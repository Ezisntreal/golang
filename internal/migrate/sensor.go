package migrate

import "api/internal/db"

func CreateSensorTable() error {
	query := `
		DROP TABLE IF EXISTS sensors CASCADE;

		CREATE TABLE sensors (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			code TEXT UNIQUE NOT NULL,
			unit TEXT,
			metric TEXT,
			position JSONB DEFAULT '{}',
			min DOUBLE PRECISION,
			max DOUBLE PRECISION,
			ctime BIGINT,
			mtime BIGINT,
			dtime BIGINT
		);`

	_, err := db.DB.Exec(query)
	return err
}