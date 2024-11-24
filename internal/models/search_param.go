package models

type SearchParams struct {
	Query      string
	EventStart string
	EventEnd   string
	Tags       []string
	Category   int
}
