package server

import (
	"context"
	"memtracker/internal/server/handlers"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewMetricServer(addr string, h handlers.MetricsHandler, ctx context.Context) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Post("/update/{mtype}/{mname}/{val}", h.UpdateHandler)
	r.Get("/value/{mtype}/{mname}", h.RetrieveMetric)
	r.Get("/", h.RetrieveMetrics)
	return &http.Server{
		Addr:        ":8080",
		Handler:     r,
		BaseContext: func(listener net.Listener) context.Context { return ctx },
	}
}
