package manager

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/job"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/queue"
)


func ManageJos(job job.Job, jobQueue queue.Queue) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		if IsAllowedTime(){
			if err:=queue.Enque(jobQueue,job);err!=nil{
				if err:=queue.ScheduleQueueJobs(jobQueue, job);err!=nil{
				ctx.JSON(http.StatusInternalServerError,gin.H{"error":"Couldn't schedule the job","details":err})
				return 
				}
			}
			ctx.JSON(http.StatusCreated,gin.H{"Details":"Job Enqueued!"})
		} else {
			if err:=queue.ScheduleTimedJobs(jobQueue, job, time.Now());err!=nil{
				ctx.JSON(http.StatusInternalServerError,gin.H{"error":"Couldn't schedule the job","details":err})
				return 
			}
			ctx.JSON(http.StatusCreated,gin.H{"Details":"Job Enqueued!"})
		}
	}
}

func IsAllowedTime() bool{
	now := time.Now()

	allowedStart:= time.Date(now.Year(),now.Month(),now.Day(),9,0,0,0,now.Location())
	allowedEnd:=time.Date(now.Year(),now.Month(),now.Day(),18,0,0,0,now.Location())

	if !now.Before(allowedStart) && !now.After(allowedEnd){
		return true
	} else{
		return false
	}
}