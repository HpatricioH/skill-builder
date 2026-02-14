package task

import (
	"errors"
	"strings"
	"time"
)

type Service struct {
	tasks  []Task
	nextID int
}

func NewService(existing []Task) *Service {
	s := &Service{
		tasks:  existing,
		nextID: 1,
	}

	s.nextID = s.computeNextID()

	return s
}

func (s *Service) computeNextID() int {
	maxID := 0
	for _, t := range s.tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

func (s *Service) AddTask(title string) (Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, errors.New("title cannot be empty")
	}

	t := Task{
		ID:        s.nextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	s.tasks = append(s.tasks, t)
	s.nextID++

	return t, nil
}

func (s *Service) ListTasks() []Task {
	// Return a copy to avoid accidental external mutation.
	out := make([]Task, len(s.tasks))
	copy(out, s.tasks)
	return out
}

func (s *Service) MarkDone(id int) error {
	for i := range s.tasks {
		if s.tasks[i].ID == id {
			if s.tasks[i].Completed {
				return errors.New("task is already completed")
			}
			s.tasks[i].Completed = true
			return nil
		}
	}
	return errors.New("task not found")
}

func (s *Service) DeleteTask(id int) error {
	for i := range s.tasks {
		if s.tasks[i].ID == id {
			// Remove element at index i:
			// - keep everything before i
			// - append everything after i
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}
