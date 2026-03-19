package worker

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

type JobType string

const (
	JobTaskCompleted JobType = "task_completed"
)

var (
	ErrProcessorStopped = errors.New("processor stopped")
	ErrQueueFull        = errors.New("job queue is full")
)

type Job struct {
	Type      JobType
	TaskID    int
	Message   string
	CreatedAt time.Time
}

type Processor struct {
	// This is a channel of job
	jobs    chan Job
	workers int
	wg      sync.WaitGroup

	mu      sync.RWMutex
	stopped bool
}

func NewProcessor(buffer int, workers int) *Processor {
	return &Processor{
		// this create a buffered channel. buffer is the jobs can be queued
		jobs:    make(chan Job, buffer),
		workers: workers,
	}
}

func (p *Processor) Start(ctx context.Context) {
	for i := 1; i <= p.workers; i++ {
		workerID := i

		p.wg.Add(1)

		go func() {
			defer p.wg.Done()

			for {
				select {
				case <-ctx.Done():
					log.Printf("worker=%d cancelled", workerID)
					return

				case job, ok := <-p.jobs:
					if !ok {
						log.Printf("worker=%d stopped", workerID)
						return
					}

					log.Printf(
						"worker=%d processing job type=%s task_id=%d created_at=%s message=%q",
						workerID,
						job.Type,
						job.TaskID,
						job.CreatedAt.Format(time.RFC3339),
						job.Message,
					)

					// Simulate work
					time.Sleep(3 * time.Second)

					log.Printf(
						"worker=%d finished job type=%s task_id=%d",
						workerID,
						job.Type,
						job.TaskID,
					)
				}
			}
		}()
	}
}

func (p *Processor) Enqueue(ctx context.Context, job Job) error {
	p.mu.RLock()
	stopped := p.stopped
	p.mu.RUnlock()

	if stopped {
		return ErrProcessorStopped
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.jobs <- job:
		return nil
	default:
		return ErrQueueFull
	}
}

// this closes the channel
func (p *Processor) Stop() {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return
	}
	p.stopped = true
	close(p.jobs)
	p.mu.Unlock()

	p.wg.Wait()
}
