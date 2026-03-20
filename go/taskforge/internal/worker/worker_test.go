package worker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestProcessor_Enqueue_SucceedsWhenQueueHasRoom(t *testing.T) {
	p := NewProcessor(1, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker so jobs can be processed if needed
	p.Start(ctx)
	defer p.Stop()

	err := p.Enqueue(context.Background(), Job{
		Type:      JobTaskCompleted,
		TaskID:    1,
		Message:   "task 1 completed",
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("Enqueue() error = %v, want nil", err)
	}
}

func TestProcessor_Enqueue_ReturnsQueueFullWhenBufferIsFull(t *testing.T) {
	p := NewProcessor(1, 0)

	// Do NOT start workers
	// This keeps the buffered channel from being drained
	// so we can force queue-full behaviour deterministically

	firstErr := p.Enqueue(context.Background(), Job{
		Type:      JobTaskCompleted,
		TaskID:    1,
		Message:   "task 1 completed",
		CreatedAt: time.Now(),
	})
	if firstErr != nil {
		t.Fatalf("Enqueue() error = %v, want nil", firstErr)
	}

	secondErr := p.Enqueue(context.Background(), Job{
		Type:      JobTaskCompleted,
		TaskID:    2,
		Message:   "task 2 completed",
		CreatedAt: time.Now(),
	})
	if !errors.Is(secondErr, ErrQueueFull) {
		t.Fatalf("second enqueue() error = %v, want %v", secondErr, ErrQueueFull)
	}

	// Clean up safely without workers
	p.Stop()
}

func TestProcessor_Enqueue_ReturnsProcessorStoppedAfterStop(t *testing.T) {
	p := NewProcessor(1, 1)

	p.Stop()

	err := p.Enqueue(context.Background(), Job{
		Type:      JobTaskCompleted,
		TaskID:    1,
		Message:   "task 1 completed",
		CreatedAt: time.Now(),
	})
	if !errors.Is(err, ErrProcessorStopped) {
		t.Fatalf("Enqueue() after Stop error = %v, want %v", err, ErrProcessorStopped)
	}
}
