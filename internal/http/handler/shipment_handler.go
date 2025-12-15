package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/pablolimapy-web/go-tracker/internal/domain/shipment"
)

type ShipmentHandler struct {
	service *shipment.Service
}

func NewShipmentHandler(s *shipment.Service) *ShipmentHandler {
	return &ShipmentHandler{service: s}
}

func (h *ShipmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Code    string `json:"code"`
		Carrier string `json:"carrier"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.Code == "" || body.Carrier == "" {
		http.Error(w, "code and carrier are required", http.StatusBadRequest)
		return
	}

	// por enquanto mock; depois vem do JWT
	userID := int64(1)

	sh, err := h.service.CreateShipment(r.Context(), userID, body.Code, body.Carrier)
	if err != nil {
		if errors.Is(err, shipment.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(sh)
}

func (h *ShipmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	sh, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, shipment.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, shipment.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sh)
}
