package main

import (
	"log"
	"memtracker/internal/server/db"
	"memtracker/internal/server/handlers"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cancelChan := make(chan os.Signal, 1)

	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		handler := handlers.Handler{DB: db.DB{Storage: db.MemStoageDB()}}
		r := chi.NewRouter()
		r.Use(middleware.Logger)
		r.Post("/update/{mtype}/{mname}/{val}", handler.UpdateHandler)
		r.Get("/value/{mtype}/{mname}", handler.RetrieveMetric)
		r.Get("/", handler.RetrieveMetrics)
		server := &http.Server{
			Addr:    ":8080",
			Handler: r,
		}
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("%v", err)
		}
	}()
	sig := <-cancelChan
	log.Printf("Caught signal %v", sig)
	close(cancelChan)
}
