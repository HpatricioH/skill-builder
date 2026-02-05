package storage

import (
	"encoding/json"
	"errors"
	"os"

	"taskforge/internal/task"
)

type FileStorage struct {
	Path string
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{Path: path}
}

func (f *FileStorage) Load() ([]task.Task, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		// if file doesn't exist, that's fine: start empty
		if errors.Is(err, os.ErrNotExist) {
			return []task.Task{}, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return []task.Task{}, nil
	}

	var tasks []task.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (f *FileStorage) Save(tasks []task.Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(f.Path, data, 0644)
}
