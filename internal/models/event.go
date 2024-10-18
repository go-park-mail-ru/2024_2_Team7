package models

type Event struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	DateStart   string   `json:"date_start"`
	DateEnd     string   `json:"date_end"`
	AuthorID    int      `json:"author"`
	Tag         []string `json:"tag"`
	Description string   `json:"description"`
	ImageURL    string   `json:"image"`
}
