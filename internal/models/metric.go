package models

import "time"

type Metric struct {
	APIName   string    `bson:"api_name"`
	Status    int       `bson:"status"`
	Latency   int64     `bson:"latency"`
	Timestamp time.Time `bson:"timestamp"`
}
