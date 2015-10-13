package main

import (
	"fmt"

	"github.com/kellydunn/golang-geo"
)

func Filter(store *Store, mapquestApiKey string, origin *geo.Point, r float64) {
	// Fetch all prospects
	prospects := []Prospect{}
	store.DB.Where("LENGTH(location) > 0").Find(&prospects)

	//
	geo.SetMapquestAPIKey(mapquestApiKey)
	geocoder := &geo.MapQuestGeocoder{}
	for _, p := range prospects {
		point, err := geocoder.Geocode(p.Location)
		if err != nil {
			continue
		}
		d := origin.GreatCircleDistance(point)
		if d < r {
			fmt.Println(p)
		}
	}
}
