package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"taskforge/internal/httpapi"
	"taskforge/internal/storage"
	"taskforge/internal/task"
)

func main() {
	store := storage.NewFileStorage(filepath.Join(".", "tasks.json"))

	existing, err := store.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading tasks:", err)
		os.Exit(1)
	}

	svc := task.NewService(existing)

	mux := httpapi.NewServer(svc, store)
	handler := httpapi.WithMiddleware(mux)

	addr := ":8080"
	fmt.Println("TaskForge API listening on", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		fmt.Fprintln(os.Stderr, "Server error:", err)
		os.Exit(1)
	}
}
