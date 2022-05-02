package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/niubrandon/nebula-server/handlers"
	"github.com/rs/cors"
)

func main() {

	// TO-DO get OpenAPI Swagger
	r := mux.NewRouter()

	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/users/login", handlers.Login).Methods("POST")

	fmt.Println("Starting up on 8080")

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://foo.com:8080"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"set-cookie"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	handler := c.Handler(r)
	http.ListenAndServe(":8080", handler)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
