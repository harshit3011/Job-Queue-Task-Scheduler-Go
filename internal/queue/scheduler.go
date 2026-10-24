package queue

import (
	"container/heap"
	"sync"
	"time"

	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/job"
)

type Scheduler struct {
	TimedJobs   *TimePriorityQueue
	WaitingJobs *JobPriorityQueue
	mu          sync.Mutex
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		TimedJobs:   NewTimePriorityQueue(),
		WaitingJobs: NewJobPriorityQueue(),
	}
}

func (s *Scheduler) ScheduleQueueJobs(jobQueue Queue, job job.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	heap.Push(s.WaitingJobs, &job)

	return nil
}

func (s *Scheduler) ScheduleTimedJobs(jobQueue Queue, job job.Job, t time.Time) error {
	now := time.Now()
	var scheduledAt time.Time
	allowedStart := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
	if t.Before(allowedStart) {
		scheduledAt = allowedStart
	} else {
		tomorrow := now.AddDate(0, 0, 1)
		newTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 9, 0, 0, 0, tomorrow.Location())
		scheduledAt = newTime
	}
	timedJob := TimedJob{
		Job:         job,
		ScheduledAt: scheduledAt,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	heap.Push(s.TimedJobs, &timedJob)
	return nil
}

func (s *Scheduler) Run(jobQueue Queue) {
	for range jobQueue.capacityAvailable {
		for {
			s.mu.Lock()

			if s.WaitingJobs.Len() == 0 {
				s.mu.Unlock()
				break
			}

			j := heap.Pop(s.WaitingJobs).(job.Job)
			s.mu.Unlock()
			if err := Enque(jobQueue, j); err != nil {
				s.mu.Lock()
				heap.Push(s.WaitingJobs, j)
				s.mu.Unlock()
				break
			}
		}
	}
}
