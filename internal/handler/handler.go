package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/timurkash/queue2/internal/biz"
)

type Handler struct {
	queueSvc       *biz.Service
	defaultTimeout time.Duration
}

func New(queueSvc *biz.Service, defaultTimeout time.Duration) *Handler {
	return &Handler{
		queueSvc:       queueSvc,
		defaultTimeout: defaultTimeout,
	}
}

type PutRequest struct {
	Message string `json:"message"`
}

func getQueue(r *http.Request) string {
	queueName := r.PathValue("queue")
	if queueName == "" {
		queueName = r.URL.Path[7:]
	}
	return queueName
}

func (h *Handler) PutQueue(w http.ResponseWriter, r *http.Request) {
	var req PutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := h.queueSvc.Put(getQueue(r), req.Message); err != nil {
		switch {
		case errors.Is(err, biz.ErrQueueLimit):
			http.Error(w, "biz limit exceeded", http.StatusServiceUnavailable)
		case errors.Is(err, biz.ErrQueueFull):
			http.Error(w, "biz is full", http.StatusServiceUnavailable)
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

	msg, err := h.queueSvc.Get(ctx, getQueue(r))
	if err != nil {
		if errors.Is(err, biz.ErrTimeout) {
			http.Error(w, "not found", http.StatusNotFound)
		} else if errors.Is(err, biz.ErrQueueNotExists) {
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
