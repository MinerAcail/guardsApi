package database

import (
	"log"

	"github.com/mineracail/guardApi/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {

	// Migrate the schema
	err := db.AutoMigrate(
		&models.Student{},
		&models.Staff{},
		&models.Calendar{},
		&models.HomeArrival{},
		&models.SchoolArrival{},
		&models.Parent{},
		&models.Message{},
	)
	if err != nil {
		log.Fatal("Error migrating schema:", err)
	}

}
