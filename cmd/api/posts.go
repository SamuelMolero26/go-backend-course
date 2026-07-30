package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/samuelmolero26/go-backend-course/internal/storage"
)

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
		return
	}

	user := r.Context().Value(userContextKey).(*storage.User)

	post := &storage.Post{
		UserID:  user.ID,
		Content: payload.Content,
	}

	if err := app.store.Posts().Create(r.Context(), post); err != nil {
		http.Error(w, `{"error": "could not create post"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid post id"}`, http.StatusBadRequest)
		return
	}

	post, err := app.store.Posts().GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error": "post not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (app *application) getUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid user id"}`, http.StatusBadRequest)
		return
	}

	posts, err := app.store.Posts().GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "could not get posts"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := r.Context().Value(postContextKey).(*storage.Post)

	var payload struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
		return
	}

	post.Content = payload.Content

	if err := app.store.Posts().Update(r.Context(), post); err != nil {
		http.Error(w, `{"error": "could not update post"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid user id"}`, http.StatusBadRequest)
		return
	}

	user, err := app.store.Users().GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
