package Roster_Management_

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

type Roster struct {
	Id int `json:"id"`
	DriverName string `json:"DriverName"`
	Rate int `json:"Rate"`
}


func AddDriverToRoster(driver Roster) error {


	rosterSchema := Database_Management.Database{
		DbName: Database_Management.RosterDBPath,
		Query:  "INSERT INTO roster (Drivername, Rate) VALUES " +
			"('" + driver.DriverName + "'" +
			",'"  + strconv.Itoa(driver.Rate) + "')",
	}

	err := rosterSchema.ExecDB()

	if err != nil {
		return err
	} else {
		return nil
	}

}

func RemoveDriverFromRoster(driverName string) error {

	rosterSchema := Database_Management.Database{
		DbName: Database_Management.RosterDBPath,
		Query:  "DELETE FROM roster WHERE Drivername=('" + driverName + "')",
	}

	err := rosterSchema.ExecDB()

	if err != nil {
		return err
	} else {
		return nil
	}

}

func GetAllDriversFromRoster() []Roster {

	rosterSchema := Database_Management.Database{
		DbName: Database_Management.RosterDBPath,
		Query:  "SELECT id, Drivername, Rate FROM roster",
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


