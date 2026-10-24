package queue

import (
	"fmt"

	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/job"
)

type Queue struct {
	queue    chan job.Job
	capacity int
}

func CreateQueue(cap int) Queue{
	return Queue{
		queue: make(chan job.Job, cap),
		capacity: cap,
	}
}

func QueueCapacity(queue Queue) int{
	return queue.capacity
}

func Enque(queue Queue, job job.Job) error{
	select{
	case queue.queue <- job:
		return nil
	default:
		return fmt.Errorf("Failed: Channel buffer is full! Cannot enter without blocking.")
	}
}

