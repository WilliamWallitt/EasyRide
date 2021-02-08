package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

var dbName string = "roster_management"

//"CREATE TABLE IF NOT EXISTS roster (id INTEGER PRIMARY KEY, Drivername TEXT, Rate INTEGER )"


func CreateRoster(dbName string, query string) *sql.DB {
	db, err := sql.Open("sqlite3", "./" + dbName)
	if err != nil {
		log.Println(err)
	}
	statement, _ := db.Prepare(query)
	_, err = statement.Exec()
	if err != nil {
		log.Println(err)
	}
	return db
}

func AddDriverToRoster(dbName string, driver Roster) {
	db, err := sql.Open("sqlite3", "./" + dbName)
	if err != nil {
		log.Println(err)
	}
	statement, err := db.Prepare("INSERT INTO roster (id, Drivername, Rate) VALUES (?, ?, ?)")
	if err != nil {
		log.Println(err)
	}
	_, err = statement.Exec(driver.Id, driver.DriverName, driver.Rate)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Added driver to roster")
	}

}

func RemoveDriverFromRoster(dbName string, id int) {

	db, err := sql.Open("sqlite3", "./" + dbName)
	if err != nil {
		log.Println(err)
	}
	statement, err := db.Prepare("DELETE FROM roster WHERE id = ?")
	if err != nil {
		log.Println(err)
	}
	_, err = statement.Exec(id)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Deleted driver from roster")
	}

}

func GetAllDriversFromRoster(dbName string) []Roster {
	db, err := sql.Open("sqlite3", "./" + dbName)
	if err != nil {
		log.Println(err)
	}
	statement, err := db.Prepare("SELECT id, Drivername, Rate FROM roster")
	if err != nil {
		log.Println(err)
	}
	var id, rate int
	var driverName string
	var roster []Roster

	rows, _ := statement.Query()
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

	err := json.NewEncoder(w).Encode(GetAllDriversFromRoster(dbName))
	if err != nil {
		log.Println(err)
	}

}

func addDriverToRosterHandler(w http.ResponseWriter, r *http.Request) {
	var driver Roster
	err := json.NewDecoder(r.Body).Decode(&driver)
	if err != nil {
		log.Println(err)
	}

	AddDriverToRoster(dbName, driver)
}


func removeDriverFromRoster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driver_id := vars["id"]
	int, err := strconv.Atoi(driver_id)
	if err != nil {
		fmt.Println(err)
		return
	}
	RemoveDriverFromRoster(dbName, int)
}


