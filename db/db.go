package db

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

func getUser(email string) (User, error) {
	var uu User
	err := db.QueryRow("SELECT email, password from users WHERE email=$1", email).Scan(&uu.Email, &uu.Password)

	return uu, err

}

func createUser(username string, email string, hashedPwd string, phone string) (User, error) {

	result, err := db.Exec("INSERT INTO users (username, email, password, phone) VALUES ($1, $2, $3, $4)", username, email, hashedPwd, phone)

	return result, err
}
