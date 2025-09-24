package migrate

import (
	"fmt"
)

func RunMigrations() {
	if err := createUserTable(); err != nil {
		panic(fmt.Sprintf("❌ Failed to migrate users: %v", err))
	}

	if err := createSensorTable(); err != nil {
		panic(fmt.Sprintf("❌ Failed to migrate sensors: %v", err))
	}

	fmt.Println("✅ All migrations done")
}
