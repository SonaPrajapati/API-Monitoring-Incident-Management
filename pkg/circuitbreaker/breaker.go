package circuitbreaker

import (
	"time"

	"github.com/sony/gobreaker"
)

var Breaker *gobreaker.CircuitBreaker

func InitBreaker() {

	settings := gobreaker.Settings{
		Name: "API-Monitor",

		MaxRequests: 3,

		Interval: 60 * time.Second,

		Timeout: 30 * time.Second,

		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
	}

	Breaker = gobreaker.NewCircuitBreaker(settings)
}
