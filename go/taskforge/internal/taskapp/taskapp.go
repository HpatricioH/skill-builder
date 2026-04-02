package taskapp

import (
	"context"

	"taskforge/internal/task"
	"taskforge/internal/taskrepo"
)

type App struct {
	repo       *taskrepo.Repository
	dispatcher Dispatcher
}

func New(repo *taskrepo.Repository, dispatcher Dispatcher) *App {
	return &App{
		repo:       repo,
		dispatcher: dispatcher,
	}
}

func (a *App) CreateTask(ctx context.Context, title string) (task.Task, error) {
	return a.repo.Create(ctx, title)
}

func (a *App) MarkDone(ctx context.Context, id int) ([]Event, error) {
	if err := a.repo.MarkDone(ctx, id); err != nil {
		return nil, err
	}

	events := []Event{
		{
			Type:   EventTaskCompleted,
			TaskID: id,
		},
	}

	if a.dispatcher != nil {
		a.dispatcher.Dispatch(ctx, events)
	}

	return events, nil
}

func (a *App) DeleteTask(ctx context.Context, id int) error {
	return a.repo.Delete(ctx, id)
}

func (a *App) ListTasks(ctx context.Context) ([]task.Task, error) {
	return a.repo.List(ctx)
}

func (a *App) GetTaskByID(ctx context.Context, id int) (task.Task, error) {
	return a.repo.GetByID(ctx, id)
}
