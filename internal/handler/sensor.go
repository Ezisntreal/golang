package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"api/internal/db"
)

type Position struct {
	Top  *float64 `json:"top,omitempty"`
	Left *float64 `json:"left,omitempty"`
	Lat  *float64 `json:"lat,omitempty"`
	Lng  *float64 `json:"lng,omitempty"`
}

type Sensor struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Code     string    `json:"code"`
	Unit     string    `json:"unit"`
	Metric   string    `json:"metric"`
	Position *Position `json:"position"`
	Min      *float64  `json:"min,omitempty"`
	Max      *float64  `json:"max,omitempty"`
	Ctime    *int64    `json:"ctime"`
	Mtime    *int64    `json:"mtime,omitempty"`
	DTime    *int64    `json:"dtime,omitempty"`
}

// 🔹 GET /sensors
func GetSensorsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	rows, err := db.DB.QueryContext(ctx,
		`SELECT id, name, code, unit, metric, position, min, max, ctime, mtime, dtime 
		 FROM sensors WHERE dtime IS NULL`,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var sensors []Sensor
	for rows.Next() {
		var s Sensor
		var posJSON []byte
		if err := rows.Scan(&s.ID, &s.Name, &s.Code, &s.Unit, &s.Metric, &posJSON,
			&s.Min, &s.Max, &s.Ctime, &s.Mtime, &s.DTime); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(posJSON) > 0 {
			_ = json.Unmarshal(posJSON, &s.Position)
		}
		sensors = append(sensors, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sensors)
}

// 🔹 POST /sensors
func CreateSensorHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var s Sensor
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	now := time.Now().UnixMilli()
	s.Ctime = &now
	s.Mtime = &now

	posJSON, err := json.Marshal(s.Position)
	if err != nil {
		http.Error(w, "invalid position data", http.StatusBadRequest)
		return
	}

	err = db.DB.QueryRowContext(ctx,
		`INSERT INTO sensors(name, code, unit, metric, position, min, max, ctime, mtime) 
		 VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		s.Name, s.Code, s.Unit, s.Metric, posJSON, s.Min, s.Max, s.Ctime, s.Mtime,
	).Scan(&s.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// 🔹 PUT /sensors?id=1
func UpdateSensorHandler(w http.ResponseWriter, r *http.Request) {
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

	var s Sensor
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	now := time.Now().UnixMilli()
	s.Mtime = &now

	posJSON, err := json.Marshal(s.Position)
	if err != nil {
		http.Error(w, "invalid position data", http.StatusBadRequest)
		return
	}

	_, err = db.DB.ExecContext(ctx,
		`UPDATE sensors 
		 SET name = $1, code = $2, unit = $3, metric = $4, 
		     position = $5, min = $6, max = $7, mtime = $8
		 WHERE id = $9 AND dtime IS NULL`,
		s.Name, s.Code, s.Unit, s.Metric, posJSON, s.Min, s.Max, s.Mtime, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "✅ updated successfully"}`))
}

// 🔹 DELETE /sensors?id=1
func DeleteSensorHandler(w http.ResponseWriter, r *http.Request) {
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
		"UPDATE sensors SET dtime = $1 WHERE id = $2 AND dtime IS NULL",
		now, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "🗑️ deleted successfully"}`))
}
