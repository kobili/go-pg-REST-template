package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/go-chi/chi/v5"

	"server/db"
	"server/handlers"
)

/*
Serves a static React Router app. Handles the case where the client refreshes on a URL
set client side by the React Router Javascript.
*/
func serveReactRouterApp(fs http.Handler) http.HandlerFunc {
	// Regex to check if a string ends in a file extension (.js, .css, etc)
	fileMatcher := regexp.MustCompile(`\.[a-zA-Z]*$`)

	fn := func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if !fileMatcher.MatchString(req.URL.Path) {
			// if the request is not for a specific file, it is probably a reload on a url set by react router
			// in this case, just serve the entry point into the react app and let the client side js handle the redirect
			http.ServeFile(w, req, "./static/index.html")
		} else {
			// This is still necessary as the frontend will still expect javascript and css to be served
			fs.ServeHTTP(w, req)
		}
	}

	return http.HandlerFunc(fn)
}

func main() {

	db := db.ConnectToDB()
	defer db.Close()

	router := chi.NewRouter()

	fs := http.FileServer(http.Dir("./static"))
	router.HandleFunc("/*", serveReactRouterApp(fs))

	router.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello world"))
		})

		r.Get("/users", handlers.ListUsersHandler(db))
		r.Get("/users/{userId}", handlers.RetrieveUserHandler(db))
		r.Post("/users", handlers.CreateUserHandler(db))
		r.Patch("/users/{userId}", handlers.UpdateUserHandler(db))
		r.Delete("/users/{userId}", handlers.DeleteUserHandler(db))
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
