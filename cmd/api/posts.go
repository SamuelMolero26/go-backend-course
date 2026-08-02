package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/minio/minio-go/v7"
	"github.com/samuelmolero26/go-backend-course/internal/storage"
)

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)

	//multipart
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, `{"error": "invalid multipart form"}`, http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	file, _, err := r.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			file = nil
		} else {
			http.Error(w, `{"error": "could not get image"}`, http.StatusBadRequest)
			return
		}
	}

	user := r.Context().Value(userContextKey).(*storage.User)

	post := &storage.Post{
		UserID:  user.ID,
		Content: content,
	}

	// si viene imagen: validar, subir a MinIO y guardar la URL
	if file != nil {
		// 1. detectar el tipo real (no confiar en el header)
		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		contentType := http.DetectContentType(buf[:n])

		// 2. rebobinar: el Read avanzó la posición
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			http.Error(w, `{"error": "could not rewind image"}`, http.StatusInternalServerError)
			return
		}

		// 3. nombre único generado por el servidor (nunca el del cliente)
		keyBytes := make([]byte, 16)
		if _, err := rand.Read(keyBytes); err != nil {
			http.Error(w, `{"error": "could not generate object name"}`, http.StatusInternalServerError)
			return
		}
		objectName := fmt.Sprintf("posts/%s", hex.EncodeToString(keyBytes))

		// 4. subir a MinIO
		_, err = app.minioClient.PutObject(
			r.Context(),
			app.config.minioBucket,
			objectName,
			file,
			-1,
			minio.PutObjectOptions{ContentType: contentType},
		)
		if err != nil {
			http.Error(w, `{"error": "could not upload image"}`, http.StatusInternalServerError)
			return
		}

		// 5. guardar la URL pública del objeto
		url := fmt.Sprintf("http://localhost:9000/%s/%s", app.config.minioBucket, objectName)
		post.ImageURL = &url
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


func (app *application) getFeedHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := app.store.Posts().GetFeed(r.Context())
	if err != nil {
		http.Error(w, `{"error": "could not get feed"}`, http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)

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
