package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mineracail/guardApi/resolvers"
	"gorm.io/gorm"
)

func StaffRoute(db *gorm.DB, r *chi.Mux) {
	// Define routes for CRUD operations
	r.Post("/staffs", func(w http.ResponseWriter, r *http.Request) {
		resolvers.CreateStaff(db, w, r)
	})
	r.Get("/staffs/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetStaffByID(db, w, r)
	})
	r.Get("/staffs/all", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetAllStaffs(db, w, r)
	})
	r.Put("/staffs/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.UpdateStaffByID(db, w, r)
	})
	r.Delete("/staffs/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.DeleteStaffByID(db, w, r)
	})
}
