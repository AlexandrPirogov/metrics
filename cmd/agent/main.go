package main

import (
	"log"
	"memtracker/internal/memtrack"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go func() {
		memtracker := memtrack.NewHTTPMemTracker()
		log.Printf("Started agent on %s, poll: %d, report: %d", memtracker.Host, memtracker.PollInterval, memtracker.ReportInterval)
		memtracker.ReadAndSend()
	}()

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

	<-cancelChan
	log.Printf("os.Interrupt-- shutting down...\n")

	defer os.Exit(0)
	close(cancelChan)
}
