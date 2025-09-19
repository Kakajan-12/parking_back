package main

import (
	"backend/config"
	"backend/database"
)

func main() {

	// Load configuration
	config.LoadConfig()
	database.InitDb(config.AppConfig.DatabaseURL)
}
