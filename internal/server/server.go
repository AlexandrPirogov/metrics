package server

import (
	"context"
	"memtracker/internal/server/handlers/api"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewMetricServer(addr string, h api.MetricsHandler, ctx context.Context) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("text/plain"))
		r.Post("/update/{mtype}/{mname}/{val}", h.UpdateHandler)
		r.Get("/value/{mtype}/{mname}", h.RetrieveMetric)
		r.Get("/", h.RetrieveMetrics)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/update/", h.UpdateHandlerJSON)
		r.Post("/value/", h.RetrieveMetricJSON)
	})

	return &http.Server{
		Addr:        ":8080",
		Handler:     r,
		BaseContext: func(listener net.Listener) context.Context { return ctx },
	}
}
