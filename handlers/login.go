package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/niubrandon/nebula-server/db"
)

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=> Post request /users/login")
	var u User
	// decode the payload
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	res, err := db.GetUser(u.Email)
	if err != nil {
		http.Error(w, "user not found in db", 400)
	}

	fmt.Println("=> found user infor", res.Email, res.Password)
	if comparePasswords(res.Password, []byte(u.Password)) == true {
		// Create the Claims
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
	} else {
		http.Error(w, "wrong password", 401)
		return
	}

}
