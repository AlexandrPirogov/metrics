package server

import (
	"context"
	"memtracker/internal/config/server"
	"memtracker/internal/server/middlewares"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type MetricsHandler interface {
	RetrieveMetrics(w http.ResponseWriter, r *http.Request)
	RetrieveMetric(w http.ResponseWriter, r *http.Request)
	UpdateHandler(w http.ResponseWriter, r *http.Request)
	PingHandler(w http.ResponseWriter, r *http.Request)

	RetrieveMetricJSON(w http.ResponseWriter, r *http.Request)
	UpdateHandlerJSON(w http.ResponseWriter, r *http.Request)
	UpdatesHandlerJSON(w http.ResponseWriter, r *http.Request)
}

func NewMetricServer(h MetricsHandler, ctx context.Context) *http.Server {
	cfg := server.ServerCfg
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Group(func(r chi.Router) {
		r.Use(middlewares.GZIPer)
		r.Use(middleware.AllowContentType("text/plain"))
		r.Post("/update/{mtype}/{mname}/{val}", h.UpdateHandler)
		r.Get("/value/{mtype}/{mname}", h.RetrieveMetric)
		r.Get("/", h.RetrieveMetrics)
	})

	r.Group(func(r chi.Router) {
		r.Use(middlewares.GZIPer)
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/update/", h.UpdateHandlerJSON)
		r.Post("/value/", h.RetrieveMetricJSON)
		r.Post("/updates/", h.UpdatesHandlerJSON)
	})
	r.Get("/ping", h.PingHandler)
	r.Mount("/debug", middleware.Profiler())
	return &http.Server{
		Addr:        cfg.Address,
		Handler:     r,
		BaseContext: func(listener net.Listener) context.Context { return ctx },
	}
}
