package collector

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"net/http/httputil"
)

type IanRecordCollector struct {
	mealGaugeVec   prometheus.GaugeVec
	costCountVec   prometheus.CounterVec
	weightGaugeVec prometheus.GaugeVec
	client         http.Client
	address        string
	interval       int
}

func (i *IanRecordCollector) GetData() error {
	return nil
}

// Collect 收集指标
func (i *IanRecordCollector) Collect(metrics chan<- prometheus.Metric) {
	var (
		err error
		t   = &TRecord{}
	)
	defer func() { i.mealGaugeVec.Collect(metrics) }()
	request, err := http.NewRequest("GET", i.address, nil)
	if err != nil {
		fmt.Println("Request failed:", err)
		i.mealGaugeVec.Collect(metrics)
		return
	}

	resp, err := i.client.Do(request)
	if resp == nil {
		fmt.Println("Request failed:", err)
		i.mealGaugeVec.Collect(metrics)
		return
	}
	if resp.StatusCode != 200 {
		fmt.Println("Request failed:", err)
	}
	res, err := httputil.DumpResponse(resp, true)
	fmt.Printf("%s\n", res)
	err = json.Unmarshal(res, t)
	if err != nil {
		fmt.Println("Unmarshal failed:", err)
		return
	}
	if t.Weight > 0 {
		i.weightGaugeVec.WithLabelValues(TimePeriod(int64(t.UpdateTime))).Set(t.Weight)
	}
	if t.Cost > 0 {
		i.costCountVec.WithLabelValues().Add(float64(t.Cost))
	}

}

// Describe 向prometheus注册指标
func (i *IanRecordCollector) Describe(descs chan<- *prometheus.Desc) {
	i.mealGaugeVec.Describe(descs)
}

func NewIanRecordCollector(address string) *IanRecordCollector {
	return &IanRecordCollector{
		mealGaugeVec: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:        "ian_test_meal",
			Help:        "help ian control the energy input",
			ConstLabels: prometheus.Labels{},
		}, []string{"type", "cost"}),
		costCountVec: SumMoney,
		weightGaugeVec: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "weight",
			Help: "Ian's weight records, divided into morning, afternoon, and evening measurements",
		}, []string{"Late Night", "Dawn", "Morning", "Noon", "Afternoon", "Evening", "Late Night"}),

		client:  http.Client{},
		address: address,
	}
}

type TRecordToMetrics struct {
}

var Weight = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "weight",
	Help: "Ian's weight records, divided into morning, afternoon, and evening measurements",
}, []string{"time"})

var SumMoney = *prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "sum_money",
	Help: "Ian's sum money records",
}, []string{"time"})
