package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mineracail/guardApi/middleware"
	"github.com/mineracail/guardApi/resolvers"
	"gorm.io/gorm"
)

func StudentRoute(db *gorm.DB, r *chi.Mux) {
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
	r.Delete("/students/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.DeleteStudentByID(db, w, r)
	})
		// Login route for authentication
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			middleware.Login(db, w, r)
		})
}
