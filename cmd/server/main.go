package main

import (
	"context"
	"log"
	"memtracker/internal/server/db"
	"memtracker/internal/server/handlers"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	handler := handlers.Handler{DB: &db.DB{Storage: db.MemStoageDB()}}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/update/{mtype}/{mname}/{val}", handler.UpdateHandler)
	r.Get("/value/{mtype}/{mname}", handler.RetrieveMetric)
	r.Get("/", handler.RetrieveMetrics)
	server := &http.Server{
		Addr:        ":8080",
		Handler:     r,
		BaseContext: func(listener net.Listener) context.Context { return ctx },
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("%v", err)
		}

	}()

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

	<-cancelChan
	log.Printf("os.Interrupt-- shutting down...\n")
	go func() {
		<-cancelChan
		log.Fatalf("os.Kill -- terminating...\n")
	}()

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Printf("Shutdown error %v\n", err)
		defer os.Exit(1)
		return
	} else {
		log.Printf("Server shutdowned\n")
	}
	cancel()
	defer os.Exit(0)
	close(cancelChan)
}
