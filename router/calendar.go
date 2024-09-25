package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mineracail/guardApi/resolvers"
	"gorm.io/gorm"
)

func CalendarRoute(db *gorm.DB, r *chi.Mux) {
	// Define routes for CRUD operations
	r.Post("/calendars", func(w http.ResponseWriter, r *http.Request) {
		resolvers.CreateCalendar(db, w, r)
	})
	r.Get("/calendars/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetCalendarByID(db, w, r)
	})
	r.Get("/calendars/all", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetAllCalendars(db, w, r)
	})
	r.Put("/calendars/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.UpdateCalendarByID(db, w, r)
	})
	r.Delete("/calendars/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.DeleteCalendarByID(db, w, r)
	})
}
