package main

import (
	"encoding/json"
	"github.com/samuelmolero26/go-backend-course/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

	var payload struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//Decode json
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
		return
	}

	//mandatory validate
	if payload.Username == "" || payload.Email == "" || payload.Password == "" {
		http.Error(w, `{"error": "username, email, and password are required"}`, http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error": "could not hash password"}`, http.StatusInternalServerError)
		return
	}
	// creater username using store
	user := &storage.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: string(hashedPassword),
		RoleID:   1,
	}

	if err := app.store.Users().Create(r.Context(), user); err != nil {
		http.Error(w, `{"error": "could not create user"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

}
