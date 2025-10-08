package main

import (
	"api/internal/db"
	"api/internal/migrate"
)

func main() {
	db.InitDB()
	migrate.RunMigrations()
}
