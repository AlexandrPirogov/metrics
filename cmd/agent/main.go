package main

import (
	"memtracker/internal/memtrack"
	"net/http"
	"time"
)

// Better use .env
var host string = "localhost"
var port string = ":8080"

func main() {
	var client http.Client = http.Client{Timeout: time.Second / 2}
	memtracker := memtrack.NewHttpMemTracker(client, host+port)
	memtracker.ReadAndSend(time.Second*2, time.Second*10)
}
