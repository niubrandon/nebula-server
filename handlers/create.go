package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/niubrandon/nebula-server/db"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=> Post request /users")
	var u User
	// decode the payload
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	// check if email is in the db

	result := db.FindUser(u.Email)

	if result == true {
		// send error code
		http.Error(w, "user already in database", 400)
		return

	} else {
		// create user
		// encrpt password
		hashedPwd := hashAndSalt([]byte(u.Password))
		//hash password
		res := db.AddUser(u.Username, u.Email, hashedPwd, u.Password)
		if res == true {
			// add user success
			claims := MyCustomClaims{
				u.Email,
				jwt.StandardClaims{
					ExpiresAt: 15000,
					Issuer:    "test",
				},
			}
			// generate jwt token
			mySigningKey := []byte("AllYourBase")
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
		}
		http.Error(w, "adding new user to database failed", 400)
		return
	}

}
