package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/samuelmolero26/go-backend-course/internal/storage"
)

type application struct {
	config config
	store  *storage.Store
}

type config struct {
	addr string
	db   string
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

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
			r.Get("/users/{id}/posts", app.getUserPostsHandler)

			//ownership group
			r.Group(func(r chi.Router) {
				r.Use(app.requirePostOwnership)
				r.Delete("/posts/{id}", app.deletePostHandler)
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
