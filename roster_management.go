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

type Roster struct {
	Id int `json:"id"`
	DriverName string `json:"DriverName"`
	Rate int `json:"Rate"`
}


func AddDriverToRoster(driver Roster) error {

	rosterSchema := Database_Management.Database{
		DbName: "./roster_management",
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

func RemoveDriverFromRoster(id int) error {

	rosterSchema := Database_Management.Database{
		DbName: "./roster_management",
		Query:  "DELETE FROM roster WHERE id=('" + strconv.Itoa(id) + "')",
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
		DbName: "./roster_management",
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


func getCurrentRosterHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(GetAllDriversFromRoster())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

}

func addDriverToRosterHandler(w http.ResponseWriter, r *http.Request) {
	var driver Roster
	err := json.NewDecoder(r.Body).Decode(&driver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = AddDriverToRoster(driver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		http.Error(w, "Driver added to roster", http.StatusOK)
		return
	}
}


func removeDriverFromRoster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driver_id := vars["id"]
	int, err := strconv.Atoi(driver_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = RemoveDriverFromRoster(int)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		http.Error(w, "Driver removed from roster", http.StatusOK)
		return
	}
}


