package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/producer"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/queue"
	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/storage/mongo"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("No .env file found, using environment variables")
	}
	client := mongo.ConnectDB()
	mongo.Client = client
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal("MongoDB couldn't be connected!!")
	}
	log.Println("MongoDB connected successfully!")

	jobQueue:=queue.CreateQueue(50)
	
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	router.POST("/jobs", producer.ProduceJobs(jobQueue))
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server couldn't be connected!!")
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
