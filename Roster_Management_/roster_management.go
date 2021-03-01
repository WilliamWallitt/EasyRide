package main

import (
	"encoding/json"
	"app/Libraries/Database_Management"
	"app/Libraries/Error_Management"
	"app/Libraries/Middleware"
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
)


type Roster struct {
	Id int `json:"id"`
	DriverName string `json:"DriverName"`
	Rate int `json:"Rate"`
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

	// checking that query didn't lead to any DB errors
	if err != nil {
		return err
	} else {
		return nil
	}

}

// AddDriverToRoster gets all drivers from the roster DB and returns then in a []Roster struct
func GetAllDriversFromRoster() []Roster {

	rosterSchema := Database_Management.Database{
		DbName: Database_Management.RosterDBPath,
		Query:  "SELECT id, Drivername, Rate FROM roster ORDER BY Rate LIMIT 1",
	}

	rows, _ := rosterSchema.QueryDB()

	if rows == nil {
		return []Roster{}
	}

	var id, rate int
	var driverName string
	var roster []Roster

	for rows.Next() {
		err := rows.Scan(&id, &driverName, &rate)
		if err != nil {
			log.Println(err)
		}
		roster = append(roster, Roster{
			Id:         id,
			DriverName: driverName,
			Rate:       rate,
		})
	}

	return roster
}


func GetCurrentRosterHandler(w http.ResponseWriter, r *http.Request) {

	err := json.NewEncoder(w).Encode(GetAllDriversFromRoster())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func AddDriverToRosterHandler(w http.ResponseWriter, r *http.Request) {

	var driverJson Error_Management.Driver
	err := json.NewDecoder(r.Body).Decode(&driverJson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	model, e := Error_Management.FormValidationHandler(driverJson)

	if e.ResponseCode != http.StatusOK {
		w.WriteHeader(e.ResponseCode)
		err := json.NewEncoder(w).Encode(e)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	m := *model
	driverJson = m.(Error_Management.Driver)

	var driver Roster

	driver.DriverName = fmt.Sprintf("%v", context.Get(r, "driverName"))
	driver.Rate = driverJson.Rate

	err = AddDriverToRoster(driver)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		return
	}
}


func RemoveDriverFromRosterHandler(w http.ResponseWriter, r *http.Request) {

	err := RemoveDriverFromRoster(fmt.Sprintf("%v", context.Get(r, "driverName")))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}


func main () {

	err := Database_Management.CreateDatabase(Database_Management.RosterDBPath, Database_Management.RosterDBInit)
	if err != nil {
		log.Fatal(err)
	}

	authRouter := mux.NewRouter().StrictSlash(true)

	// get all drivers in the roster (GET)
	// curl -v -X GET http://localhost:8082/rosters
	authRouter.HandleFunc("/rosters", GetCurrentRosterHandler).Methods("GET")

	// add driver to roster (POST)
	//curl -H "Content-Type: application/json" -X POST -d '{"Rate":11}' http://localhost:8082/rosters
	authRouter.Handle("/rosters", Middleware.AuthMiddleware(AddDriverToRosterHandler)).Methods("POST")

	// remove driver from roster (DELETE)
	//curl -X DELETE http://localhost:8082/rosters
	authRouter.Handle("/rosters", Middleware.AuthMiddleware(RemoveDriverFromRosterHandler)).Methods("DELETE")

	err = http.ListenAndServe(":8082", authRouter)
	if err != nil {
		log.Fatal(err)
	}

}


