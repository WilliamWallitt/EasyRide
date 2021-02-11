package main

import (
	"enterprise_computing_cw/Database_Management"
	"enterprise_computing_cw/Driver_Allocation"
	"enterprise_computing_cw/Driver_Authentication"
	"enterprise_computing_cw/Driver_Management_"
	"enterprise_computing_cw/Roster_Management_"
	"enterprise_computing_cw/Trip_Mangement"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)


func main() {

	err := Database_Management.CreateDatabases()
	if err != nil {
		log.Print(err)
	}

	authRouter := mux.NewRouter().StrictSlash(true)

	authRouter.Use(Driver_Authentication.AuthMiddleware)

	// get all users (GET)
	// curl -v -X GET localhost:10000/auth/users
	authRouter.HandleFunc("/auth/users", Driver_Authentication.GetAllUsers).Methods("GET")

	// user login (POST)
	//curl -H "Content-Type: application/json" -X POST -d '{"Username":"root","Password":"root"}' http://localhost:10000/auth/login
	authRouter.HandleFunc("/auth/login", Driver_Authentication.SignIn).Methods("POST")

	// user signup
	//curl -H "Content-Type: application/json" -X POST -d '{"Username":"root","Password":"root"}' http://localhost:10000/auth/signup
	authRouter.HandleFunc("/auth/signup", Driver_Authentication.SignUp).Methods("POST")

	// get all drivers (GET)
	// curl -v -X GET localhost:10000/drivers
	authRouter.HandleFunc("/drivers", Driver_Management_.GetAllDriversHandler).Methods("GET")

	// create driver (POST)
	//curl -H "Content-Type: application/json" -X POST -d '{"DriverName":"root","Rate":11}' http://localhost:10000/drivers
	authRouter.HandleFunc("/drivers", Driver_Management_.CreateDriverHandler).Methods("POST")

	// get driver (GET)
	// curl -v -X GET localhost:10000/drivers/{id}
	authRouter.HandleFunc("/drivers/{id}", Driver_Management_.GetDriverHandler).Methods("GET")

	// update driver (PUT)
	// curl -X PUT -H "Content-Type: application/json" -d '{"Rate":11}' http://localhost:1000/drivers/{id}
	authRouter.HandleFunc("/drivers/{id}", Driver_Management_.UpdateDriverHandler).Methods("PUT")

	// get all drivers in the roster (GET)
	// curl -v -X GET http://localhost:10000/rosters
	authRouter.HandleFunc("/rosters", Roster_Management_.GetCurrentRosterHandler).Methods("GET")

	// add driver to roster (POST)
	//curl -H "Content-Type: application/json" -X POST -d '{"DriverName":"root","Rate":11}' http://localhost:10000/rosters
	authRouter.HandleFunc("/rosters", Roster_Management_.AddDriverToRosterHandler).Methods("POST")

	// remove driver from roster (DELETE)
	//curl -X DELETE http://localhost:10000/rosters/{id}
	authRouter.HandleFunc("/rosters/{id}", Roster_Management_.RemoveDriverFromRosterHandler).Methods("DELETE")

	// get trip directions (POST)
	//curl -H "Content-Type: application/json" -X POST -d '{"Origin":"London","Destination":"Exeter"}' http://localhost:10000/directions
	authRouter.HandleFunc("/directions", trip_mangement.DirectionsHandler).Methods("POST")
	// can create function that each route uses to check cookies (maybe...)

	// get best driver (POST)
	//curl -H "Content-Type: application/json" -X POST -d '{"Origin":"London","Destination":"Exeter"}' http://localhost:10000/allocation
	authRouter.HandleFunc("/allocation", Driver_Allocation.GetBestDriverHandler).Methods("POST")


	err = http.ListenAndServe(":10000", authRouter)
	if err != nil {
		log.Fatal(err)
	}

}