package main

import (
	"context"
	"log"
	"memtracker/internal/server"
	"memtracker/internal/server/db"
	"memtracker/internal/server/handlers/api"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	handler := &api.DefaultHandler{DB: &db.DB{Storage: db.MemStoageDB()}}
	server := server.NewMetricServer(":8080", handler, ctx)
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
