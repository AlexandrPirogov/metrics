package main

import (
	"context"
	"log"
	config "memtracker/internal/config/server"
	"memtracker/internal/server"
	"memtracker/internal/server/handlers/api"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildDate)
	log.Printf("Build commit: %s", buildCommit)

	config.Exec()
	ctx, cancel := context.WithCancel(context.Background())

	handler := api.NewHandler()
	server := server.NewMetricServer(handler, ctx)
	handler.DB.Start()

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("%v", err)
		}

	}()

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	log.Printf("started server on %s\n", server.Conf.Address)
	sig := <-cancelChan
	log.Printf("Got signal %v\n", sig)

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Printf("Shutdown error %v\n", err)
		defer os.Exit(1)
		return
	}

	log.Printf("Server shutdowned\n")

	cancel()
	defer os.Exit(0)
	close(cancelChan)
}
