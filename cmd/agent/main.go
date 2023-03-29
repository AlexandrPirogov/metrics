package main

import (
	"log"
	"memtracker/internal/memtrack"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Better use .env
var host string = "localhost"
var port string = ":8080"

func main() {
	go func() {
		var client = http.Client{Timeout: time.Second / 2}
		memtracker := memtrack.NewHTTPMemTracker(client, host+port)
		memtracker.ReadAndSend()
	}()

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

	<-cancelChan
	log.Printf("os.Interrupt-- shutting down...\n")
	//We don't need to clear http.Client and stop it properly like server
	go func() {
		<-cancelChan
		log.Fatalf("os.Kill -- terminating...\n")
	}()
	defer os.Exit(0)
	close(cancelChan)
}
