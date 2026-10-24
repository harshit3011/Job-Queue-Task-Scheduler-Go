package queue

import (
	"container/heap"
	"time"

	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/job"
)

var pMap = map[string]int{
	"HIGH":   3,
	"MEDIUM": 2,
	"LOW":    1,
}

type TimedJob struct{
	Job job.Job
	ScheduledAt time.Time
}

type JobPriorityQueue []*job.Job

func (pq JobPriorityQueue) Len() int{
	return len(pq)
}

func(pq JobPriorityQueue) Less(i, j int) bool{
	
	return pMap[string(pq[i].Priority)] > pMap[string(pq[j].Priority)]
}
func (pq JobPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *JobPriorityQueue) Push(x any) {
job := x.(*job.Job)
*pq = append(*pq, job)
}

func (pq *JobPriorityQueue) Pop() any {
old := *pq
n := len(old)
job := old[n-1]
*pq = old[0:n-1]
return job
}

type TimePriorityQueue []*TimedJob

func (pq TimePriorityQueue) Len() int{
	return len(pq)
}

func(pq TimePriorityQueue) Less(i, j int) bool{
	
	return pq[i].ScheduledAt.Before(pq[j].ScheduledAt)
}
func (pq TimePriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *TimePriorityQueue) Push(x any) {
job := x.(*TimedJob)
*pq = append(*pq, job)
}

func (pq *TimePriorityQueue) Pop() any {
old := *pq
n := len(old)
job := old[n-1]
*pq = old[0:n-1]
return job
}

func NewJobPriorityQueue() *JobPriorityQueue{
	pq:=JobPriorityQueue{}
	heap.Init(&pq)
	return &pq
}

func NewTimePriorityQueue() *TimePriorityQueue{
	pq:=TimePriorityQueue{}
	heap.Init(&pq)
	return &pq
}