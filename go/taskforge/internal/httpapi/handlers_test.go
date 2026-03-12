package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"taskforge/internal/storage"
	"taskforge/internal/task"
)

func TestGETTasks_Empty(t *testing.T) {
	// Arrange
	svc := task.NewService(nil)
	store := storage.NewFileStorage("test_tasks.json") // not used by GET
	mux := NewServer(svc, store, nil)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("status =%d, want %d", rec.Code, http.StatusOK)
	}
}
