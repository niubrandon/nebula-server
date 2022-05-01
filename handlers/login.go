package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	mySigningKey := []byte("AllYourBase")
	fmt.Println("Post request /users/login")
	var u User
	// decode the payload
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	//check with dbquery
	var uu User

	err := db.QueryRow("SELECT email, password from users WHERE email=$1", u.Email).Scan(&uu.Email, &uu.Password)
	if err != nil {
		fmt.Println("Query err", err)
		//can't find account, sent http request
		http.Error(w, "account not exists!", 401)
		return
	}
	fmt.Println("printing row result from db", uu.Email, uu.Password)

	if comparePasswords(uu.Password, []byte(u.Password)) {
		// Create the Claims
		claims := MyCustomClaims{
			u.Email,
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
