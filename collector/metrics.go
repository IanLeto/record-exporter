package collector

import "github.com/prometheus/client_golang/prometheus"

type IanMetrics struct {
	DinnerCount int `json:"dinner"`
}

var DinnerCount = prometheus.NewCounter(prometheus.CounterOpts{})
