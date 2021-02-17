package Driver_Management_

import (
	"encoding/json"
	"enterprise_computing_cw/Database_Management"
	"enterprise_computing_cw/Error_Management"
	"fmt"
	"github.com/gorilla/context"
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


func GetDriverHandler(w http.ResponseWriter, r *http.Request) {

	driverSchema := Database_Management.Database{
		DbName: Database_Management.DriverDBPath,
		Query:  "SELECT id, Drivername, Rate FROM drivers WHERE Drivername=('" + fmt.Sprintf("%v", context.Get(r, "driverName")) + "')",
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
