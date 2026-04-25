package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"taskforge/internal/storage"
	"taskforge/internal/task"
	"taskforge/internal/taskapp"
)

type Handlers struct {
	mu    sync.Mutex
	svc   *task.Service
	app   *taskapp.App
	store *storage.FileStorage
}

type errorResponse struct {
	Error string `json: "error"`
}

type messageResponse struct {
	Message string `json: "message"`
}

func (h *Handlers) save(w http.ResponseWriter) bool {
	if err := h.store.Save(h.svc.ListTasks()); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to save tasks",
		})
		return false
	}
	return true
}

func (h *Handlers) handleListTasks(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	limit := 20
	offset := 0

	if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil || parsed <= 0 {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid limit"})
			return
		}
		limit = parsed
	}

	if rawOffset := r.URL.Query().Get("offset"); rawOffset != "" {
		parsed, err := strconv.Atoi(rawOffset)
		if err != nil || parsed < 0 {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid offset"})
			return
		}
		offset = parsed
	}

	if limit > 100 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "limit cannot be greated thatn 100"})
		return
	}

	tasks, err := h.app.ListTasksPaginated(r.Context(), limit, offset)
	if err != nil {
		fmt.Println("list paginated error", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to list tasks"})
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handlers) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	var body struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}

	t, err := h.app.CreateTask(r.Context(), body.Title)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, t)
}

func (h *Handlers) handleMarkDone(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	id, ok := parseIDFromPath(r.PathValue("id"))
	if !ok {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}

	if _, err := h.app.MarkDone(r.Context(), id); err != nil {
		// simple mapping for now
		code := http.StatusNotFound
		if strings.Contains(err.Error(), "already completed") {
			code = http.StatusConflict
		}
		writeJSON(w, code, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, messageResponse{
		Message: fmt.Sprintf("task %d maked done", id),
	})
}

func (h *Handlers) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) != 2 || parts[0] != "tasks" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid path"})
		return
	}

	id, ok := parseIDFromPath(parts[1])
	if !ok {
		if err := h.store.Save(h.svc.ListTasks()); err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to save tasks"})
			return
		}

		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}

	if err := h.app.DeleteTask(r.Context(), id); err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, messageResponse{
		Message: fmt.Sprintf("task %d deleted", id),
	})
}

func (h *Handlers) handleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	id, ok := parseIDFromPath(r.PathValue("id"))
	if !ok {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}

	t, err := h.app.GetTaskByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, t)
}

func (h *Handlers) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) != 2 || parts[0] != "tasks" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid path"})
		return
	}

	id, ok := parseIDFromPath(parts[1])
	if !ok {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
	}

	var body struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}

	body.Title = strings.TrimSpace(body.Title)
	if body.Title == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "title cannot be empty"})
		return
	}

	updated, err := h.app.UpdateTaskTitle(r.Context(), id, body.Title)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func parseIDFromPath(raw string) (int, bool) {
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
