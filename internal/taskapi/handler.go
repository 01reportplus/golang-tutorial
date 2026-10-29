package taskapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
)

const maxBody = 1024 * 1024

type errorBody struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		http.Error(w, "response encoding failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	// A client disconnect after the headers cannot be converted into another response.
	_, _ = w.Write(append(data, '\n'))
}

func decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeJSON(w, http.StatusUnsupportedMediaType, errorBody{"Content-Type must be application/json"})
		return false
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBody)
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(dst)
	if err == nil {
		var extra any
		err = dec.Decode(&extra)
		if err == io.EOF {
			return true
		}
	}
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		writeJSON(w, http.StatusRequestEntityTooLarge, errorBody{"request body exceeds 1 MiB"})
	} else {
		writeJSON(w, http.StatusBadRequest, errorBody{"body must contain exactly one JSON object with known fields"})
	}
	return false
}

func taskID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, errorBody{"id must be a positive int64"})
		return 0, false
	}
	return id, true
}

func storeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		writeJSON(w, http.StatusNotFound, errorBody{ErrNotFound.Error()})
	case errors.Is(err, ErrTitle):
		writeJSON(w, http.StatusBadRequest, errorBody{ErrTitle.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, errorBody{"internal error"})
	}
}

// NewHandler builds routes on a private ServeMux. Store must be initialized by NewStore.
func NewHandler(store *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /api/tasks", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, store.List())
	})
	mux.HandleFunc("POST /api/tasks", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Title string `json:"title"`
		}
		if !decode(w, r, &req) {
			return
		}
		task, err := store.Create(req.Title)
		if err != nil {
			storeError(w, err)
			return
		}
		w.Header().Set("Location", fmt.Sprintf("/api/tasks/%d", task.ID))
		writeJSON(w, http.StatusCreated, task)
	})
	mux.HandleFunc("GET /api/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := taskID(w, r)
		if !ok {
			return
		}
		task, err := store.Get(id)
		if err != nil {
			storeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, task)
	})
	mux.HandleFunc("PATCH /api/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := taskID(w, r)
		if !ok {
			return
		}
		var req struct {
			Done *bool `json:"done"`
		}
		if !decode(w, r, &req) {
			return
		}
		if req.Done == nil {
			writeJSON(w, http.StatusBadRequest, errorBody{"done must be a boolean"})
			return
		}
		task, err := store.SetDone(id, *req.Done)
		if err != nil {
			storeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, task)
	})
	mux.HandleFunc("DELETE /api/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := taskID(w, r)
		if !ok {
			return
		}
		if err := store.Delete(id); err != nil {
			storeError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	return mux
}
