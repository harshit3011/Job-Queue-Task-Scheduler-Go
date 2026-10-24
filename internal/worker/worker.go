package worker

import (
	"fmt"

	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/dlq"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/executor"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/job"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/queue"
)

type Worker struct {
	jobQueue  queue.Queue
	scheduler *queue.Scheduler
	executor  *executor.Executor
	dlq       *dlq.DeadLetterQueue
}

func NewWorker(
	jobQueue queue.Queue,
	scheduler *queue.Scheduler,
	executor *executor.Executor,
	dlq *dlq.DeadLetterQueue,
) *Worker {
	return &Worker{
		jobQueue:  jobQueue,
		scheduler: scheduler,
		executor:  executor,
		dlq:       dlq,
	}
}

func (w *Worker) Start() {
	for {
		j := queue.Deque(w.jobQueue)
		j.JobStatus = job.StatusProcessing
		j.CurrAttempts++
		err := w.executor.Execute(j)

		if err == nil {
			j.JobStatus = job.StatusCompleted
			continue
		}
		if j.CurrAttempts < j.MaxAttempts {
			j.JobStatus = job.StatusPending
			if scheduleErr := w.scheduler.ScheduleQueueJobs(j); scheduleErr != nil {
				j.JobStatus = job.StatusFailed
				w.dlq.Add(j, fmt.Errorf(
					"execution failed: %v; retry scheduling failed: %w",
					err,
					scheduleErr,
				))
				continue
			}
			continue
		}
		j.JobStatus = job.StatusFailed
		w.dlq.Add(j, err)
	}
}
