package main

import (
	"encoding/csv"
	"os"

	"github.com/kellydunn/golang-geo"
)

func Filter(store *Store, mapquestApiKey string, origin *geo.Point, r float64) {
	// Fetch all prospects
	prospects := []Prospect{}
	store.DB.Where("LENGTH(location) > 0").Find(&prospects)

	// Geocode and output
	w := csv.NewWriter(os.Stdout)
	geo.SetMapquestAPIKey(mapquestApiKey)
	geocoder := &geo.MapQuestGeocoder{}
	for _, p := range prospects {
		point, err := geocoder.Geocode(p.Location)
		if err != nil {
			continue
		}
		d := origin.GreatCircleDistance(point)
		if d < r {
			err := w.Write([]string{
				p.Name,
				p.Email,
				p.Location,
				p.User,
				p.Repo,
			})
			if err != nil {
				panic(err)
			}
			w.Flush()
		}
	}
}
