package trip_mangement

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
	"io/ioutil"
	"net/http"
)

type Directions struct {
	Origin string `json:"Origin"`
	Destination string `json:"Destination"`
}


func GetDirections(origin string, destination string, apikey string) ([]maps.Route, error) {
	c, err := maps.NewClient(maps.WithAPIKey(apikey))
	if err != nil {
		return nil, err
	}
	r := &maps.DirectionsRequest{
		Origin: origin,
		Destination: destination,
	}
	route, _, err := c.Directions(context.Background(), r)
	if err != nil {
		return nil, err
	}
	return route, nil

}


func DirectionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var directions Directions
	err := json.NewDecoder(r.Body).Decode(&directions)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// try get request ------------------------------------------------------------------------

	resp, err := http.Get("https://maps.googleapis.com/maps/api/directions/json?origin="+
		directions.Origin+"&destination="+directions.Destination+"&key=" + "AIzaSyB2rJrmiL6i3APBb-IMOoykhj8IYqiWc6k")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	type dir struct {
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


	var test dir

	err = json.Unmarshal(body, &test)
	if err != nil {
		fmt.Println(err.Error())
		//invalid character '\'' looking for beginning of object key string
	}

	//err = json.NewDecoder(resp.Body).Decode(&test)
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	fmt.Println("bad")
	//	return
	//}
	pretty.Print(test.Routes[0].Legs)

	// ------------------------------------------------------------------------

	route, err := GetDirections(directions.Origin, directions.Destination, "AIzaSyB2rJrmiL6i3APBb-IMOoykhj8IYqiWc6k")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.NewEncoder(w).Encode(route)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}