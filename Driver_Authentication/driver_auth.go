package main

import (
	"app/Libraries/Database_Management"
	"app/Libraries/Error_Management"
	"app/Libraries/Middleware"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// secret key for the jwtKey
var jwtKey = []byte("my_secret_key")

// jwt claims struct
type Claims struct {
	Username string `json:"Username"`
	jwt.StandardClaims
}

// function to hash and salt the password
func hashAndSaltPwd(pwd string) string {
	bytePwd := []byte(pwd)
	hash, _ := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	return string(hash)
}

// function to verify a given hash is the right password (hash)
func verifyPassword(hash string, pwd []byte) bool {
	byteHash := []byte(hash)
	err := bcrypt.CompareHashAndPassword(byteHash, pwd)
	if err != nil {
		return false
	}
	return true
}

// SignUp http handler for user signup
func SignUp(w http.ResponseWriter, r *http.Request) {

	// decode the Username and Password json request into the Auth struct
	var newUser Error_Management.Auth
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check that all json fields are correct and valid
	model, e := Error_Management.FormValidationHandler(newUser)

	if e.ResponseCode != http.StatusOK {
		w.WriteHeader(e.ResponseCode)
		err := json.NewEncoder(w).Encode(e)
		if err != nil {
			fmt.Println("JSON err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	// m points to model interface
	m := *model
	// cast interface m to struct Auth
	newUser = m.(Error_Management.Auth)

	// create scheme with the query and auth db
	authSchema := Database_Management.Database{
		DbName: Database_Management.DriverAuthDBPath,
		Query:  "INSERT INTO auth (Username, Password) VALUES " + "('" + newUser.Username + "'" + ",'" + hashAndSaltPwd(newUser.Password) + "')",
	}

	// execute the query
	err = authSchema.ExecDB()
	// handle any internal db errors
	if err != nil {
		w.WriteHeader(http.StatusConflict)
	}

	w.WriteHeader(http.StatusCreated)

}

// SignIn http handler for user login
func SignIn(w http.ResponseWriter, r *http.Request) {

	// decode the Username and Password json request into an Auth struct
	var newUser Error_Management.Auth
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 	// check that all json fields are correct and valid
	model, e := Error_Management.FormValidationHandler(newUser)

	if e.ResponseCode != http.StatusOK {
		w.WriteHeader(e.ResponseCode)
		err := json.NewEncoder(w).Encode(e)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	// m points to model interface
	m := *model
	// cast interface m to struct Auth
	newUser = m.(Error_Management.Auth)

	// create schema with query and auth db
	// getting the Username and Password (hash) for that user
	authSchema := Database_Management.Database{
		DbName: Database_Management.DriverAuthDBPath,
		Query:  "SELECT Username, Password FROM auth WHERE Username=('" + newUser.Username + "')",
	}

	// query the db
	rows, err := authSchema.QueryDB()
	// handle any interal db errors
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// handle if the user doesnt exist
	if rows == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// variables to store the user's Username and Password
	var Username string
	var Password string
	// iterate over the rows returned
	for rows.Next() {
		// store that row's Username and Password
		err = rows.Scan(&Username, &Password)
		// check if the Username and Password is correct
		if Username == newUser.Username && verifyPassword(Password, []byte(newUser.Password)) {
			// add 24 hr expiration date to token
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
			// token -> signed string(token)
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

	// if the password and username is invalid, then we send an unauthorized response back to the user
	w.WriteHeader(http.StatusUnauthorized)
	http.Error(w, "Login failed, username or password incorrect", http.StatusUnauthorized)
	return
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fmt.Println(f.Name())
	}

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


func main() {

	Database_Management.CreateDatabase(Database_Management.DriverAuthDBPath, Database_Management.DriverAuthDBInit)

	authRouter := mux.NewRouter().StrictSlash(true)

	// get all users (GET) - remove
	// curl -b 'token=<token from user login here>' -X GET localhost:3000/auth/users
	authRouter.Handle("/auth/users", Middleware.AuthMiddleware(GetAllUsers)).Methods("GET")

	// user login (POST)
	// curl -H "Content-Type: application/json" -X POST -i -d '{"Username":"test","Password":"test"}' http://localhost:3000/login
	authRouter.HandleFunc("/login", SignIn).Methods("POST")

	// user signup (POST)
	//curl -H "Content-Type: application/json" -X POST -d '{"Username":"root","Password":"root"}' http://localhost:3000/signup
	authRouter.HandleFunc("/signup", SignUp).Methods("POST")

	err := http.ListenAndServe(":8080", authRouter)
	if err != nil {
		log.Fatal(err)
	}
}

