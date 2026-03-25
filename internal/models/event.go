package models

import "time"

type APIStatusEvent struct {
	APIID        string    `json:"api_id"`
	URL          string    `json:"url"`
	Status       string    `json:"status"`
	ResponseTime int64     `json:"response_time"`
	CheckedAt    time.Time `json:"checked_at"`
}
