package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type GetOrderHandler struct {
	ParamFromRequest func(r *http.Request, key string) string
}

type Order struct {
	UUID uuid.UUID `json:"uuid"`
}

func (h *GetOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	orderUUID := h.ParamFromRequest(r, "orderUUID")
	pOrderUUID, err := uuid.Parse(orderUUID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call your application service here
	order := Order{
		UUID: pOrderUUID,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
