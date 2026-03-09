package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

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

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		fmt.Println("TaskForge API listening on", server.Addr)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintln(os.Stderr, "Server error:", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	fmt.Println("\nShutdown signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "Graceful shutdown failed:", err)
		os.Exit(1)
	}

	fmt.Println("Server stopped cleanly")
}
