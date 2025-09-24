package main

import (
	"day1/internal/db"
	"day1/internal/seed"
)

func main() {
	db.InitDB()
	seed.RunSeed()
}
