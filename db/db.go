package db

import (
	"fmt"

	_ "github.com/lib/pq"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

//exposed function need Captial first letter
func FindUser(email string) bool {
	db := createConnection()
	defer db.Close()
	var uu User
	sqlStatement := `SELECT email, password from users WHERE email=$1`
	err := db.QueryRow(sqlStatement, email).Scan(&uu.Email, &uu.Password)
	if err != nil {
		fmt.Println("User not found in database", err)
		return false
	}
	fmt.Println("Found user from database", uu)
	return true

}

func GetUser(email string) (User, error) {
	db := createConnection()
	defer db.Close()
	var uu User
	sqlStatement := `SELECT email, password from users WHERE email=$1`
	err := db.QueryRow(sqlStatement, email).Scan(&uu.Email, &uu.Password)
	if err != nil {
		fmt.Println("User not found in database", err)
	}
	fmt.Println("Found user from database", uu)
	return uu, err

}

func AddUser(username string, email string, hashedPwd string, phone string) bool {
	db := createConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO users (username, email, password, phone) VALUES ($1, $2, $3, $4)`
	result, err := db.Exec(sqlStatement, username, email, hashedPwd, phone)
	if err != nil {
		fmt.Println("Unable to create user", err)
		return false
	}
	fmt.Println("New user has been added to Database", result)
	return true
}
