package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mineracail/guardApi/resolvers"
	"gorm.io/gorm"
)

func MessageRoute( db *gorm.DB,r *chi.Mux) {
	r.Post("/messages", func(w http.ResponseWriter, r *http.Request) {
		resolvers.CreateMessage(db, w, r)
	})
	r.Post("/messages/multiple", func(w http.ResponseWriter, r *http.Request) {
		resolvers.CreateMessageToMultiple(db, w, r)
	})
	r.Get("/messages/all", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetAllMessages(db, w, r)
	})
	r.Get("/messages/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.GetMessageByID(db, w, r)
	})
	r.Put("/messages/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.UpdateMessageByID(db, w, r)
	})
	r.Delete("/messages/{id}", func(w http.ResponseWriter, r *http.Request) {
		resolvers.DeleteMessageByID(db, w, r)
	})
}
