package db

import (
	"database/sql"
	"fmt"
)

func InitSchema(db *sql.DB) error {
	query := `
  CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
	  title TEXT NOT NULL, 
	  completed BOOLEAN NOT NULL DEFAULT 0, 
	  created_at DATETIME NOT NULL
	);
	`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("create tasks table: %w", err)
	}

	return nil
}
