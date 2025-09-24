package main

import (
	"day1/internal/db"
	"day1/internal/migrate"
)

func main() {
	db.InitDB()
	migrate.RunMigrations()
}
