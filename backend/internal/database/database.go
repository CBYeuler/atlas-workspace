package database

import (
	//"atlas-workspace/backend/internal/config"
	"fmt"
	"log"

	"github.com/CBYeuler/atlas-workspace/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	cfg := config.C
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
}
