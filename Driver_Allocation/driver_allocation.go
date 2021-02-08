//package Driver_Allocation

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
	Rate int `json:"Rate"`
}


func getSurgePricingTimeHandler(currPrice float64) float64 {
	hour := time.Now().Hour()
	if hour > 23 || hour < 6 {
		fmt.Println(strconv.Itoa(hour) + ":00, Surge pricing added")
		return currPrice * 2
	}
	fmt.Println(strconv.Itoa(hour) + ":00, Surge pricing not added")
	return currPrice
}

func getSurgePricingRouteHandler(origin string, destination string, driverRate float64) float64 {

	postBody, _ := json.Marshal(map[string]string {
			"Origin": origin,
			"Destination": destination,
	})

	responseBody := bytes.NewBuffer(postBody)

	body, err := http.Post("http://localhost:10000/directions", "application/json", responseBody)
	if err != nil {
		log.Println(err)
	}

	var res Response
	err = json.NewDecoder(body.Body).Decode(&res)
	if err != nil {
		fmt.Println(err)
	}

	legs := res[0].Legs[0]
	numARoads, totalDistance := 0, float64(0)


	for i, s := range legs.Steps {
		distance, instructions := s.Distance.Text, s.HtmlInstructions
		isARoad := instructionsHelper(instructions)
		if isARoad {
			numARoads += 1
		}
		roadDistance := distanceHelper(distance)
		totalDistance += roadDistance
		if i == len(legs.Steps) - 1 {

			if i / numARoads < 2 {
				fmt.Println("Majority A roads")
				fmt.Println("Total distance: ", totalDistance, "km")
				return driverRate * totalDistance * 2
			}
			if i / numARoads >= 2 {
				fmt.Println("Minority A roads")
				fmt.Println("Total distance: ", totalDistance, "km")
				return driverRate * totalDistance
			}
		}
	}

	defer body.Body.Close()

	return 0

}


func distanceHelper(distance string) float64 {
	var number float64
	for pos, char := range distance {
		if string(char) == "m" || string(char) == "k" {
			i, err := strconv.ParseFloat(distance[0: pos - 1], 64)
			if err != nil {
				log.Println(err)
			}
			if string(char) == "m" {
				i = i / 1000
			}
			number = i
			break
		}
	}
	return number

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


func getSurgePricingRosterHandler(origin string, destination string) int {
	body, err := http.Get("http://localhost:10000/rosters")
	if err != nil {
		log.Println(err)
	}

	var roster Roster
	err = json.NewDecoder(body.Body).Decode(&roster)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(roster)
	}

	currBest := math.Inf(1)
	var driverIndex int
	for i, driver := range roster {
		currPrice := getSurgePricingTimeHandler(getSurgePricingRouteHandler(origin, destination, float64(driver.Rate)))
		if len(roster) < 5 {
			currPrice *= 2
		}
		if currPrice < currBest {
			currBest = currPrice
			driverIndex = i
		}

	}

	fmt.Println(currBest, roster[driverIndex])

	defer body.Body.Close()
	return 0
}




func main() {

	getSurgePricingRosterHandler("London", "Manchester")
}