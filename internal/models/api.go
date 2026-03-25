package models

type API struct {
	ID       string `bson:"_id,omitempty"`
	Name     string `bson:"name"`
	URL      string `bson:"url"`
	Interval int    `bson:"interval"`
}
