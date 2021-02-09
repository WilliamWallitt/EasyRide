package main

import (
	"encoding/json"
	"enterprise_computing_cw/Database_Management"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
	"github.com/dgrijalva/jwt-go"

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


func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}


func signUp(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	fmt.Println("Endpoint: signUp")
	var newUser Auth
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authSchema := Database_Management.Database{
		DbName: "./driver_auth",
		Query:  "",
	}

	authSchema.Query = "INSERT INTO auth (Username, Password) VALUES (" + "'" + newUser.Username + "'" + ", " + "'" + hashAndSaltPwd(newUser.Password) + "'" + ")"
	err = authSchema.ExecDB()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		http.Error(w, "Sign up successful", http.StatusOK)
	}

}

func signIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Endpoint: signIn")
	var newUser Auth
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authSchema := Database_Management.Database{
		DbName: "./driver_auth",
		Query:  "SELECT Username, Password FROM auth",
	}

	rows, err := authSchema.QueryDB()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rows == nil {
		http.Error(w, "No users found", http.StatusInternalServerError)
		return
	}

	var Username string
	var Password string

	for rows.Next() {

		err = rows.Scan(&Username, &Password)

		if Username == newUser.Username && verifyPassword(Password, []byte(newUser.Password)) {

			expiration := time.Now().Add(24 * time.Hour)

			//jwt new logic

			cookie := http.Cookie{
				Name: "username",
				Value: newUser.Username,
				Expires: expiration,
			}

			http.SetCookie(w, &cookie)
			http.Error(w, "Login success, cookie set", http.StatusOK)
			return
		}
	}

	http.Error(w, "Login failed, username or password incorrect", http.StatusBadRequest)
	return
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint: getAllUsers")

	authSchema := Database_Management.Database{
		DbName: "./driver_auth",
		Query:  "SELECT id, Username, Password FROM auth",
	}

	rows, err := authSchema.QueryDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rows == nil {
		http.Error(w, "No users found", http.StatusBadRequest)
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
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		AllUsers = append(AllUsers, Users{
			Id: id,
			Username: Username,
			Password: Password,
		})
	}

	err = json.NewEncoder(w).Encode(AllUsers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}


}
