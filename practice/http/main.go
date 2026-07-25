package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	app := &api{}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      app.RegisterRoutes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Println("Server listening on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
