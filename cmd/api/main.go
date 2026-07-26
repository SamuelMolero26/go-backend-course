package main

import (
	"github.com/samuelmolero26/go-backend-course/internal/env"
	"log"
)

func main() {

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount()

	log.Fatal(app.serve(mux))

}
