package main

import (
	"log"
	"net/http"

	"github.com/thesouldev/goboxd/server"
)

func main() {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: server.NewMux(),
	}

	log.Println("Starting server on port 8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on port 8080: %v", err)
	}
}
