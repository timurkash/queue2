package app

import (
	"fmt"
	"net/http"

	"github.com/timurkash/queue2/internal/data"
	"github.com/timurkash/queue2/internal/handler"
)

type Application struct {
	server *http.Server
}

func New(cfg Config) *Application {
	qSvc := data.New(cfg.MaxQueues, cfg.QueueCapacity)
	theHandler := handler.New(qSvc, cfg.DefaultTimeout)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /queue/{queue}", theHandler.PutQueue)
	mux.HandleFunc("GET /queue/{queue}", theHandler.GetQueue)

	return &Application{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: mux,
		},
	}
}

func (a *Application) Run() error {
	return a.server.ListenAndServe()
}
