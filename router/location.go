package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mineracail/guardApi/resolvers"
	"gorm.io/gorm"
)

func LocationRoute(db *gorm.DB, r *chi.Mux) {
	// Define routes for CRUD operations
	r.Post("/locationbyparent", func(w http.ResponseWriter, r *http.Request) {
		resolvers.CreateHomeArrival(db, w, r)
	})
	r.Post("/locationbystaff", func(w http.ResponseWriter, r *http.Request) {
		resolvers.CreateSchoolArrival(db, w, r)
	})
	r.Get("/locationbyparent/all", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetAllHomeArrivalsForThatWeek(db, w, r)
	})
	r.Get("/locationbyparent/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetAllHomeArrivalsForThatWeekByParentId(db, w, r)
	})

	r.Get("/locationbyday/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetConfirmedArrivalsByParent(db, w, r)
	})
	r.Get("/stafflocationbyday/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetConfirmedArrivalsByStaff(db, w, r)
	})

	r.Get("/locationbyday/all", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetAllConfirmedArrivals(db, w, r)
	})
	r.Get("/stafflocationbyday/all", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetAllConfirmedArrivalsStaff(db, w, r)
	})

}
