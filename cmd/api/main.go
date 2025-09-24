package main

import (
	"fmt"
	"net/http"

	"day1/internal/db"
	"day1/internal/handler"
)

func main() {
	db.InitDB()

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetUsersHandler(w, r)
			return
		}
		if r.Method == http.MethodPost {
			handler.CreateUserHandler(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	fmt.Println("🚀 Server running at :8080")
	http.ListenAndServe(":8080", nil)
}
