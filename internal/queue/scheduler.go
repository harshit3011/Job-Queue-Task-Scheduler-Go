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
	wake        chan struct{}
	stop        chan struct{}
	stopOnce    sync.Once
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		TimedJobs:   NewTimePriorityQueue(),
		WaitingJobs: NewJobPriorityQueue(),
		wake:        make(chan struct{}, 1),
		stop:        make(chan struct{}),
	}
}

func (s *Scheduler) ScheduleQueueJobs(j job.Job) error {
	s.mu.Lock()
	heap.Push(s.WaitingJobs, &j)
	s.mu.Unlock()
	s.notify()
	return nil
}

func (s *Scheduler) ScheduleTimedJobs(job job.Job, t time.Time) error {
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
	heap.Push(s.TimedJobs, &timedJob)
	s.mu.Unlock()
	s.notify()
	return nil
}

func (s *Scheduler) Run(jobQueue Queue) {
	for {
		now := time.Now()
		s.mu.Lock()
		for s.TimedJobs.Len() > 0 && !s.TimedJobs.Peek().ScheduledAt.After(now) {
			timedJob := heap.Pop(s.TimedJobs).(*TimedJob)
			heap.Push(s.WaitingJobs, &timedJob.Job)
		}
		s.mu.Unlock()

		for {
			s.mu.Lock()
			if s.WaitingJobs.Len() == 0 {
				s.mu.Unlock()
				break
			}
			j := *heap.Pop(s.WaitingJobs).(*job.Job)
			s.mu.Unlock()

			if err := Enque(jobQueue, j); err != nil {
				s.mu.Lock()
				heap.Push(s.WaitingJobs, &j)
				s.mu.Unlock()
				break
			}
		}

		var timer *time.Timer
		var timerC <-chan time.Time
		s.mu.Lock()
		if s.TimedJobs.Len() > 0 {
			delay := time.Until(s.TimedJobs.Peek().ScheduledAt)
			if delay < 0 {
				delay = 0
			}
			timer = time.NewTimer(delay)
			timerC = timer.C
		}
		s.mu.Unlock()

		select {
		case <-jobQueue.capacityAvailable:
		case <-s.wake:
		case <-timerC:
		case <-s.stop:
			if timer != nil {
				timer.Stop()
			}
			return
		}

		if timer != nil && !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
	}
}

func (s *Scheduler) Stop() {
	s.stopOnce.Do(func() {
		close(s.stop)
	})
}

func (s *Scheduler) notify() {
	select {
	case s.wake <- struct{}{}:
	default:
	}
}
