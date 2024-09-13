	package database		

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"	
	"github.com/joho/godotenv"
)

// ConnectDB establishes a connection to the database
func ConnectDB() *gorm.DB {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {		
		log.Fatal("Error loading .env file")
	}

	// Get environment variables
	user := os.Getenv("GO_DB_USER")
	password := os.Getenv("GO_DB_PASSWORD")
	dbname := os.Getenv("GO_DB_NAME")	
	host := os.Getenv("GO_DB_HOST")
	port := os.Getenv("GO_DB_PORT")

	// Define the connection string with service name `postgres`
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Open a connection to the database using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	// defer func() {
	// 	if err := sqlDB.Close(); err != nil {
	// 		log.Fatal("Error closing the database connection:", err)
	// 	}		
	// }()

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database using GORM!")

	return db
}
