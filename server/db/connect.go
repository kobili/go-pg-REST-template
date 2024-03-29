package db

import (
	"context"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MONGO_DB_NAME = "Go_REST"
const MONGO_COLLECTION_USER = "users"

func ConnectToMongoDB() *mongo.Client {
	uri := os.Getenv("MONGO_DB_URI")
	if uri == "" {
		log.Fatalf("'MONGO_DB_URI' environment variable must be set")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB server at %s: %v", uri, err)
	}

	log.Printf("Successfully connected to MongoDB server")

	return client
}

func getUserCollection(client mongo.Client) *mongo.Collection {
	return client.Database(MONGO_DB_NAME).Collection(MONGO_COLLECTION_USER)
}
