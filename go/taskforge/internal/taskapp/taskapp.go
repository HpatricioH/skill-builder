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

func (a *App) MarkDone(ctx context.Context, id int) ([]Event, error) {
	if err := a.service.MarkDone(id); err != nil {
		return nil, err
	}

	if err := a.store.Save(a.service.ListTasks()); err != nil {
		return nil, err
	}

	events := []Event{
		{
			Type:   EventTaskCompleted,
			TaskID: id,
		},
	}

	a.dispatchEvents(ctx, events)

	return events, nil
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

func (a *App) dispatchEvents(ctx context.Context, events []Event) {
	for _, event := range events {
		switch event.Type {
		case EventTaskCompleted:
			a.enqueueTaskCompleted(ctx, event.TaskID)
		}
	}
}

func (a *App) enqueueTaskCompleted(ctx context.Context, taskID int) {
	if a.processor == nil {
		return
	}

	err := a.processor.Enqueue(ctx, worker.Job{
		Type:      worker.JobTaskCompleted,
		TaskID:    taskID,
		Message:   fmt.Sprintf("task %d completed", taskID),
		CreatedAt: time.Now(),
	})
	if err != nil {
		fmt.Printf("warning: failed to enqueue background job for task %d: %v\n", taskID, err)
	}
}
