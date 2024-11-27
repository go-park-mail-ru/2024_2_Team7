package models

type SearchParams struct {
	Query        string
	EventStart   string
	EventEnd     string
	Tags         []string
	Category     int
	LatitudeMin  float64
	LatitudeMax  float64
	LongitudeMin float64
	LongitudeMax float64
}
