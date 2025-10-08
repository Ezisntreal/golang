package migrate

import "api/internal/db"

func CreateBoxTable() error {
	query := `
	DROP TABLE IF EXISTS boxes CASCADE;
	
	CREATE TABLE boxes (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		sensors TEXT[],
		ctime BIGINT,
		dtime BIGINT
	);`
	
	_, err := db.DB.Exec(query)
	return err
}