package database

import (
	"log"

	"github.com/mineracail/guardApi/models"
	"gorm.io/gorm"
)

func AutoMigrate(	) {

	// Migrate the schema
	err := db.AutoMigrate(&models.Student{})
	if err != nil {
		log.Fatal("Error migrating schema:", err)
	}

}
