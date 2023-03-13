package main

import (
	"memtracker/internal/server/handlers"
	"net/http"
)

func main() {

	http.HandleFunc("/update/", handlers.UpdateHandler)
	http.HandleFunc("/", handlers.NotImplementedHandler)
	server := &http.Server{
		Addr: ":8080",
	}
	server.ListenAndServe()
}
