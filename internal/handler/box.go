package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"api/internal/db"
)

type Box struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Sensors  []string `json:"sensors"`
	Ctime    *int64   `json:"ctime"`
	DTime    *int64   `json:"dtime,omitempty"`
}

// 🔹 GET /boxes
func GetBoxesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	rows, err := db.DB.QueryContext(ctx,
		"SELECT id, name, sensors, ctime, dtime FROM boxes WHERE dtime IS NULL",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var boxes []Box
	for rows.Next() {
		var b Box
		var sensorsJSON []byte
		if err := rows.Scan(&b.ID, &b.Name, &sensorsJSON, &b.Ctime, &b.DTime); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = json.Unmarshal(sensorsJSON, &b.Sensors)
		boxes = append(boxes, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boxes)
}

// 🔹 POST /boxes
func CreateBoxHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var b Box
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	now := time.Now().UnixMilli()
	b.Ctime = &now

	sensorsJSON, err := json.Marshal(b.Sensors)
	if err != nil {
		http.Error(w, "invalid sensors data", http.StatusBadRequest)
		return
	}

	err = db.DB.QueryRowContext(ctx,
		"INSERT INTO boxes(name, sensors, ctime) VALUES($1, $2, $3) RETURNING id",
		b.Name, sensorsJSON, b.Ctime,
	).Scan(&b.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
}

// 🔹 PUT /boxes?id=1
func UpdateBoxHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id param", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var b Box
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	sensorsJSON, err := json.Marshal(b.Sensors)
	if err != nil {
		http.Error(w, "invalid sensors data", http.StatusBadRequest)
		return
	}

	_, err = db.DB.ExecContext(ctx,
		`UPDATE boxes 
		 SET name = $1, sensors = $2 
		 WHERE id = $3 AND dtime IS NULL`,
		b.Name, sensorsJSON, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "✅ updated successfully"}`))
}

// 🔹 DELETE /boxes?id=1
func DeleteBoxHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id param", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	now := time.Now().UnixMilli()
	_, err = db.DB.ExecContext(ctx,
		"UPDATE boxes SET dtime = $1 WHERE id = $2 AND dtime IS NULL",
		now, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "🗑️ deleted successfully"}`))
}