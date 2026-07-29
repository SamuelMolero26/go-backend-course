package main

import (
	"context"
	"log"

	"github.com/samuelmolero26/go-backend-course/internal/env"
	"github.com/samuelmolero26/go-backend-course/internal/storage"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db:   env.GetString("DB_DSN", ""),
	}

	ctx := context.Background()

	store, err := storage.NewStorage(ctx, cfg.db)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.serve(mux))
}
