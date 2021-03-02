package main

import (
	"app/Libraries/Error_Management"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
	"net/url"
)

type Mapping struct {
	MajorityARoads bool `json:"MajorityARoads"`
	TotalKilometres float64 `json:"TotalKilometres"`
	TimeSurgePricing bool `json:"TimeSurgePricing"`
}

// for the json response of the google maps service
type Response struct {
	Routes []struct {
		Legs []struct{
			Steps []struct {
				Distance struct {
					Text string `json:"text"`
				}
				HtmlInstructions string `json:"html_instructions"`
			}
		}
	}
}



// function that uses the origin and destination with the google maps API
// to extract the total distance and pricing (with / without surge pricing) for a driver
func getSurgePricingRouteHandler(origin string, destination string) (*Mapping, error) {

	// need to handle spaces


	resp, err := http.Get("https://maps.googleapis.com/maps/api/directions/json?units=metric&region=UK&origin="+
		url.QueryEscape(origin)+"&destination="+url.QueryEscape(destination)+"&key=" + "AIzaSyB2rJrmiL6i3APBb-IMOoykhj8IYqiWc6k")

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if body == nil {
		return nil, nil
	}

	// we are taking our json maps response and extracting the information we need into the Response struct
	var directions Response
	err = json.Unmarshal(body, &directions)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// check that there are directions to extract the distance and instructions from
	if len(directions.Routes) < 1 {
		return nil, err
	}

	// create Mapping struct

	// get current legs of the route
	legs := directions.Routes[0].Legs
	// storing the number of A roads and total distance
	numARoads, totalDistance := 0, float64(0)
	
	
	mapping := Mapping{
		MajorityARoads:   false,
		TotalKilometres:  0,
		TimeSurgePricing: false,
	}

	// iterate over each step in legs[0].Steps
	for i, s := range legs[0].Steps {
		distance, instructions := s.Distance.Text, s.HtmlInstructions
		// check if the instructions contain an A road
		isARoad := instructionsHelper(instructions)
		// if there is an A road, add 1 to numARoads
		if isARoad {
			numARoads += 1
		}
		// get current distance in km for that step
		roadDistance, err := distanceHelper(distance)
		if err != nil {
			return nil, err
		}
		// add that distance to the total distance
		totalDistance += roadDistance
		// if we have reached the end of our steps
		if i == len(legs[0].Steps) - 1 {
			// if we have no A roads, no surge pricing
			if numARoads == 0 {
				mapping.MajorityARoads = false
				mapping.TotalKilometres = totalDistance
				mapping.TimeSurgePricing = getSurgePricingTimeHandler()
				return &mapping, nil
			}
			// if i divided by the number of A roads is less than 2, then majority A roads
			// surge pricing applies
			if i / numARoads < 2 {
				mapping.MajorityARoads = true
				mapping.TotalKilometres = totalDistance
				mapping.TimeSurgePricing = getSurgePricingTimeHandler()
				return &mapping, nil
			}
			// otherwise surge pricing doesnt apply
			if i / numARoads >= 2 {
				mapping.MajorityARoads = false
				mapping.TotalKilometres = totalDistance
				mapping.TimeSurgePricing = getSurgePricingTimeHandler()
				return &mapping, nil
			}
		}
	}
	// in case we don't return anything in the above for loop
	return nil, err

}

// parses the google maps distance string into an float value in km
func distanceHelper(distance string) (float64, error) {
	var number float64
	// for each character in the distance string
	for pos, char := range distance {
		// if we get to an "m" or "k" ("km")
		if string(char) == "m" || string(char) == "k" {
			// we will convert the previous chars to an float
			i, err := strconv.ParseFloat(distance[0: pos - 1], 64)
			// handle errors if this doenst work
			if err != nil {
				return 0, err
			}
			// convert m to km
			if string(char) == "m" {
				i = i / 1000
			}
			// store the distance in the number variable
			number = i
			// exit the for loop
			break
		}
	}
	// return the distance
	return number, nil

}

// parses the google maps html instruction string returns a boolean
// True if the road is an A road, False otherwise
func instructionsHelper(instructions string) bool {
	substrings := []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9"}

	for _, substr := range substrings {
		if strings.Contains(instructions, substr) {
			return true
		}
	}
	return false

}


// http handler for the allocation route - getting the best driver for the trip
func GetMappingHandler(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "application/json")

	// decode the user's json request into the Trip struct (origin and destination)
	var trip Error_Management.Trip
	err := json.NewDecoder(r.Body).Decode(&trip)
	// handles if any errors have occured doing this
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

	// get Mapping struct

	trip_mapping, err := getSurgePricingRouteHandler(trip.Origin, trip.Destination)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	if trip_mapping == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Println(trip_mapping)

	// otherwise encode the best driver information as the response
	err = json.NewEncoder(w).Encode(trip_mapping)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

}

// checks the current time and returns a boolean (True if surge pricing applies, False otherwise)
func getSurgePricingTimeHandler() bool {
	hour := time.Now().Hour()
	if hour > 23 || hour < 6 {
		return true
	}
	return false
}

func main() {

	// trailing slash is allowed for any route ie /allocation/ allowed and allocation/ allowed
	tripRouter := mux.NewRouter().StrictSlash(true)

	// get mapping (POST)
	//curl -H "Content-Type: application/json" -X POST -d '{"Origin":"London","Destination":"Exeter"}' http://localhost:8083/mapping

	tripRouter.HandleFunc("/mapping", GetMappingHandler).Methods("POST")

	err := http.ListenAndServe(":8083", tripRouter)
	if err != nil {
		log.Fatal(err)
	}

}



