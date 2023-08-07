// Package server provides an API for the http.server instance
package http

import (
	"context"
	"memtracker/internal/config/server"
	"memtracker/internal/server/http/api"
	"memtracker/internal/server/http/middlewares"

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

type metricServer struct {
	http *http.Server
	Conf server.ServerConfig
}

func (m *metricServer) ListenAndServe() error {
	return m.Conf.Run(m.http)
}

func (m *metricServer) Shutdown(ctx context.Context) error {
	return m.http.Shutdown(ctx)
}

func BuildHTPP() *metricServer {
	ctx := context.Background()
	cfg := server.ServerCfg
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	if cfg.Subnet != "" {
		r.Use(middlewares.SubnetValidate)
	}

	h := api.NewHandler()
	group(r, h)

	h.DB.Start()
	return &metricServer{
		&http.Server{
			Addr:        cfg.Address,
			Handler:     r,
			BaseContext: func(listener net.Listener) context.Context { return ctx },
		},
		*cfg,
	}
}

func group(r *chi.Mux, h MetricsHandler) {
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
		r.Post("/updates/", h.UpdatesHandlerJSON)
		r.Post("/update/", h.UpdateHandlerJSON)
		r.Post("/value/", h.RetrieveMetricJSON)
	})

	r.Get("/ping", h.PingHandler)
	r.Mount("/debug", middleware.Profiler())
}
