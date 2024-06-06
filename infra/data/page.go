package data

import "time"

type Page struct {
	Url         string         `json:"url" bson:"url" db:"url"`
	Links       []string       `json:"links" bson:"links" db:"links"`
	Title       string         `json:"title" bson:"title" db:"title"`
	Description string         `json:"description" bson:"description" db:"description"`
	Meta        []string       `json:"meta" bson:"meta" db:"meta"`
	Visited     bool           `json:"visited" bson:"visited" db:"visited"`
	Timestamp   time.Time      `json:"timestamp" bson:"timestamp" db:"timestamp"`
	Words       map[string]int `json:"words" bson:"words" db:"words"`
}
type PageVisited struct {
	Url     string `json:"url" bson:"url"`
	Visited bool   `json:"visited" bson:"visited"`
}

type PageSearch struct {
	Url   string `json:"url" bson:"url"`
	Title string `json:"title" bson:"title"`
}
type PageSearchWithFrequency struct {
	PageSearch
	Frequency int `json:"frequency" bson:"frequency"`
}
