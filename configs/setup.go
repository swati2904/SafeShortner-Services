package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client                // Global database client
var URLCollection *mongo.Collection // Global URL collection variable

// ConnectDB connects to the MongoDB database and initializes global variables
func ConnectDB() *mongo.Client {
	// Create a new MongoDB client using the URI from the environment
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
	}

	// Create a context with a 10-second timeout for connecting to the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect the client to the MongoDB server
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Ping the database to ensure a successful connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB")

	// Assign the client to the global DB variable
	DB = client

	// Initialize the URL collection globally
	URLCollection = GetCollection(client, "urls") // "urls" is collection name

	return client
}

// GetCollection retrieves a specific collection from the database
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("safeShortnerAPI").Collection(collectionName) // "safeShortnerAPI" database name
}
