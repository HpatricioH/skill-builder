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

	"taskforge/internal/db"
	"taskforge/internal/httpapi"
	"taskforge/internal/storage"
	"taskforge/internal/task"
	"taskforge/internal/taskrepo"
	"taskforge/internal/worker"
)

func main() {
	store := storage.NewFileStorage(filepath.Join(".", "tasks.json"))

	existing, err := store.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading tasks:", err)
		os.Exit(1)
	}

	svc := task.NewService(existing)

	database, err := db.Open("taskforge.db")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening database:", err)
		os.Exit(1)
	}
	defer database.Close()

	if err := db.InitSchema(database); err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing schema:", err)
		os.Exit(1)
	}

	repo := taskrepo.New(database)
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	processor := worker.NewProcessor(10, 3)
	processor.Start(workerCtx)
	defer processor.Stop()

	mux := httpapi.NewServer(svc, store, processor, repo)
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

	// Stop workers first
	workerCancel()
	processor.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "Graceful shutdown failed:", err)
		os.Exit(1)
	}

	fmt.Println("Server stopped cleanly")
}
