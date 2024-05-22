package data

import "time"

type Page struct {
	Url         string    `json:"url" bson:"url"`
	Links       []string  `json:"links" bson:"links"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Meta        []string  `json:"meta" bson:"meta"`
	Visited     bool      `json:"visited" bson:"visited"`
	Timestamp   time.Time `json:"timestamp" bson:"timestamp"`
}
type PageVisited struct {
	Url     string `json:"url" bson:"url"`
	Visited bool   `json:"visited" bson:"visited"`
}
