package Driver_Allocation

import (
	"encoding/json"
	"enterprise_computing_cw/Error_Management"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)


// for the list of current drivers in the roster
type Roster []struct {
	Id int `json:"id"`
	DriverName string `json:"DriverName"`
	Rate float64 `json:"Rate"`
}

// for the best driver selected from the roster
type driver struct {
	Id int `json:"id"`
	DriverName string `json:"DriverName"`
	Rate float64 `json:"Rate"`
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

	var directions Response
	err = json.Unmarshal(body, &directions)
	if err != nil {
		return 0, err
	}

	err = resp.Body.Close()
	if err != nil {
		return 0, err
	}

	if len(directions.Routes) < 1 {
		return 0, err
	}


	legs := directions.Routes[0].Legs
	numARoads, totalDistance := 0, float64(0)

	for i, s := range legs[0].Steps {
		distance, instructions := s.Distance.Text, s.HtmlInstructions
		isARoad := instructionsHelper(instructions)
		if isARoad {
			numARoads += 1
		}

		roadDistance, err := distanceHelper(distance)
		if err != nil {
			return 0, err
		}
		totalDistance += roadDistance
		if i == len(legs[0].Steps) - 1 {

			if numARoads == 0 {
				return driverRate * totalDistance, nil
			}

			if i / numARoads < 2 {
				return driverRate * totalDistance * 2, nil
			}
			if i / numARoads >= 2 {
				return driverRate * totalDistance, nil
			}
		}
	}

	return 0, err


}

// parses the google maps distance string into an float value in km
func distanceHelper(distance string) (float64, error) {
	var number float64
	for pos, char := range distance {
		if string(char) == "m" || string(char) == "k" {
			i, err := strconv.ParseFloat(distance[0: pos - 1], 64)
			if err != nil {
				return 0, err
			}
			if string(char) == "m" {
				i = i / 1000
			}
			number = i
			break
		}
	}
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


	body, err := http.Get("http://localhost:10000/rosters")
	if err != nil {
		return nil, err
	}

	var roster Roster
	err = json.NewDecoder(body.Body).Decode(&roster)
	if err != nil {
		return  nil, err
	}
	if len(roster) < 1 {
		return  nil, err
	}

	currBest := math.Inf(1)
	var driverIndex int

	for i, driver := range roster {
		routePrice, err := getSurgePricingRouteHandler(origin, destination, driver.Rate)
		if err != nil || routePrice == 0 {
			return nil, err
		}
		currPrice := getSurgePricingTimeHandler(routePrice)

		if len(roster) < 5 {
			currPrice *= 2
		}

		if currPrice < currBest {
			currBest = currPrice
			driverIndex = i
		}
	}

	err = body.Body.Close()
	if err != nil {
		return nil, err
	}
	// convert pence/km to pound/km
	roster[driverIndex].Rate = currBest / 100
	bestDriver := driver {
		Id: roster[driverIndex].Id,
		DriverName: roster[driverIndex].DriverName,
		Rate: roster[driverIndex].Rate,
	}

	return &bestDriver, nil

}

// http handler for the allocation route - getting the best driver for the trip
func GetBestDriverHandler(w http.ResponseWriter, r *http.Request){


	var trip Error_Management.Trip
	err := json.NewDecoder(r.Body).Decode(&trip)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	model, e := Error_Management.FormValidationHandler(trip)
	if e.ResponseCode != http.StatusOK {
		w.WriteHeader(e.ResponseCode)
		err := json.NewEncoder(w).Encode(e)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	m := *model
	trip = m.(Error_Management.Trip)

	bestDriver, err := getSurgePricingRosterHandler(trip.Origin, trip.Destination)

	if err != nil {

		w.WriteHeader(http.StatusNoContent)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	if bestDriver == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = json.NewEncoder(w).Encode(bestDriver)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

}

