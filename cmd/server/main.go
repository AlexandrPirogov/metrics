package main

import (
	"memtracker/internal/server/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/update/{mtype}/{mname}/{val}", handlers.UpdateHandler)
	r.Get("/value/{mtype}/{mname}", handlers.RetrieveMetric)
	r.Get("/", handlers.RetrieveMetrics)
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	server.ListenAndServe()
}
