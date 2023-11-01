package collector

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"net/http/httputil"
)

type IanRecordCollector struct {
	mealGaugeVec prometheus.GaugeVec
	client       http.Client
	address      string
}

func (i *IanRecordCollector) Describe(descs chan<- *prometheus.Desc) {
	i.mealGaugeVec.Describe(descs)
}

func (i *IanRecordCollector) Collect(metrics chan<- prometheus.Metric) {
	var (
		err error
		t   = &TRecord{}
	)
	i.mealGaugeVec.WithLabelValues("BF").Set(0)
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
	}
	if t.BF != "" {
		i.mealGaugeVec.WithLabelValues("BF").Set(1)
	} else {
		i.mealGaugeVec.WithLabelValues("BF").Set(0)
	}
	if t.LUN != "" {
		i.mealGaugeVec.WithLabelValues("LUN").Set(1)
	} else {
		i.mealGaugeVec.WithLabelValues("LUN").Set(0)
	}
	if t.DIN != "" {
		i.mealGaugeVec.WithLabelValues("DIN").Set(1)
	} else {
		i.mealGaugeVec.WithLabelValues("DIN").Set(0)
	}
	i.mealGaugeVec.Collect(metrics)

}

func NewIanRecordCollector(address string) *IanRecordCollector {
	return &IanRecordCollector{
		mealGaugeVec: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "ian_test_meal",
			Help: "help ian control the energy input",
		}, []string{"time"}),
		client:  http.Client{},
		address: address,
	}
}
