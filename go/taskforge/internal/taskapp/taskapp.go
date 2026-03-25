package taskapp

import (
	"context"
	"fmt"
	"time"

	"taskforge/internal/storage"
	"taskforge/internal/task"
	"taskforge/internal/worker"
)

type App struct {
	service   *task.Service
	store     *storage.FileStorage
	processor *worker.Processor
}

func New(service *task.Service, store *storage.FileStorage, processor *worker.Processor) *App {
	return &App{
		service:   service,
		store:     store,
		processor: processor,
	}
}

func (a *App) CreateTask(ctx context.Context, title string) (task.Task, error) {
	t, err := a.service.AddTask(title)
	if err != nil {
		return task.Task{}, err
	}

	if err := a.store.Save(a.service.ListTasks()); err != nil {
		return task.Task{}, err
	}

	return t, nil
}

func (a *App) MarkDone(ctx context.Context, id int) error {
	if err := a.service.MarkDone(id); err != nil {
		return err
	}

	if err := a.store.Save(a.service.ListTasks()); err != nil {
		return err
	}

	if a.processor != nil {
		_ = a.processor.Enqueue(ctx, worker.Job{
			Type:      worker.JobTaskCompleted,
			TaskID:    id,
			Message:   fmt.Sprintf("task %d completed", id),
			CreatedAt: time.Now(),
		})
	}

	return nil
}

func (a *App) DeleteTask(ctx context.Context, id int) error {
	if err := a.service.DeleteTask(id); err != nil {
		return err
	}

	if err := a.store.Save(a.service.ListTasks()); err != nil {
		return err
	}

	return nil
}
