package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"taskforge/internal/storage"
	"taskforge/internal/task"
)

func main() {
	// store tasks.json in the current working directory
	store := storage.NewFileStorage(filepath.Join(".", "tasks.json"))

	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		exitErr(nil, 1)
	}

	// Load tasks at the start of every run
	existing, err := store.Load()
	if err != nil {
		exitErr(fmt.Errorf("loading tasks: %w", err), 1)
	}

	svc := task.NewService(existing)

	switch args[0] {
	case "add":
		title := strings.TrimSpace(strings.Join(args[1:], " "))
		if title == "" {
			exitErr(fmt.Errorf("missing task title"), 1)
		}

		t, err := svc.AddTask(title)
		if err != nil {
			exitErr(err, 1)
		}

		// save after changes
		saveOrExit(store, svc)
		fmt.Printf("Added: #%d %s\n", t.ID, t.Title)

	case "list":
		printTasks(svc.ListTasks())

	case "done":
		id := parseIDOrExit(args, 1)
		if err := svc.MarkDone(id); err != nil {
			exitErr(err, 1)
		}

		saveOrExit(store, svc)
		fmt.Printf("Marked task #%d as completed\n", id)

	case "delete":
		id := parseIDOrExit(args, 1)
		if err := svc.DeleteTask(id); err != nil {
			exitErr(err, 1)
		}

		saveOrExit(store, svc)
		fmt.Printf("Deleted task #%d\n", id)

	default:
		printUsage()
		exitErr(fmt.Errorf("unknown command: %s", args[0]), 1)
	}
}

func parseIDOrExit(args []string, index int) int {
	if len(args) <= index {
		printUsage()
		exitErr(fmt.Errorf("missing task ID"), 1)
	}

	id, err := strconv.Atoi(args[index])
	if err != nil || id <= 0 {
		exitErr(fmt.Errorf("invalid task ID: %q", args[index]), 1)
	}

	return id
}

func saveOrExit(store *storage.FileStorage, svc *task.Service) {
	if err := store.Save(svc.ListTasks()); err != nil {
		exitErr(fmt.Errorf("saving tasks %w", err), 1)
	}
}

func printTasks(tasks []task.Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks yet.")
		return
	}
	for _, t := range tasks {
		status := " "
		if t.Completed {
			status = "x"
		}
		fmt.Printf("[%s] #%d %s\n", status, t.ID, t.Title)
	}
}

func exitErr(err error, code int) {
	if err != nil {
		fmt.Fprint(os.Stderr, "Error:", err)
	}
	os.Exit(code)
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" taskforge add <title>")
	fmt.Println(" taskforge list")
	fmt.Println(" taskforge done <id>")
	fmt.Println(" taskforge delete <id>")
}
