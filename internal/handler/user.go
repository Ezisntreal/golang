package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"day1/internal/db"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Phone    string `json:"phone"`
}

// GET /users
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	rows, err := db.DB.QueryContext(ctx,
		"SELECT id, username, fullname, phone FROM users",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Fullname, &u.Phone); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// POST /users
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err := db.DB.QueryRowContext(ctx,
		"INSERT INTO users(username, fullname, phone) VALUES($1, $2, $3) RETURNING id",
		u.Username, u.Fullname, u.Phone,
	).Scan(&u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}
