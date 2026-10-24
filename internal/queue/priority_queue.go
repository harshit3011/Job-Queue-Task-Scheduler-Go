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