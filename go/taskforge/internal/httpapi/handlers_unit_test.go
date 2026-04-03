package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"taskforge/internal/db"
	"taskforge/internal/task"
	"taskforge/internal/taskapp"
	"taskforge/internal/taskrepo"
)

func newHandlers(t *testing.T) (*Handlers, string) {
	t.Helper()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	t.Cleanup(func() {
		_ = database.Close()
	})

	if err := db.InitSchema(database); err != nil {
		t.Fatalf("init schema: %v", err)
	}

	repo := taskrepo.New(database)
	app := taskapp.New(repo, nil)

	return &Handlers{
		app: app,
	}, dbPath
}

func newTestRepo(t *testing.T) *taskrepo.Repository {
	t.Helper()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	t.Cleanup(func() {
		_ = database.Close()
	})

	if err := db.InitSchema(database); err != nil {
		t.Fatalf("init schema: %v", err)
	}

	return taskrepo.New(database)
}

func TestHandlers_Flow_Create_List_Done_Delete(t *testing.T) {
	h, _ := newHandlers(t)

	// POST /tasks
	{
		body := []byte(`{"title":"Buy milk"}`)
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.handleCreateTask(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("POST status=%d want=%d body=%s", rec.Code, http.StatusCreated, rec.Body.String())
		}
	}

	// GET /tasks -> expect 1 task
	{
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		rec := httptest.NewRecorder()

		h.handleListTasks(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET status=%d want=%d body=%s", rec.Code, http.StatusOK, rec.Body.String())
		}

		var list []task.Task
		if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
			t.Fatalf("decode GET list error=%v", err)
		}
		if len(list) != 1 {
			t.Fatalf("len(list)=%d want=1 list=%v", len(list), list)
		}
		if list[0].Completed {
			t.Fatalf("Completed=%v want=false", list[0].Completed)
		}
	}

	// PATCH /tasks/{id}/done  (we inject the path param directly)
	{
		req := httptest.NewRequest(http.MethodPatch, "/tasks/1/done", nil)
		req.SetPathValue("id", "1")
		rec := httptest.NewRecorder()

		h.handleMarkDone(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("PATCH status=%d want=%d body=%s", rec.Code, http.StatusOK, rec.Body.String())
		}
	}

	// GET /tasks -> expect same task now Completed=true
	{
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		rec := httptest.NewRecorder()

		h.handleListTasks(rec, req)

		var list []task.Task
		if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
			t.Fatalf("decode GET list error=%v body=%s", err, rec.Body.String())
		}
		if len(list) != 1 || !list[0].Completed {
			t.Fatalf("after done, tasks=%v; want one task with Completed=true", list)
		}
	}

	// DELETE /tasks/{id}
	{
		req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
		req.SetPathValue("id", "1")
		rec := httptest.NewRecorder()

		h.handleDeleteTask(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("DELETE status=%d want=%d body=%s", rec.Code, http.StatusOK, rec.Body.String())
		}
	}

	// GET /tasks -> expect empty
	{
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		rec := httptest.NewRecorder()

		h.handleListTasks(rec, req)

		var list []task.Task
		if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
			t.Fatalf("decode GET list error=%v", err)
		}
		if len(list) != 0 {
			t.Fatalf("len(list)=%d want=0 list=%v", len(list), list)
		}
	}
}
