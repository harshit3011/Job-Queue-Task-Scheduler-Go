package queue

import (
	"container/heap"
	"testing"
	"time"

	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/job"
)

func TestScheduleQueueJobsUsesPriorityHeap(t *testing.T) {
	scheduler := NewScheduler()
	if err := scheduler.ScheduleQueueJobs(job.Job{
		ID:       "medium",
		Priority: job.MediumPriority,
	}); err != nil {
		t.Fatalf("ScheduleQueueJobs returned error: %v", err)
	}
	if err := scheduler.ScheduleQueueJobs(job.Job{
		ID:       "high",
		Priority: job.HighPriority,
	}); err != nil {
		t.Fatalf("ScheduleQueueJobs returned error: %v", err)
	}

	scheduler.mu.Lock()
	got := heap.Pop(scheduler.WaitingJobs).(*job.Job)
	scheduler.mu.Unlock()
	if got.ID != "high" {
		t.Fatalf("highest-priority job = %q, want %q", got.ID, "high")
	}
}

func TestRunEnqueuesTimedJobWhenDue(t *testing.T) {
	jobQueue := CreateQueue(1)
	scheduler := NewScheduler()
	defer scheduler.Stop()

	heap.Push(scheduler.TimedJobs, &TimedJob{
		Job:         job.Job{ID: "timed"},
		ScheduledAt: time.Now().Add(10 * time.Millisecond),
	})
	go scheduler.Run(jobQueue)

	select {
	case got := <-jobQueue.queue:
		if got.ID != "timed" {
			t.Fatalf("enqueued job = %q, want %q", got.ID, "timed")
		}
	case <-time.After(time.Second):
		t.Fatal("timed job was not enqueued")
	}
}
