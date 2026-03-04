package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
}

func TestPOSTTasks_InvalidJSON(t *testing.T) {
	mux, _ := newTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.Newreader([]byte(`{bad json`)))
}
