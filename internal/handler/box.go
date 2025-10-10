package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"api/internal/db"
	"github.com/lib/pq"
)

type Box struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Sensors []string `json:"sensors"`
	Ctime   *int64   `json:"ctime"`
	DTime   *int64   `json:"dtime,omitempty"`
}

// 🔹 GET /box
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
		var sensorsStr *string // đọc text[] từ PostgreSQL thành chuỗi
		if err := rows.Scan(&b.ID, &b.Name, &sensorsStr, &b.Ctime, &b.DTime); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// ✅ Parse thủ công thành []string
		if sensorsStr != nil && *sensorsStr != "" {
			s := strings.Trim(*sensorsStr, "{}[]") // bỏ {}, [] nếu có
			s = strings.ReplaceAll(s, `"`, "")     // bỏ dấu "
			if s != "" {
				b.Sensors = strings.Split(s, ",")
				for i := range b.Sensors {
					b.Sensors[i] = strings.TrimSpace(b.Sensors[i])
				}
			}
		} else {
			b.Sensors = []string{}
		}

		boxes = append(boxes, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boxes)
}

// 🔹 POST /box
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

	// ❌ bỏ json.Marshal
	err := db.DB.QueryRowContext(ctx,
		"INSERT INTO boxes(name, sensors, ctime) VALUES($1, $2, $3) RETURNING id",
		b.Name, pq.Array(b.Sensors), *b.Ctime,
	).Scan(&b.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
}

// 🔹 PUT /box?id=1
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

	// ✅ Không dùng json.Marshal — dùng pq.Array vì cột là TEXT[]
	_, err = db.DB.ExecContext(ctx,
		`UPDATE boxes 
		 SET name = $1, sensors = $2 
		 WHERE id = $3 AND dtime IS NULL`,
		b.Name, pq.Array(b.Sensors), id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "✅ updated successfully"}`))
}


// 🔹 DELETE /box?id=1
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
