package taskrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"taskforge/internal/task"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, title string) (task.Task, error) {
	now := time.Now()

	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO tasks (title, completed, created_at)
		VALUES (?, 0, ?)`,
		title,
		now,
	)
	if err != nil {
		return task.Task{}, fmt.Errorf("insert task: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return task.Task{}, fmt.Errorf("get last insert id: %w", err)
	}

	return task.Task{
		ID:        int(id),
		Title:     title,
		Completed: false,
		CreatedAt: now,
	}, nil
}

func (r *Repository) List(ctx context.Context) ([]task.Task, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, completed, created_at
		FROM tasks 
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list task: %w", err)
	}
	defer rows.Close()

	var tasks []task.Task

	for rows.Next() {
		var t task.Task
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Completed,
			&t.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *Repository) MarkDone(ctx context.Context, id int) error {
	t, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if t.Completed {
		return fmt.Errorf("task already completed")
	}

	res, err := r.db.ExecContext(
		ctx,
		`UPDATE tasks SET completed = 1 WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("mark done: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(
		ctx,
		`DELETE FROM tasks WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (task.Task, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, title, completed, created_at
		FROM tasks
		WHERE id = ?
		`, id)

	var t task.Task
	if err := row.Scan(
		&t.ID,
		&t.Title,
		&t.Completed,
		&t.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return task.Task{}, fmt.Errorf("task not found")
		}
		return task.Task{}, fmt.Errorf("get task by id: %w", err)
	}

	return t, nil
}

func (r *Repository) UpdateTitle(ctx context.Context, id int, title string) (task.Task, error) {
	res, err := r.db.ExecContext(
		ctx,
		`UPDATE tasks SET title = ? WHERE id = ?`,
		title,
		id,
	)
	if err != nil {
		return task.Task{}, fmt.Errorf("update task title: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return task.Task{}, fmt.Errorf("rows affected: %w", err)
	}

	if affected == 0 {
		return task.Task{}, fmt.Errorf("task no found")
	}

	updated, err := r.GetByID(ctx, id)
	if err != nil {
		return task.Task{}, err
	}

	return updated, nil
}

func (r *Repository) ListPaginated(ctx context.Context, limit, offset int) ([]task.Task, error) {
	rows, err := r.db.QueryContext(ctx, `
    SELECT id, title, completed, created_at
		FROM task
		ORDER BY id ASC 
		LIMIT ? OFFSET ?
		`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list paginated tasks: %w", err)
	}
	defer rows.Close()

	var tasks []task.Task

	for rows.Next() {
		var t task.Task
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Completed,
			&t.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, nil
}
