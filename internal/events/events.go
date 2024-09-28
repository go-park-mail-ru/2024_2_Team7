package events

type EventDB struct {
	Events []Event
}

func NewEventDB() *EventDB {
	eventsFeed := createEventMapWithDefaultValues()

	return &EventDB{
		Events: eventsFeed,
	}
}

func createEventMapWithDefaultValues() []Event {
	events := []Event{
		{
			ID:          "1",
			Title:       "Выставка в Третьяковской галерее",
			DateStart:   "2024-09-27",
			DateEnd:     "2024-10-3",
			Tag:         []string{"Popular"},
			Description: "Творчество Малевича",
			ImageURL:    "/static/images/Events1.jpg",
		},
		{
			ID:          "2",
			Title:       "Экскурсии по центру",
			DateStart:   "2024-09-26",
			DateEnd:     "2024-11-5",
			Tag:         []string{"Popular"},
			Description: "Прогулка по Кремлю",
			ImageURL:    "/static/images/Events2.jpg",
		},
		{
			ID:          "3",
			Title:       "Концерт в Концертном зале",
			DateStart:   "2024-10-15",
			DateEnd:     "2024-10-15",
			Tag:         []string{"Popular"},
			Description: "Музыкальное выступление известной группы",
			ImageURL:    "/static/images/PetrovConcert.jpg",
		},
		{
			ID:          "4",
			Title:       "Фестиваль искусств",
			DateStart:   "2024-12-01",
			DateEnd:     "2024-12-05",
			Tag:         []string{"Festival"},
			Description: "Мультимедиальный фестиваль современного искусства",
			ImageURL:    "/static/images/SemenovArtFestival.jpg",
		},
	}
	return events
}
