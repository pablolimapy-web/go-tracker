package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/pablolimapy-web/go-tracker/internal/http/handler"
)

func New(h *handler.ShipmentHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/shipments", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/{id}", h.GetByID)
	})

	return r
}
