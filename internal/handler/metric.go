package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"api/internal/db"
)

// Metric model
type Metric struct {
	ID    int     `json:"id"`
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Ctime *int64  `json:"ctime"`
	DTime *int64  `json:"dtime,omitempty"`
}

// 🔹 GET /metrics
func GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	rows, err := db.DB.QueryContext(ctx,
		"SELECT id, code, name, ctime, dtime FROM metrics WHERE dtime IS NULL",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var metrics []Metric
	for rows.Next() {
		var m Metric
		if err := rows.Scan(&m.ID, &m.Code, &m.Name, &m.Ctime, &m.DTime); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		metrics = append(metrics, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// 🔹 POST /metrics
func CreateMetricHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var m Metric
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	now := time.Now().UnixMilli()
	m.Ctime = &now

	err := db.DB.QueryRowContext(ctx,
		"INSERT INTO metrics(code, name, ctime) VALUES($1, $2, $3) RETURNING id",
		m.Code, m.Name, m.Ctime,
	).Scan(&m.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

// 🔹 PUT /metrics?id=1
func UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
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

	var m Metric
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	_, err = db.DB.ExecContext(ctx,
		`UPDATE metrics 
		 SET code = $1, name = $2 
		 WHERE id = $3 AND dtime IS NULL`,
		m.Code, m.Name, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "✅ updated successfully"}`))
}

// 🔹 DELETE /metrics?id=1
func DeleteMetricHandler(w http.ResponseWriter, r *http.Request) {
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
		"UPDATE metrics SET dtime = $1 WHERE id = $2 AND dtime IS NULL",
		now, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "🗑️ deleted successfully"}`))
}
