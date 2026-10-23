package job

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/storage/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ProduceJobs() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var job JobCreation

		err:=ctx.ShouldBindJSON(&job)
		if err!=nil{
			ctx.JSON(http.StatusBadRequest,gin.H{"error":"Invalid request payload","details":err})
			return 
		}
		jobsCollection:=mongo.OpenCollection("jobs",mongo.Client)
		if jobsCollection==nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":"No such collection exists","details":err})
			return 
		}
		var newJob Job
		newJob.ID= bson.NewObjectID().Hex()
		newJob.Title=job.Title
		newJob.Type=job.Type
		newJob.JobStatus=job.JobStatus
		newJob.Priority=job.Priority
		newJob.Payload=json.RawMessage(job.Payload.String())
		newJob.CurrAttempts = 0
		newJob.MaxAttempts = 5
		newJob.CreatedAt = time.Now()
		c, cancel:=context.WithTimeout(context.Background(),3*time.Second)
		defer cancel()
		result,err:=jobsCollection.InsertOne(c,newJob)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a new job","details":err})
			return
		}
		ctx.JSON(http.StatusCreated, result)
	}
}