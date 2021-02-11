package trip_mangement

import (
	"context"
	"encoding/json"
	"googlemaps.github.io/maps"
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
	route, err := GetDirections(directions.Origin, directions.Destination, "AIzaSyB2rJrmiL6i3APBb-IMOoykhj8IYqiWc6k")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.NewEncoder(w).Encode(route)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}