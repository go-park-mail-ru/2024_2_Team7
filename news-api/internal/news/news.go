package news

import (
	"encoding/json"
	"net/http"
)

var newsFeed = []News{
	{ID: "1", Title: "News Title 1", Date: "2024-09-27", Description: "Description 1", ImageURL: "/static/images/news1.jpg"},
	{ID: "2", Title: "News Title 2", Date: "2024-09-26", Description: "Description 2", ImageURL: "/static/images/news2.jpg"},
}

func GetNewsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(newsFeed)
}
