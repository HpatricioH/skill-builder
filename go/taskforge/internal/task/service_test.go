package task

import "testing"

func TestService_AddTask(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		wantErr   bool
		wantID    int
		wantTitle string
	}{
		{
			name:      "valid title",
			title:     "Buy milk",
			wantErr:   false,
			wantID:    1,
			wantTitle: "Buy milk",
		},
		{
			name:    "empty title",
			title:   "   ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(nil)

			got, err := svc.AddTask(tt.title)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AddTask() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if got.ID != tt.wantID {
				t.Fatalf("AddTask() got ID = %d, want %d", got.ID, tt.wantID)
			}
			if got.Title != tt.wantTitle {
				t.Fatalf("AddTask() got Title = %q, want %q", got.Title, tt.wantTitle)
			}

			// Also verify the service actually stores the task
			tasks := svc.ListTasks()
			if len(tasks) != 1 {
				t.Fatalf("ListTasks() len = %d, want 1", len(tasks))
			}
			if tasks[0].ID != got.ID {
				t.Fatalf("stored task ID = %d, want %d", tasks[0].ID, got.ID)
			}
		})
	}
}

func TestService_MarkDone(t *testing.T) {
	svc := NewService(nil)

	t1, err := svc.AddTask("Task 1")
	if err != nil {
		t.Fatalf("AddTask() error = %v", err)
	}

	// Mark it done
	if err := svc.MarkDone(t1.ID); err != nil {
		t.Fatalf("MarkDone() error = %v", err)
	}

	// Verify it is done
	tasks := svc.ListTasks()
	if len(tasks) != 1 {
		t.Fatalf("ListTasks() len = %d, want 1", len(tasks))
	}
	if !tasks[0].Completed {
		t.Fatalf("task Completed = %v, want true", tasks[0].Completed)
	}

	// Marking again should error (based on our current rules)
	if err := svc.MarkDone(t1.ID); err == nil {
		t.Fatalf("MarkDone() second time error = nil, want error")
	}

	// Non-existent ID should error
	if err := svc.MarkDone(999); err == nil {
		t.Fatalf("MarkDone() non-existent error = nil, want error")
	}
}

func TestService_DeleteTask(t *testing.T) {
	svc := NewService(nil)

	t1, _ := svc.AddTask("Task 1")
	t2, _ := svc.AddTask("Task 2")
	t3, _ := svc.AddTask("Task 3")

	// Delete middle task
	if err := svc.DeleteTask(t2.ID); err != nil {
		t.Fatalf("DeleteTask() error = %v", err)
	}

	tasks := svc.ListTasks()
	if len(tasks) != 2 {
		t.Fatalf("ListTasks() len = %d, want 2", len(tasks))
	}

	// Ensure remaining IDs are t1 and t3
	if tasks[0].ID != t1.ID || tasks[1].ID != t3.ID {
		t.Fatalf("remaining IDs = [%d, %d], want [%d, %d]", tasks[0].ID, tasks[1].ID, t1.ID, t3.ID)
	}

	// Deleting non-existent should error
	if err := svc.DeleteTask(999); err == nil {
		t.Fatalf("DeleteTask() non-existent error = nil, want error")
	}
}
