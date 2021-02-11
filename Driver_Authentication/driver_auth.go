package Driver_Authentication

import (
	"encoding/json"
	"enterprise_computing_cw/Database_Management"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

type Auth struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

// jwt stuff
var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"Username"`
	jwt.StandardClaims
}


func hashAndSaltPwd(pwd string) string {
	byte_pwd := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(byte_pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func verifyPassword(hash string, pwd []byte) bool {
	byteHash := []byte(hash)
	err := bcrypt.CompareHashAndPassword(byteHash, pwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}


func SignUp(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	fmt.Println("Endpoint: signUp")
	var newUser Auth
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authSchema := Database_Management.Database{
		DbName: Database_Management.DriverAuthDBPath,
		Query:  "INSERT INTO auth (Username, Password) VALUES (" + "'" + newUser.Username + "'" + ", " + "'" + hashAndSaltPwd(newUser.Password) + "'" + ")",
	}

	err = authSchema.ExecDB()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Endpoint: signIn")
	var newUser Auth
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authSchema := Database_Management.Database{
		DbName: Database_Management.DriverAuthDBPath,
		Query:  "SELECT Username, Password FROM auth",
	}

	rows, err := authSchema.QueryDB()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rows == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var Username string
	var Password string

	for rows.Next() {

		err = rows.Scan(&Username, &Password)

		if Username == newUser.Username && verifyPassword(Password, []byte(newUser.Password)) {

			expiration := time.Now().Add(24 * time.Hour)

			// create claims struct with username, expiration date
			claims := &Claims{
				Username: newUser.Username,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expiration.Unix(),
				},
			}

			// create token with server secret + claims data
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			// token -> string(token)
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// create client cookie, which stores the token
			http.SetCookie(w, &http.Cookie{
				Name: "token",
				Path: "/",
				Value: tokenString,
				Expires: expiration,
			})
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
	http.Error(w, "Login failed, username or password incorrect", http.StatusUnauthorized)
	return
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint: getAllUsers")

	authSchema := Database_Management.Database{
		DbName: Database_Management.DriverAuthDBPath,
		Query:  "SELECT id, Username, Password FROM auth",
	}

	rows, err := authSchema.QueryDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rows == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	type Users struct {
		Id int
		Username string
		Password string
	}

	var AllUsers []Users
	var id int
	var Username string
	var Password string

	for rows.Next() {
		err := rows.Scan(&id, &Username, &Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		AllUsers = append(AllUsers, Users{
			Id: id,
			Username: Username,
			Password: Password,
		})
	}

	err = json.NewEncoder(w).Encode(AllUsers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


}
