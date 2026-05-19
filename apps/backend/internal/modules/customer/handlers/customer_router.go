package handlers

import (
	"github.com/go-chi/chi/v5"
)

func (h *CustomerHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Deactivate)
	return r
}
