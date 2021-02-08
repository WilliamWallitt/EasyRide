package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

type Driver struct {
	Id int `json:"id"`
	DriverName string `json:"DriverName"`
	Rate int `json:"Rate"`
}

func initDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./driver_management")
	if err != nil {
		log.Println(err)
	}
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS drivers (id INTEGER PRIMARY KEY, Drivername TEXT, Rate INTEGER )")
	_, err = statement.Exec()
	if err != nil {
		log.Println(err)
	}
	return db
}



func getAllDriversHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := initDb()
	var id int
	var DriverName string
	var Rate int

	// create array
	var drivers []Driver

	rows, _ := db.Query("SELECT id, Drivername, Rate FROM drivers")

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
	err := json.NewEncoder(w).Encode(drivers)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Get drivers")
		_ = db.Close()

	}
}


func updateDriverHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// parse path params
	vars := mux.Vars(r)
	driverName := vars["driver"]
	var driver Driver
	err := json.NewDecoder(r.Body).Decode(&driver)
	if err != nil {
		log.Println(err)
		//http.Error(w, err.Error(), http.StatusBadRequest)
		//return
	} else {
		log.Println(driver.Rate)
	}

	db := initDb()

	query := "UPDATE drivers SET Rate = ? WHERE Drivername = ?"
	statement, _ := db.Prepare(query)
	_, err = statement.Exec(driver.Rate, driverName)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Updated driver")
		_ = db.Close()
	}

}


func getDriverHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driverName := vars["driver"]
	db := initDb()
	query :=  "SELECT id, Drivername, Rate FROM drivers WHERE Drivername = ?"
	statement, _ := db.Prepare(query)

	rows, err := statement.Query(driverName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var id, rate int
	var drivername string
	var driver Driver

	for rows.Next() {
		err := rows.Scan(&id, &drivername, &rate)
		if err != nil {
			log.Println(err)
		}
		driver = Driver{
			Id: id,
			DriverName: driverName,
			Rate: rate,
		}
		fmt.Println(driver)
		err = json.NewEncoder(w).Encode(driver)
		if err != nil {
			log.Println(err)
		}

	}



}


func createDriverHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var driver Driver
	err := json.NewDecoder(r.Body).Decode(&driver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := initDb()

	statement, err := db.Prepare("INSERT INTO drivers (Drivername, Rate) VALUES (?, ?)")
	if err != nil {
		log.Println(err)
	}
	_, err = statement.Exec(driver.DriverName, driver.Rate)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Created new driver")
		_ = db.Close()

	}

}
