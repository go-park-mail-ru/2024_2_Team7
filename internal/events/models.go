package events

type Event struct {
	ID          string
	Title       string
	DateStart   string
	DateEnd     string
	Tag         []string // чтобы в дальнейшем делать подборки по тегам
	Description string
	ImageURL    string
}
