package taskrepo

import (
	"context"
	"path/filepath"
	"testing"

	"taskforge/internal/db"
)

func newTestRepo(t *testing.T) *Repository {
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

	return New(database)
}

// Test Create + List

func TestRepository_CreateAndList(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	created, err := repo.Create(ctx, "Buy milk")
	if err != nil {
		t.Fatalf("Create() err = %v", err)
	}

	if created.ID <= 0 {
		t.Fatalf("Create() ID = %d", created.ID)
	}

	if created.Title != "Buy milk" {
		t.Fatalf("Create() title = %q, want %q", created.Title, "Buy milk")
	}

	if created.Completed {
		t.Fatalf("Create() Completed = %v, want false", created.Completed)
	}

	tasks, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List() err = %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("List() len = %d, want 1", len(tasks))
	}
	if tasks[0].Title != "Buy milk" {
		t.Fatalf("List()[0].Title = %q, want %q", tasks[0].Title, "Buy milk")
	}
}

// GetByID Test

func TestRepository_GetByID(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	created, err := repo.Create(ctx, "Study Go")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if got.ID != created.ID {
		t.Fatalf("GetByID() ID = %d, want %d", got.ID, created.ID)
	}

	if got.Title != created.Title {
		t.Fatalf("GetByID() Title = %q, want %q", got.Title, created.Title)
	}
}

// MarkDone Test
func TestRepository_MarkDone(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	created, err := repo.Create(ctx, "Finish task")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := repo.MarkDone(ctx, created.ID); err != nil {
		t.Fatalf("MarkDone() error = %v", err)
	}

	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if !got.Completed {
		t.Fatalf("Completed = %v, want true", got.Completed)
	}
}

// Delete Test
func TestRepository_Delete(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	created, err := repo.Create(ctx, "Temporary Task")
	if err != nil {
		t.Fatalf("Created() error = %v", err)
	}

	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete() error = %q", err)
	}

	tasks, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(tasks) != 0 {
		t.Fatalf("List() len = %d, want 0", len(tasks))
	}
}

// GetByID not found test
func TestRepository_GetByID_NotFound(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, 999)
	if err == nil {
		t.Fatalf("GetByID() error = nil, watn error")
	}
}

// MarkDone already completed
func TestRepository_MarkDone_AlreadyCompleted(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	created, err := repo.Create(ctx, "Already done")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := repo.MarkDone(ctx, created.ID); err != nil {
		t.Fatalf("first MarkDone() error = %v", err)
	}

	err = repo.MarkDone(ctx, created.ID)
	if err == nil {
		t.Fatalf("second MarkDone() error = nil, want error")
	}
}

// Delete not found
func TestRepository_Delete_NotFound(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	err := repo.Delete(ctx, 999)
	if err == nil {
		t.Fatal("Delete() error = nil, want error")
	}
}

// Update title test
func TestRepository_UpdateTitle(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	created, err := repo.Create(ctx, "Old title")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	updated, err := repo.UpdateTitle(ctx, created.ID, "New Title")
	if err != nil {
		t.Fatalf("UpdateTitle() error = %v", err)
	}

	if updated.Title != "New title" {
		t.Fatalf("updated.Title = %q, want %q", updated.Title, "New Title")
	}

	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID() err = %v", err)
	}

	if got.Title != "New title" {
		t.Fatalf("stored title = %q, want %q", got.Title, "New title")
	}
}
