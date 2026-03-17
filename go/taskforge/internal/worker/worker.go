package worker

import (
	"log"
	"sync"
	"time"
)

type JobType string

const (
	JobTaskCompleted JobType = "task_completed"
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
}

func NewProcessor(buffer int, workers int) *Processor {
	return &Processor{
		// this create a buffered channel. buffer is the jobs can be queued
		jobs:    make(chan Job, buffer),
		workers: workers,
	}
}

func (p *Processor) Start() {
	for i := 1; i <= p.workers; i++ {
		workerID := i

		p.wg.Add(1)

		go func() {
			defer p.wg.Done()

			for job := range p.jobs {
				log.Printf(
					"worker=%d processing job type=%s task_id=%d created_at=%s message=%q",
					workerID,
					job.Type,
					job.TaskID,
					job.CreatedAt.Format(time.RFC3339),
					job.Message,
				)

				// Simulate work
				time.Sleep(500 * time.Millisecond)

				log.Printf(
					"worker=%d finished job type=%s task_id=%d",
					workerID,
					job.Type,
					job.TaskID,
				)
			}

			log.Printf("worker=%d stopeed", workerID)
		}()
	}
}

// this sends a job into the channel
func (p *Processor) Enqueue(job Job) {
	p.jobs <- job
}

// this closes the channel
func (p *Processor) Stop() {
	close(p.jobs)
	p.wg.Wait()
}
