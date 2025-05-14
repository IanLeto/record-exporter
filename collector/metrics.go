package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

type IanCollector struct {
	beatGaugeVec   prometheus.GaugeVec
	eventGaugeVec  prometheus.GaugeVec
	handleGaugeVec prometheus.GaugeVec
	harvesterGauge prometheus.GaugeVec

	outputGaugeVec prometheus.GaugeVec
	pipelineGauge  prometheus.GaugeVec

	up      *prometheus.Desc
	client  http.Client
	address string
	current float64
}

type IanMetrics struct {
	DinnerCount int `json:"dinner"`
}

var DinnerCount = prometheus.NewCounter(prometheus.CounterOpts{})

func init() {

}
