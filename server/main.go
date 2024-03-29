package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"server/db"
)

func main() {

	mongoClient := db.ConnectToMongoDB()
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to close MongoDB client: %v", err)
		}
	}()

	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	router.Route("/users", func(r chi.Router) {
		// r.Get("/", ListUsersHandler(pg))
		// r.Get("/{userId}", RetrieveUserHandler(pg))
		r.Post("/", CreateUserHandler(mongoClient))
		// r.Patch("/{userId}", UpdateUserHandler(pg))
		// r.Delete("/{userId}", DeleteUserHandler(pg))
	})

	serverPort := os.Getenv("SERVER_PORT")

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", serverPort),
		Handler: router,
	}

	fmt.Printf("Starting server on localhost:%s\n", serverPort)

	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		fmt.Println("Server shutting down")
	}
	if err != nil {
		fmt.Println(err)
	}
}
