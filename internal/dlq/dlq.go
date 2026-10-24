package dlq

import (
	"sync"
	"time"

	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/job"
)

type Entry struct {
	job      job.Job
	Reason   string
	FailedAt time.Time
}

type DeadLetterQueue struct {
	entries []Entry
	mu      sync.Mutex
}

func NewDeadLetterQueue() *DeadLetterQueue{
	return &DeadLetterQueue{
		entries: make([]Entry, 0),
	}
}

func (d *DeadLetterQueue) Add(job job.Job, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.entries = append(d.entries, Entry{
		job: job,
		Reason: err.Error(),
		FailedAt: time.Now(),
	})
}

func (d *DeadLetterQueue) Entries() []Entry {
	d.mu.Lock()
	defer d.mu.Unlock()

	return append([]Entry(nil), d.entries...)
}