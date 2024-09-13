package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mineracail/guardApi/database"
	"github.com/mineracail/guardApi/resolvers"
)

func main() {
	r := chi.NewRouter()
	db := database.ConnectDB()

	// Define routes for CRUD operations
	r.Post("/students", func(w http.ResponseWriter, r *http.Request) {
		resolvers.CreateStudent(db, w, r)
	})
	r.Get("/students/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetStudentByID(db, w, r)
	})
	r.Get("/students/all", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetAllStudents(db, w, r)
	})
	r.Put("/students/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.UpdateStudentByID(db, w, r)	
	})
	r.Delete("/`students`/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.DeleteStudentByID(db, w, r)
	})

	log.Println("Starting server on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
