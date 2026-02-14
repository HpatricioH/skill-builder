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
		os.Exit(1)
	}

	// Load tasks at the start of every run
	existing, err := store.Load()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		os.Exit(1)
	}

	svc := task.NewService(existing)

	switch args[0] {
	case "add":
		if len(args) < 2 {
			fmt.Println("Missing task title.")
			printUsage()
			os.Exit(1)
		}
		title := strings.Join(args[1:], " ")
		t, err := svc.AddTask(title)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		// save after changes
		if err := store.Save(svc.ListTasks()); err != nil {
			fmt.Println("Error saving tasks:", err)
			os.Exit(1)
		}

		fmt.Printf("Added: #%d %s\n", t.ID, t.Title)

	case "list":
		tasks := svc.ListTasks()
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

	case "done":
		if len(args) < 2 {
			fmt.Println("Missing task ID.")
			printUsage()
			os.Exit(1)
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid task ID:")
			os.Exit(1)
		}

		if err := svc.MarkDone(id); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		if err := store.Save(svc.ListTasks()); err != nil {
			fmt.Println("Error saving tasks:", err)
			os.Exit(1)
		}

		fmt.Printf("Marked task #%d as completed\n", id)

	case "delete":
		if len(args) < 2 {
			fmt.Println("Missing task ID.")
			printUsage()
			os.Exit(1)
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid task ID:")
			os.Exit(1)
		}

		if err := svc.DeleteTask(id); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		if err := store.Save(svc.ListTasks()); err != nil {
			fmt.Println("Error saving tasks:", err)
			os.Exit(1)
		}

		fmt.Printf("Deleted task #%d\n", id)

	default:
		fmt.Printf("Unknown command: %s\n", args[0])
		printUsage()
		os.Exit(1)
	}

}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" taskforge add <title>")
	fmt.Println(" taskforge list")
	fmt.Println(" taskforge done <id>")
	fmt.Println(" taskforge delete <id>")
}
