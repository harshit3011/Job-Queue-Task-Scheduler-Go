package worker

import (
	"fmt"

	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/dlq"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/executor"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/job"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/queue"
)

type Worker struct{
	jobQueue *queue.Queue
	scheduler *queue.Scheduler
	executor *executor.Executor
	dlq *dlq.DeadLetterQueue
}

func NewWorker(
	jobQueue *queue.Queue,
	scheduler *queue.Scheduler,
	executor *executor.Executor,
	dlq *dlq.DeadLetterQueue,
) *Worker{
	return &Worker{
		jobQueue: jobQueue,
		scheduler: scheduler,
		executor: executor,
		dlq: dlq,
	}
}

func (w *Worker) Start(){
	job := queue.Deque(*w.jobQueue)
	 
}