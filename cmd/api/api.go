package main

import (
	"log"
	"net/http"
)

type application struct {
	config config
}

type config struct {
	addr string
}

func (app *application) serve() error {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    app.config.addr,
		Handler: mux,
	}

	log.Printf("Server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}
