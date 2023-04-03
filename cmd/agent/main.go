package main

import (
	"log"
	"memtracker/internal/config/agent"
	"memtracker/internal/memtrack"
	"os"
	"os/signal"
	"syscall"
)

// Better use .env
var host string = "localhost"
var port string = ":8080"

func main() {
	agent.Exec()
	go func() {
		memtracker := memtrack.NewHTTPMemTracker(host + port)
		log.Printf("Started agent on %s, poll: %d, report: %d", memtracker.Host, memtracker.PollInterval, memtracker.ReportInterval)
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
