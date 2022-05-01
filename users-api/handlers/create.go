package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost     int = 4  // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword
)

type MyCustomClaims struct {
	Foo string `json:"foo"`
	jwt.StandardClaims
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	mySigningKey := []byte("AllYourBase")
	fmt.Println("Post request /users")
	var u User
	// decode the payload
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	// check if email is in the db
	var uu User
	err := db.QueryRow("SELECT email from users WHERE email=$1", u.Email).Scan(&uu.Email)
	if err != nil {
		fmt.Println("Ready to create account", err)
		// can't find account, create new account
		// encrpt password
		hashedPwd := hashAndSalt([]byte(u.Password))
		// add to db
		result, err := db.Exec("INSERT INTO users (username, email, password, phone) VALUES ($1, $2, $3, $4)", u.Username, u.Email, hashedPwd, u.Phone)
		if err != nil {
			http.Error(w, "creating user failed in db", 400)
			fmt.Println("db execution error", err)
			return
		}
		fmt.Println("New user created:", result)
		// send jwt token
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

		http.Error(w, "email exists", 400)
		return
	}

}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
