package Driver_Allocation

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)


type Response []struct{
	Legs []struct {
		Steps []struct {
			HtmlInstructions string `json:"html_instructions"`
			Distance struct {
				Text string `json:"text"`
			}
		}
	}

}

type Roster []struct {
	Id int `json:"id"`
	DriverName string `json:"DriverName"`
	Rate float64 `json:"Rate"`
}


func getSurgePricingTimeHandler(currPrice float64) float64 {
	hour := time.Now().Hour()
	if hour > 23 || hour < 6 {
		return currPrice * 2
	}
	return currPrice
}

func getSurgePricingRouteHandler(origin string, destination string, driverRate float64) (float64, error) {

	postBody, _ := json.Marshal(map[string]string {
			"Origin": origin,
			"Destination": destination,
	})

	responseBody := bytes.NewBuffer(postBody)

	body, err := http.Post("http://localhost:10000/directions", "application/json", responseBody)
	if err != nil {
		return 0, err
	}

	var res Response
	err = json.NewDecoder(body.Body).Decode(&res)
	if err != nil {
		return 0, err
	}

	legs := res[0].Legs[0]
	numARoads, totalDistance := 0, float64(0)


	for i, s := range legs.Steps {
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
		if i == len(legs.Steps) - 1 {

			if i / numARoads < 2 {
				return driverRate * totalDistance * 2, nil
			}
			if i / numARoads >= 2 {
				return driverRate * totalDistance, nil
			}
		}
	}

	err = body.Body.Close()
	if err != nil {
		return 0, err
	}
	return 0, err

}


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

func instructionsHelper(instructions string) bool {
	substrings := []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9"}

	for _, substr := range substrings {
		if strings.Contains(instructions, substr) {
			return true
		}
	}
	return false

}


func getSurgePricingRosterHandler(origin string, destination string) (struct {
	Id         int    `json:"id"`
	DriverName string `json:"DriverName"`
	Rate       float64    `json:"Rate"`
}, error) {
	body, err := http.Get("http://localhost:10000/rosters")
	if err != nil {
		return struct {
	Id         int    `json:"id"`
	DriverName string `json:"DriverName"`
	Rate       float64    `json:"Rate"`
	}{0, "", 0}, err
	}

	var roster Roster
	err = json.NewDecoder(body.Body).Decode(&roster)
	if err != nil {
		return roster[0], err
	}

	currBest := math.Inf(1)
	var driverIndex int
	for i, driver := range roster {
		routePrice, err := getSurgePricingRouteHandler(origin, destination, float64(driver.Rate))
		if err != nil {
			return roster[0], err
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
		return roster[0], err
	}
	roster[driverIndex].Rate = currBest
	return roster[driverIndex], nil

}


func GetBestDriver(w http.ResponseWriter, r *http.Request) {
	var journee struct {
		Origin string `json:"Origin"`
		Destination string `json:"Destination"`
	}

	err := json.NewDecoder(r.Body).Decode(&journee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bestDriver, err := getSurgePricingRosterHandler(journee.Origin, journee.Destination)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(bestDriver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		http.Error(w, "", http.StatusOK)
	}

}

