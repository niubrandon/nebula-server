package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	db_password := goDotEnvVariable("DB_PASSWORD")
	// db connection
	const (
		host = "localhost"
		port = 5432
		// use dbuser instead of user
		dbuser = "postgres"
		dbname = "nebula"
	)

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, dbuser, db_password, dbname)

	//open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	//close database
	defer db.Close()

	//check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")

	const (
		MinCost     int = 4  // the minimum allowable cost as passed in to GenerateFromPassword
		MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
		DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword
	)

	type user struct {
		ID       int    `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	mySigningKey := []byte("AllYourBase")

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		jwt.StandardClaims
	}

	// TO-DO get OpenAPI Swagger
	r := mux.NewRouter()

	// POST on /login
	r.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Post request /users")
		var u user
		// decode the payload
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "bad request", 400)
			return
		}
		//check with dbquery
		var uu user

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

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
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
