package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/timurkash/queue2/internal/data"
)

type Handler struct {
	queueSvc       *data.Service
	defaultTimeout time.Duration
}

func New(queueSvc *data.Service, defaultTimeout time.Duration) *Handler {
	return &Handler{
		queueSvc:       queueSvc,
		defaultTimeout: defaultTimeout,
	}
}

type PutRequest struct {
	Message string `json:"message"`
}

func (h *Handler) PutQueue(w http.ResponseWriter, r *http.Request) {
	var req PutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	queueName := r.PathValue("data")
	if err := h.queueSvc.Put(queueName, req.Message); err != nil {
		switch {
		case errors.Is(err, data.ErrQueueLimit):
			http.Error(w, "data limit exceeded", http.StatusServiceUnavailable)
		case errors.Is(err, data.ErrQueueFull):
			http.Error(w, "data is full", http.StatusServiceUnavailable)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetQueue(w http.ResponseWriter, r *http.Request) {
	timeout := h.defaultTimeout
	if t := r.URL.Query().Get("timeout"); t != "" {
		if td, err := time.ParseDuration(t + "s"); err == nil {
			timeout = td
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	queueName := r.PathValue("data")
	msg, err := h.queueSvc.Get(ctx, queueName)
	if err != nil {
		if errors.Is(err, data.ErrTimeout) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"message": msg.Data})
	if err != nil {
		return
	}
}
