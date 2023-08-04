package main

import (
	"context"
	"log"
	cfg "memtracker/internal/config/server"
	"memtracker/internal/server"
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
	cfg.Exec()
	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildDate)
	log.Printf("Build commit: %s", buildCommit)

	_, cancel := context.WithCancel(context.Background())

	server := server.New()

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("%v", err)
		}

	}()

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	//log.Printf("started server on %s\n", server.Conf.Address)
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
