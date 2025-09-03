package database

import (
	"learning-platform/backend/internal/config"
	"learning-platform/backend/internal/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase(cfg *config.Config) {
	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Course{},
		&models.Enrollment{},
		&models.Lesson{},
		&models.Assignment{},
		&models.TestCase{},
		&models.Solution{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connection established")
}
