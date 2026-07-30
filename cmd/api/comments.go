package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/samuelmolero26/go-backend-course/internal/storage"

)

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid post id"}`, http.StatusBadRequest)
		return
	}

	var payload  struct {
		Content string `json:"content"`
	}

	if err:= json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w,`{"error": "invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if payload.Content == "" {
		http.Error(w, `{"error": "content is required"}`, http.StatusBadRequest)
		return
	}

	user := r.Context().Value(userContextKey).(*storage.User)

	comment := &storage.Comment{
		PostID:  postID,
		UserID:  user.ID,
		Content: payload.Content,
	}
	
	if err := app.store.Comments().Create(r.Context(), comment); err != nil {
		http.Error(w, `{"error": "could not create comment"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)

}

func (app *application) getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid post id"}`, http.StatusBadRequest)
		return
	}

	comments, err := app.store.Comments().GetByPostID(r.Context(), postID)
	if err != nil {
		http.Error(w, `{"error": "could not get comments"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	comment := r.Context().Value(commentContextKey).(*storage.Comment)

	if err := app.store.Comments().Delete(r.Context(), comment.ID, comment.UserID); err != nil {
		http.Error(w, `{"error": "could not delete comment"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}