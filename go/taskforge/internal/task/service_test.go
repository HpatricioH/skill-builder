package task

import "testing"

func TestService_AddTask(t *testing.T) {
	test := []struct {
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

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(nil)

			got, err := svc.AddTask(tt.title)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AddTask() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
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
				t.Fatalf("ListTasks() got %d tasks, want 1", len(tasks))
			}

			if tasks[0].ID != got.ID {
				t.Fatalf("stored task ID = %d, want %d", tasks[0].ID, got.ID)
			}
		})
	}
}
