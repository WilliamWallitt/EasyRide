package main

import (
	trip_mangement "enterprise_computing_cw/Trip_Mangement"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)


func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// logic
		cookie, err := r.Cookie("username")
		log.Println(r.RequestURI)
		log.Println(cookie)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else {
			next.ServeHTTP(w, r)
		}

	})
}

func main() {


	CreateRoster(dbName, "CREATE TABLE IF NOT EXISTS roster (id INTEGER PRIMARY KEY, Drivername TEXT, Rate INTEGER )")

	authRouter := mux.NewRouter().StrictSlash(true)


	authRouter.HandleFunc("/auth", authPage)
	authRouter.HandleFunc("/auth/users", getAllUsers)
	authRouter.HandleFunc("/auth/login", signIn).Methods("POST")
	authRouter.HandleFunc("/auth/signup", signUp).Methods("POST")

	// get all drivers
	authRouter.HandleFunc("/drivers", getAllDriversHandler)

	// get driver
	authRouter.HandleFunc("/drivers/{driver}", getDriverHandler)

	// create driver (POST)
	authRouter.HandleFunc("/driver", createDriverHandler).Methods("POST")

	// update driver
	authRouter.HandleFunc("/drivers/{driver}", updateDriverHandler).Methods("PUT")

	// get all drivers in the roster
	authRouter.HandleFunc("/rosters", getCurrentRosterHandler)

	// add driver to roster (POST)
	authRouter.HandleFunc("/roster", addDriverToRosterHandler).Methods("POST")

	// remove driver from roster (DELETE)
	authRouter.HandleFunc("/roster/{id}", removeDriverFromRoster).Methods("DELETE")


	authRouter.HandleFunc("/directions", trip_mangement.DirectionsHandler).Methods("POST")
	// can create function that each route uses to check cookies (maybe...)

	//authRouter.HandleFunc("/directions", Driver_Allocation.RouteHandler)

	log.Fatal(http.ListenAndServe(":10000", authRouter))

}