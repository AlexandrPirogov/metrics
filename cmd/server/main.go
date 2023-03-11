package main

import (
	"memtracker/internal/handlers"
	"net/http"
)

func main() {

	http.HandleFunc("/update/", handlers.UpdateHandler)
	http.HandleFunc("/", http.NotFound)
	server := &http.Server{
		Addr: ":8080",
	}
	server.ListenAndServe()
}
