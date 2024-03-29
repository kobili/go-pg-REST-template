package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"server/db"
)

func main() {

	db := db.ConnectToDB()
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
