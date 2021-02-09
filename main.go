package main

import (
	"enterprise_computing_cw/Database_Management"
	"enterprise_computing_cw/Driver_Allocation"
	trip_mangement "enterprise_computing_cw/Trip_Mangement"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)


func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonAuthRoutes := []string{"/",  "/auth/users", "/auth/login", "/auth/signup"}
		currentRoute := r.RequestURI
		requireAuth := true
		for _, route := range nonAuthRoutes {
			if route  == currentRoute {
				requireAuth = false
			}
		}
		fmt.Print(requireAuth, nonAuthRoutes, currentRoute)
		if requireAuth {
			// needs authentication
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else {
			next.ServeHTTP(w, r)
		}
		//cookie, err := r.Cookie("username")
		//log.Println(r.RequestURI)
		//log.Println(cookie)
		//if err != nil {
		//	fmt.Println(err)
		//	http.Error(w, "Forbidden", http.StatusForbidden)
		//} else {
		//	next.ServeHTTP(w, r)
		//}

	})
}





func main() {

	// init db's
	Database_Management.CreateDatabases()

	authRouter := mux.NewRouter().StrictSlash(true)
	authRouter.Use(authMiddleware)

	authRouter.HandleFunc("/", redirectToLogin)
	authRouter.HandleFunc("/auth/users", getAllUsers)
	authRouter.HandleFunc("/auth/login", signIn).Methods("POST")
	authRouter.HandleFunc("/auth/signup", signUp).Methods("POST")

	// get all drivers
	authRouter.HandleFunc("/drivers", getAllDriversHandler)

	// get driver
	authRouter.HandleFunc("/drivers/{id}", getDriverHandler)

	// create driver (POST)
	authRouter.HandleFunc("/driver", createDriverHandler).Methods("POST")

	// update driver
	authRouter.HandleFunc("/driver/{id}", updateDriverHandler).Methods("PUT")

	// get all drivers in the roster
	authRouter.HandleFunc("/rosters", getCurrentRosterHandler)

	// add driver to roster (POST)
	authRouter.HandleFunc("/roster", addDriverToRosterHandler).Methods("POST")

	// remove driver from roster (DELETE)
	authRouter.HandleFunc("/roster/{id}", removeDriverFromRoster).Methods("DELETE")


	authRouter.HandleFunc("/directions", trip_mangement.DirectionsHandler).Methods("POST")
	// can create function that each route uses to check cookies (maybe...)

	authRouter.HandleFunc("/allocation", Driver_Allocation.GetBestDriver).Methods("POST")

	//authRouter.HandleFunc("/directions", Driver_Allocation.RouteHandler)

	log.Fatal(http.ListenAndServe(":10000", authRouter))

}