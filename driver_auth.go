package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// dummy data
type Auth struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type Error struct {
	Message string `json:"Error"`
}


// password handling

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


// run after the req handler
func authPage(w http.ResponseWriter, r *http.Request) {
	// would render the authentication page
	fmt.Println("Endpoint: authPage")
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

	db, err := sql.Open("sqlite3", "./driver_auth")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS auth (id INTEGER PRIMARY KEY, Username TEXT, Password TEXT)")
	statement.Exec()
	statement, _ = db.Prepare("INSERT INTO auth (Username, Password) VALUES (?, ?)")
	statement.Exec(newUser.Username, hashAndSaltPwd(newUser.Password))

	err = db.Close()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("DB successfully closed")
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

	db, err := sql.Open("sqlite3", "./driver_auth")
	if err != nil {
		panic(err)
	}

	rows, _ := db.Query("SELECT Username, Password FROM auth")

	var Username string
	var Password string
	for rows.Next() {
		err = rows.Scan(&Username, &Password)
		if err != nil {
			log.Println(err)
		}

		if Username == newUser.Username && verifyPassword(Password, []byte(newUser.Password)) {

			expiration := time.Now().Add(24 * time.Hour)
			cookie := http.Cookie{
				Name: "username",
				Value: newUser.Username,
				Expires: expiration,
			}

			http.SetCookie(w, &cookie)


			if err != nil {
				log.Println(err)
			}

			err = json.NewEncoder(w).Encode(Error{
				Message: "Login success, cookie set",
			})

			if err != nil {
				log.Println(err)
			}

			return
		}
	}

	err = json.NewEncoder(w).Encode(Error{
		Message: "Username or password incorrect",
	})

	if err != nil {
		log.Println(err)
	}

	err = db.Close()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("DB successfully closed")
	}

}

func getAllUsers(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint: getAllUsers")

	db, err := sql.Open("sqlite3", "./driver_auth")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS auth (id INTEGER PRIMARY KEY, Username TEXT, Password TEXT)")
	statement.Exec()

	rows, _ := db.Query("SELECT id, Username, Password FROM auth")
	var id int
	var Username string
	var Password string
	for rows.Next() {
		err := rows.Scan(&id, &Username, &Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		fmt.Fprint(w, "id:" + strconv.Itoa(id) + " Username : " + Username + " Password : " +  Password + "\n")
	}
	if err != nil {
		panic(err)
	}

	err = db.Close()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("DB successfully closed")
	}
}
