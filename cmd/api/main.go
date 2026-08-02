package main

import (
	"context"
	"log"

	"github.com/samuelmolero26/go-backend-course/internal/env"
	"github.com/samuelmolero26/go-backend-course/internal/storage"
	"github.com/samuelmolero26/go-backend-course/internal/ratelimit"

)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db:   env.GetString("DB_DSN", ""),
		minioBucket: env.GetString("MINIO_BUCKET", "photos"),
	}

	ctx := context.Background()

	store, err := storage.NewStorage(ctx, cfg.db)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	minioClient, err :=newMinioClient()
	if err != nil {
		log.Fatal(err)
	}
	
	if err := ensureBucket(ctx, minioClient, cfg.minioBucket); err != nil{
		log.Fatal(err)
	}

	if err := setPublicReadPolicy(ctx, minioClient, cfg.minioBucket); err != nil {
		log.Fatal(err)
	}

	app := &application{
		config: cfg,
		store:  store,
		rateLimiter: ratelimit.New(10,2),
		minioClient: minioClient,
	}

	mux := app.mount()

	log.Fatal(app.serve(mux))
}
