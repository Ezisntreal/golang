package seed

import (
	"day1/internal/db"
	"fmt"
)

// RunSeed chèn dữ liệu mẫu vào DB
func RunSeed() {
	seedUsers()
	seedSensors()
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

func seedSensors() {
	sensors := []struct {
		code string
		name string
	}{
		{"SEN001", "Temperature Sensor"},
		{"SEN002", "Humidity Sensor"},
		{"SEN003", "Pressure Sensor"},
	}

	for _, s := range sensors {
		_, err := db.DB.Exec(
			"INSERT INTO sensors(code, name) VALUES($1, $2) ON CONFLICT (code) DO NOTHING",
			s.code, s.name,
		)
		if err != nil {
			fmt.Println("⚠️ Seed sensor error:", err)
		}
	}
}
