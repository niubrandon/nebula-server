package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		var u user
		// decode the payload
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "bad request", 400)
			return
		}
		// check if user is in database
		for _, uu := range users {
			fmt.Println("user is", uu.Email, uu.Password)
			if u.Email == "niubrandon@nebula.com" {
				fmt.Println("found email", uu.Email)
				// check password
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
					fmt.Printf("%v %v", ss, err)
					// send jwt as cookie
					http.SetCookie(w, &http.Cookie{
						Name:  "token",
						Value: ss,
					})
					return
				}

			}
		}

		// vars := mux.Vars(r)
		//	email := r.FormValue("email")
		//	password := r.FormValue("password")
		//check password

		fmt.Println(u.Password)
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
