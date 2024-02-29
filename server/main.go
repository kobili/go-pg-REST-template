package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"server/db"
)

func main() {
	db := db.ConnectToDB("localhost", 5432, "postgres", "password", "go_test")
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

	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		fmt.Println("Server shutting down")
	}
	if err != nil {
		fmt.Println(err)
	}
}
