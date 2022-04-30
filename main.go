package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
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

	mySigningKey := []byte("AllYourBase")

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		jwt.StandardClaims
	}

	// Create the Claims

	// TO-DO get OpenAPI
	r := mux.NewRouter()

	// POST on /users
	r.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Post request /users")
		// vars := mux.Vars(r)
		var u user
		// decode the payload
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "bad request", 400)
			return
		}
		// check if user is in database
		for _, uu := range users {

			if u.Email == "niubrandon@nebula.com" {
				fmt.Println("found user email:", uu.Email)
				// found user now check password
				if uu.Password == u.Password {
					fmt.Println("password is correct")

					// Create the Claims
					claims := MyCustomClaims{
						//"bar",
						uu.Email,
						jwt.StandardClaims{
							ExpiresAt: 15000,
							Issuer:    "test",
						},
					}
					// generate jwt token
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
					ss, err := token.SignedString(mySigningKey)
					expiresAt := time.Now().Add(1200 * time.Second)
					fmt.Printf("%v %v", ss, err)
					// send jwt as cookie

					http.SetCookie(w, &http.Cookie{
						Name:    "session_token",
						Value:   ss,
						Expires: expiresAt,
					})
					return
				} else {
					http.Error(w, "wrong password", 401)
					return
				}
			}
		}

		http.Error(w, "user not found", 401)
	}).Methods("POST")

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
