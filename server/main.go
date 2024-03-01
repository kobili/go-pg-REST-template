package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"

	"server/db"
)

func main() {
	db_host := os.Getenv("DB_HOST")
	db_port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Environment variable DB_PORT is supposed to be an integer. Got: %v", os.Getenv("DB_PORT"))
	}
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")

	db := db.ConnectToDB(db_host, db_port, db_user, db_password, db_name)
	defer db.Close()

	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	router.Route("/users", func(r chi.Router) {
		r.Get("/", ListUsersHandler(db))
		r.Get("/{userId}", RetrieveUserHandler(db))
		r.Post("/", CreateUserHandler(db))
		r.Patch("/{userId}", UpdateUserHandler(db))
		r.Delete("/{userId}", DeleteUserHandler(db))
	})

	server := http.Server{
		Addr:    ":4321",
		Handler: router,
	}

	fmt.Printf("Starting server on %s\n", "localhost:4321")

	err = server.ListenAndServe()
	if err == http.ErrServerClosed {
		fmt.Println("Server shutting down")
	}
	if err != nil {
		fmt.Println(err)
	}
}
