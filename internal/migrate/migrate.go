package migrate

import (
	"fmt"
)

func RunMigrations() {
	fmt.Println("🚀 Starting migration...")

	if err := CreateUserTable(); err != nil {
		panic(fmt.Sprintf("❌ Failed to migrate users: %v", err))
	}

	if err := CreateMetricTable(); err != nil {
		panic(fmt.Sprintf("❌ Failed to migrate metrics: %v", err))
	}

	if err := CreateBoxTable(); err != nil {
		panic(fmt.Sprintf("❌ Failed to migrate boxs: %v", err))
	}

	if err := CreateSensorTable(); err != nil {
		panic(fmt.Sprintf("❌ Failed to migrate sensors: %v", err))
	}

	fmt.Println("✅ All migrations done")
}
