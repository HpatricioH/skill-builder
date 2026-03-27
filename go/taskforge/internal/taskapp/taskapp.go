package taskapp

import (
	"context"

	"taskforge/internal/storage"
	"taskforge/internal/task"
)

type App struct {
	service    *task.Service
	store      *storage.FileStorage
	dispatcher Dispatcher
}

func New(service *task.Service, store *storage.FileStorage, dispatcher Dispatcher) *App {
	return &App{
		service:    service,
		store:      store,
		dispatcher: dispatcher,
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

	if a.dispatcher != nil {
		a.dispatcher.Dispatch(ctx, events)
	}

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
