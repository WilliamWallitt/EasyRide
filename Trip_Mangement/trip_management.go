package trip_mangement

import (
	"context"
	"encoding/json"
	"googlemaps.github.io/maps"
	"log"
	"net/http"
)

type Directions struct {
	Origin string `json:"Origin"`
	Destination string `json:"Destination"`
}



func Getdirections(origin string, destination string, apikey string) []maps.Route {
	c, err := maps.NewClient(maps.WithAPIKey(apikey))
	if err != nil {
		log.Fatal("Error", err)
	}
	r := &maps.DirectionsRequest{
		Origin: origin,
		Destination: destination,
	}
	route, _, err := c.Directions(context.Background(), r)
	if err != nil {
		log.Fatal("Error", err)
	}
	return route

}


func DirectionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var directions Directions
	err := json.NewDecoder(r.Body).Decode(&directions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = json.NewEncoder(w).Encode(Getdirections(directions.Origin, directions.Destination, "AIzaSyB2rJrmiL6i3APBb-IMOoykhj8IYqiWc6k"))
	if err != nil {
		log.Fatal(err)
	}
}