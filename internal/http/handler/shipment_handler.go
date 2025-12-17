package handler

import (
	"encoding/json"
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

type createShipmentRequest struct {
	Code    string `json:"code"`
	Carrier string `json:"carrier"`
}

func (h *ShipmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createShipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Sem auth por enquanto: user fixo (depois vira JWT)
	userID := int64(1)

	created, err := h.service.Create(r.Context(), userID, req.Code, req.Carrier)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(created)
}

func (h *ShipmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	item, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if err == shipment.ErrNotFound {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(item)
}

func (h *ShipmentHandler) List(w http.ResponseWriter, r *http.Request) {
	// Sem auth por enquanto: user fixo (depois vira JWT)
	userID := int64(1)

	items, err := h.service.ListByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}
