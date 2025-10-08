package seed

import (
	"api/internal/db"
	"fmt"
	"time"
	"github.com/lib/pq"
)

// RunSeed chèn dữ liệu mẫu vào DB
func RunSeed() {
	seedUsers()
	seedMetrics()
	seedSensors()
	seedBoxes()
	fmt.Println("✅ Seeding done")
}

func seedUsers() {
	users := []struct {
		username string
		fullname string
		phone    string
	}{
		{"alice01", "Alice Wonderland", "0123456789"},
		{"bob02", "Bob Marley", "0987654321"},
		{"charlie03", "Charlie Brown", "0112233445"},
	}

	for _, u := range users {
		_, err := db.DB.Exec(
			"INSERT INTO users(username, fullname, phone) VALUES($1, $2, $3) ON CONFLICT (username) DO NOTHING",
			u.username, u.fullname, u.phone,
		)
		if err != nil {
			fmt.Println("⚠️ Seed user error:", err)
		}
	}
}

func seedMetrics() {
	metrics := []struct {
		code string
		name string
	}{
		{"RA", "Lượng mưa"},
		{"WA", "Mực nước"},
		{"DR", "Độ mở cống"},
	}

	now := time.Now().UnixMilli() // thời gian hiện tại (kiểu int64)

	for _, m := range metrics {
		_, err := db.DB.Exec(
			`INSERT INTO metrics (code, name, ctime)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (code) DO NOTHING`,
			m.code, m.name, now,
		)
		if err != nil {
			fmt.Println("⚠️ Seed metric error:", err)
		}
	}
}

func seedSensors() {
	sensors := []struct {
		name   string
		code   string
		unit   string
		metric string
	}{
		{"Độ mở cống 1", "DR1", "m", "DR"},
		{"Độ mở cống 2", "DR2", "m", "DR"},
		{"Độ mở cống 3", "DR3", "m", "DR"},
		{"Lượng mưa", "RA1", "mm", "RA"},
	}

	now := time.Now().UnixMilli() // thời gian hiện tại (int64)

	for _, s := range sensors {
		_, err := db.DB.Exec(
			`INSERT INTO sensors (name, code, unit, metric, ctime)
			 VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (code) DO NOTHING`,
			s.name, s.code, s.unit, s.metric, now,
		)
		if err != nil {
			fmt.Println("⚠️ Seed sensor error:", err)
		}
	}
}

func seedBoxes() {
	boxes := []struct {
		name         string
		sensors  []string
	}{
		{"Trạm đo mở cống", []string{"DR1", "DR2", "DR3"}},
		{"Trạm đo mưa", []string{"RA1"}},
	}

	now := time.Now().UnixMilli() // thời gian hiện tại (int64)

	for _, b := range boxes {
		_, err := db.DB.Exec(
			`INSERT INTO boxes (name, sensors, ctime)
			 VALUES ($1, $2, $3)`,
			b.name, pq.Array(b.sensors), now,
		)
		if err != nil {
			fmt.Println("⚠️ Seed boxs error:", err)
		}
	}
}
