package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"server/utils"
)

func ConnectToDB() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := utils.GetNumericEnv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatalf("sql.Open: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("sql.Ping: %v", err)
	}

	return db
}

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
