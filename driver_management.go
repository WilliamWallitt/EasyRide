package main

import (
	"encoding/json"
	"enterprise_computing_cw/Database_Management"
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


func getAllDriversHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var id int
	var DriverName string
	var Rate int
	var drivers []Driver


	driverSchema := Database_Management.Database{
		DbName: "./driver_management",
		Query:  "SELECT id, Drivername, Rate FROM drivers",
	}

	//driverSchema.Query = "SELECT id, Drivername, Rate FROM drivers"
	rows, err := driverSchema.QueryDB()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rows == nil {
		http.Error(w, "No users found", http.StatusOK)
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
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

}


func updateDriverHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	driverId := vars["id"]

	var driver Driver
	err := json.NewDecoder(r.Body).Decode(&driver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	driverSchema := Database_Management.Database{
		DbName: "./driver_management",
		Query:  "UPDATE drivers SET Rate = (" + "'" + strconv.Itoa(driver.Rate) +"'"+ ") " +
			"WHERE id = (" + "'" + driverId + "'" + ")",
	}

	err = driverSchema.ExecDB()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		http.Error(w, "Driver updated", http.StatusOK)
		return
	}
}


func getDriverHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driverId := vars["id"]

	driverSchema := Database_Management.Database{
		DbName: "./driver_management",
		Query:  "SELECT id, Drivername, Rate FROM drivers WHERE id=('" + driverId + "')",
	}

	rows, err := driverSchema.QueryDB()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rows == nil {
		http.Error(w, "No driver found", http.StatusOK)
		return
	}

	var id, rate int
	var driverName string
	//var driverName string
	var driver Driver

	for rows.Next() {
		err := rows.Scan(&id, &driverName, &rate)

		driver = Driver{
			Id: id,
			DriverName: driverName,
			Rate: rate,
		}

		err = json.NewEncoder(w).Encode(driver)
		if err != nil {
			http.Error(w, "Json encode error", http.StatusOK)
			return
		}
	}

	http.Error(w, "No Driver found", http.StatusOK)
	return

}


func createDriverHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var driver Driver
	err := json.NewDecoder(r.Body).Decode(&driver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	driverSchema := Database_Management.Database{
		DbName: "./driver_management",
		Query:  "INSERT INTO drivers (Drivername, Rate) VALUES ('" + driver.DriverName + "' , '" + strconv.Itoa(driver.Rate) + "')",
	}

	err = driverSchema.ExecDB()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		http.Error(w, "Driver added", http.StatusOK)
		return
	}

}
