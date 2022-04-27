package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	type user struct {
		ID       int    `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var users = []user{
		{ID: 1, Email: "niubrandon@nebula.com", Password: "superman"},
		{ID: 2, Email: "admin@nebula.com", Password: "superman"},
	}

	//TO-DO get OpenAPI
	r := mux.NewRouter()

	// POST on /users
	r.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "You are trying to login")
		var u user
		//decode the payload
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "bad request", 400)
			return
		}
		//check if user is in database

		//send jwt

		// vars := mux.Vars(r)
		//	email := r.FormValue("email")
		//	password := r.FormValue("password")
		//check password

		fmt.Println(u)
		fmt.Println(users)
		//send jwt token
	}).Methods("POST")

	fmt.Println("Starting up on 8080")

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://foo.com:8080"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	handler := c.Handler(r)
	http.ListenAndServe(":8080", handler)
}

func Hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Hello world!")
}
