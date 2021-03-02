package main

import (
	"encoding/json"
	"app/Libraries/Error_Management"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	Price float64
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

// checks the current time and returns a boolean (True if surge pricing applies, False otherwise)
func getSurgePricingTimeHandler(currPrice float64) float64 {
	hour := time.Now().Hour()
	if hour > 23 || hour < 6 {
		return currPrice * 2
	}
	return currPrice
}

// function that uses the origin and destination with the google maps API
// to extract the total distance and pricing (with / without surge pricing) for a driver
func getSurgePricingRouteHandler(origin string, destination string, driverRate float64) (float64, error) {

	resp, err := http.Get("https://maps.googleapis.com/maps/api/directions/json?units=metric&region=UK&origin="+
		origin+"&destination="+destination+"&key=" + "AIzaSyB2rJrmiL6i3APBb-IMOoykhj8IYqiWc6k")

	if err != nil {
		return 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if body == nil {
		return 0, nil
	}

	// we are taking our json maps response and extracting the information we need into the Response struct
	var directions Response
	err = json.Unmarshal(body, &directions)
	if err != nil {
		return 0, err
	}

	err = resp.Body.Close()
	if err != nil {
		return 0, err
	}

	// check that there are directions to extract the distance and instructions from
	if len(directions.Routes) < 1 {
		return 0, err
	}

	// get current legs of the route
	legs := directions.Routes[0].Legs
	// storing the number of A roads and total distance
	numARoads, totalDistance := 0, float64(0)

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
			return 0, err
		}
		// add that distance to the total distance
		totalDistance += roadDistance
		// if we have reached the end of our steps
		if i == len(legs[0].Steps) - 1 {
			// if we have no A roads, no surge pricing
			if numARoads == 0 {
				return driverRate * totalDistance, nil
			}
			// if i divided by the number of A roads is less than 2, then majority A roads
			// surge pricing applies
			if i / numARoads < 2 {
				return driverRate * totalDistance * 2, nil
			}
			// otherwise surge pricing doesnt apply
			if i / numARoads >= 2 {
				return driverRate * totalDistance, nil
			}
		}
	}
	// in case we don't return anything in the above for loop
	return 0, err

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

// combines the time surge pricing handler and route surge pricing handler to calculate the final rate
// for each driver, and returns a sturct of the best driver (id, driver name, final rate)
func getSurgePricingRosterHandler(origin string, destination string) (*driver, error) {

	// get request to the rosters service
	body, err := http.Get("http://host.docker.internal:3002/rosters")
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

	routePrice, err := getSurgePricingRouteHandler(origin, destination, roster.Rate)
	if err != nil || routePrice == 0 {
		return nil, err
	}

	// add the time surge price if applicable

	currPrice := getSurgePricingTimeHandler(routePrice)

	// add the roster surge price if applicable
	if roster.Total < 5 {
		currPrice *= 2
	}

	// close request
	err = body.Body.Close()
	if err != nil {
		return nil, err
	}

	// convert pence/km to pound/km and round to 2 decimal places
	price, err := strconv.ParseFloat(fmt.Sprintf("%.2f", currPrice / 100), 64)
	if err != nil {
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
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	// as the form validation helper function returns an interface we need to convert
	// it back into the correct struct

	m := *model
	trip = m.(Error_Management.Trip)

	// apply all surge pricing if applicable and return the trip price
	bestDriver, err := getSurgePricingRosterHandler(trip.Origin, trip.Destination)
	// handle any errors that might of occured
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
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
		_ = json.NewEncoder(w).Encode(err)
		return
	}

}


func main () {


	// trailing slash is allowed for any route ie /allocation/ allowed and allocation/ allowed
	authRouter := mux.NewRouter().StrictSlash(true)

	// get best driver (POST)
	//curl -H "Content-Type: application/json" -X POST -d '{"Origin":"London","Destination":"Exeter"}' http://localhost:8083/allocation

	authRouter.HandleFunc("/allocation", GetBestDriverHandler).Methods("POST")
	err := http.ListenAndServe(":8083", authRouter)
	if err != nil {
		log.Fatal(err)
	}

}

