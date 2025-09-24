package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	godotenv.Load()
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("❌ DB_URL is not set in environment")
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("cannot connect db:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("cannot ping db:", err)
	}

	fmt.Println("✅ Connected to DB")
}
