package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/samuelmolero26/go-backend-course/internal/storage"
)

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid user id"}`, http.StatusBadRequest)
		return
	}

	currentUser := r.Context().Value(userContextKey).(*storage.User)

	if err := app.store.Followers().Follow(r.Context(), userID, currentUser.ID); err != nil {
		log.Printf("ERROR follow user: %v", err)
		http.Error(w, `{"error": "could not follow user"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "followed"})
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid user id"}`, http.StatusBadRequest)
		return
	}

	currentUser := r.Context().Value(userContextKey).(*storage.User)

	if err := app.store.Followers().Unfollow(r.Context(), userID, currentUser.ID); err != nil {
		log.Printf("ERROR unfollow user: %v", err)
		http.Error(w, `{"error": "could not unfollow user"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "unfollowed"})
}

func (app *application) getFollowersHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid user id"}`, http.StatusBadRequest)
		return
	}

	followers, err := app.store.Followers().GetFollowers(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "could not get followers"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

func (app *application) getFollowingHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid user id"}`, http.StatusBadRequest)
		return
	}

	following, err := app.store.Followers().GetFollowing(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "could not get following"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(following)
}
