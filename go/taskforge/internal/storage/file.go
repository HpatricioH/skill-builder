package storage

import (
	"encoding/json"
	"os"

	"taskforge/internal/task"
)

type FileStorage struct {
	Path string
}

func (f *FileStorage) Save(tasks []task.Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(f.Path, data, 0644)
}
