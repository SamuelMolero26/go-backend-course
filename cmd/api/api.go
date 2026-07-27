package main

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

type application struct {
	config config
	user   []User
}

type config struct {
	addr string
}

var users = []User{}

func (a *application) getusersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *application) createUsersHandler(w http.ResponseWriter, r *http.Request) {
	//decode request body to User Struct
	var payload User
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	u := User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		UserName:  payload.UserName,
	}

	if err = insertUser(u); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
	}

	w.WriteHeader(http.StatusCreated)

}

func insertUser(u User) error {

	if u.FirstName == "" {
		return errors.New("First name is required")

	}

	if u.LastName == "" {
		return errors.New("Last name is required")
	}

	//sotring validation
	for _, user := range users {
		if user.UserName == u.UserName {
			return errors.New("Username already exists")
		}
	}

	users = append(users, u)
	return nil

}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthChecker)
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
