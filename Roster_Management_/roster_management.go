package main

import (
	"app/Libraries/Database_Management"
	"app/Libraries/Error_Management"
	"app/Libraries/Middleware"
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
)


type Roster struct {
	DriverName string `json:"DriverName"`
	Rate int `json:"Rate"`
	Total int
}

// AddDriverToRoster takes in a Roster struct and creates a query to store a new driver in the Roster_DB and executes that query
func AddDriverToRoster(driver Roster) error {

	// create Schema using the db path and query (using driver information)
	rosterSchema := Database_Management.Database{
		DbName: Database_Management.RosterDBPath,
		Query:  "INSERT INTO roster (Drivername, Rate) VALUES " +
			"('" + driver.DriverName + "'" +
			",'"  + strconv.Itoa(driver.Rate) + "')",
	}

	// execute query
	err := rosterSchema.ExecDB()

	// checking that query didn't lead to any DB errors
	if err != nil {
		return err
	} else {
		return nil
	}

}

// RemoveDriverFromRoster takes in the driver name and deletes that driver from the roster (all drivernames are unique)
func RemoveDriverFromRoster(driverName string) error {

	// create schema using db name and query (using driver name)
	rosterSchema := Database_Management.Database{
		DbName: Database_Management.RosterDBPath,
		Query:  "DELETE FROM roster WHERE Drivername=('" + driverName + "')",
	}

	// execute query
	err := rosterSchema.ExecDB()

	// check if query caused any DB errors
	if err != nil {
		return err
	} else {
		return nil
	}

}

// GetBestDriverFromRoster queries the database to get the best driver with the lowest rate
// it also make another query to get how many drivers are in the current roster
// returning a Roster struct containing the driver's name, rate and total number of drivers in the roster
func GetBestDriverFromRoster() (*Roster, error) {

	// get the lowest rate driver
	rosterSchema := Database_Management.Database{
		DbName: Database_Management.RosterDBPath,
		Query:  "SELECT Drivername, Rate FROM roster ORDER by rate LIMIT 1",
	}

	// execute query
	rows, err := rosterSchema.QueryDB()

	// check if query returned an error
	if err != nil {
		return nil, nil
	}
	// check if roster is empty (rows = nil)
	if rows == nil {
		return nil, nil
	}

	// variables to store driverName and rate from query
	var rate int
	var driverName string
	// creating roster struct to return
	var roster Roster

	// go over the rows returned and populate the roster struct
	for rows.Next() {
		err := rows.Scan(&driverName, &rate)
		if err != nil {
			return nil, err
		}
		roster = Roster{
			DriverName: driverName,
			Rate:       rate,
		}
	}

	// get the number of drivers in the roster
	rosterSchema = Database_Management.Database{
		DbName: Database_Management.RosterDBPath,
		Query:  "SELECT COUNT(1) FROM roster",
	}

	// query the db
	rows, err =  rosterSchema.QueryDB()
	// check if db returns an error
	if err != nil {
		return nil, err
	}
	// variable to store the count returned
	var count int
	// iterate over the returned rows, populate the count variable
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return nil, err
		}
	}
	// update the roster.Total attribute of the struct to the driver count
	roster.Total = count
	// return that roster struct
	return &roster, nil
}

// GetBestDriverFromRosterHandler is the http handler for getting the best driver driver from the roster
func GetBestDriverFromRosterHandler(w http.ResponseWriter, r *http.Request) {
	// get the roster struct from the GetBestDriverFromRoster function
	roster, err := GetBestDriverFromRoster()
	// make sure the function didnt return an error
	if err != nil {
		// if an error is returned there are no drivers in the roster
		w.WriteHeader(http.StatusNoContent)

	}
	// encode the roster to json as the response
	err = json.NewEncoder(w).Encode(roster)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// AddDriverToRosterHandler is the http handler for adding a driver to the roster
func AddDriverToRosterHandler(w http.ResponseWriter, r *http.Request) {

	// we are validating the request the user sent to make sure it is correct
	var driverJson Error_Management.Driver
	// decode the json request into the Driver struct
	err := json.NewDecoder(r.Body).Decode(&driverJson)
	// make sure we don't have any problems with this
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// check that the json request is of the correct format
	model, e := Error_Management.FormValidationHandler(driverJson)
	// handle if the json request is not of the correct format errors
	if e.ResponseCode != http.StatusOK {
		w.WriteHeader(e.ResponseCode)
		err := json.NewEncoder(w).Encode(e)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	// as an interface is returned we need to convert it back into the correct struct
	m := *model
	driverJson = m.(Error_Management.Driver)

	// we now use the Rate from the user request and the driver name from the decoded JWT token
	// to populate a Roster struct and pass it through the the AddDriverToRoster function
	// to add it to the roster db
	var driver Roster
	// get driver name from the context returned by the auth middleware
	driver.DriverName = fmt.Sprintf("%v", context.Get(r, "driverName"))
	// get the driver's rate from the user request
	driver.Rate = driverJson.Rate
	// add the driver to the roster
	err = AddDriverToRoster(driver)
	// check that it was successful
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		return
	}
}

// RemoveDriverFromRosterHandler http handler uses the current JWT token's drivername to delete that driver
// from the roster
func RemoveDriverFromRosterHandler(w http.ResponseWriter, r *http.Request) {
	// use RemoveDriverFromRoster function to handle all the logic for deleting that driver
	err := RemoveDriverFromRoster(fmt.Sprintf("%v", context.Get(r, "driverName")))
	// check if any errors have occurred doing this
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}


func main () {

	// set up database when this service is started
	err := Database_Management.CreateDatabase(Database_Management.RosterDBPath, Database_Management.RosterDBInit)
	if err != nil {
		log.Fatal(err)
	}

	// trailing slash is allowed for any route ie /rosters/ allowed and rosters/ allowed
	authRouter := mux.NewRouter().StrictSlash(true)

	// get all drivers in the roster (GET)
	// curl -v -X GET http://localhost:8082/rosters
	authRouter.HandleFunc("/rosters", GetBestDriverFromRosterHandler).Methods("GET")

	// add driver to roster (POST)
	//curl -H "Content-Type: application/json" -X POST -d '{"Rate":11}' http://localhost:8082/rosters
	authRouter.Handle("/rosters", Middleware.AuthMiddleware(AddDriverToRosterHandler)).Methods("POST")

	// remove driver from roster (DELETE)
	//curl -X DELETE http://localhost:8082/rosters
	authRouter.Handle("/rosters", Middleware.AuthMiddleware(RemoveDriverFromRosterHandler)).Methods("DELETE")

	// start server on port :8082, handle any server errors that may occur when starting the server
	err = http.ListenAndServe(":8082", authRouter)
	if err != nil {
		log.Fatal(err)
	}

}


