package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type api struct {
	user []User
}

var users = []User{}

func (a *api) getusersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) createUsersHandler(w http.ResponseWriter, r *http.Request) {
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

func (a *api) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", a.index)
	mux.HandleFunc("GET /users", a.getusersHandler)
	mux.HandleFunc("POST /users", a.createUsersHandler)

	return mux
}

func (a *api) index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/":
			w.Write([]byte("Index page"))
			return
		}

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error"))

	}

}

func (a *api) users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/users":
			w.Write([]byte("Users  page"))
			return
		}

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error"))

	}

}
