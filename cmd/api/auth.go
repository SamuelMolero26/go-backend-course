package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/samuelmolero26/go-backend-course/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {

	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
		return
	}

	//sanity check - check for empty
	if payload.Email == "" || payload.Password == "" {
		http.Error(w, `{"error": "email and password are required"}`, http.StatusBadRequest)
		return
	}

	user, err := app.store.Users().GetByEmail(r.Context(), payload.Email)
	if err != nil {
		http.Error(w, `{"error": "Invalid Credentials"}`, http.StatusUnauthorized)
		return
	}

	//compare agains bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		http.Error(w, `{"error": "Invalid Credentials"}`, http.StatusUnauthorized)
		return
	}

	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)

	if err != nil {
		http.Error(w, `{"error": "could not generate token"}`, http.StatusInternalServerError)
		return
	}

	token := hex.EncodeToString(tokenBytes)

	authToken := &storage.AuthToken{

		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // expira en 24h
	}

	if err := app.store.AuthTokens().Create(r.Context(), authToken); err != nil {
		http.Error(w, `{"error": "could not save token"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})

}
