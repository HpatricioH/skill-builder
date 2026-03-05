package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"taskforge/internal/storage"
	"taskforge/internal/task"
)

func newTestServer(t *testing.T) (*http.ServeMux, string) {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	store := storage.NewFileStorage(path)

	// Start with empty tasks from file
	existing, err := store.Load()
	if err != nil {
		t.Fatalf("store.Load() error =%v", err)
	}

	svc := task.NewService(existing)
	mux := NewServer(svc, store)

	return mux, path
}

func decodeJSON(t *testing.T, body *bytes.Buffer, v any) {
	t.Helper()
	if err := json.NewDecoder(body).Decode(v); err != nil {
		t.Fatalf("decode json error =%v", err)
	}
}

func TestTasks_Flow_Create_List_Done_Delete(t *testing.T) {
	mux, filePath := newTestServer(t)

	// 1) POST /tasks
	postBody := []byte(`{"title": "Buy milk"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(postBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /tasks status = %d, want %d, body=%s", rec.Code, http.StatusCreated, rec.Body)
	}

	var created task.Task
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("decode created task error = %v", err)
	}
	if created.ID != 1 {
		t.Fatalf("created.ID = %d, want 1", created.ID)
	}
	if created.Title != "Buy milk" {
		t.Fatalf("created.Title = %q, want %q", created.Title, "Buy milk")
	}
	if created.Completed {
		t.Fatalf("created.Completed = %v, want false", created.Completed)
	}

	// Verify persistence file exists and is non-empty
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("tasks.json stat error = %v", err)
	}
	if info.Size() == 0 {
		t.Fatalf("tasks.json size = 0, want > 0")
	}

	// 2) GET /tasks
	req = httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /tasks status =%d, want %d", rec.Code, http.StatusOK)
	}

	var list []task.Task
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatalf("decode list error = %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("len(list) = %d, want 1", len(list))
	}
	if list[0].ID != 1 {
		t.Fatalf("list[0].ID = %d, want 1", list[0].ID)
	}

	// 3) PATCH /tasks/1/done
	req = httptest.NewRequest(http.MethodPatch, "/tasks/1/done", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH /tasks/1/done status =%d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	// GET again, verify done
	req = httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /tasks (after done) status = %d, want %d", rec.Code, http.StatusOK)
	}

	list = nil
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatalf("decode list after done error = %v", err)
	}
	if len(list) != 1 || !list[0].Completed {
		t.Fatalf("after done, tasks=%v; want one tasks with Completed= true", list)
	}

	// 4) DELETE /tasks/1
	req = httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("DELETE /tasks/1 status =%d, want %d, body %s", rec.Code, http.StatusNoContent, rec.Body.String())
	}

	// GET again, should be empty
	req = httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /tasks (after done) status = %d, want %d", rec.Code, http.StatusOK)
	}

	list = nil
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatalf("decode list after done error = %v", err)
	}
	if len(list) != 1 || !list[0].Completed {
		t.Fatalf("after done, tasks=%v; want one tasks with Completed= true", list)
	}
}

func TestPOSTTasks_InvalidJSON(t *testing.T) {
	mux, _ := newTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader([]byte(`{bad json`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("POST /task invalid json status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}
