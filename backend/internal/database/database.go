package database

import (
	//"atlas-workspace/backend/internal/config"
	"fmt"
	"log"

	"github.com/CBYeuler/atlas-workspace/backend/internal/config"
	"github.com/CBYeuler/atlas-workspace/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	cfg := config.C
	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed toconnect to Atlas database: %v", err)
	}

	fmt.Println("Connected to Atlas database successfully")

	// Auto-migrate models
	DB.AutoMigrate(
		&models.User{},
		&models.Workspace{},
		&models.WorkspaceMember{},
		&models.Session{},
	)
}
