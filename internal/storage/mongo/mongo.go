package mongo

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDB() (*mongo.Client, error) {
	_ = godotenv.Load("../../.env")
	mongoURL := os.Getenv("MONGODB_URI")
	if mongoURL == "" {
		return nil, fmt.Errorf("MONGODB_URI is not set")
	}

	mongoOptions := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(mongoOptions)
	if err != nil {
		return nil, fmt.Errorf("connect to MongoDB: %w", err)
	}
	return client, nil
}

var Client *mongo.Client

func OpenCollection(collection string, client *mongo.Client) (*mongo.Collection, error) {
	if client == nil {
		return nil, fmt.Errorf("MongoDB client is not initialized")
	}
	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		return nil, fmt.Errorf("DATABASE_NAME is not set")
	}
	if collection == "" {
		return nil, fmt.Errorf("collection name is empty")
	}
	return client.Database(dbName).Collection(collection), nil
}
