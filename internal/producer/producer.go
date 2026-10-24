package producer

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/job"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/manager"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/queue"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/storage/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ProduceJobs(jobQueue queue.Queue) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var jobRequest job.JobCreation

		err := ctx.ShouldBindJSON(&jobRequest)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err})
			return
		}
		jobsCollection := mongo.OpenCollection("jobs", mongo.Client)
		if jobsCollection == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No such collection exists", "details": err})
			return
		}
		var newJob job.Job
		newJob.ID = bson.NewObjectID().Hex()
		newJob.Title = jobRequest.Title
		newJob.Type = jobRequest.Type
		newJob.JobStatus = jobRequest.JobStatus
		newJob.Priority = jobRequest.Priority
		newJob.Payload = json.RawMessage(jobRequest.Payload.String())
		newJob.CurrAttempts = 0
		newJob.MaxAttempts = 5
		newJob.CreatedAt = time.Now()
		c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		result, err := jobsCollection.InsertOne(c, newJob)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a new job", "details": err})
			return
		}
		ctx.JSON(http.StatusCreated, result)
		manageHandler := manager.ManageJos(newJob,jobQueue)
		manageHandler(ctx)
	}
}
