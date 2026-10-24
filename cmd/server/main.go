package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/dlq"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/executor"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/producer"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/queue"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/storage/mongo"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/worker"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	client, err := mongo.ConnectDB()
	if err != nil {
		log.Fatalf("MongoDB client couldn't be created: %v", err)
	}
	mongo.Client = client
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Ping(pingCtx, nil)
	pingCancel()
	if err != nil {
		log.Printf("MongoDB couldn't be connected: %v", err)
		return
	}
	log.Println("MongoDB connected successfully!")

	jobQueue := queue.CreateQueue(50)
	scheduler := queue.NewScheduler()
	defer scheduler.Stop()
	executorInstance := executor.NewExecutor()
	deadLetterQueue := dlq.NewDeadLetterQueue()

	if err := executorInstance.RegisterExecutor("PRINT", func(rm json.RawMessage) error {
		var data struct {
			Message string `json:"message"`
		}

		if err := json.Unmarshal(rm, &data); err != nil {
			return err
		}

		fmt.Printf("Processed the job: %s\n", data.Message)
		return nil
	}); err != nil {
		log.Printf("Couldn't register PRINT executor: %v", err)
		return
	}

	go scheduler.Run(jobQueue)

	workerCount := 3
	for i := 0; i < workerCount; i++ {
		w := worker.NewWorker(
			jobQueue,
			scheduler,
			executorInstance,
			deadLetterQueue,
		)

		go w.Start()
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	router.POST("/jobs", producer.ProduceJobs(jobQueue, scheduler))

	serverErrors := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("Shutting down server...")
	case err := <-serverErrors:
		log.Printf("Server couldn't start: %v", err)
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
}
