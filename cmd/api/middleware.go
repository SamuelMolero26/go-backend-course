package main

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/samuelmolero26/go-backend-course/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

type contextKey string

const userContextKey contextKey = "user"
const postContextKey contextKey = "post"

func (app *application) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}

		token := parts[1]
		authToken, err := app.store.AuthTokens().GetByToken(r.Context(), token)
		if err != nil {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}

		user, err := app.store.Users().GetByID(r.Context(), authToken.UserID)
		if err != nil {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	//take the token out of the header
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	//Erase from DB
	if err := app.store.AuthTokens().Delete(r.Context(), token); err != nil {
		http.Error(w, `{"error": "could not logout"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "logged out"})
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	//get post out of the context
	post := r.Context().Value(postContextKey).(*storage.Post)

	//Erase from DB
	if err := app.store.Posts().Delete(r.Context(), post.ID, post.UserID); err != nil {
		http.Error(w, `{"error": "could not delete post"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) requirePostOwnership(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Take Post ID out of the URL
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

		if err != nil {
			http.Error(w, `{"error": "Invalid ID, post does not exist"}`, http.StatusBadRequest)
			return
		}
		//search for post on DB
		post, err := app.store.Posts().GetByID(r.Context(), id)

		if err != nil {
			http.Error(w, `{"error": "post not found"}`, http.StatusNotFound)
			return
		}

		// take use out of the context
		user := r.Context().Value(userContextKey).(*storage.User)

		//compare against the owner
		if post.UserID != user.ID {
			http.Error(w, `{"error": "forbidden"}`, http.StatusForbidden)
			return
		}

		//put the post on the context (to search for it again)
		ctx := context.WithValue(r.Context(), postContextKey, post)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
