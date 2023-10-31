package collector

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"net/http/httputil"
	"time"
)

type IanMetrics struct {
	DinnerCount int `json:"dinner"`
}

var DinnerCount = prometheus.NewCounter(prometheus.CounterOpts{})

type TRecord struct {
	Name       string  `json:"name"`
	Weight     float32 `json:"weight"`
	BF         string  `json:"bf"`
	LUN        string  `json:"lun"`
	DIN        string  `json:"din"`
	EXTRA      string  `json:"extra"`
	Core       int     `json:"core"`
	Runner     int     `json:"runner"`
	Support    int     `json:"support"`
	Squat      int     `json:"squat"`
	EasyBurpee int     `json:"easy_burpee"`
	Chair      int     `json:"chair"`
	Stretch    int     `json:"stretch"`
	Vol1       string  `json:"vol1"`
	Vol2       string  `json:"vol2"`
	Vol3       string  `json:"vol3"`
	Vol4       string  `json:"vol4"`
	Content    string  `json:"content"`
	Region     string  `json:"region"`
}

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
		Name: "weight",
		Help: "Ian's weight records, divided into morning, afternoon, and evening measurements",
		ConstLabels: prometheus.Labels{
			"morning":   t.BF,
			"afternoon": t.LUN,
			"evening":   t.DIN,
		},
	})
	weightVecGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{}, []string{})
	weightGauge.Set(float64(t.Weight))
	weightVecGauge.WithLabelValues("morning").Set(float64(t.Weight))
	prometheus.MustRegister(weightGauge)
	prometheus.MustRegister(weightVecGauge)
}

func init() {
	toExporter(&TRecord{
		Weight: 70.0,
	})
	prometheus.MustRegister(DinnerCount)
}
