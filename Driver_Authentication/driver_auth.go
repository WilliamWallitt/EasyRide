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

type Error struct {
	Errors []string
}

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"Username"`
	jwt.StandardClaims
}


func hashAndSaltPwd(pwd string) string {
	bytePwd := []byte(pwd)
	hash, _ := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	return string(hash)
}

func verifyPassword(hash string, pwd []byte) bool {
	byteHash := []byte(hash)
	err := bcrypt.CompareHashAndPassword(byteHash, pwd)
	if err != nil {
		return false
	}
	return true
}


func SignUp(w http.ResponseWriter, r *http.Request) {

	var newUser Error_Management.Auth
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	authSchema := Database_Management.Database{
		DbName: Database_Management.DriverAuthDBPath,
		Query:  "INSERT INTO auth (Username, Password) VALUES (" + "'" + newUser.Username + "'" + ", " + "'" + hashAndSaltPwd(newUser.Password) + "'" + ")",
	}

	err = authSchema.ExecDB()

	fmt.Println("DB err", err)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func SignIn(w http.ResponseWriter, r *http.Request) {

	var newUser Error_Management.Auth
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	model, e := Error_Management.FormValidationHandler(newUser)

	if e.ResponseCode != http.StatusOK {
		w.WriteHeader(e.ResponseCode)
		err := json.NewEncoder(w).Encode(e)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// m points to model interface
	m := *model
	// cast interface m to struct Auth
	newUser = m.(Error_Management.Auth)

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
		w.WriteHeader(http.StatusNoContent)
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

	fmt.Println(authSchema)

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
	// curl -v -X GET localhost:10000/auth/users
	authRouter.Handle("/auth/users", Middleware.AuthMiddleware(GetAllUsers)).Methods("GET")

	// user login (POST)
	//curl -H "Content-Type: application/json" -X POST -d '{"Username":"root","Password":"root"}' http://localhost:8080/login
	authRouter.HandleFunc("/login", SignIn).Methods("POST")

	// user signup (POST)
	//curl -H "Content-Type: application/json" -X POST -d '{"Username":"root","Password":"root"}' http://localhost:8080/signup
	authRouter.HandleFunc("/signup", SignUp).Methods("POST")

	err := http.ListenAndServe(":8080", authRouter)
	if err != nil {
		log.Fatal(err)
	}
}

