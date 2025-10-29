package orders

import (
	"encoding/json"
	"gophermart/internal/auth"
	"io"
	"net/http"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uint)
	body, _ := io.ReadAll(r.Body)
	number := strings.TrimSpace(string(body))

	if number == "" {
		http.Error(w, "empty order number", http.StatusBadRequest)
		return
	}

	status, err := h.service.CreateOrder(userID, number)
	if err != nil {
		if status == "conflict" {
			http.Error(w, "order belongs to another user", http.StatusConflict)
			return
		}
		if err.Error() == "invalid number" {
			http.Error(w, "invalid number", http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if status == "own" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uint)
	orders, err := h.service.GetUserOrders(userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
