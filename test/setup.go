package test

import (
	"log"
	"os"
	"os/exec"

	"gorm.io/gorm"

	"backend/config"
	"backend/database"
)

var TestDB *gorm.DB // global DB for all tests

// InitTestDB initializes DB and runs migrations for test env
func InitTestDB() {
	// Set test environment
	if err := os.Setenv("APP_ENV", "test"); err != nil {
		panic(err)
	}

	config.LoadConfig()

	// Connect to database (so tests can use it)
	var err error
	TestDB, err = database.ConnectDB(config.AppConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}

	// Run Atlas migrations for test environment
	runAtlasMigrateUp()
}

// CleanupTestDB drops all tables at the end for test env
func CleanupTestDB() {
	runAtlasMigrateDown()
}

// runAtlasMigrateUp runs `atlas migrate apply --env test`
func runAtlasMigrateUp() {
	cmd := exec.Command("atlas", "migrate", "apply", "--env", "test")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run atlas migrate apply: %v", err)
	}
}

// runAtlasMigrateDown drops all tables using Atlas for test env
func runAtlasMigrateDown() {
	cmd := exec.Command("atlas", "migrate", "drop", "-yes", "--env", "test")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run atlas migrate drop: %v", err)
	}
}
