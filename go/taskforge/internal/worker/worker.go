package worker

import (
	"fmt"
	"log"
)

type Processor struct {
	// This is a channel of strings
	jobs chan string
}

func NewProcessor(buffer int) *Processor {
	return &Processor{
		// this create a buffered channel. buffer is the jobs can be queued
		jobs: make(chan string, buffer),
	}
}

func (p *Processor) Start() {
	go func() {
		for job := range p.jobs {
			log.Printf("processing job: %s\n", job)
			fmt.Printf("background job done: %s\n", job)
		}
	}()
}

// this sends a job into the channel
func (p *Processor) Enqueue(job string) {
	p.jobs <- job
}

// this closes the channel
func (p *Processor) Stop() {
	close(p.jobs)
}
