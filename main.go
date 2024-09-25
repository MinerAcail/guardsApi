package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mineracail/guardApi/database"

	"github.com/mineracail/guardApi/router"
)

func main() {
	r := chi.NewRouter()
	// Apply the middleware to all routes
	//r.Use(middleware.Middleware)
	

	db := database.ConnectDB()
	// Migrate the schema	
	database.AutoMigrate(db)

	// Define routes for CRUD operations
	router.StudentRoute(db,r)
	router.StaffRoute(db,r)
	router.CalendarRoute(db,r)	
	router.ParentRoute(db,r)	
	router.LocationRoute(db,r)	
	
	log.Println("Starting server on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
