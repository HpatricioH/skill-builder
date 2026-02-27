package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"taskforge/internal/storage"
	"taskforge/internal/task"
)

type Handlers struct {
	svc   *task.Service
	store *storage.FileStorage
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
	writeJSON(w, http.StatusOK, h.svc.ListTasks())
}

func (h *Handlers) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return

	}

	t, err := h.svc.AddTask(body.Title)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if !h.save(w) {
		return
	}

	writeJSON(w, http.StatusCreated, t)
}

func (h *Handlers) handleMarkDone(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDFromPath(r.PathValue("id"))
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	if err := h.svc.MarkDone(id); err != nil {
		// simple mapping for now
		code := http.StatusNotFound
		if strings.Contains(err.Error(), "already completed") {
			code = http.StatusConflict
		}
		writeJSON(w, code, map[string]string{"error": err.Error()})
		return
	}

	if !h.save(w) {
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": fmt.Sprintf("task %d maked done", id),
	})
}

func (h *Handlers) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDFromPath("id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	if err := h.svc.DeleteTask(id); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	if !h.save(w) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
