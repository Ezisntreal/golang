package main

import (
	"fmt"
	"net/http"

	"api/internal/db"
	"api/internal/router"
)

func main() {
	db.InitDB()

	router.RegisterRoutes()

	fmt.Println("🚀 Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
