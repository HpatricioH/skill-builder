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

	if created.Title != "Buy Milk" {
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
