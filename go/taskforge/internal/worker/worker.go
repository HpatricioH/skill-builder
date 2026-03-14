package worker

import (
	"log"
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
	// This is a channel of strings
	jobs chan Job
}

func NewProcessor(buffer int) *Processor {
	return &Processor{
		// this create a buffered channel. buffer is the jobs can be queued
		jobs: make(chan Job, buffer),
	}
}

func (p *Processor) Start() {
	go func() {
		for job := range p.jobs {
			log.Printf(
				"processing job type=%s task_id=%d created_at=%s message=%q",
				job.Type,
				job.TaskID,
				job.CreatedAt.Format(time.RFC3339),
				job.Message,
			)

			// Simulate work
			time.Sleep(500 * time.Millisecond)

			log.Printf(
				"finished job type=%s task_id=%d",
				job.Type,
				job.TaskID,
			)
		}
	}()
}

// this sends a job into the channel
func (p *Processor) Enqueue(job Job) {
	p.jobs <- job
}

// this closes the channel
func (p *Processor) Stop() {
	close(p.jobs)
}
