package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mineracail/guardApi/resolvers"
	"gorm.io/gorm"
)

func ParentRoute(db *gorm.DB, r *chi.Mux) {
	// Define routes for CRUD operations for Parent
	r.Post("/parents", func(w http.ResponseWriter, r *http.Request) {
		resolvers.CreateParent(db, w, r)
	})
	r.Get("/parents/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetParentByID(db, w, r)
	})
	r.Get("/parentschild/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetChildByParentID(db, w, r)
	})
	r.Get("/parents/all", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetAllParents(db, w, r)
	})
	r.Put("/parents/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.UpdateParentByID(db, w, r)
	})
	r.Delete("/parents/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.DeleteParentByID(db, w, r)
	})
}

