package collector

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"net/http/httputil"
	"time"
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

func NewWatch() {
	var (
		client = http.Client{}
	)
	interval := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-interval.C:
			request, err := http.NewRequest("GET", "http://example.com", nil)
			if err != nil {
				fmt.Println("Request failed:", err)
				continue
			}
			resp, err := client.Do(request)
			res, err := httputil.DumpResponse(resp, true)
			fmt.Printf("%s\n", res)
			err = json.Unmarshal(res, &TRecord{})
			if err != nil {
				fmt.Println("Unmarshal failed:", err)
			}
		}
	}

}

func toExporter(t *TRecord) {
	weightGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "weight",
		Help:        "Ian's weight records, divided into morning, afternoon, and evening measurements",
		ConstLabels: prometheus.Labels{},
	})
	weightVecGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{}, []string{})
	weightGauge.Set(float64(t.Weight))
	weightVecGauge.WithLabelValues("morning").Set(float64(t.Weight))
	prometheus.MustRegister(weightGauge)
	prometheus.MustRegister(weightVecGauge)
}

func init() {
	//toExporter(&TRecord{
	//	Weight: 70.0,
	//})
	//prometheus.MustRegister(DinnerCount)
}

func NewIanExporter(address string) *IanCollector {
	return &IanCollector{}
}
