package database

import (
	"backend/contrib/models"
	"backend/contrib/repository"
	"log"
	"os"

	"gorm.io/gorm"
)

func initSuperUser(db *gorm.DB) {
	repo := repository.NewUserRepository(db)

	superUsername := os.Getenv("FIRST_SUPER_USERNAME")
	superPassword := os.Getenv("FIRST_SUPERUSER_PASSWORD")

	if superUsername == "" || superPassword == "" {
		log.Fatal("Environment variables FIRST_SUPER_USERNAME and FIRST_SUPERUSER_PASSWORD must be set")
	}
	includeSuperuser := true
	exists, err := repo.UserExist(&repository.UserRepoFilter{
		Username:         &superUsername,
		IncludeSuperuser: &includeSuperuser,
	})
	if err != nil {
		log.Fatalf("Failed to check if superuser exists: %v", err)
	}

	if exists {
		log.Println("Superuser already exists, skipping creation")
		return
	}

	isSuper := true
	// Create superuser
	user, err := repo.UserCreate(
		superUsername,
		superPassword,
		"Admin User",
		models.AdminRole,
		true,
		nil,
		&isSuper,
	)
	if err != nil {
		log.Fatalf("Failed to create superuser: %v", err)
	}

	log.Printf("Superuser successfully created: ID %d", user.ID)
}

func initMacUser(db *gorm.DB) {

	macUsername := os.Getenv("MACROSCOP_USERNAME")
	macPassword := os.Getenv("MACROSCOP_PASSWORD")

	if macUsername == "" || macPassword == "" {
		log.Fatal("Environment variables MACROSCOP_USERNAME and MACROSCOP_PASSWORD must be set")
	}

	repo := repository.NewMacUserRepository(db)

	exists, err := repo.MacUserExist(1)
	if err != nil {
		log.Fatalf("Failed to check if macuser exists: %v", err)
	}

	if exists {
		log.Println("MacUser already exists, skipping creation")
		return
	}

	// Create MacUser
	newMacUser := &models.MacUser{
		MacUsername: macUsername,
		MacPassword: macPassword,
	}

	if err := repo.MacUserCreate(newMacUser); err != nil {
		log.Fatalf("Failed to create MacUser: %v", err)
	}

	log.Println("Mac user successfully created")
}

func InitDb(dsn string) {
	// Check if superuser exists

	db, err := ConnectDB(dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	initSuperUser(db)
	initMacUser(db)
	log.Println("Database initialization complete")
}
