package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"taskforge/internal/storage"
	"taskforge/internal/task"
)

func TestRequestIDMiddleware_AddsHeader(t *testing.T) {
	svc := task.NewService(nil)
	store := storage.NewFileStorage("test_tasks.json")
	repo := newTestRepo(t)

	mux := NewServer(svc, store, nil, repo)
	handler := WithMiddleware(mux)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	requestID := rec.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Fatal("expected X-Request-ID header to be set")
	}
}
