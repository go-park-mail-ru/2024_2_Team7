package models

type Test struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Questions []Question `json:"question"`
}

type Question struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}
