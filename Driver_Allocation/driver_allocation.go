package main

import (
	"bytes"
	"encoding/json"
	"app/Libraries/Error_Management"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)


// for the list of current drivers in the roster
type Roster struct {
	DriverName string `json:"DriverName"`
	Rate float64 `json:"Rate"`
	Total int `json:"Total"`
}

// for the best driver selected from the roster
type driver struct {
	DriverName string `json:"DriverName"`
	Rate float64 `json:"Rate"`
	Price float64 `json:"Price"`
}

type Mapping struct {
	MajorityARoads bool `json:"MajorityARoads"`
	TotalKilometres float64 `json:"TotalKilometres"`
	TimeSurgePricing bool `json:"TimeSurgePricing"`
}


// combines the time surge pricing handler and route surge pricing handler to calculate the final rate
// for each driver, and returns a sturct of the best driver (id, driver name, final rate)
func getSurgePricingRosterHandler(origin string, destination string) (*driver, error) {

	// get request to the rosters service
	body, err := http.Get("http://host.docker.internal:3001/rosters")
	// check an error occured making the request
	if err != nil {
		return nil, err
	}
	// decode the best driver into the Roster struct
	var roster Roster
	err = json.NewDecoder(body.Body).Decode(&roster)
	if err != nil {
		return  nil, err
	}

	// get the route price (applies route surge pricing if applicable)
	var trip_mapping Mapping

	requestBody, err := json.Marshal(map[string]string{
		"Origin" : origin,
		"Destination": destination,
	})

	if err != nil {
		fmt.Println(err, "1")
		return nil, err
	}


	// MAPPING POST REQ
	response, err := http.Post("http://host.docker.internal:3003/mapping", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, err
	}

	err =  json.NewDecoder(response.Body).Decode(&trip_mapping)
	if err != nil {
		fmt.Println(err, "3")
		return nil, err
	}

	err = response.Body.Close()
	if err != nil {
		return nil, err
	}

	currPrice := trip_mapping.TotalKilometres * roster.Rate

	// add the time surge price if applicable
	if trip_mapping.TimeSurgePricing {
		currPrice *= 2
	}

	if trip_mapping.MajorityARoads {
		currPrice *= 2
	}

	// add the roster surge price if applicable
	if roster.Total < 5 {
		currPrice *= 2
	}

	// close request
	err = body.Body.Close()
	if err != nil {
		fmt.Println(err, "4")
		return nil, err
	}

	// convert pence/km to pound/km and round to 2 decimal places
	price, err := strconv.ParseFloat(fmt.Sprintf("%.2f", currPrice / 100), 64)
	if err != nil {
		fmt.Println(err, "5")
		return nil, err
	}

	// create a Driver struct with the driver's name, rate and total price of trip
	bestDriver := driver{
		DriverName: roster.DriverName,
		Rate:       roster.Rate,
		Price:      price,
	}
	// return the Driver struct
	return &bestDriver, nil

}

// http handler for the allocation route - getting the best driver for the trip
func GetBestDriverHandler(w http.ResponseWriter, r *http.Request){

	// decode the user's json request into the Trip struct (origin and destination)
	var trip Error_Management.Trip
	err := json.NewDecoder(r.Body).Decode(&trip)
	// handle decoding error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	// form validation, check that origin and destination fields are correctly filled in

	model, e := Error_Management.FormValidationHandler(trip)
	if e.ResponseCode != http.StatusOK {
		w.WriteHeader(e.ResponseCode)
		err := json.NewEncoder(w).Encode(e)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	// as the form validation helper function returns an interface we need to convert
	// it back into the correct struct

	m := *model
	trip = m.(Error_Management.Trip)

	// apply all surge pricing if applicable and return the trip price
	bestDriver, err := getSurgePricingRosterHandler(trip.Origin, trip.Destination)
	// handle any errors that might of occured in the function
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// handle if the best driver is nil (there are no drivers in the roster)
	if bestDriver == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// otherwise encode the best driver information as the response
	err = json.NewEncoder(w).Encode(bestDriver)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}


func main () {

	// trailing slash is allowed for any route ie /allocation/ allowed and allocation/ allowed
	allocationRouter := mux.NewRouter().StrictSlash(true)

	// get best driver (POST)
	//curl -H "Content-Type: application/json" -X POST -d '{"Origin":"London","Destination":"Exeter"}' http://localhost:3002/allocation

	allocationRouter.HandleFunc("/allocation", GetBestDriverHandler).Methods("POST")
	err := http.ListenAndServe(":8082", allocationRouter)
	if err != nil {
		log.Fatal(err)
	}

}

