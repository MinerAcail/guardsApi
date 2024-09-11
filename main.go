package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"github.com/joho/godotenv"
)

func main() {
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
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port)

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database!")
}
	