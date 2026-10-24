package manager

import (
	"fmt"
	"time"

	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/job"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/queue"
)

func ManageJob(job job.Job, jobQueue queue.Queue, scheduler *queue.Scheduler) error {
	if scheduler == nil {
		return fmt.Errorf("scheduler is not initialized")
	}
	if IsAllowedTime() {
		if err := queue.Enque(jobQueue, job); err != nil {
			return scheduler.ScheduleQueueJobs(job)
		}
		return nil
	}
	return scheduler.ScheduleTimedJobs(job, time.Now())
}

func IsAllowedTime() bool {
	now := time.Now()

	allowedStart := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
	allowedEnd := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())

	if !now.Before(allowedStart) && !now.After(allowedEnd) {
		return true
	} else {
		return false
	}
}
