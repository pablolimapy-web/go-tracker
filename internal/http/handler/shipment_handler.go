// internal/http/handler/shipment_handler.go
package handler

import (
	"encoding/json"
	"net/http"
	_ "strconv"

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

	// aqui você pegaria o userID do JWT, por enquanto mock:
	userID := int64(1)

	sh, err := h.service.CreateShipment(r.Context(), userID, body.Code, body.Carrier)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(sh)
}
