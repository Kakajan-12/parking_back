package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	env := os.Getenv("APP_ENV")
	logLevel := logger.Silent

	if env == "development" || env == "test" {
		logLevel = logger.Info // log queries in dev and test
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Successfully connected to PostgreSQL")
	return db, nil
}
