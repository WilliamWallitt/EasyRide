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

// Driver struct to store a driver information
type Driver struct {
	Id int `json:"id"`
	DriverName string `json:"DriverName"`
	Rate int `json:"Rate"`
}

func GetAllDriversHandler(w http.ResponseWriter, r *http.Request) {

	var id int
	var DriverName string
	var Rate int
	var drivers []Driver

	driverSchema := Database_Management.Database{
		DbName: Database_Management.DriverDBPath,
		Query:  "SELECT id, Drivername, Rate FROM drivers",
	}

	rows, err := driverSchema.QueryDB()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rows == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	for rows.Next() {
		err := rows.Scan(&id, &DriverName, &Rate)
		if err != nil {
			log.Println(err)
		}
		drivers = append(drivers, Driver{
			Id: id,
			DriverName: DriverName,
			Rate:       Rate,
		})
	}

	err = json.NewEncoder(w).Encode(drivers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// UpdateDriverHandler http handler updates the current drivers rate
// it uses the driver token to make sure the driver is updating it's own rate
// and the json rate sent via PUT request for the new rate
func UpdateDriverHandler(w http.ResponseWriter, r *http.Request) {

	// get the driverName from the token
	var driverName = context.Get(r, "driverName")
	// check that the json request is of the correct format and is valid
	var driverJson Error_Management.Driver
	// decode the json request into a Error_Management.Driver struct
	err := json.NewDecoder(r.Body).Decode(&driverJson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// check that this json request's format is valid
	model, e := Error_Management.FormValidationHandler(driverJson)
	// if not, the user has made a bad request
	if e.ResponseCode != http.StatusOK {
		w.WriteHeader(e.ResponseCode)
		err := json.NewEncoder(w).Encode(e)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	// as the Validation handler returns an interface, cast it into the required struct
	m := *model
	driverJson = m.(Error_Management.Driver)
	// Create scheme with query and driver db
	driverSchema := Database_Management.Database{
		DbName: Database_Management.DriverDBPath,
		Query:  "UPDATE drivers SET Rate = (" + "'" + strconv.Itoa(driverJson.Rate) +"'"+ ") " +
			"WHERE Drivername = (" + "'" + fmt.Sprintf("%v", driverName) + "'" + ")",
	}
	// execute query
	err = driverSchema.ExecDB()
	// check that there is not db errors from this request
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// CreateDriverHandler http handler creates a new driver
// it uses the driver token to get the driver's name
// and the json rate sent via POST request for rate
func CreateDriverHandler(w http.ResponseWriter, r *http.Request) {
	// decode the request into an Error_Management.Driver struct
	var driverJson Error_Management.Driver
	err := json.NewDecoder(r.Body).Decode(&driverJson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// check that the json request is valid and of the correct format
	model, e := Error_Management.FormValidationHandler(driverJson)
	// if not, handle the error
	if e.ResponseCode != http.StatusOK {
		w.WriteHeader(e.ResponseCode)
		err := json.NewEncoder(w).Encode(e)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	// as the Validation handler returns an interface, cast it into the required struct
	m := *model
	driverJson = m.(Error_Management.Driver)
	// create schema with query and driver db
	driverSchema := Database_Management.Database{
		DbName: Database_Management.DriverDBPath,
		Query:  "INSERT INTO drivers (Drivername, Rate) VALUES " +
			"('" + fmt.Sprintf("%v", context.Get(r, "driverName")) +
			"' , '" + strconv.Itoa(driverJson.Rate) + "')",
	}
	// execute query
	err = driverSchema.ExecDB()
	// handle any db errors
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func main() {

	// when service is started, create a driver database
	err := Database_Management.CreateDatabase(Database_Management.DriverDBPath, Database_Management.DriverDBInit)
	if err != nil {
		log.Fatal(err)
	}
	// trailing slash is allowed for any route ie /rosters/ allowed and rosters/ allowed
	authRouter := mux.NewRouter().StrictSlash(true)

	// get all drivers (GET) - remove
	// curl -v -X GET localhost:10000/drivers
	authRouter.Handle("/drivers", Middleware.AuthMiddleware(GetAllDriversHandler)).Methods("GET")

	// create driver (POST)
	//curl -H "Content-Type: application/json" -X POST -d '{"Rate":11}' http://localhost:10000/drivers
	authRouter.Handle("/drivers", Middleware.AuthMiddleware(CreateDriverHandler)).Methods("POST")

	// update driver (PUT)
	// curl -X PUT -H "Content-Type: application/json" -d '{"Rate":11}' http://localhost:8081/drivers
	authRouter.Handle("/drivers", Middleware.AuthMiddleware(UpdateDriverHandler)).Methods("PUT")

	err = http.ListenAndServe(":8081", authRouter)

	if err != nil {
		log.Fatal(err)
	}
}



