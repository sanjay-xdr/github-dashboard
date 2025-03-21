package database

import (
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const mongoURI = "mongodb://localhost:27017"

var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
)

// InitMongo initializes the MongoDB client (Singleton Pattern)
func InitMongo() {
	clientOnce.Do(func() {
		var err error
		clientInstance, err = mongo.Connect(options.Client().ApplyURI(mongoURI))
		if err != nil {
			log.Fatal("MongoDB Connection Error:", err)
		}
		fmt.Println("Connected to MongoDB")
	})
}

// GetMongoClient returns the singleton MongoDB client
func GetMongoClient() *mongo.Client {
	if clientInstance == nil {
		InitMongo()
	}
	return clientInstance
}

// GetPRCollection returns the PR collection
func GetPRCollection() *mongo.Collection {
	return GetMongoClient().Database("github_dashboard").Collection("pull_requests")
}

func GetWorkflowCollection() *mongo.Collection {
	return GetMongoClient().Database("github_dashboard").Collection("workflow")
}
