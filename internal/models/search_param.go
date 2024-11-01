package models

type SearchParams struct {
	Str        string
	EventStart string
	EventEnd   string
	Tags       []string
	Category   int
}
