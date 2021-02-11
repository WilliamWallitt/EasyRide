package Driver_Management_

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


func GetAllDriversHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var id int
	var DriverName string
	var Rate int
	var drivers []Driver


	driverSchema := Database_Management.Database{
		DbName: Database_Management.DriverDBPath,
		Query:  "SELECT id, Drivername, Rate FROM drivers",
	}

	//driverSchema.Query = "SELECT id, Drivername, Rate FROM drivers"
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

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	driverId := vars["id"]

	var driver Driver
	err := json.NewDecoder(r.Body).Decode(&driver)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	driverSchema := Database_Management.Database{
		DbName: Database_Management.DriverDBPath,
		Query:  "UPDATE drivers SET Rate = (" + "'" + strconv.Itoa(driver.Rate) +"'"+ ") " +
			"WHERE id = (" + "'" + driverId + "'" + ")",
	}

	err = driverSchema.ExecDB()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}


func GetDriverHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driverId := vars["id"]

	driverSchema := Database_Management.Database{
		DbName: Database_Management.DriverDBPath,
		Query:  "SELECT id, Drivername, Rate FROM drivers WHERE id=('" + driverId + "')",
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

}


func CreateDriverHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var driver Driver
	err := json.NewDecoder(r.Body).Decode(&driver)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	driverSchema := Database_Management.Database{
		DbName: Database_Management.DriverDBPath,
		Query:  "INSERT INTO drivers (Drivername, Rate) VALUES ('" + driver.DriverName + "' , '" + strconv.Itoa(driver.Rate) + "')",
	}

	err = driverSchema.ExecDB()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		return
	}

}
