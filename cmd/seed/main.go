package main

import (
	"api/internal/db"
	"api/internal/seed"
)

func main() {
	db.InitDB()
	seed.RunSeed()
}
