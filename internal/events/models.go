package events

type Event struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	DateStart   string   `json:"date_start"`
	DateEnd     string   `json:"date_end"`
	Tag         []string `json:"tag"`
	Description string   `json:"description"`
	ImageURL    string   `json:"image"`
}

type EventDB struct {
	Events []Event
}

type Handler struct {
	EventDB EventDB
}
