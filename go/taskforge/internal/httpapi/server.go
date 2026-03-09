package httpapi

import (
	"net/http"

	"taskforge/internal/storage"
	"taskforge/internal/task"
)

func NewServer(svc *task.Service, store *storage.FileStorage) *http.ServeMux {
	h := &Handlers{svc: svc, store: store}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handleHealth)

	// Collection routes
	mux.HandleFunc("GET /tasks", h.handleListTasks)
	mux.HandleFunc("POST /tasks", h.handleCreateTask)

	// Item routes
	mux.HandleFunc("PATCH /tasks/{id}/done", h.handleMarkDone)
	mux.HandleFunc("DELETE /tasks/{id}", h.handleDeleteTask)

	return mux
}

func WithMiddleware(handler http.Handler) http.Handler {
	handler = LoggingMiddleware(handler)
	handler = RequestIDMiddleware(handler)

	return handler
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
