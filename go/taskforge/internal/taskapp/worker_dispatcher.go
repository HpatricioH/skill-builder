package taskapp

import (
	"context"
	"fmt"
	"time"

	"taskforge/internal/worker"
)

type WorkerDispatcher struct {
	processor *worker.Processor
}

func NewWorkerDispatcher(processor *worker.Processor) *WorkerDispatcher {
	return &WorkerDispatcher{
		processor: processor,
	}
}

func (d *WorkerDispatcher) Dispatch(ctx context.Context, events []Event) {
	if d.processor == nil {
		return
	}

	for _, event := range events {
		switch event.Type {
		case EventTaskCompleted:
			err := d.processor.Enqueue(ctx, worker.Job{
				Type:      worker.JobTaskCompleted,
				TaskID:    event.TaskID,
				Message:   fmt.Sprintf("task %d completed", event.TaskID),
				CreatedAt: time.Now(),
			})
			if err != nil {
				fmt.Printf("warking: failed to enqueue background job for task %d: %v\n", event.TaskID, err)
			}
		}
	}
}
