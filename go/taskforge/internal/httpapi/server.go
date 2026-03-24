package httpapi

import (
	"net/http"

	"taskforge/internal/storage"
	"taskforge/internal/task"
	"taskforge/internal/taskapp"
	"taskforge/internal/worker"
)

func NewServer(svc *task.Service, store *storage.FileStorage, processor *worker.Processor) *http.ServeMux {
	app := taskapp.New(svc, store, processor)

	h := &Handlers{
		svc:   svc,
		app:   app,
		store: store,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handleHealth)

	// Collection routes
	mux.HandleFunc("GET /tasks", h.handleListTasks)
	mux.HandleFunc("GET /tasks/{id}", h.handleGetTaskByID)
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
