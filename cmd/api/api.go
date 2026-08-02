package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/samuelmolero26/go-backend-course/internal/storage"
	"github.com/samuelmolero26/go-backend-course/internal/ratelimit"
	"github.com/minio/minio-go/v7"
)

type application struct {
	config 		config
	store  		*storage.Store
	rateLimiter *ratelimit.Limiter
	minioClient	*minio.Client
}

type config struct {
	addr string
	db   string
	minioBucket string
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(app.ratelimiterMiddleware)

	r.Route("/v1", func(r chi.Router) {
		// Public
		r.Get("/health", app.healthChecker)
		r.Post("/users", app.createUserHandler)
		r.Post("/login", app.loginHandler)

		// Protected
		r.Group(func(r chi.Router) {
			r.Use(app.authMiddleware)
			r.Post("/logout", app.logoutHandler)
			r.Post("/posts", app.createPostHandler)
			r.Get("/posts/{id}", app.getPostHandler)
			r.Get("/users/{id}", app.getUserHandler)
			r.Get("/users/{id}/posts", app.getUserPostsHandler)
			r.Get("/posts", app.getFeedHandler)
			r.Get("/users/{id}/followers", app.getFollowersHandler)
			r.Get("/users/{id}/following", app.getFollowingHandler)
			r.Post("/users/{id}/follow", app.followUserHandler)
			r.Delete("/users/{id}/follow", app.unfollowUserHandler)
			r.Post("/posts/{id}/comments", app.createCommentHandler)
			r.Get("/posts/{id}/comments", app.getCommentsHandler)

			//ownership group
			r.Group(func(r chi.Router) {
				r.Use(app.requirePostOwnership)
				r.Delete("/posts/{id}", app.deletePostHandler)
				r.Patch("/posts/{id}", app.updatePostHandler)
			})

			//comments
			r.Group(func(r chi.Router) {
				r.Use(app.requireCommentOwnership)
				r.Delete("/comments/{id}", app.deleteCommentHandler)
			})
		})
	})

	return r
}

func (app *application) serve(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}
