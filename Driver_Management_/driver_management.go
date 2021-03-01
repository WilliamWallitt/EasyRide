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


func UpdateDriverHandler(w http.ResponseWriter, r *http.Request) {

	var driverName = context.Get(r, "driverName")
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

	driverSchema := Database_Management.Database{
		DbName: Database_Management.DriverDBPath,
		Query:  "UPDATE drivers SET Rate = (" + "'" + strconv.Itoa(driverJson.Rate) +"'"+ ") " +
			"WHERE Drivername = (" + "'" + fmt.Sprintf("%v", driverName) + "'" + ")",
	}

	err = driverSchema.ExecDB()

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}


func CreateDriverHandler(w http.ResponseWriter, r *http.Request) {


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

	driverSchema := Database_Management.Database{
		DbName: Database_Management.DriverDBPath,
		Query:  "INSERT INTO drivers (Drivername, Rate) VALUES " +
			"('" + fmt.Sprintf("%v", context.Get(r, "driverName")) +
			"' , '" + strconv.Itoa(driverJson.Rate) + "')",
	}

	err = driverSchema.ExecDB()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func main() {


	err := Database_Management.CreateDatabase(Database_Management.DriverDBPath, Database_Management.DriverDBInit)
	if err != nil {
		log.Fatal(err)
	}

	authRouter := mux.NewRouter().StrictSlash(true)

	//authRouter.Use(Middleware.AuthMiddleware)

	// get all drivers (GET)
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



