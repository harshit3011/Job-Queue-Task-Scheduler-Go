package mongo

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDB() *mongo.Client{
	err:=godotenv.Load("../../.env")
	if err!=nil{
		log.Fatalf("Failed to load the enviroment variables, %v",err)
	}
	mongo_url:=os.Getenv("MONGODB_URI")
	if mongo_url==""{
		log.Fatal("Mongo URL not set")
	}
	mongoOptions:=options.Client().ApplyURI(mongo_url)
	client, err:=mongo.Connect(mongoOptions)
	if err!=nil{
		log.Fatalf("Failed the mongo client: %v",err)
	}
	return client
}

var Client *mongo.Client

func OpenCollection(collection string,client *mongo.Client) *mongo.Collection{
	err:=godotenv.Load(".env")
	if err!=nil{
		log.Println("Failed to load the enviroment variables")
	}
	db_name:=os.Getenv("DATABASE_NAME")
	if db_name==""{
		log.Fatal("Database name not set or Database doesn't exist")
	}
	cln:=client.Database(db_name).Collection(collection)
	if cln==nil{
		return nil
	}
	return cln
}